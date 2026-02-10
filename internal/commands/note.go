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
		content := fs.String("content", "", "Note Content")
		tags := fs.String("tags", "", "Tags (comma-sep)")
		fs.Parse(argsToParse)

		if *pid == 0 || *content == "" {
			log.Fatal("Project ID and Content required")
		}
		id, err := cli.AddNote(db, *pid, *content, *tags)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Note added: %d\n", id)

	case "list":
		fs := flag.NewFlagSet("note list", flag.ExitOnError)
		pid := fs.Int64("project", 0, "Project ID")
		fs.Parse(argsToParse)
		if *pid == 0 {
			log.Fatal("Project ID required")
		}

		notes, err := cli.ListNotes(db, *pid, role)
		if err != nil {
			log.Fatal(err)
		}
		for _, n := range notes {
			fmt.Printf("[%d] %s [Tags: %s]\n", n.ID, n.Content, n.Tags)
		}
	}
}
