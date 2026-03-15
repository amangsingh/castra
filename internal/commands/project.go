package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
	"strconv"
)

type ProjectAddCommand struct{ lastID int64 }

func (c *ProjectAddCommand) Name() string        { return "add" }
func (c *ProjectAddCommand) Description() string { return "Add a new project" }
func (c *ProjectAddCommand) Usage() string {
	return "castra project add --name <name> [--desc <description>] [--notes <notes>]"
}

func (c *ProjectAddCommand) AllowedRoles() []string { return []string{"architect"} }

func (c *ProjectAddCommand) Execute(ctx *Context) error {

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
	c.lastID = id
	fmt.Printf("Project added with ID: %d\n", id)
	return nil
}

func (c *ProjectAddCommand) AuditInfo() (string, int64, string) {
	return "project", c.lastID, "project.add"
}

type ProjectListCommand struct{}

func (c *ProjectListCommand) Name() string        { return "list" }
func (c *ProjectListCommand) Description() string { return "List all projects" }
func (c *ProjectListCommand) Usage() string       { return "castra project list [--archived] [--deleted]" }

func (c *ProjectListCommand) ReadInfo() (string, string) {
	return "project", "project.list"
}

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

type ProjectDeleteCommand struct{ lastID int64 }

func (c *ProjectDeleteCommand) Name() string        { return "delete" }
func (c *ProjectDeleteCommand) Description() string { return "Delete a project (soft delete)" }
func (c *ProjectDeleteCommand) Usage() string       { return "castra project delete <id>" }

func (c *ProjectDeleteCommand) AllowedRoles() []string { return []string{"architect"} }

func (c *ProjectDeleteCommand) Execute(ctx *Context) error {

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
	c.lastID = id
	fmt.Println("Project soft deleted.")
	return nil
}

func (c *ProjectDeleteCommand) AuditInfo() (string, int64, string) {
	return "project", c.lastID, "project.delete"
}

type ProjectViewCommand struct{}

func (c *ProjectViewCommand) Name() string        { return "view" }
func (c *ProjectViewCommand) Description() string { return "View project details" }
func (c *ProjectViewCommand) Usage() string       { return "castra project view <id>" }

func (c *ProjectViewCommand) ReadInfo() (string, string) {
	return "project", "project.view"
}

func (c *ProjectViewCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("project view", flag.ExitOnError)
	fs.Parse(ctx.Args)

	if len(fs.Args()) < 1 {
		return fmt.Errorf("project ID required")
	}
	id, _ := strconv.ParseInt(fs.Args()[0], 10, 64)

	p, err := cli.GetProject(ctx.DB, id)
	if err != nil {
		return err
	}

	fmt.Printf("--- Project [%d]: %s ---\n", p.ID, p.Name)
	fmt.Printf("Status: %s\n", p.Status)
	fmt.Printf("\nDescription:\n%s\n", p.Description)
	fmt.Printf("\nNotes:\n%s\n", p.Notes)

	fmt.Println("\n--- Next Actions ---")
	printProjectNextActions(ctx.DB, p.Status, ctx.Role, id)
	return nil
}

type ProjectUpdateCommand struct{ lastID int64 }

func (c *ProjectUpdateCommand) Name() string        { return "update" }
func (c *ProjectUpdateCommand) Description() string { return "Update project status" }
func (c *ProjectUpdateCommand) Usage() string       { return "castra project update --status <active|archived> <id>" }

func (c *ProjectUpdateCommand) AllowedRoles() []string { return []string{"architect"} }

func (c *ProjectUpdateCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("project update", flag.ExitOnError)
	status := fs.String("status", "", "New Status (active|archived)")
	fs.Parse(ctx.Args)

	if len(fs.Args()) < 1 {
		return fmt.Errorf("project ID required")
	}
	id, _ := strconv.ParseInt(fs.Args()[0], 10, 64)

	if *status == "" {
		return fmt.Errorf("status required")
	}

	if err := cli.UpdateProjectStatus(ctx.DB, id, *status, ctx.Role); err != nil {
		return err
	}
	c.lastID = id
	fmt.Printf("Project %d updated to %s\n", id, *status)
	return nil
}

func (c *ProjectUpdateCommand) AuditInfo() (string, int64, string) {
	return "project", c.lastID, "project.update"
}
