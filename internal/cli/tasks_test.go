package cli

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}

	query := `
	CREATE TABLE projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	CREATE TABLE tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER NOT NULL,
		sprint_id INTEGER,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'todo',
		priority TEXT DEFAULT 'medium',
		qa_approved BOOLEAN DEFAULT FALSE,
		security_approved BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME
	);
	CREATE TABLE project_notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		tags TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}
	return db
}

func TestListNotes_RoleFiltering(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed notes
	// 1: general (no tags)
	// 2: engineering specific (tags: junior-engineer)
	// 3: qa specific (tags: qa-functional)
	_, err := db.Exec(`INSERT INTO project_notes (content, project_id, tags) VALUES 
		('Note 1', 1, ''),
		('Note 2', 1, 'junior-engineer'),
		('Note 3', 1, 'qa-functional')`)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		role          string
		expectedCount int
	}{
		{"architect", 3},       // Sees all
		{"doc-writer", 3},      // Sees all (new requirement)
		{"junior-engineer", 1}, // Sees 'junior-engineer' only (Note 2) - Wait, logic says contains OR nothing?
		// Logic in notes.go:
		// query := `SELECT ...` -> gets ALL
		// then filtered by `containsRole`.
		// Note 1 has empty tags. `containsRole("", "junior-engineer")` -> false.
		// so it sees 1 note?
		// Let's check `containsRole` logic.
	}

	// Update expectations based on logic reading: filtering happens in Go.
	// If tags contain role => include.
	// Note 1: tags="" => exclude.
	// Note 2: tags="junior-engineer" => include.
	// Note 3: tags="qa-functional" => exclude.

	// Wait, architect/doc-writer bypass filter.

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			notes, err := ListNotes(db, 1, tt.role)
			if err != nil {
				t.Fatalf("ListNotes failed: %v", err)
			}
			if len(notes) != tt.expectedCount {
				t.Errorf("Role %s: expected %d notes, got %d", tt.role, tt.expectedCount, len(notes))
			}
		})
	}
}

func TestListTasks_RoleFiltering(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed data
	// 1: todo (visible to engineer, architect, doc-writer)
	// 2: review (visible to qa, security, architect, doc-writer)
	// 3: done (visible to architect, doc-writer)
	// 4: blocked (visible to engineer, architect, doc-writer)
	_, err := db.Exec(`INSERT INTO tasks (title, project_id, status) VALUES 
		('Task 1', 1, 'todo'),
		('Task 2', 1, 'review'),
		('Task 3', 1, 'done'),
		('Task 4', 1, 'blocked')`)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		role          string
		expectedCount int
	}{
		{"architect", 4},
		{"doc-writer", 4},
		{"junior-engineer", 2}, // todo, blocked
		{"senior-engineer", 2}, // todo, blocked
		{"qa-functional", 1},   // review
		{"security-ops", 1},    // review
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			tasks, err := ListTasks(db, 1, nil, false, tt.role)
			if err != nil {
				t.Fatalf("ListTasks failed: %v", err)
			}
			if len(tasks) != tt.expectedCount {
				t.Errorf("Role %s: expected %d tasks, got %d", tt.role, tt.expectedCount, len(tasks))
			}
		})
	}
}

func TestUpdateTaskStatus_RoleRestrictions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create task
	res, _ := db.Exec(`INSERT INTO tasks (title, project_id, status) VALUES ('Task 1', 1, 'todo')`)
	id, _ := res.LastInsertId()

	// 1. Engineer cannot mark done
	err := UpdateTaskStatus(db, id, "done", "junior-engineer")
	if err == nil {
		t.Error("Expected error when engineer marks done, got nil")
	}

	// 2. Engineer can mark review
	err = UpdateTaskStatus(db, id, "review", "junior-engineer")
	if err != nil {
		t.Errorf("Engineer failed to mark review: %v", err)
	}

	// 3. QA/Sec can only work on review (already review)
	// Try to mark done by QA
	err = UpdateTaskStatus(db, id, "done", "qa-functional")
	if err != nil {
		t.Errorf("QA failed to approve: %v", err)
	}

	// Verify still review (needs sec)
	var status string
	var qaApp, secApp bool
	db.QueryRow("SELECT status, qa_approved, security_approved FROM tasks WHERE id = ?", id).Scan(&status, &qaApp, &secApp)
	if status != "review" {
		t.Errorf("Task should still be review, got %s", status)
	}
	if !qaApp {
		t.Error("QA approval not set")
	}

	// 4. Sec approves
	err = UpdateTaskStatus(db, id, "done", "security-ops")
	if err != nil {
		t.Errorf("Sec failed to approve: %v", err)
	}

	// Verify done
	db.QueryRow("SELECT status, qa_approved, security_approved FROM tasks WHERE id = ?", id).Scan(&status, &qaApp, &secApp)
	if status != "done" {
		t.Errorf("Task should be done, got %s", status)
	}
	if !secApp {
		t.Error("Sec approval not set")
	}
}
