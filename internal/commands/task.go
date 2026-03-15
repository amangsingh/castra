package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
	"strconv"
)

type TaskAddCommand struct{ lastID int64 }

func (c *TaskAddCommand) Name() string        { return "add" }
func (c *TaskAddCommand) Description() string { return "Add a new task" }
func (c *TaskAddCommand) Usage() string {
	return "castra task add --project <pid> [--milestone <mid>] [--sprint <sid>] [--archetype <aid>] --title <title> [--desc <description>] [--prio <low|medium|high>]"
}

func (c *TaskAddCommand) AllowedRoles() []string { return []string{"architect"} }

func (c *TaskAddCommand) Execute(ctx *Context) error {

	fs := flag.NewFlagSet("task add", flag.ExitOnError)
	pid := fs.Int64("project", 0, "Project ID")
	mid := fs.Int64("milestone", 0, "Milestone ID (optional)")
	sid := fs.Int64("sprint", 0, "Sprint ID (optional)")
	aid := fs.Int64("archetype", 0, "Archetype ID (optional)")
	title := fs.String("title", "", "Title")
	desc := fs.String("desc", "", "Description")
	prio := fs.String("prio", "medium", "Priority")
	fs.Parse(ctx.Args)

	if *pid == 0 || *title == "" {
		return fmt.Errorf("project ID and Title required")
	}

	var milestoneID *int64
	if *mid != 0 {
		milestoneID = mid
	}

	var sprintID *int64
	if *sid != 0 {
		sprintID = sid
	}

	var archetypeID *int64
	if *aid != 0 {
		archetypeID = aid
	}

	id, err := cli.AddTask(ctx.DB, *pid, milestoneID, sprintID, archetypeID, *title, *desc, *prio)
	if err != nil {
		return err
	}
	c.lastID = id
	fmt.Printf("Task created: %d\n", id)
	return nil
}

func (c *TaskAddCommand) AuditInfo() (string, int64, string) { return "task", c.lastID, "task.add" }

type TaskListCommand struct{}

func (c *TaskListCommand) Name() string        { return "list" }
func (c *TaskListCommand) Description() string { return "List tasks" }
func (c *TaskListCommand) Usage() string {
	return "castra task list --project <pid> [--milestone <mid>] [--sprint <sid>] [--backlog]"
}

func (c *TaskListCommand) ReadInfo() (string, string) {
	return "task", "task.list"
}

func (c *TaskListCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("task list", flag.ExitOnError)
	pid := fs.Int64("project", 0, "Project ID")
	mid := fs.Int64("milestone", 0, "Milestone ID")
	sid := fs.Int64("sprint", 0, "Sprint ID")
	backlog := fs.Bool("backlog", false, "Show backlog only")
	fs.Parse(ctx.Args)

	if *pid == 0 {
		return fmt.Errorf("project ID required")
	}

	var milestoneID *int64
	if *mid != 0 {
		milestoneID = mid
	}

	var sprintID *int64
	if *sid != 0 {
		sprintID = sid
	}

	tasks, err := cli.ListTasks(ctx.DB, *pid, milestoneID, sprintID, *backlog, ctx.Role)
	if err != nil {
		return err
	}
	for _, t := range tasks {
		badge := ""
		if t.QABypassed || t.SecurityBypassed {
			badge = "[QA/Sec Bypassed]"
		} else {
			if t.QAApproved {
				badge += "[QA]"
			}
			if t.SecurityApproved {
				badge += "[SEC]"
			}
		}
		fmt.Printf("[%d] %s (%s) %s\n", t.ID, t.Title, t.Status, badge)
	}
	return nil
}

type TaskViewCommand struct{}

func (c *TaskViewCommand) Name() string        { return "view" }
func (c *TaskViewCommand) Description() string { return "View task details" }
func (c *TaskViewCommand) Usage() string       { return "castra task view <id>" }

func (c *TaskViewCommand) ReadInfo() (string, string) {
	return "task", "task.view"
}

