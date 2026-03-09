package cli

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID               int64
	ProjectID        int64
	MilestoneID      *int64 // Pointer to handle NULL
	SprintID         *int64 // Pointer to handle NULL
	Title            string
	Description      string
	Status           string
	Priority         string
	QAApproved       bool
	SecurityApproved bool
}

func AddTask(db *sql.DB, projectID int64, milestoneID, sprintID *int64, title, desc, priority string) (int64, error) {
	// Defaults to Todo
	query := `INSERT INTO tasks (project_id, milestone_id, sprint_id, title, description, priority, status) VALUES (?, ?, ?, ?, ?, ?, 'todo')`
	res, err := db.Exec(query, projectID, milestoneID, sprintID, title, desc, priority)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(db *sql.DB, id int64) (*Task, error) {
	var t Task
	query := `SELECT id, project_id, milestone_id, sprint_id, title, COALESCE(description, ''), status, priority, qa_approved, security_approved FROM tasks WHERE id = ? AND deleted_at IS NULL`
	err := db.QueryRow(query, id).Scan(&t.ID, &t.ProjectID, &t.MilestoneID, &t.SprintID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.QAApproved, &t.SecurityApproved)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return &t, nil
}

func ListTasks(db *sql.DB, projectID int64, milestoneID, sprintID *int64, backlogOnly bool, role string) ([]Task, error) {
	query := `SELECT id, project_id, milestone_id, sprint_id, title, COALESCE(description, ''), status, priority, qa_approved, security_approved FROM tasks WHERE project_id = ? AND deleted_at IS NULL`
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

	// Filter by Role Context
	// 1. Architect: Everything
	// 2. Engineer: todo, doing, blocked, pending
	// 3. QA & Security: review
	switch role {
	case "junior-engineer", "senior-engineer", "designer":
		query += ` AND status IN ('todo', 'doing', 'blocked', 'pending')`
	case "qa-functional", "security-ops":
		query += ` AND status = 'review'`
	case "doc-writer":
		// Doc writer sees all, no status filter needed
	case "architect":
		// No filter
	default:
		return nil, fmt.Errorf("unknown role: %s", role)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.MilestoneID, &t.SprintID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.QAApproved, &t.SecurityApproved); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func UpdateTaskStatus(db *sql.DB, id int64, newStatus string, role string) error {
	// 1. Fetch current task state
	var currentStatus string
	var qaApp, secApp bool
	var sprintID sql.NullInt64
	err := db.QueryRow(`SELECT status, qa_approved, security_approved, sprint_id FROM tasks WHERE id = ?`, id).Scan(&currentStatus, &qaApp, &secApp, &sprintID)
	if err != nil {
		return err
	}

	// 2. Validate Transitions by Role
	switch role {
	case "architect":
		// Can do anything
	case "junior-engineer", "senior-engineer", "designer":
		if newStatus == "done" {
			return fmt.Errorf("engineer/designer cannot mark task as done (must be approved by qa & security)")
		}
	case "doc-writer":
		return fmt.Errorf("doc-writer cannot update tasks")
	}

	// 3. Handle Approval Logic (Gates)
	targetStatus, proceed, err := handleTaskApprovals(db, id, currentStatus, newStatus, role, qaApp, secApp)
	if err != nil {
		return err
	}
	if !proceed {
		return nil
	}

	// 4. Update status if allowed
	_, err = db.Exec(`UPDATE tasks SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, targetStatus, id)
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
	query := `SELECT id, project_id, milestone_id, sprint_id, title, COALESCE(description, ''), status, priority, qa_approved, security_approved FROM tasks WHERE deleted_at IS NULL`
	var args []interface{}

	switch role {
	case "junior-engineer", "senior-engineer", "designer":
		query += ` AND status IN ('todo', 'doing', 'blocked', 'pending')`
	case "qa-functional", "security-ops":
		query += ` AND status = 'review'`
	case "doc-writer", "architect":
		// No filter
	default:
		return nil, fmt.Errorf("unknown role: %s", role)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.MilestoneID, &t.SprintID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.QAApproved, &t.SecurityApproved); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
