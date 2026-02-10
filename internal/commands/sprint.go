package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
	"log"
	"os"
)

func HandleSprint(role string) {
	subCmdIdx := GetSubcommandIndex()
	if subCmdIdx == -1 {
		return
	}
	cmd := os.Args[subCmdIdx]

	if role != "architect" && cmd != "list" {
		log.Fatal("Only architect can manage sprints")
	}

	db := GetDB()
	defer db.Close()
	argsToParse := FilterArgs(os.Args[subCmdIdx+1:])

	switch cmd {
	case "add":
		fs := flag.NewFlagSet("sprint add", flag.ExitOnError)
		pid := fs.Int64("project", 0, "Project ID")
		name := fs.String("name", "", "Sprint Name")
		start := fs.String("start", "", "YYYY-MM-DD")
		end := fs.String("end", "", "YYYY-MM-DD")
		fs.Parse(argsToParse)
		if *pid == 0 || *name == "" {
			log.Fatal("Project ID and Name required")
		}
		id, err := cli.AddSprint(db, *pid, *name, *start, *end)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sprint added: %d\n", id)

	case "list":
		fs := flag.NewFlagSet("sprint list", flag.ExitOnError)
		pid := fs.Int64("project", 0, "Project ID")
		fs.Parse(argsToParse)

		if *pid == 0 {
			log.Fatal("Project ID required")
		}
		sprints, err := cli.ListSprints(db, *pid)
		if err != nil {
			log.Fatal(err)
		}
		for _, s := range sprints {
			fmt.Printf("[%d] %s (%s)\n", s.ID, s.Name, s.Status)
		}
	}
}
