package commands

import (
	"castra/internal/db"
	"castra/internal/git"
	"database/sql"
	"log"
	"os"
	"path/filepath"
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
	dbPath := "workspace.db"

	// Try to find the git common directory to share state across worktrees
	if commonDir, err := git.DiscoverCommonDir(); err == nil {
		targetPath := filepath.Join(commonDir, "..", "workspace.db")
		legacyPath := "workspace.db"

		absTarget, _ := filepath.Abs(targetPath)
		absLegacy, _ := filepath.Abs(legacyPath)

		if absTarget != absLegacy {
			// Check if legacy file exists
			if _, err := os.Stat(legacyPath); err == nil {
				// Legacy exists. Check if target exists.
				if _, err := os.Stat(targetPath); os.IsNotExist(err) {
					log.Printf("Migrating legacy ledger from %s to %s", absLegacy, absTarget)
					if err := os.Rename(absLegacy, absTarget); err != nil {
						log.Printf("Warning: failed to migrate legacy ledger: %v", err)
					} else {
						// Also try to migrate journals/WAL files if they exist
						for _, suffix := range []string{"-journal", "-wal", "-shm"} {
							if _, err := os.Stat(legacyPath + suffix); err == nil {
								os.Rename(legacyPath+suffix, targetPath+suffix)
							}
						}
					}
				} else {
					log.Printf("Warning: legacy ledger found at %s, but central ledger already exists at %s. Manual merge required.", absLegacy, absTarget)
				}
			}
		}
		dbPath = targetPath
	}

	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to DB at %s: %v", dbPath, err)
	}
	return database
}
