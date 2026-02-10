package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func HandleProject(role string) {
	subCmdIdx := GetSubcommandIndex()
	if subCmdIdx == -1 {
		fmt.Println("Usage: castra project --role <role> <add|list|delete> ...")
		return
	}
	cmd := os.Args[subCmdIdx]

	if role != "architect" && cmd != "list" {
		log.Fatal("Only architect can modify projects")
	}

	db := GetDB()
	defer db.Close()

	// Parse flags starting after subcommand
	argsToParse := FilterArgs(os.Args[subCmdIdx+1:])

	switch cmd {
	case "add":
		fs := flag.NewFlagSet("project add", flag.ExitOnError)
		name := fs.String("name", "", "Project Name")
		desc := fs.String("desc", "", "Description")
		notes := fs.String("notes", "", "Notes/Docs")

		fs.Parse(argsToParse)

		if *name == "" {
			log.Fatal("Name required")
		}
		id, err := cli.AddProject(db, *name, *desc, *notes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Project added with ID: %d\n", id)

	case "list":
		fs := flag.NewFlagSet("project list", flag.ExitOnError)
		archived := fs.Bool("archived", false, "Show archived")
		deleted := fs.Bool("deleted", false, "Show deleted")
		fs.Parse(argsToParse)

		projects, err := cli.ListProjects(db, *archived, *deleted)
		if err != nil {
			log.Fatal(err)
		}
		for _, p := range projects {
			fmt.Printf("[%d] %s (%s)\n", p.ID, p.Name, p.Status)
		}

	case "delete":
		fs := flag.NewFlagSet("project delete", flag.ExitOnError)
		fs.Parse(argsToParse)

		idParsed := fs.Args() // Remaining args after flags
		if len(idParsed) < 1 {
			log.Fatal("ID required")
		}
		id, _ := strconv.ParseInt(idParsed[0], 10, 64)
		if err := cli.SoftDeleteProject(db, id); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Project soft deleted.")
	}
}
