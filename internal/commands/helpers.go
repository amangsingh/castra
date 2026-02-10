package commands

import (
	"castra/internal/db"
	"database/sql"
	"log"
	"os"
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

// Helper to find subcommand index, skipping flags
func GetSubcommandIndex() int {
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--role" {
			i++ // Skip value
			continue
		}
		// Found non-flag arg, likely subcommand
		if len(os.Args[i]) > 0 && os.Args[i][0] != '-' {
			return i
		}
	}
	return -1
}
