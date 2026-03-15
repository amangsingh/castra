package cli

import (
	"castra/internal/persona"
	"database/sql"
	"fmt"
	"strings"
)

type Task struct {
	ID               int64
	ProjectID        int64
	MilestoneID      *int64 // Pointer to handle NULL
	SprintID         *int64 // Pointer to handle NULL
	ArchetypeID      *int64 // Pointer to handle NULL
	Title            string
	Description      string
	Status           string
	Priority         string
	QAApproved       bool
	SecurityApproved bool
	QABypassed       bool
	SecurityBypassed bool
}

func AddTask(db *sql.DB, projectID int64, milestoneID, sprintID, archetypeID *int64, title, desc, priority string) (int64, error) {
	// 1. Determine the initial status for this task (AC2)
	initialStatus := "todo"

	// If archetype is not provided, try to inherit from milestone
	if archetypeID == nil && milestoneID != nil {
		var msAid sql.NullInt64
		err := db.QueryRow(`SELECT archetype_id FROM milestones WHERE id = ?`, *milestoneID).Scan(&msAid)
		if err == nil && msAid.Valid {
			id := msAid.Int64
			archetypeID = &id
		}
	}

	if archetypeID != nil {
		var raw sql.NullString
		err := db.QueryRow(`SELECT statuses FROM archetypes WHERE id = ?`, *archetypeID).Scan(&raw)
		if err == nil && raw.Valid && raw.String != "" {
			statuses := parseStatuses(raw.String)
			if len(statuses) > 0 {
				initialStatus = statuses[0]
			}
		}
	}

	query := `INSERT INTO tasks (project_id, milestone_id, sprint_id, archetype_id, title, description, priority, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(query, projectID, milestoneID, sprintID, archetypeID, title, desc, priority, initialStatus)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(db *sql.DB, id int64) (*Task, error) {
	var t Task
	query := `SELECT id, project_id, milestone_id, sprint_id, archetype_id, title, COALESCE(description, ''), status, priority, qa_approved, security_approved, COALESCE(qa_bypassed, 0), COALESCE(security_bypassed, 0) FROM tasks WHERE id = ? AND deleted_at IS NULL`
	err := db.QueryRow(query, id).Scan(&t.ID, &t.ProjectID, &t.MilestoneID, &t.SprintID, &t.ArchetypeID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.QAApproved, &t.SecurityApproved, &t.QABypassed, &t.SecurityBypassed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return &t, nil
}

func ListTasks(db *sql.DB, projectID int64, milestoneID, sprintID *int64, backlogOnly bool, role string) ([]Task, error) {
	query := `SELECT id, project_id, milestone_id, sprint_id, archetype_id, title, COALESCE(description, ''), status, priority, qa_approved, security_approved, COALESCE(qa_bypassed, 0), COALESCE(security_bypassed, 0) FROM tasks WHERE project_id = ? AND deleted_at IS NULL`
	args := []interface{}{projectID}

	if backlogOnly {
		query += ` AND sprint_id IS NULL AND milestone_id IS NULL`
	} else {
		if milestoneID != nil {
			query += ` AND milestone_id = ?`
			args = append(args, *milestoneID)
		}
		if sprintID != nil {
			query += ` AND sprint_id = ?`
			args = append(args, *sprintID)
		}
	}

	// Filter by Role Context (Dynamic)
	visibleStatuses, seeAll, err := GetVisibleStatuses(db, role)
	if err != nil {
		return nil, err
	}
	if !seeAll {
		if len(visibleStatuses) == 0 {
			// If no visible statuses, query a condition that matches nothing
			query += " AND 1=0"
		} else {
			placeholders := make([]string, len(visibleStatuses))
			for i, s := range visibleStatuses {
				placeholders[i] = "?"
				args = append(args, s)
			}
			query += fmt.Sprintf(" AND status IN (%s)", strings.Join(placeholders, ","))
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.MilestoneID, &t.SprintID, &t.ArchetypeID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.QAApproved, &t.SecurityApproved, &t.QABypassed, &t.SecurityBypassed); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func UpdateTaskStatus(db *sql.DB, id int64, newStatus string, desc string, role string, breakGlass bool, reason string) error {
	// --- Persona Audit (Gate 3: The Self) ---
	// Core engine enforcement of the Doctrine of Command.
	linter := persona.NewLinter()
	
	// First-pass: Is this a recognized system persona with any operational capability?
	if err := linter.PersonaAudit(role, persona.CapabilityLifecycle); err != nil {
		if errAudit := linter.PersonaAudit(role, persona.CapabilityAudit); errAudit != nil {
			// Persona Non-Compliance: reset to start of pipeline
			_ = AddAuditEntry(db, "system", 0, "persona_non_compliance", role, "[engine] "+err.Error())
			
			statuses := GetTaskStatuses(db, id)
			initialStatus := "todo"
			if len(statuses) > 0 {
				initialStatus = statuses[0]
			}

			_, _ = db.Exec(`UPDATE tasks SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, initialStatus, id)
			LogTaskAction(db, id, "rejected", "system", fmt.Sprintf("Persona Non-Compliance: Task reset to '%s' due to engine-level jurisdictional violation.", initialStatus))
			return err
		}
	}

	// 1. Fetch current task state
	var currentStatus string
	var qaApp, secApp bool
	var sprintID sql.NullInt64
	err := db.QueryRow(`SELECT status, qa_approved, security_approved, sprint_id FROM tasks WHERE id = ?`, id).Scan(&currentStatus, &qaApp, &secApp, &sprintID)
	if err != nil {
		return err
	}

	// Handle description-only update (no status change requested)
	if newStatus == "" && desc != "" {
		_, err := db.Exec(`UPDATE tasks SET description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, desc, id)
		if err != nil {
			return err
		}
		LogTaskAction(db, id, "description_update", role, "Task description updated.")
		return nil
	}

	// 2. Load Pipeline & Validate Capabilities
	statuses := GetTaskStatuses(db, id)
	if !breakGlass {
		// Only Audit roles (or Architect) can move to "done"
		if (newStatus == "done" || newStatus == "closed") {
			if err := linter.PersonaAudit(role, persona.CapabilityAudit); err != nil {
				return err
			}
		}

		// Only Lifecycle roles (or Architect) can move tasks intentionally
		// (QA/Security can also move to start of pipeline on rejection, which is handled in handleTaskApprovals)
		initialStatus := "todo"
		if len(statuses) > 0 {
			initialStatus = statuses[0]
		}
		if newStatus != "" && newStatus != currentStatus && newStatus != initialStatus {
			if err := linter.PersonaAudit(role, persona.CapabilityLifecycle); err != nil {
				// Special check: Auditors can move to review status OR terminal statuses
				isTerminal := newStatus == "done" || newStatus == "closed"
				if (newStatus == "review" || isTerminal) && linter.PersonaAudit(role, persona.CapabilityAudit) == nil {
					// Permitted
				} else {
					return err
				}
			}
		}

		// Archetype Role Enforcement: check if task is claimed by the correct role
		if currentStatus == "todo" && newStatus == "doing" {
			var defaultRole sql.NullString
			_ = db.QueryRow(`SELECT a.default_role FROM tasks t JOIN archetypes a ON t.archetype_id = a.id WHERE t.id = ?`, id).Scan(&defaultRole)
			if defaultRole.Valid && defaultRole.String != "" && role != defaultRole.String {
				return fmt.Errorf("role mismatch: this task is designated for %s archetypes, but you are %s", defaultRole.String, role)
			}
		}
	}

	// 2b. Validate archetype status transition (Task 80)
	// Architects are exempt — they are only blocked at the dual-gate approval level.
	if !breakGlass && role != "architect" {
		if err := ValidateTransition(statuses, currentStatus, newStatus); err != nil {
			return err
		}
	}

	// 3. Handle Approval Logic (Gates)
	targetStatus := newStatus
	proceed := true
	targetStatus, proceed, err = handleTaskApprovals(db, id, currentStatus, newStatus, role, qaApp, secApp, breakGlass, reason, statuses)
	if err != nil {
		return err
	}
	if !proceed {
		return nil
	}

	// 4. Update status (and optionally description)
	if desc != "" {
		_, err = db.Exec(`UPDATE tasks SET status = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, targetStatus, desc, id)
	} else {
		_, err = db.Exec(`UPDATE tasks SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, targetStatus, id)
	}
	if err != nil {
		return err
	}

	// 5. Finalize action
	LogTaskAction(db, id, "status_change", role, fmt.Sprintf("Status changed from '%s' to '%s'", currentStatus, targetStatus))
	handleSprintAutomation(db, id, targetStatus, role, sprintID)

	return nil
}

func MoveTaskToSprint(db *sql.DB, id int64, sprintID int64) error {
	_, err := db.Exec(`UPDATE tasks SET sprint_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, sprintID, id)
	return err
}

