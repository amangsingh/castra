package cli

import (
	"castra/internal/persona"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

// DefaultStatuses is the fallback pipeline for tasks with no archetype.
var DefaultStatuses = []string{"todo", "doing", "review", "done"}

type Archetype struct {
	ID          int64
	ProjectID   *int64 // Optional FK for project-specific archetypes
	Name        string
	Description string
	DefaultRole string
	Statuses    []string // ordered valid statuses (JSON-encoded in DB)
}

// statusesJSON encodes a string slice to JSON for DB storage.
func statusesJSON(statuses []string) string {
	b, _ := json.Marshal(statuses)
	return string(b)
}

// parseStatuses decodes a JSON string from the DB into a string slice.
func parseStatuses(raw string) []string {
	if raw == "" {
		return DefaultStatuses
	}
	var s []string
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return DefaultStatuses
	}
	return s
}

func AddArchetype(db *sql.DB, projectID *int64, name, description, defaultRole string, statuses []string) (int64, error) {
	query := `INSERT INTO archetypes (project_id, name, description, default_role, statuses) VALUES (?, ?, ?, ?, ?)`
	res, err := db.Exec(query, projectID, name, description, defaultRole, statusesJSON(statuses))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetArchetype(db *sql.DB, id int64) (*Archetype, error) {
	var a Archetype
	var raw string
	query := `SELECT id, project_id, name, COALESCE(description, ''), default_role, statuses FROM archetypes WHERE id = ? AND deleted_at IS NULL`
	err := db.QueryRow(query, id).Scan(&a.ID, &a.ProjectID, &a.Name, &a.Description, &a.DefaultRole, &raw)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("archetype not found")
		}
		return nil, err
	}
	a.Statuses = parseStatuses(raw)
	return &a, nil
}

// GetTaskStatuses returns the ordered status pipeline for a task.
// Falls back to DefaultStatuses if no archetype is assigned.
func GetTaskStatuses(db *sql.DB, taskID int64) []string {
	var raw sql.NullString
	err := db.QueryRow(`
		SELECT a.statuses FROM tasks t
		LEFT JOIN archetypes a ON t.archetype_id = a.id
		WHERE t.id = ?`, taskID).Scan(&raw)
	if err != nil || !raw.Valid || raw.String == "" {
		return DefaultStatuses
	}
	return parseStatuses(raw.String)
}

// NextStatus returns the next valid status in the pipeline after current.
// Returns "" if current is the last status or not found.
func NextStatus(statuses []string, current string) string {
	for i, s := range statuses {
		if s == current && i+1 < len(statuses) {
			return statuses[i+1]
		}
	}
	return ""
}

// GetReviewStatus returns the status where QA/Security gates are enforced.
// This is typically the status immediately preceding the terminal (done) status.
func GetReviewStatus(statuses []string) string {
	if len(statuses) < 2 {
		return ""
	}
	return statuses[len(statuses)-2]
}

// ValidateTransition checks if newStatus is a legal transition from current,
// given the archetype's ordered statuses. Rejections (first status) are always allowed.
func ValidateTransition(statuses []string, current, newStatus string) error {
	if len(statuses) == 0 {
		return nil
	}
	firstStatus := statuses[0]
	// Rejection (returning to start) is always valid
	if newStatus == firstStatus {
		return nil
	}
	// Find next valid status
	next := NextStatus(statuses, current)
	if next == "" {
		return fmt.Errorf("task is already in terminal status '%s'", current)
	}
	if newStatus != next {
		validOptions := []string{firstStatus, next}
		return fmt.Errorf("invalid status transition '%s'→'%s': valid options are %s",
			current, newStatus, strings.Join(validOptions, " or "))
	}
	return nil
}

func ListArchetypes(db *sql.DB, projectID *int64) ([]Archetype, error) {
	query := `SELECT id, project_id, name, COALESCE(description, ''), default_role, statuses FROM archetypes WHERE (project_id IS NULL OR project_id = ?) AND deleted_at IS NULL`
	rows, err := db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var archetypes []Archetype
	for rows.Next() {
		var a Archetype
		var raw string
		if err := rows.Scan(&a.ID, &a.ProjectID, &a.Name, &a.Description, &a.DefaultRole, &raw); err != nil {
			return nil, err
		}
		a.Statuses = parseStatuses(raw)
		archetypes = append(archetypes, a)
	}
	return archetypes, nil
}

func SoftDeleteArchetype(db *sql.DB, id int64) error {
	query := `UPDATE archetypes SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetVisibleStatuses returns the list of statuses a role is permitted to see across all archetypes.
func GetVisibleStatuses(db *sql.DB, role string) ([]string, bool, error) {
	linter := persona.NewLinter()

	// Architects and Doc Writers see everything
	if linter.PersonaAudit(role, persona.CapabilityManagement) == nil || role == "doc-writer" {
		return nil, true, nil
	}

	canLifecycle := linter.PersonaAudit(role, persona.CapabilityLifecycle) == nil
	canAudit := linter.PersonaAudit(role, persona.CapabilityAudit) == nil

	// Fetch all unique status sets from archetypes
	rows, err := db.Query(`SELECT DISTINCT statuses FROM archetypes WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	statusMap := make(map[string]bool)
	
	// Process each archetype's pipeline
	processPipeline := func(statuses []string) {
		if len(statuses) == 0 {
			return
		}
		
		// If role has lifecycle capability, they see todo -> doing (everything before review)
		// Default pipeline: ["todo", "doing", "review", "done"]
		// Lifecycle: see [0, len-3]
		if canLifecycle {
			// Always allow visibility of common lifecycle states even if not in the specific pipeline
			statusMap["todo"] = true
			statusMap["doing"] = true
			statusMap["blocked"] = true
			statusMap["pending"] = true
			
			for i := 0; i <= len(statuses)-3; i++ {
				if i >= 0 && i < len(statuses) {
					statusMap[statuses[i]] = true
				}
			}
		}

		// If role has audit capability, they see "review" status
		// Review is index len-2
		if canAudit {
			if len(statuses) >= 2 {
				statusMap[statuses[len(statuses)-2]] = true
			}
		}
	}

	for rows.Next() {
		var raw sql.NullString
		if err := rows.Scan(&raw); err != nil {
			return nil, false, err
		}
		if raw.Valid && raw.String != "" {
			processPipeline(parseStatuses(raw.String))
		}
	}

	// Also process DefaultStatuses for tasks with no archetype
	processPipeline(DefaultStatuses)

	var result []string
	for s := range statusMap {
		result = append(result, s)
	}

	// If no statuses found and not seeAll, it's effectively an unknown/restricted role
	if len(result) == 0 && !canLifecycle && !canAudit {
		return nil, false, fmt.Errorf("unknown role: %s", role)
	}

	return result, false, nil
}
