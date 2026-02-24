package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
	"log"
	"os"
)

func HandleLog(role string) {
	subCmdIdx := GetSubcommandIndex()
	if subCmdIdx == -1 {
		fmt.Println("Usage: castra log --role <role> <add|list>")
		return
	}
	cmd := os.Args[subCmdIdx]

	db := GetDB()
	defer db.Close()
	argsToParse := FilterArgs(os.Args[subCmdIdx+1:])

	switch cmd {
	case "add":
		fs := flag.NewFlagSet("log add", flag.ExitOnError)
		msg := fs.String("msg", "", "Log message")
		entityType := fs.String("type", "", "Entity type (project, sprint, task)")
		entityID := fs.Int64("entity", 0, "Entity ID")
		fs.Parse(argsToParse)

		if *msg == "" {
			log.Fatal("Message required (--msg)")
		}

		err := cli.AddAuditEntry(db, *entityType, *entityID, *msg, role, "")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Audit entry logged.")

	case "list":
		fs := flag.NewFlagSet("log list", flag.ExitOnError)
		entityType := fs.String("type", "", "Filter by entity type")
		entityID := fs.Int64("entity", 0, "Filter by entity ID")
		fs.Parse(argsToParse)

		var eid *int64
		if *entityID != 0 {
			eid = entityID
		}

		entries, err := cli.ListAuditEntries(db, *entityType, eid)
		if err != nil {
			log.Fatal(err)
		}
		for _, e := range entries {
			fmt.Printf("[%s] %s/%d: %s (%s)\n", e.Timestamp, e.EntityType, e.EntityID, e.Action, e.Role)
		}
	}
}
