package cli

import (
	"castra/internal/persona"
	"database/sql"
	"fmt"
)

// handleTaskApprovals logic for dual-gate approval (QA/Security)
// Returns: updated newStatus, whether to proceed with status update, and error
func handleTaskApprovals(db *sql.DB, id int64, currentStatus, newStatus, role string, qaApp, secApp bool, breakGlass bool, reason string, statuses []string) (string, bool, error) {
	linter := persona.NewLinter()
	
	// 1. Jurisdictional Audit (Gate 3: The Self)
	// Only valid personas with lifecycle or audit capabilities can initiate this logic.
	if err := linter.PersonaAudit(role, persona.CapabilityLifecycle); err != nil {
		if errAudit := linter.PersonaAudit(role, persona.CapabilityAudit); errAudit != nil {
			handlePersonaNonCompliance(db, id, role, err.Error())
			return newStatus, false, err
		}
	}

	if breakGlass {
		if err := linter.PersonaAudit(role, persona.CapabilityBreakGlass); err != nil {
			handlePersonaNonCompliance(db, id, role, err.Error())
			return newStatus, false, err
		}
		LogTaskAction(db, id, "status_change.break_glass", role, fmt.Sprintf("Architect used BREAK-GLASS to force status to '%s'", newStatus))
		fmt.Printf("!!! BREAK-GLASS Protocol used by Architect to force status to %s !!!\n", newStatus)
		// If forcing done, set bypass flags and auto-create post-incident review task.
		if newStatus == "done" {
			_, err := db.Exec(`UPDATE tasks SET qa_bypassed = TRUE, security_bypassed = TRUE WHERE id = ?`, id)
			if err != nil {
				return newStatus, false, err
			}
			// Auto-create post-incident review task
			var projectID int64
			_ = db.QueryRow(`SELECT project_id FROM tasks WHERE id = ?`, id).Scan(&projectID)
			title := fmt.Sprintf("Post-Incident Review: Task [%d]", id)
			desc := fmt.Sprintf("Architect bypassed QA & Security for Task [%d]. Review the changes.", id)
			_, _ = db.Exec(
				`INSERT INTO tasks (project_id, title, description, status, priority) VALUES (?, ?, ?, 'review', 'low')`,
				projectID, title, desc,
			)
			fmt.Printf("Post-incident review task created for Task [%d].\n", id)
		}
		return newStatus, true, nil
	}

	// State guard: Reviewers (QA/Security) can only act on tasks that are in the review/gatekeeper status.
	// Architects are exempt from this guard as they oversee the entire lifecycle.
	reviewStatus := GetReviewStatus(statuses)
	isGatekeeper := linter.PersonaAudit(role, persona.CapabilityAudit) == nil
	isArchitect := linter.PersonaAudit(role, persona.CapabilityBreakGlass) == nil

	if isGatekeeper && !isArchitect && currentStatus != reviewStatus {
		return newStatus, false, fmt.Errorf("%s can only act on tasks in '%s' status (current: '%s')", role, reviewStatus, currentStatus)
	}

	if newStatus == "todo" && isGatekeeper {
		// Rejection: reset BOTH approval flags
		_, err := db.Exec(`UPDATE tasks SET qa_approved = false, security_approved = false, status = 'todo', updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
		if err != nil {
			return newStatus, false, err
		}
		logMsg := "Task rejected from review to todo. All approvals reset."
		if reason != "" {
			logMsg = fmt.Sprintf("Task rejected. Reason: %s. All approvals reset.", reason)
		}
		LogTaskAction(db, id, "rejected", role, logMsg)
		fmt.Printf("Task rejected by %s. All approvals reset. Task returned to todo.\n", role)
		return "todo", false, nil // Already updated status here
	}

	if newStatus == statuses[len(statuses)-1] || newStatus == "done" {
		// Determine required gates from archetype statuses.
		var defaultRole string
		_ = db.QueryRow(`SELECT a.default_role FROM tasks t LEFT JOIN archetypes a ON t.archetype_id = a.id WHERE t.id = ?`, id).Scan(&defaultRole)
		requiredQA := defaultRole != "security-ops"
		requiredSec := true

		// Update approvals based on Persona Capabilities (Gate Audit)
		updated := false
		if err := linter.PersonaAudit(role, persona.CapabilityAuditFunctional); err == nil {
			if !qaApp {
				qaApp = true
				updated = true
			}
		}
		if err := linter.PersonaAudit(role, persona.CapabilityAuditSecurity); err == nil {
			if !secApp {
				secApp = true
				updated = true
			}
		}

		if updated {
			_, err2 := db.Exec(`UPDATE tasks SET qa_approved = ?, security_approved = ? WHERE id = ?`, qaApp, secApp, id)
			if err2 != nil {
				return newStatus, false, err2
			}
		}

		// Check if all required gates are satisfied
		if (requiredQA && !qaApp) || (requiredSec && !secApp) {
			if updated {
				LogTaskAction(db, id, "approved", role, "Approval granted. Waiting for other gate.")
				fmt.Printf("Task approved by %s. Waiting for other required approval(s) to mark DONE.\n", role)
			} else {
				// If not an auditor, maybe they are trying to mark DONE without permission?
				// Engineers cannot mark as DONE (enforced in UpdateTaskStatus too)
				if err := linter.PersonaAudit(role, persona.CapabilityAudit); err != nil {
					fmt.Printf("Unauthorized attempt to mark DONE by %s: missing required approvals. Use --break-glass to override if Architect.\n", role)
				}
			}
			return newStatus, false, nil // Don't proceed to 'done' yet
		}

		if updated {
			LogTaskAction(db, id, "approved", role, "Final approval granted. All required gates passed.")
		}
		return "done", true, nil // Proceed to update status to 'done'
	}

	return newStatus, true, nil
}

// handlePersonaNonCompliance resets task to its initial status and logs the audit failure
func handlePersonaNonCompliance(db *sql.DB, id int64, role string, errMsg string) {
	_ = AddAuditEntry(db, "system", 0, "persona_non_compliance", role, "[handleTaskApprovals] "+errMsg)
	
	// Reset to archetype's first status
	statuses := GetTaskStatuses(db, id)
	initialStatus := "todo"
	if len(statuses) > 0 {
		initialStatus = statuses[0]
	}

	_, _ = db.Exec(`UPDATE tasks SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, initialStatus, id)
	LogTaskAction(db, id, "rejected", "system", fmt.Sprintf("Persona Non-Compliance: Task reset to '%s' due to jurisdictional violation.", initialStatus))
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
