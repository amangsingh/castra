package cli

import (
	"castra/internal/db"
	"database/sql"
	"testing"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}
	return database
}

func TestListNotes_RoleFiltering(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed notes
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
		{"doc-writer", 3},      // Sees all
		{"junior-engineer", 1}, // Sees 'junior-engineer' tagged only
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			notes, err := ListNotes(db, 1, nil, tt.role)
			if err != nil {
				t.Fatalf("ListNotes failed: %v", err)
			}
			if len(notes) != tt.expectedCount {
				t.Errorf("Role %s: expected %d notes, got %d", tt.role, tt.expectedCount, len(notes))
			}
		})
	}
}

func TestListNotes_TaskFiltering(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Seed: project-level note and task-level note
	db.Exec(`INSERT INTO project_notes (content, project_id, task_id, tags) VALUES 
		('Project note', 1, NULL, ''),
		('Task note', 1, 42, 'qa-functional')`)

	taskID := int64(42)
	notes, err := ListNotes(db, 1, &taskID, "architect")
	if err != nil {
		t.Fatalf("ListNotes failed: %v", err)
	}
	if len(notes) != 1 {
		t.Errorf("Expected 1 task-level note, got %d", len(notes))
	}
}

func TestListTasks_RoleFiltering(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

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
			tasks, err := ListTasks(db, 1, nil, nil, false, tt.role)
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

	// 3. QA approves
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

func TestUpdateTaskStatus_RejectionResetsApprovals(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	res, _ := db.Exec(`INSERT INTO tasks (title, project_id, status) VALUES ('Task 1', 1, 'review')`)
	id, _ := res.LastInsertId()

	// QA approves
	UpdateTaskStatus(db, id, "done", "qa-functional")

	var qaApp, secApp bool
	db.QueryRow("SELECT qa_approved, security_approved FROM tasks WHERE id = ?", id).Scan(&qaApp, &secApp)
	if !qaApp {
		t.Fatal("QA approval should be set")
	}

	// Security REJECTS
	err := UpdateTaskStatus(db, id, "todo", "security-ops")
	if err != nil {
		t.Fatalf("Security rejection failed: %v", err)
	}

	// Verify BOTH flags are reset
	var status string
	db.QueryRow("SELECT status, qa_approved, security_approved FROM tasks WHERE id = ?", id).Scan(&status, &qaApp, &secApp)
	if status != "todo" {
		t.Errorf("Task should be todo after rejection, got %s", status)
	}
	if qaApp {
		t.Error("QA approval should be reset after rejection")
	}
	if secApp {
		t.Error("Security approval should be reset after rejection")
	}
}
