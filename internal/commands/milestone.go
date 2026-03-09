package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
	"strconv"
)

type MilestoneAddCommand struct{}

func (c *MilestoneAddCommand) Name() string        { return "add" }
func (c *MilestoneAddCommand) Description() string { return "Add a new milestone" }
func (c *MilestoneAddCommand) Usage() string {
	return "castra milestone add --project <pid> --name <name>"
}

func (c *MilestoneAddCommand) Execute(ctx *Context) error {
	if ctx.Role != "architect" {
		return fmt.Errorf("only architect can manage milestones")
	}

	fs := flag.NewFlagSet("milestone add", flag.ExitOnError)
	pid := fs.Int64("project", 0, "Project ID")
	name := fs.String("name", "", "Milestone Name")
	fs.Parse(ctx.Args)

	if *pid == 0 || *name == "" {
		return fmt.Errorf("project ID and name required")
	}
	id, err := cli.AddMilestone(ctx.DB, *pid, *name)
	if err != nil {
		return err
	}
	fmt.Printf("Milestone added: %d\n", id)
	return nil
}

type MilestoneListCommand struct{}

func (c *MilestoneListCommand) Name() string        { return "list" }
func (c *MilestoneListCommand) Description() string { return "List milestones for a project" }
func (c *MilestoneListCommand) Usage() string       { return "castra milestone list --project <pid>" }

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
		fmt.Printf("[%d] %s (%s)\n", m.ID, m.Name, m.Status)
	}
	return nil
}

type MilestoneUpdateCommand struct{}

func (c *MilestoneUpdateCommand) Name() string        { return "update" }
func (c *MilestoneUpdateCommand) Description() string { return "Update milestone status" }
func (c *MilestoneUpdateCommand) Usage() string {
	return "castra milestone update --status <open|completed> <id>"
}

func (c *MilestoneUpdateCommand) Execute(ctx *Context) error {
	if ctx.Role != "architect" {
		return fmt.Errorf("only architect can manage milestones")
	}

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
	return cli.UpdateMilestoneStatus(ctx.DB, id, *status, ctx.Role)
}

type MilestoneDeleteCommand struct{}

func (c *MilestoneDeleteCommand) Name() string        { return "delete" }
func (c *MilestoneDeleteCommand) Description() string { return "Delete a milestone (soft delete)" }
func (c *MilestoneDeleteCommand) Usage() string       { return "castra milestone delete <id>" }

func (c *MilestoneDeleteCommand) Execute(ctx *Context) error {
	if ctx.Role != "architect" {
		return fmt.Errorf("only architect can manage milestones")
	}

	fs := flag.NewFlagSet("milestone delete", flag.ExitOnError)
	fs.Parse(ctx.Args)

	if len(fs.Args()) < 1 {
		return fmt.Errorf("ID required")
	}
	id, _ := strconv.ParseInt(fs.Args()[0], 10, 64)
	return cli.SoftDeleteMilestone(ctx.DB, id)
}
