package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func HandleMilestone(role string) {
	subCmdIdx := GetSubcommandIndex()
	if subCmdIdx == -1 {
		fmt.Println("Usage: castra milestone --role <role> <add|list|update|delete> ...")
		return
	}
	cmd := os.Args[subCmdIdx]

	if role != "architect" && cmd != "list" {
		log.Fatal("Only architect can manage milestones")
	}

	db := GetDB()
	defer db.Close()
	argsToParse := FilterArgs(os.Args[subCmdIdx+1:])

	switch cmd {
	case "add":
		fs := flag.NewFlagSet("milestone add", flag.ExitOnError)
		pid := fs.Int64("project", 0, "Project ID")
		name := fs.String("name", "", "Milestone Name")
		fs.Parse(argsToParse)
		if *pid == 0 || *name == "" {
			log.Fatal("Project ID and Name required")
		}
		id, err := cli.AddMilestone(db, *pid, *name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Milestone added: %d\n", id)

	case "list":
		fs := flag.NewFlagSet("milestone list", flag.ExitOnError)
		pid := fs.Int64("project", 0, "Project ID")
		fs.Parse(argsToParse)

		if *pid == 0 {
			log.Fatal("Project ID required")
		}
		milestones, err := cli.ListMilestones(db, *pid, role)
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range milestones {
			fmt.Printf("[%d] %s (%s)\n", m.ID, m.Name, m.Status)
		}

	case "update":
		fs := flag.NewFlagSet("milestone update", flag.ExitOnError)
		status := fs.String("status", "", "New Status (open|completed)")
		fs.Parse(argsToParse)

		idParsed := fs.Args()
		if len(idParsed) < 1 {
			log.Fatal("ID required")
		}
		id, _ := strconv.ParseInt(idParsed[0], 10, 64)

		if *status == "" {
			log.Fatal("Status required")
		}
		if err := cli.UpdateMilestoneStatus(db, id, *status, role); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Milestone status updated.")

	case "delete":
		fs := flag.NewFlagSet("milestone delete", flag.ExitOnError)
		fs.Parse(argsToParse)

		idParsed := fs.Args()
		if len(idParsed) < 1 {
			log.Fatal("ID required")
		}
		id, _ := strconv.ParseInt(idParsed[0], 10, 64)
		if err := cli.SoftDeleteMilestone(db, id); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Milestone deleted.")
	}
}
