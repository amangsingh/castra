package db

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

// TestFreshDatabase verifies that running migrations on an empty database
// creates all tables and sets the version to the latest.
func TestFreshDatabase(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	if err := RunMigrations(db); err != nil {
		t.Fatalf("RunMigrations failed on fresh DB: %v", err)
	}

	// Verify all tables exist
	tables := []string{"projects", "sprints", "tasks", "project_notes", "audit_log", "schema_version"}
	for _, table := range tables {
		var name string
		err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&name)
		if err != nil {
			t.Errorf("Table %s not found after migration: %v", table, err)
		}
	}

	// Verify version is at latest
	latestVersion := Migrations[len(Migrations)-1].Version
	version, err := getCurrentVersion(db)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}
	if version != latestVersion {
		t.Errorf("Expected version %d, got %d", latestVersion, version)
	}
}

// TestPreMigrationDatabase simulates an existing workspace.db created before
// the migration system existed, and verifies it upgrades correctly.
func TestPreMigrationDatabase(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	// Simulate a pre-migration database: tables exist but no schema_version.
	// Use the OLD schema that lacks task_id on project_notes and role on audit_log.
	oldSchema := `
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
		content TEXT NOT NULL,
		tags TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);
	CREATE TABLE audit_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		entity_type TEXT,
		entity_id INTEGER,
		action TEXT NOT NULL,
		payload TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(oldSchema); err != nil {
		t.Fatalf("Failed to create old schema: %v", err)
	}

	// Now run migrations — should detect pre-migration DB, skip v1, apply v2
	if err := RunMigrations(db); err != nil {
		t.Fatalf("RunMigrations failed on pre-migration DB: %v", err)
	}

	// Verify task_id column exists on project_notes
	_, err := db.Exec(`INSERT INTO project_notes (project_id, task_id, content, tags) VALUES (1, NULL, 'test', '')`)
	if err != nil {
		t.Errorf("task_id column missing from project_notes after migration: %v", err)
	}

	// Verify role column exists on audit_log
	_, err = db.Exec(`INSERT INTO audit_log (entity_type, entity_id, action, role, payload) VALUES ('task', 1, 'test', 'architect', 'test')`)
	if err != nil {
		t.Errorf("role column missing from audit_log after migration: %v", err)
	}

	// Verify version is at latest
	latestVersion := Migrations[len(Migrations)-1].Version
	version, err := getCurrentVersion(db)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}
	if version != latestVersion {
		t.Errorf("Expected version %d, got %d", latestVersion, version)
	}
}

// TestIdempotency verifies that running migrations multiple times is safe.
func TestIdempotency(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	if err := RunMigrations(db); err != nil {
		t.Fatalf("First RunMigrations failed: %v", err)
	}

	if err := RunMigrations(db); err != nil {
		t.Fatalf("Second RunMigrations failed (not idempotent): %v", err)
	}

	latestVersion := Migrations[len(Migrations)-1].Version
	version, err := getCurrentVersion(db)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}
	if version != latestVersion {
		t.Errorf("Expected version %d after re-run, got %d", latestVersion, version)
	}
}

// TestPartialMigration verifies that only pending migrations are applied
// when the DB is already at an intermediate version.
func TestPartialMigration(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	// Manually set up schema_version at v1 with the v1 tables already created
	if _, err := db.Exec(Migrations[0].SQL); err != nil {
		t.Fatalf("Failed to apply v1 schema manually: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_version (version INTEGER NOT NULL DEFAULT 0)`); err != nil {
		t.Fatalf("Failed to create schema_version: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO schema_version (version) VALUES (1)`); err != nil {
		t.Fatalf("Failed to set version: %v", err)
	}

	// Run migrations — should only apply v2+
	if err := RunMigrations(db); err != nil {
		t.Fatalf("RunMigrations failed on partial migration: %v", err)
	}

	latestVersion := Migrations[len(Migrations)-1].Version
	version, err := getCurrentVersion(db)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}
	if version != latestVersion {
		t.Errorf("Expected version %d, got %d", latestVersion, version)
	}
}

// TestInitDB verifies the full InitDB flow creates a working database.
func TestInitDB(t *testing.T) {
	database, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer database.Close()

	// Should be able to insert and query across all tables
	_, err = database.Exec(`INSERT INTO projects (name, description) VALUES ('Test', 'A test project')`)
	if err != nil {
		t.Errorf("Failed to insert project: %v", err)
	}

	_, err = database.Exec(`INSERT INTO tasks (project_id, title, status) VALUES (1, 'Task 1', 'todo')`)
	if err != nil {
		t.Errorf("Failed to insert task: %v", err)
	}

	_, err = database.Exec(`INSERT INTO project_notes (project_id, task_id, content, tags) VALUES (1, 1, 'note', 'architect')`)
	if err != nil {
		t.Errorf("Failed to insert note with task_id: %v", err)
	}

	_, err = database.Exec(`INSERT INTO audit_log (entity_type, entity_id, action, role, payload) VALUES ('task', 1, 'create', 'architect', 'created')`)
	if err != nil {
		t.Errorf("Failed to insert audit entry with role: %v", err)
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}
	return db
}
