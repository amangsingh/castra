package cli

import (
	"database/sql"
	"fmt"
)

// handleTaskApprovals logic for dual-gate approval (QA/Security)
// Returns: updated newStatus, whether to proceed with status update, and error
func handleTaskApprovals(db *sql.DB, id int64, currentStatus, newStatus, role string, qaApp, secApp bool) (string, bool, error) {
	if role == "qa-functional" || role == "security-ops" {
		if currentStatus != "review" {
			return newStatus, false, fmt.Errorf("%s can only process tasks in 'review' status", role)
		}

		if newStatus == "todo" {
			// Rejection: reset BOTH approval flags
			_, err := db.Exec(`UPDATE tasks SET qa_approved = false, security_approved = false, status = 'todo', updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
			if err != nil {
				return newStatus, false, err
			}
			LogTaskAction(db, id, "rejected", role, "Task rejected from review to todo. All approvals reset.")
			fmt.Printf("Task rejected by %s. All approvals reset. Task returned to todo.\n", role)
			return "todo", false, nil // Already updated status here
		}

		if newStatus == "done" {
			if role == "qa-functional" {
				qaApp = true
			}
			if role == "security-ops" {
				secApp = true
			}

			_, err := db.Exec(`UPDATE tasks SET qa_approved = ?, security_approved = ? WHERE id = ?`, qaApp, secApp, id)
			if err != nil {
				return newStatus, false, err
			}

			if !qaApp || !secApp {
				LogTaskAction(db, id, "approved", role, "Approval granted. Waiting for other gate.")
				fmt.Printf("Task approved by %s. Waiting for other approval to mark DONE.\n", role)
				return newStatus, false, nil // Don't proceed to 'done' yet
			}
			LogTaskAction(db, id, "approved", role, "Final approval granted. Both gates passed.")
			return "done", true, nil // Proceed to update status to 'done'
		}
	}
	return newStatus, true, nil
}

// handleSprintAutomation triggers sprint state changes based on task state
func handleSprintAutomation(db *sql.DB, taskID int64, newStatus, role string, sprintID sql.NullInt64) {
	if !sprintID.Valid {
		return
	}

	sid := sprintID.Int64

	// Auto-start
	if newStatus == "doing" {
		var sprintStatus string
		err := db.QueryRow(`SELECT status FROM sprints WHERE id = ?`, sid).Scan(&sprintStatus)
		if err == nil && sprintStatus == "planning" {
			_, errUpdate := db.Exec(`UPDATE sprints SET status = 'in progress' WHERE id = ?`, sid)
			if errUpdate == nil {
				payload := fmt.Sprintf("[%s] Sprint automatically started by task %d status change", role, taskID)
				_ = AddAuditEntry(db, "sprint", sid, "status_change", role, payload)
				fmt.Printf("Sprint %d automatically started.\n", sid)
			}
		}
	}

	// Auto-complete
	if newStatus == "done" {
		var pendingCount int
		err := db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE sprint_id = ? AND status != 'done' AND deleted_at IS NULL`, sid).Scan(&pendingCount)
		if err == nil && pendingCount == 0 {
			_, errUpdate := db.Exec(`UPDATE sprints SET status = 'done' WHERE id = ?`, sid)
			if errUpdate == nil {
				payload := fmt.Sprintf("[%s] Sprint automatically completed by task %d completion", role, taskID)
				_ = AddAuditEntry(db, "sprint", sid, "status_change", role, payload)
				fmt.Printf("Sprint %d automatically completed.\n", sid)
			}
		}
	}
}
