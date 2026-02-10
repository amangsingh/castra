package cli

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID               int64
	ProjectID        int64
	SprintID         *int64 // Pointer to handle NULL
	Title            string
	Description      string
	Status           string
	Priority         string
	QAApproved       bool
	SecurityApproved bool
}

func AddTask(db *sql.DB, projectID int64, sprintID *int64, title, desc, priority string) (int64, error) {
	// Defaults to Todo
	query := `INSERT INTO tasks (project_id, sprint_id, title, description, priority, status) VALUES (?, ?, ?, ?, ?, 'todo')`
	res, err := db.Exec(query, projectID, sprintID, title, desc, priority)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ListTasks(db *sql.DB, projectID int64, sprintID *int64, backlogOnly bool, role string) ([]Task, error) {
	query := `SELECT id, project_id, sprint_id, title, description, status, priority, qa_approved, security_approved FROM tasks WHERE project_id = ? AND deleted_at IS NULL`
	args := []interface{}{projectID}

	if backlogOnly {
		query += ` AND sprint_id IS NULL`
	} else if sprintID != nil {
		query += ` AND sprint_id = ?`
		args = append(args, *sprintID)
	}

	// Filter by Role Context
	// 1. Architect: Everything
	// 2. Engineer: todo, doing, blocked, pending
	// 3. QA & Security: review
	switch role {
	case "engineer":
		query += ` AND status IN ('todo', 'doing', 'blocked', 'pending')`
	case "qa", "security":
		query += ` AND status = 'review'`
	case "doc-writer":
		query += ` AND status = 'done'`
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
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.SprintID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.QAApproved, &t.SecurityApproved); err != nil {
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
	err := db.QueryRow(`SELECT status, qa_approved, security_approved FROM tasks WHERE id = ?`, id).Scan(&currentStatus, &qaApp, &secApp)
	if err != nil {
		return err
	}

	// 2. Validate Transitions by Role
	switch role {
	case "architect":
		// Can do anything
	case "engineer":
		// Cannot set to 'done'
		if newStatus == "done" {
			return fmt.Errorf("engineer cannot mark task as done (must be approved by qa & security)")
		}
	case "qa", "security":
		// Can only pick up from 'review' and mark as 'done' (conditional)
		if currentStatus != "review" {
			return fmt.Errorf("%s can only process tasks in 'review' status", role)
		}

		if newStatus == "done" {
			// Register Approval
			if role == "qa" {
				qaApp = true
			}
			if role == "security" {
				secApp = true
			}

			// Update Approval Flags
			_, err := db.Exec(`UPDATE tasks SET qa_approved = ?, security_approved = ? WHERE id = ?`, qaApp, secApp, id)
			if err != nil {
				return err
			}

			// Check Lock: Both must be true to transition to DONE
			if !qaApp || !secApp {
				fmt.Printf("Task approved by %s. Waiting for other approval to mark DONE.\n", role)
				return nil // Not an error, just didn't transition status yet
			}
			// Fallthrough to update status to done
		}
	case "doc-writer":
		return fmt.Errorf("doc-writer cannot update tasks")
	}

	_, err = db.Exec(`UPDATE tasks SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, newStatus, id)
	return err
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
