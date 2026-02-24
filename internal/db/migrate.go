package db

import (
	"database/sql"
	"fmt"
	"log"
)

// Migration represents a single schema migration step.
type Migration struct {
	Version     int
	Description string
	SQL         string              // Raw SQL to execute (used if Func is nil)
	Func        func(*sql.Tx) error // Programmatic migration (used if SQL is empty)
}

// Migrations is the ordered list of all schema migrations.
// New migrations MUST be appended to the end with incrementing version numbers.
var Migrations = []Migration{
	{
		Version:     1,
		Description: "Initial schema: projects, sprints, tasks, project_notes, audit_log",
		SQL: `
		CREATE TABLE projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL DEFAULT 'active',
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME
		);

		CREATE TABLE sprints (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			start_date DATE,
			end_date DATE,
			status TEXT NOT NULL DEFAULT 'planning',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			FOREIGN KEY(project_id) REFERENCES projects(id)
		);

		CREATE TABLE tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			sprint_id INTEGER,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL DEFAULT 'todo',
			priority TEXT DEFAULT 'medium',
			context_ref TEXT,
			qa_approved BOOLEAN DEFAULT FALSE,
			security_approved BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			FOREIGN KEY(project_id) REFERENCES projects(id),
			FOREIGN KEY(sprint_id) REFERENCES sprints(id)
		);

		CREATE TABLE project_notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			task_id INTEGER,
			content TEXT NOT NULL,
			tags TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			FOREIGN KEY(project_id) REFERENCES projects(id),
			FOREIGN KEY(task_id) REFERENCES tasks(id)
		);

		CREATE TABLE audit_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			entity_type TEXT,
			entity_id INTEGER,
			action TEXT NOT NULL,
			role TEXT,
			payload TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	},
	{
		Version:     2,
		Description: "Add task_id to project_notes and role to audit_log (for pre-migration DBs)",
		Func: func(tx *sql.Tx) error {
			if err := addColumnIfNotExists(tx, "project_notes", "task_id", "INTEGER REFERENCES tasks(id)"); err != nil {
				return err
			}
			if err := addColumnIfNotExists(tx, "audit_log", "role", "TEXT"); err != nil {
				return err
			}
			return nil
		},
	},
	{
		Version:     3,
		Description: "Add milestones table and milestone_id to tasks",
		Func: func(tx *sql.Tx) error {
			// Create milestones table
			_, err := tx.Exec(`
			CREATE TABLE IF NOT EXISTS milestones (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				project_id INTEGER NOT NULL,
				name TEXT NOT NULL,
				status TEXT NOT NULL DEFAULT 'open',
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				deleted_at DATETIME,
				FOREIGN KEY(project_id) REFERENCES projects(id)
			)`)
			if err != nil {
				return err
			}

			// Add milestone_id to tasks
			if err := addColumnIfNotExists(tx, "tasks", "milestone_id", "INTEGER REFERENCES milestones(id)"); err != nil {
				return err
			}
			return nil
		},
	},
}

// RunMigrations applies all pending migrations to the database.
func RunMigrations(db *sql.DB) error {
	// 1. Ensure the schema_version table exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_version (
		version INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		return fmt.Errorf("failed to create schema_version table: %w", err)
	}

	// 2. Get current version
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	// 3. Detect pre-migration database (tables exist but no version tracking)
	if currentVersion == 0 && coreTablesExist(db) {
		// This is an existing workspace.db from before the migration system.
		// The tables already exist, so we skip migration v1 (CREATE TABLEs).
		// We mark it as v1 and let v2+ ALTER TABLEs bring it current.
		log.Println("Detected pre-migration database. Marking as version 1.")
		if err := setVersion(db, 1); err != nil {
			return fmt.Errorf("failed to set initial version for pre-migration db: %w", err)
		}
		currentVersion = 1
	}

	// 4. Apply pending migrations
	for _, m := range Migrations {
		if m.Version <= currentVersion {
			continue
		}

		log.Printf("Applying migration v%d: %s", m.Version, m.Description)

		if err := applyMigration(db, m); err != nil {
			return fmt.Errorf("migration v%d failed: %w", m.Version, err)
		}
	}

	return nil
}

// applyMigration runs a single migration inside a transaction.
func applyMigration(db *sql.DB, m Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if m.Func != nil {
		if err := m.Func(tx); err != nil {
			return fmt.Errorf("func error: %w", err)
		}
	} else if m.SQL != "" {
		if _, err := tx.Exec(m.SQL); err != nil {
			return fmt.Errorf("SQL error: %w", err)
		}
	}

	if _, err := tx.Exec(`UPDATE schema_version SET version = ?`, m.Version); err != nil {
		return err
	}

	return tx.Commit()
}

// getCurrentVersion reads the current schema version. Returns 0 if no version is set.
func getCurrentVersion(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM schema_version`).Scan(&count)
	if err != nil {
		return 0, err
	}

	if count == 0 {
		// No version row yet — initialize it
		_, err := db.Exec(`INSERT INTO schema_version (version) VALUES (0)`)
		if err != nil {
			return 0, err
		}
		return 0, nil
	}

	var version int
	err = db.QueryRow(`SELECT version FROM schema_version`).Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

// setVersion directly sets the schema version (used for pre-migration DB bootstrapping).
func setVersion(db *sql.DB, version int) error {
	_, err := db.Exec(`UPDATE schema_version SET version = ?`, version)
	return err
}

// coreTablesExist checks if the main Castra tables already exist in the database.
// Used to detect pre-migration workspace.db files.
func coreTablesExist(db *sql.DB) bool {
	var name string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='tasks'`).Scan(&name)
	return err == nil && name == "tasks"
}

// addColumnIfNotExists adds a column to a table only if it doesn't already exist.
// SQLite lacks ALTER TABLE ADD COLUMN IF NOT EXISTS, so we check pragma table_info.
func addColumnIfNotExists(tx *sql.Tx, table, column, colDef string) error {
	rows, err := tx.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, colType string
		var notNull int
		var dfltValue *string
		var pk int
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return err
		}
		if name == column {
			return nil // Column already exists
		}
	}

	_, err = tx.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, colDef))
	return err
}
