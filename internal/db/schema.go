package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// configurePragmas sets journal mode with a sandbox-safe fallback.
// WAL mode requires shared memory (shm/mmap) which sandboxed environments
// (e.g. Antigravity) restrict, causing SQLITE_CANTOPEN (error 14).
// If WAL cannot be enabled, we fall back to DELETE journal mode.
func configurePragmas(db *sql.DB) {
	// Attempt WAL mode for better concurrency.
	var journalMode string
	if err := db.QueryRow("PRAGMA journal_mode=WAL;").Scan(&journalMode); err != nil {
		log.Printf("WAL mode request failed: %v; falling back to DELETE journal mode", err)
		journalMode = ""
	}

	// SQLite may silently refuse WAL (returning the current mode instead).
	// If we didn't get WAL, explicitly set DELETE mode as a safe fallback.
	if journalMode != "wal" {
		log.Printf("WAL mode unavailable (got %q); using DELETE journal mode for sandbox compatibility", journalMode)
		if _, err := db.Exec("PRAGMA journal_mode=DELETE;"); err != nil {
			log.Printf("Error setting DELETE journal mode: %v", err)
		}
	}

	if _, err := db.Exec(`
		PRAGMA synchronous=NORMAL;
		PRAGMA busy_timeout=5000;
	`); err != nil {
		log.Printf("Error setting PRAGMAs: %v", err)
	}
}

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	configurePragmas(db)

	if err := RunMigrations(db); err != nil {
		log.Printf("Error running migrations: %v", err)
		return nil, err
	}

	return db, nil
}
