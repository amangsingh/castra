package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type ArchetypeAddCommand struct{ lastID int64 }

func (c *ArchetypeAddCommand) Name() string        { return "add" }
func (c *ArchetypeAddCommand) Description() string { return "Add a new task archetype" }
func (c *ArchetypeAddCommand) Usage() string {
	return "castra archetype add --name <name> [--project <pid>] [--desc <description>] --role <default_role> --statuses <s1,s2,s3>"
}

func (c *ArchetypeAddCommand) Execute(ctx *Context) error {
	if ctx.Role != "architect" {
		return fmt.Errorf("only architect can add archetypes")
	}

	fs := flag.NewFlagSet("archetype add", flag.ExitOnError)
	name := fs.String("name", "", "Archetype Name")
	pid := fs.Int64("project", 0, "Project ID (optional)")
	desc := fs.String("desc", "", "Description")
	role := fs.String("role", "", "Default Role")
	statuses := fs.String("statuses", "", "Comma-separated ordered status list (e.g. todo,doing,review,done)")
	fs.Parse(ctx.Args)

	if *name == "" || *statuses == "" {
		return fmt.Errorf("name and statuses are required")
	}

	statusList := strings.Split(*statuses, ",")
	for i, s := range statusList {
		statusList[i] = strings.TrimSpace(s)
	}

	var projectID *int64
	if *pid != 0 {
		projectID = pid
	}

	id, err := cli.AddArchetype(ctx.DB, projectID, *name, *desc, *role, statusList)
	if err != nil {
		return err
	}
	c.lastID = id
	fmt.Printf("Archetype created: %d\n", id)
	return nil
}

func (c *ArchetypeAddCommand) AuditInfo() (string, int64, string) {
	return "archetype", c.lastID, "archetype.add"
}

type ArchetypeListCommand struct{}

func (c *ArchetypeListCommand) Name() string        { return "list" }
func (c *ArchetypeListCommand) Description() string { return "List all task archetypes" }
func (c *ArchetypeListCommand) Usage() string       { return "castra archetype list" }

func (c *ArchetypeListCommand) ReadInfo() (string, string) {
	return "archetype", "archetype.list"
}

func (c *ArchetypeListCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("archetype list", flag.ExitOnError)
	pid := fs.Int64("project", 0, "Project ID (optional)")
	fs.Parse(ctx.Args)

	var projectID *int64
	if *pid != 0 {
		projectID = pid
	}

	archetypes, err := cli.ListArchetypes(ctx.DB, projectID)
	if err != nil {
		return err
	}

	if len(archetypes) == 0 {
		fmt.Println("No archetypes found.")
		return nil
	}

	for _, a := range archetypes {
		scope := "global"
		if a.ProjectID != nil {
			scope = fmt.Sprintf("project:%d", *a.ProjectID)
		}
		fmt.Printf("[%d] %s (%s, Role: %s, Statuses: %s)\n", a.ID, a.Name, scope, a.DefaultRole, strings.Join(a.Statuses, "→"))
		if a.Description != "" {
			fmt.Printf("    %s\n", a.Description)
		}
	}
	return nil
}

type ArchetypeDeleteCommand struct{ lastID int64 }

func (c *ArchetypeDeleteCommand) Name() string        { return "delete" }
func (c *ArchetypeDeleteCommand) Description() string { return "Delete an archetype (soft delete)" }
func (c *ArchetypeDeleteCommand) Usage() string       { return "castra archetype delete <id>" }

func (c *ArchetypeDeleteCommand) Execute(ctx *Context) error {
	if ctx.Role != "architect" {
		return fmt.Errorf("only architect can delete archetypes")
	}

	fs := flag.NewFlagSet("archetype delete", flag.ExitOnError)
	fs.Parse(ctx.Args)

	if len(fs.Args()) < 1 {
		return fmt.Errorf("archetype ID required")
	}
	id, _ := strconv.ParseInt(fs.Args()[0], 10, 64)

	err := cli.SoftDeleteArchetype(ctx.DB, id)
	if err != nil {
		return err
	}
	c.lastID = id
	fmt.Printf("Archetype %d soft deleted.\n", id)
	return nil
}

func (c *ArchetypeDeleteCommand) AuditInfo() (string, int64, string) {
	return "archetype", c.lastID, "archetype.delete"
}
