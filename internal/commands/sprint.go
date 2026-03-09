package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
)

type SprintAddCommand struct{}

func (c *SprintAddCommand) Name() string        { return "add" }
func (c *SprintAddCommand) Description() string { return "Add a new sprint" }
func (c *SprintAddCommand) Usage() string {
	return "castra sprint add --project <pid> --name <name> [--start <yyyy-mm-dd>] [--end <yyyy-mm-dd>]"
}

func (c *SprintAddCommand) Execute(ctx *Context) error {
	if ctx.Role != "architect" {
		return fmt.Errorf("only architect can manage sprints")
	}

	fs := flag.NewFlagSet("sprint add", flag.ExitOnError)
	pid := fs.Int64("project", 0, "Project ID")
	name := fs.String("name", "", "Sprint Name")
	start := fs.String("start", "", "YYYY-MM-DD")
	end := fs.String("end", "", "YYYY-MM-DD")
	fs.Parse(ctx.Args)

	if *pid == 0 || *name == "" {
		return fmt.Errorf("project ID and name required")
	}
	id, err := cli.AddSprint(ctx.DB, *pid, *name, *start, *end)
	if err != nil {
		return err
	}
	fmt.Printf("Sprint added: %d\n", id)
	return nil
}

type SprintListCommand struct{}

func (c *SprintListCommand) Name() string        { return "list" }
func (c *SprintListCommand) Description() string { return "List sprints for a project" }
func (c *SprintListCommand) Usage() string       { return "castra sprint list --project <pid>" }

func (c *SprintListCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("sprint list", flag.ExitOnError)
	pid := fs.Int64("project", 0, "Project ID")
	fs.Parse(ctx.Args)

	if *pid == 0 {
		return fmt.Errorf("project ID required")
	}
	sprints, err := cli.ListSprints(ctx.DB, *pid)
	if err != nil {
		return err
	}
	for _, s := range sprints {
		fmt.Printf("[%d] %s (%s)\n", s.ID, s.Name, s.Status)
	}
	return nil
}
