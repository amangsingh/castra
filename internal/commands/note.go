package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
	"log"
	"os"
)

func HandleNote(role string) {
	subCmdIdx := GetSubcommandIndex()
	if subCmdIdx == -1 {
		return
	}
	cmd := os.Args[subCmdIdx]

	db := GetDB()
	defer db.Close()
	argsToParse := FilterArgs(os.Args[subCmdIdx+1:])

	switch cmd {
	case "add":
		fs := flag.NewFlagSet("note add", flag.ExitOnError)
		pid := fs.Int64("project", 0, "Project ID")
		tid := fs.Int64("task", 0, "Task ID (optional, for task-level notes)")
		content := fs.String("content", "", "Note Content")
		tags := fs.String("tags", "", "Tags (comma-sep)")
		fs.Parse(argsToParse)

		if *pid == 0 || *content == "" {
			log.Fatal("Project ID and Content required")
		}

		var taskID *int64
		if *tid != 0 {
			taskID = tid
		}

		id, err := cli.AddNote(db, *pid, taskID, *content, *tags)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Note added: %d\n", id)

	case "list":
		fs := flag.NewFlagSet("note list", flag.ExitOnError)
		pid := fs.Int64("project", 0, "Project ID")
		tid := fs.Int64("task", 0, "Task ID (optional, filter by task)")
		fs.Parse(argsToParse)
		if *pid == 0 {
			log.Fatal("Project ID required")
		}

		var taskID *int64
		if *tid != 0 {
			taskID = tid
		}

		notes, err := cli.ListNotes(db, *pid, taskID, role)
		if err != nil {
			log.Fatal(err)
		}
		for _, n := range notes {
			taskLabel := ""
			if n.TaskID != nil {
				taskLabel = fmt.Sprintf(" (Task: %d)", *n.TaskID)
			}
			fmt.Printf("[%d]%s %s [Tags: %s]\n", n.ID, taskLabel, n.Content, n.Tags)
		}
	}
}
