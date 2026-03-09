package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
)

type LogAddCommand struct{}

func (c *LogAddCommand) Name() string        { return "add" }
func (c *LogAddCommand) Description() string { return "Add an audit log entry" }
func (c *LogAddCommand) Usage() string {
	return "castra log add --msg <message> [--type <type>] [--entity <id>]"
}

func (c *LogAddCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("log add", flag.ExitOnError)
	msg := fs.String("msg", "", "Log message")
	entityType := fs.String("type", "", "Entity type (project, sprint, task)")
	entityID := fs.Int64("entity", 0, "Entity ID")
	fs.Parse(ctx.Args)

	if *msg == "" {
		return fmt.Errorf("message required (--msg)")
	}

	return cli.AddAuditEntry(ctx.DB, *entityType, *entityID, *msg, ctx.Role, "")
}

type LogListCommand struct{}

func (c *LogListCommand) Name() string        { return "list" }
func (c *LogListCommand) Description() string { return "List audit log entries" }
func (c *LogListCommand) Usage() string       { return "castra log list [--type <type>] [--entity <id>]" }

func (c *LogListCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("log list", flag.ExitOnError)
	entityType := fs.String("type", "", "Filter by entity type")
	entityID := fs.Int64("entity", 0, "Filter by entity ID")
	fs.Parse(ctx.Args)

	var eid *int64
	if *entityID != 0 {
		eid = entityID
	}

	entries, err := cli.ListAuditEntries(ctx.DB, *entityType, eid)
	if err != nil {
		return err
	}
	for _, e := range entries {
		fmt.Printf("[%s] %s/%d: %s (%s)\n", e.Timestamp, e.EntityType, e.EntityID, e.Action, e.Role)
	}
	return nil
}
