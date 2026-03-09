package commands_test

import (
	"fmt"
	"strings"
	"testing"

	"castra/internal/commands"
)

// TestNewTestDB verifies that NewTestDB returns a ready-to-use, fully-migrated
// in-memory database, and that separate calls return isolated databases.
func TestNewTestDB(t *testing.T) {
	db1 := commands.NewTestDB(t)
	db2 := commands.NewTestDB(t)

	// Both DBs should be pingable.
	if err := db1.Ping(); err != nil {
		t.Fatalf("db1 is not usable: %v", err)
	}
	if err := db2.Ping(); err != nil {
		t.Fatalf("db2 is not usable: %v", err)
	}

	// Verify core tables exist (migration ran correctly).
	tables := []string{"projects", "milestones", "sprints", "tasks", "project_notes", "audit_log"}
	for _, table := range tables {
		var name string
		err := db1.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&name)
		if err != nil {
			t.Errorf("table %q missing from test DB: %v", table, err)
		}
	}

	// Isolation check — write to db1, confirm db2 is unaffected.
	if _, err := db1.Exec(`INSERT INTO projects (name) VALUES ('onlyin1')`); err != nil {
		t.Fatalf("insert into db1 failed: %v", err)
	}
	var count int
	if err := db2.QueryRow(`SELECT COUNT(*) FROM projects WHERE name='onlyin1'`).Scan(&count); err != nil {
		t.Fatalf("query db2 failed: %v", err)
	}
	if count != 0 {
		t.Errorf("isolation broken: db2 sees db1's data (count=%d)", count)
	}
}

// TestNewTestCtx verifies that NewTestCtx correctly assembles a *Context.
func TestNewTestCtx(t *testing.T) {
	db := commands.NewTestDB(t)

	ctx := commands.NewTestCtx(db, "architect", []string{"--name", "Acme"})

	if ctx.Role != "architect" {
		t.Errorf("expected role %q, got %q", "architect", ctx.Role)
	}
	if ctx.DB != db {
		t.Errorf("ctx.DB does not point to the provided database")
	}
	if len(ctx.Args) != 2 || ctx.Args[0] != "--name" || ctx.Args[1] != "Acme" {
		t.Errorf("unexpected args: %v", ctx.Args)
	}
}

// TestCaptureOutput verifies that CaptureOutput intercepts fmt.Print* calls.
func TestCaptureOutput(t *testing.T) {
	output := commands.CaptureOutput(t, func() {
		// fmt.Println writes to os.Stdout, which CaptureOutput redirects via a pipe.
		fmt.Println("hello from test")
	})

	if !strings.Contains(output, "hello from test") {
		t.Errorf("CaptureOutput did not capture stdout; got: %q", output)
	}
}

// TestAssertNoError verifies AssertNoError does not fail when err is nil.
func TestAssertNoError(t *testing.T) {
	// If this panics/calls t.Fatal, the test itself will fail — which is the proof.
	commands.AssertNoError(t, nil)
}

// TestAssertError verifies AssertError does not fail when err is non-nil.
func TestAssertError(t *testing.T) {
	commands.AssertError(t, fmt.Errorf("some error"))
}

// TestAssertOutputContains verifies the substring helper works correctly.
func TestAssertOutputContains(t *testing.T) {
	commands.AssertOutputContains(t, "Project added with ID: 42", "ID: 42")
}
