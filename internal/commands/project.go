package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
	"strconv"
)

type ProjectAddCommand struct{}

func (c *ProjectAddCommand) Name() string        { return "add" }
func (c *ProjectAddCommand) Description() string { return "Add a new project" }
func (c *ProjectAddCommand) Usage() string {
	return "castra project add --name <name> [--desc <description>] [--notes <notes>]"
}

func (c *ProjectAddCommand) Execute(ctx *Context) error {
	if ctx.Role != "architect" {
		return fmt.Errorf("only architect can modify projects")
	}

	fs := flag.NewFlagSet("project add", flag.ExitOnError)
	name := fs.String("name", "", "Project Name")
	desc := fs.String("desc", "", "Description")
	notes := fs.String("notes", "", "Notes/Docs")

	fs.Parse(ctx.Args)

	if *name == "" {
		return fmt.Errorf("name required")
	}
	id, err := cli.AddProject(ctx.DB, *name, *desc, *notes)
	if err != nil {
		return err
	}
	fmt.Printf("Project added with ID: %d\n", id)
	return nil
}

type ProjectListCommand struct{}

func (c *ProjectListCommand) Name() string        { return "list" }
func (c *ProjectListCommand) Description() string { return "List all projects" }
func (c *ProjectListCommand) Usage() string       { return "castra project list [--archived] [--deleted]" }

func (c *ProjectListCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("project list", flag.ExitOnError)
	archived := fs.Bool("archived", false, "Show archived")
	deleted := fs.Bool("deleted", false, "Show deleted")
	fs.Parse(ctx.Args)

	projects, err := cli.ListProjects(ctx.DB, *archived, *deleted)
	if err != nil {
		return err
	}
	for _, p := range projects {
		fmt.Printf("[%d] %s (%s)\n", p.ID, p.Name, p.Status)
	}
	return nil
}

type ProjectDeleteCommand struct{}

func (c *ProjectDeleteCommand) Name() string        { return "delete" }
func (c *ProjectDeleteCommand) Description() string { return "Delete a project (soft delete)" }
func (c *ProjectDeleteCommand) Usage() string       { return "castra project delete <id>" }

func (c *ProjectDeleteCommand) Execute(ctx *Context) error {
	if ctx.Role != "architect" {
		return fmt.Errorf("only architect can modify projects")
	}

	fs := flag.NewFlagSet("project delete", flag.ExitOnError)
	fs.Parse(ctx.Args)

	idParsed := fs.Args()
	if len(idParsed) < 1 {
		return fmt.Errorf("ID required")
	}
	id, _ := strconv.ParseInt(idParsed[0], 10, 64)
	if err := cli.SoftDeleteProject(ctx.DB, id); err != nil {
		return err
	}
	fmt.Println("Project soft deleted.")
	return nil
}