func MoveTaskToBacklog(db *sql.DB, id int64) error {
	_, err := db.Exec(`UPDATE tasks SET sprint_id = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}

func SoftDeleteTask(db *sql.DB, id int64) error {
	_, err := db.Exec(`UPDATE tasks SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}

func ListAllTasksForRole(db *sql.DB, role string) ([]Task, error) {
	query := `SELECT id, project_id, milestone_id, sprint_id, archetype_id, title, COALESCE(description, ''), status, priority, qa_approved, security_approved, COALESCE(qa_bypassed, 0), COALESCE(security_bypassed, 0) FROM tasks WHERE deleted_at IS NULL`
	var args []interface{}

	// Filter by Role Context (Dynamic)
	visibleStatuses, seeAll, err := GetVisibleStatuses(db, role)
	if err != nil {
		return nil, err
	}
	if !seeAll {
		if len(visibleStatuses) == 0 {
			query += " AND 1=0"
		} else {
			placeholders := make([]string, len(visibleStatuses))
			for i, s := range visibleStatuses {
				placeholders[i] = "?"
				args = append(args, s)
			}
			query += fmt.Sprintf(" AND status IN (%s)", strings.Join(placeholders, ","))
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.MilestoneID, &t.SprintID, &t.ArchetypeID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.QAApproved, &t.SecurityApproved, &t.QABypassed, &t.SecurityBypassed); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