func (c *TaskViewCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("task view", flag.ExitOnError)
	fs.Parse(ctx.Args)

	if len(fs.Args()) < 1 {
		return fmt.Errorf("task ID required")
	}
	id, _ := strconv.ParseInt(fs.Args()[0], 10, 64)

	task, err := cli.GetTask(ctx.DB, id)
	if err != nil {
		return err
	}

	approvals := ""
	if task.QABypassed || task.SecurityBypassed {
		approvals = "[QA/Sec Bypassed]"
	} else {
		if task.QAApproved {
			approvals += "[QA Approved] "
		}
		if task.SecurityApproved {
			approvals += "[Security Approved]"
		}
	}

	fmt.Printf("--- Task [%d]: %s ---\n", task.ID, task.Title)
	fmt.Printf("Status:   %s %s\n", task.Status, approvals)
	fmt.Printf("Priority: %s\n", task.Priority)
	fmt.Printf("Project:  %d\n", task.ProjectID)
	if task.MilestoneID != nil {
		fmt.Printf("Milestone: %d\n", *task.MilestoneID)
	}
	if task.SprintID != nil {
		fmt.Printf("Sprint:   %d\n", *task.SprintID)
	}
	fmt.Printf("\nDescription:\n%s\n", task.Description)

	fmt.Println("\n--- Notes ---")
	notes, err := cli.ListNotes(ctx.DB, task.ProjectID, &id, ctx.Role)
	if err == nil {
		if len(notes) == 0 {
			fmt.Println("No notes found for this role.")
		} else {
			for _, n := range notes {
				fmt.Printf("[Tags: %s] %s\n", n.Tags, n.Content)
			}
		}
	}

	fmt.Println("\n--- Audit Log ---")
	logs, err := cli.ListAuditEntries(ctx.DB, "task", &id)
	if err == nil {
		if len(logs) == 0 {
			fmt.Println("No logs found.")
		} else {
			for _, l := range logs {
				fmt.Printf("[%s] %s: %s\n", l.Timestamp, l.Role, l.Payload)
			}
		}
	}

	fmt.Println("\n--- Next Actions ---")
	printTaskNextActions(ctx.DB, task.Status, ctx.Role, id)
	return nil
}

type TaskUpdateCommand struct{ lastID int64 }

func (c *TaskUpdateCommand) Name() string        { return "update" }
func (c *TaskUpdateCommand) Description() string { return "Update task status" }
func (c *TaskUpdateCommand) Usage() string {
	return "castra task update [--status <status>] [--desc <description>] [--reason <text>] [--break-glass|--force] <id>"
}

func (c *TaskUpdateCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("task update", flag.ContinueOnError)
	status := fs.String("status", "", "New Status")
	reason := fs.String("reason", "", "Rejection reason (required when setting status to todo as qa/security)")
	desc := fs.String("desc", "", "New description text (optional)")

	// Pre-scan for --break-glass or --force alias because Go's flag.Parse stops at the first
	// positional argument (the task ID), so placing --break-glass after the ID
	// would silently ignore it. We strip it out before handing args to fs.Parse.
	breakGlass := false
	filteredArgs := ctx.Args[:0:len(ctx.Args)]
	for _, a := range ctx.Args {
		if a == "--break-glass" || a == "-break-glass" || a == "--force" || a == "-force" {
			breakGlass = true
		} else {
			filteredArgs = append(filteredArgs, a)
		}
	}
	fs.Parse(filteredArgs)

	if len(fs.Args()) < 1 {
		return fmt.Errorf("ID required")
	}
	id, _ := strconv.ParseInt(fs.Args()[0], 10, 64)

	// --status is optional only when --desc is provided
	if *status == "" && *desc == "" {
		return fmt.Errorf("--status or --desc is required")
	}

	// Enforce rejection reason for QA/Security roles reverting to todo
	if *status == "todo" && (*reason == "") && (ctx.Role == "qa-functional" || ctx.Role == "security-ops") {
		return fmt.Errorf("--reason is required when rejecting a task (setting status to 'todo') as %s", ctx.Role)
	}

	if err := cli.UpdateTaskStatus(ctx.DB, id, *status, *desc, ctx.Role, breakGlass, *reason); err != nil {
		return err
	}
	c.lastID = id
	return nil
}

func (c *TaskUpdateCommand) AuditInfo() (string, int64, string) {
	return "task", c.lastID, "task.update"
}

type TaskDeleteCommand struct{ lastID int64 }

func (c *TaskDeleteCommand) Name() string        { return "delete" }
func (c *TaskDeleteCommand) Description() string { return "Delete a task (soft delete)" }
func (c *TaskDeleteCommand) Usage() string       { return "castra task delete <id>" }

func (c *TaskDeleteCommand) AllowedRoles() []string { return []string{"architect"} }

func (c *TaskDeleteCommand) Execute(ctx *Context) error {

	fs := flag.NewFlagSet("task delete", flag.ExitOnError)
	fs.Parse(ctx.Args)

	if len(fs.Args()) < 1 {
		return fmt.Errorf("ID required")
	}
	id, _ := strconv.ParseInt(fs.Args()[0], 10, 64)
	if err := cli.SoftDeleteTask(ctx.DB, id); err != nil {
		return err
	}
	c.lastID = id
	return nil
}

func (c *TaskDeleteCommand) AuditInfo() (string, int64, string) {
	return "task", c.lastID, "task.delete"
}
