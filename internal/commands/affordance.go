package commands

import (
	"castra/internal/cli"
	"castra/internal/persona"
	"database/sql"
	"fmt"
)

type action struct {
	cmd  string
	desc string
}

func printTaskNextActions(db *sql.DB, status, role string, id int64) {
	var actions []action

	statuses := cli.GetTaskStatuses(db, id)
	if len(statuses) == 0 {
		fmt.Println("No status pipeline defined.")
		return
	}

	firstStatus := statuses[0]
	next := cli.NextStatus(statuses, status)

	// Fetch Archetype Metadata
	var defaultRole sql.NullString
	_ = db.QueryRow(`
		SELECT a.default_role 
		FROM tasks t 
		JOIN archetypes a ON t.archetype_id = a.id 
		WHERE t.id = ?`, id).Scan(&defaultRole)

	linter := persona.NewLinter()
	isLifecycle := linter.PersonaAudit(role, persona.CapabilityLifecycle) == nil
	isAuditor := linter.PersonaAudit(role, persona.CapabilityAudit) == nil
	isArchitect := linter.PersonaAudit(role, persona.CapabilityBreakGlass) == nil

	if isLifecycle {
		isAssignedRole := !defaultRole.Valid || defaultRole.String == "" || role == defaultRole.String
		if isArchitect || isAssignedRole {
			if next != "" {
				actions = append(actions, action{
					cmd:  fmt.Sprintf("castra task update --role %s --status %s %d", role, next, id),
					desc: fmt.Sprintf("Move to %s", next),
				})
			}
			if next != "" && status != "blocked" {
				actions = append(actions, action{
					cmd:  fmt.Sprintf("castra task update --role %s --status blocked %d", role, id),
					desc: "Block task",
				})
			}
			if status == "blocked" {
				actions = append(actions, action{
					cmd:  fmt.Sprintf("castra task update --role %s --status doing %d", role, id),
					desc: "Resume (doing)",
				})
			}
		}
	}

	if isAuditor {
		reviewStatus := cli.GetReviewStatus(statuses)
		if status == reviewStatus {
			if next != "" {
				actions = append(actions, action{
					cmd:  fmt.Sprintf("castra task update --role %s --status %s %d", role, next, id),
					desc: "Approve",
				})
			}
			actions = append(actions, action{
				cmd:  fmt.Sprintf("castra task update --role %s --status %s --reason \"...\" %d", role, firstStatus, id),
				desc: "Reject (return to start)",
			})
		}
	}

	displayActions(actions)
}

func printProjectNextActions(db *sql.DB, status, role string, id int64) {
	var actions []action

	linter := persona.NewLinter()
	isArchitect := linter.PersonaAudit(role, persona.CapabilityBreakGlass) == nil

	actions = append(actions, action{
		cmd:  fmt.Sprintf("castra project list --role %s", role),
		desc: "List all projects",
	})
	actions = append(actions, action{
		cmd:  fmt.Sprintf("castra milestone add --role %s --project %d --name \"...\"", role, id),
		desc: "Add milestone to this project",
	})
	actions = append(actions, action{
		cmd:  fmt.Sprintf("castra milestone list --role %s --project %d", role, id),
		desc: "List milestones for this project",
	})
	actions = append(actions, action{
		cmd:  fmt.Sprintf("castra task list --role %s --project %d", role, id),
		desc: "List all tasks for this project",
	})

	if isArchitect {
		if status == "active" {
			actions = append(actions, action{
				cmd:  fmt.Sprintf("castra project update --role %s --status archived %d", role, id),
				desc: "Archive project",
			})
		} else {
			actions = append(actions, action{
				cmd:  fmt.Sprintf("castra project update --role %s --status active %d", role, id),
				desc: "Restore project (activate)",
			})
		}
		actions = append(actions, action{
			cmd:  fmt.Sprintf("castra project delete --role %s %d", role, id),
			desc: "Delete project (soft)",
		})
	}

	displayActions(actions)
}

func printMilestoneNextActions(db *sql.DB, status, role string, projectID, id int64) {
	var actions []action

	linter := persona.NewLinter()
	isArchitect := linter.PersonaAudit(role, persona.CapabilityBreakGlass) == nil

	actions = append(actions, action{
		cmd:  fmt.Sprintf("castra milestone list --role %s --project %d", role, projectID),
		desc: "List milestones for this project",
	})
	actions = append(actions, action{
		cmd:  fmt.Sprintf("castra task list --role %s --project %d --milestone %d", role, projectID, id),
		desc: "List tasks for this milestone",
	})
	actions = append(actions, action{
		cmd:  fmt.Sprintf("castra project view --role %s %d", role, projectID),
		desc: "View parent project",
	})

	if isArchitect {
		if status == "open" {
			actions = append(actions, action{
				cmd:  fmt.Sprintf("castra milestone update --role %s --status completed %d", role, id),
				desc: "Mark as completed",
			})
		} else {
			actions = append(actions, action{
				cmd:  fmt.Sprintf("castra milestone update --role %s --status open %d", role, id),
				desc: "Reopen milestone",
			})
		}
		actions = append(actions, action{
			cmd:  fmt.Sprintf("castra milestone delete --role %s %d", role, id),
			desc: "Delete milestone (soft)",
		})
	}

	displayActions(actions)
}

func displayActions(actions []action) {
	if len(actions) == 0 {
		fmt.Println("No actions available for this role in current status.")
	} else {
		for _, a := range actions {
			fmt.Printf("  → %s\n    (%s)\n", a.cmd, a.desc)
		}
	}
}
