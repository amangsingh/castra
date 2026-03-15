package cli

import (
	"database/sql"
	"strings"
)

type Note struct {
	ID        int64
	ProjectID int64
	TaskID    *int64
	Content   string
	Tags      string
}

func AddNote(db *sql.DB, projectID int64, taskID *int64, content, tags string) (int64, error) {
	query := `INSERT INTO project_notes (project_id, task_id, content, tags) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, projectID, taskID, content, tags)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ListNotes(db *sql.DB, projectID int64, taskID *int64, role string) ([]Note, error) {
	// Build base query
	query := `SELECT id, project_id, task_id, COALESCE(content, ''), COALESCE(tags, '') FROM project_notes WHERE project_id = ? AND deleted_at IS NULL`
	args := []interface{}{projectID}

	// Filter by task if provided
	if taskID != nil {
		query += ` AND task_id = ?`
		args = append(args, *taskID)
	}

	query += ` ORDER BY created_at DESC`

	// Architect and Doc Writer see all notes
	if role == "architect" || role == "doc-writer" {
		return queryNotes(db, query, args...)
	}

	// Others see notes that are either untagged (public) OR tagged for their role.
	// Notes exclusively tagged for a different specific role are hidden.
	allNotes, err := queryNotes(db, query, args...)
	if err != nil {
		return nil, err
	}

	var filtered []Note
	for _, n := range allNotes {
		if isVisibleToRole(n.Tags, role) {
			filtered = append(filtered, n)
		}
	}
	return filtered, nil
}

func queryNotes(db *sql.DB, query string, args ...interface{}) ([]Note, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.ProjectID, &n.TaskID, &n.Content, &n.Tags); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

// knownRoles is the set of role tags that make a note role-specific.
var knownRoles = map[string]bool{
	"architect":       true,
	"senior-engineer": true,
	"junior-engineer": true,
	"qa-functional":   true,
	"security-ops":    true,
	"designer":        true,
	"doc-writer":      true,
}

// isVisibleToRole returns true if a note should be shown to the given role.
// A note is visible if:
//   - it has no role-specific tags (it is public/generic), OR
//   - its tags include the calling role.
func isVisibleToRole(tags, role string) bool {
	if tags == "" {
		return true // no tags = public
	}
	parts := strings.Split(tags, ",")
	hasRoleTag := false
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t == role {
			return true // explicitly tagged for this role
		}
		if knownRoles[t] {
			hasRoleTag = true // tagged for a different role
		}
	}
	// If none of the tags are a known role, the note is generic → public
	return !hasRoleTag
}
