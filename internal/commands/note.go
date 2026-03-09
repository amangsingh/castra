package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
)

type NoteAddCommand struct{}

func (c *NoteAddCommand) Name() string        { return "add" }
func (c *NoteAddCommand) Description() string { return "Add a new note" }
func (c *NoteAddCommand) Usage() string {
	return "castra note add --project <pid> [--task <tid>] --content <content> [--tags <tags>]"
}

func (c *NoteAddCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("note add", flag.ExitOnError)
	pid := fs.Int64("project", 0, "Project ID")
	tid := fs.Int64("task", 0, "Task ID (optional)")
	content := fs.String("content", "", "Note Content")
	tags := fs.String("tags", "", "Tags (comma-sep)")
	fs.Parse(ctx.Args)

	if *pid == 0 || *content == "" {
		return fmt.Errorf("project ID and content required")
	}

	var taskID *int64
	if *tid != 0 {
		taskID = tid
	}

	id, err := cli.AddNote(ctx.DB, *pid, taskID, *content, *tags)
	if err != nil {
		return err
	}
	fmt.Printf("Note added: %d\n", id)
	return nil
}

type NoteListCommand struct{}

func (c *NoteListCommand) Name() string        { return "list" }
func (c *NoteListCommand) Description() string { return "List notes for a project or task" }
func (c *NoteListCommand) Usage() string       { return "castra note list --project <pid> [--task <tid>]" }

func (c *NoteListCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("note list", flag.ExitOnError)
	pid := fs.Int64("project", 0, "Project ID")
	tid := fs.Int64("task", 0, "Task ID (optional)")
	fs.Parse(ctx.Args)

	if *pid == 0 {
		return fmt.Errorf("project ID required")
	}

	var taskID *int64
	if *tid != 0 {
		taskID = tid
	}

	notes, err := cli.ListNotes(ctx.DB, *pid, taskID, ctx.Role)
	if err != nil {
		return err
	}
	for _, n := range notes {
		taskLabel := ""
		if n.TaskID != nil {
			taskLabel = fmt.Sprintf(" (Task: %d)", *n.TaskID)
		}
		fmt.Printf("[%d]%s %s [Tags: %s]\n", n.ID, taskLabel, n.Content, n.Tags)
	}
	return nil
}
