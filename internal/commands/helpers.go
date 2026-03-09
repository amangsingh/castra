package commands

import (
	"castra/internal/db"
	"database/sql"
	"log"
)

// Helper to remove global flags from args for subcommands
func FilterArgs(args []string) []string {
	var filtered []string
	skipNext := false
	for _, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}
		if arg == "--role" {
			skipNext = true
			continue
		}
		filtered = append(filtered, arg)
	}
	return filtered
}

func GetDB() *sql.DB {
	database, err := db.InitDB("workspace.db")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	return database
}
