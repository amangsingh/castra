package cli

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestPersonaAudit(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Initialize schema
	schema := `
	CREATE TABLE tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER,
		status TEXT,
		qa_approved BOOLEAN DEFAULT FALSE,
		security_approved BOOLEAN DEFAULT FALSE,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE audit_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		entity_type TEXT,
		entity_id INTEGER,
		action TEXT,
		role TEXT,
		payload TEXT,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatal(err)
	}

	// Insert a test task
	res, _ := db.Exec(`INSERT INTO tasks (project_id, status) VALUES (1, 'review')`)
	taskID, _ := res.LastInsertId()

	statuses := []string{"todo", "doing", "review", "done"}

	t.Run("Unauthorized Lifecycle Action", func(t *testing.T) {
		// doc-writer tries to move task
		_, proceed, _ := handleTaskApprovals(db, taskID, "todo", "doing", "doc-writer", false, false, false, "", statuses)
		if proceed {
			t.Errorf("expected proceed to be false for unauthorized doc-writer")
		}
		// Note: The actual reset logic for non-compliance in handleTaskApprovals.breakGlass 
		// returns the error. Lifecycle checks in handleTaskApprovals are implicitly 
		// guarded by isGatekeeper. But engine-level reset is in UpdateTaskStatus.
	})

	t.Run("Unauthorized Break-Glass", func(t *testing.T) {
		_, proceed, err := handleTaskApprovals(db, taskID, "review", "done", "junior-engineer", false, false, true, "I want to be architect", statuses)
		if err == nil {
			t.Error("expected error for junior-engineer break-glass")
		}
		if proceed {
			t.Error("expected proceed=false for unauthorized break-glass")
		}

		// Verify reset to todo in log
		var lastStatus string
		_ = db.QueryRow(`SELECT status FROM tasks WHERE id = ?`, taskID).Scan(&lastStatus)
		if lastStatus != "todo" {
			t.Errorf("expected status reset to todo, got %s", lastStatus)
		}
	})

	t.Run("Authorized Auditor Action", func(t *testing.T) {
		newStatus, proceed, err := handleTaskApprovals(db, taskID, "review", "done", "qa-functional", false, false, false, "", statuses)
		if err != nil {
			t.Errorf("unexpected error for qa-functional: %v", err)
		}
		if proceed {
			t.Error("expected proceed=false (waiting for security)")
		}
		if newStatus != "done" {
			t.Errorf("expected newStatus=done, got %s", newStatus)
		}
	})
}
