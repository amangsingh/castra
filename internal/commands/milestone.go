package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
	"strconv"
)

type MilestoneAddCommand struct{ lastID int64 }

func (c *MilestoneAddCommand) Name() string        { return "add" }
func (c *MilestoneAddCommand) Description() string { return "Add a new milestone" }
func (c *MilestoneAddCommand) Usage() string {
	return "castra milestone add --project <pid> [--parent <mid>] [--archetype <aid>] --name <name> [--desc <description>]"
}

func (c *MilestoneAddCommand) AllowedRoles() []string { return []string{"architect"} }

func (c *MilestoneAddCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("milestone add", flag.ExitOnError)
	pid := fs.Int64("project", 0, "Project ID")
	parent := fs.Int64("parent", 0, "Parent Milestone ID (optional)")
	aid := fs.Int64("archetype", 0, "Archetype ID (optional)")
	name := fs.String("name", "", "Milestone Name")
	desc := fs.String("desc", "", "Description")
	fs.Parse(ctx.Args)

	if *pid == 0 || *name == "" {
		return fmt.Errorf("project ID and name required")
	}

	var parentID *int64
	if *parent != 0 {
		parentID = parent
	}

	var archetypeID *int64
	if *aid != 0 {
		archetypeID = aid
	}

	id, err := cli.AddMilestone(ctx.DB, *pid, parentID, archetypeID, *name, *desc)
	if err != nil {
		return err
	}
	c.lastID = id
	fmt.Printf("Milestone added: %d\n", id)
	return nil
}

func (c *MilestoneAddCommand) AuditInfo() (string, int64, string) {
	return "milestone", c.lastID, "milestone.add"
}

type MilestoneListCommand struct{}

func (c *MilestoneListCommand) Name() string        { return "list" }
func (c *MilestoneListCommand) Description() string { return "List milestones for a project" }
func (c *MilestoneListCommand) Usage() string       { return "castra milestone list --project <pid>" }

func (c *MilestoneListCommand) ReadInfo() (string, string) {
	return "milestone", "milestone.list"
}

func (c *MilestoneListCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("milestone list", flag.ExitOnError)
	pid := fs.Int64("project", 0, "Project ID")
	fs.Parse(ctx.Args)

	if *pid == 0 {
		return fmt.Errorf("project ID required")
	}
	milestones, err := cli.ListMilestones(ctx.DB, *pid, ctx.Role)
	if err != nil {
		return err
	}
	for _, m := range milestones {
		prefix := ""
		if m.ParentID != nil {
			prefix = "  ↳ "
		}
		fmt.Printf("%s[%d] %s (%s)\n", prefix, m.ID, m.Name, m.Status)
		if m.Description != "" {
			fmt.Printf("%s    %s\n", prefix, m.Description)
		}
	}
	return nil
}

type MilestoneUpdateCommand struct{ lastID int64 }

func (c *MilestoneUpdateCommand) Name() string        { return "update" }
func (c *MilestoneUpdateCommand) Description() string { return "Update milestone status" }
func (c *MilestoneUpdateCommand) Usage() string {
	return "castra milestone update --status <open|completed> <id>"
}

func (c *MilestoneUpdateCommand) AllowedRoles() []string { return []string{"architect", "senior-engineer"} }

func (c *MilestoneUpdateCommand) Execute(ctx *Context) error {

	fs := flag.NewFlagSet("milestone update", flag.ExitOnError)
	status := fs.String("status", "", "New Status (open|completed)")
	fs.Parse(ctx.Args)

	if len(fs.Args()) < 1 {
		return fmt.Errorf("ID required")
	}
	id, _ := strconv.ParseInt(fs.Args()[0], 10, 64)

	if *status == "" {
		return fmt.Errorf("status required")
	}
	if err := cli.UpdateMilestoneStatus(ctx.DB, id, *status, ctx.Role); err != nil {
		return err
	}
	c.lastID = id
	return nil
}

func (c *MilestoneUpdateCommand) AuditInfo() (string, int64, string) {
	return "milestone", c.lastID, "milestone.update"
}

type MilestoneDeleteCommand struct{ lastID int64 }

func (c *MilestoneDeleteCommand) Name() string        { return "delete" }
func (c *MilestoneDeleteCommand) Description() string { return "Delete a milestone (soft delete)" }
func (c *MilestoneDeleteCommand) Usage() string       { return "castra milestone delete <id>" }

func (c *MilestoneDeleteCommand) AllowedRoles() []string { return []string{"architect"} }

func (c *MilestoneDeleteCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("milestone delete", flag.ExitOnError)
	fs.Parse(ctx.Args)

	if len(fs.Args()) < 1 {
		return fmt.Errorf("ID required")
	}
	id, _ := strconv.ParseInt(fs.Args()[0], 10, 64)
	if err := cli.SoftDeleteMilestone(ctx.DB, id); err != nil {
		return err
	}
	c.lastID = id
	return nil
}

func (c *MilestoneDeleteCommand) AuditInfo() (string, int64, string) {
	return "milestone", c.lastID, "milestone.delete"
}

type MilestoneViewCommand struct{}

func (c *MilestoneViewCommand) Name() string        { return "view" }
func (c *MilestoneViewCommand) Description() string { return "View milestone details" }
func (c *MilestoneViewCommand) Usage() string       { return "castra milestone view <id>" }

func (c *MilestoneViewCommand) ReadInfo() (string, string) {
	return "milestone", "milestone.view"
}

func (c *MilestoneViewCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("milestone view", flag.ExitOnError)
	fs.Parse(ctx.Args)

	if len(fs.Args()) < 1 {
		return fmt.Errorf("milestone ID required")
	}
	id, _ := strconv.ParseInt(fs.Args()[0], 10, 64)

	m, err := cli.GetMilestone(ctx.DB, id)
	if err != nil {
		return err
	}

	fmt.Printf("--- Milestone [%d]: %s ---\n", m.ID, m.Name)
	fmt.Printf("Status:      %s\n", m.Status)
	fmt.Printf("Project:     %d\n", m.ProjectID)
	if m.ParentID != nil {
		fmt.Printf("Parent:      %d\n", *m.ParentID)
	}
	if m.ArchetypeID != nil {
		fmt.Printf("Archetype:   %d\n", *m.ArchetypeID)
	}
	if m.Description != "" {
		fmt.Printf("Description: %s\n", m.Description)
	}

	fmt.Println("\n--- Next Actions ---")
	printMilestoneNextActions(ctx.DB, m.Status, ctx.Role, m.ProjectID, id)
	return nil
}
