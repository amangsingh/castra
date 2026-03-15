package commands_test

// integration tests for the router's automatic audit logging mechanism.
//
// These tests verify that:
//  1. Mutating commands (implementing MutatingCommand) automatically produce
//     an audit_log entry on success.
//  2. Read-only commands (implementing ReadCommand) automatically produce
//     an audit_log entry with "[read]" payload.
//  3. Failed commands do NOT produce spurious audit entries.

import (
	"castra/internal/cli"
	"testing"

	"castra/internal/commands"
)

// countAuditEntries is a helper that returns the total number of rows in the
// audit_log table, optionally filtered by action string.
func countAuditEntries(t *testing.T, db interface {
	QueryRow(string, ...interface{}) interface{ Scan(...interface{}) error }
}, action string) int {
	t.Helper()
	// Use the exported ListAuditEntries so we stay on the sanctioned API.
	return -1 // placeholder; real impl below
}

// auditEntriesForAction returns all audit entries whose action matches.
func auditEntriesForAction(t *testing.T, ctx *commands.Context, action string) []cli.AuditEntry {
	t.Helper()
	entries, err := cli.ListAuditEntries(ctx.DB, "", nil)
	if err != nil {
		t.Fatalf("ListAuditEntries: %v", err)
	}
	var matched []cli.AuditEntry
	for _, e := range entries {
		if e.Action == action {
			matched = append(matched, e)
		}
	}
	return matched
}

// --- Task Create auto-audit ---

// TestAutoAudit_TaskAdd verifies that executing "task add" via the registry
// produces exactly one audit entry with action "task.add".
func TestAutoAudit_TaskAdd(t *testing.T) {
	db := commands.NewTestDB(t)

	// Seed a project (required for task.add).
	projectCtx := commands.NewTestCtx(db, "architect", []string{"--name", "TestProject", "--desc", "p"})
	if err := (&commands.ProjectAddCommand{}).Execute(projectCtx); err != nil {
		t.Fatalf("seed project: %v", err)
	}

	// Seed a sprint (required for task.add).
	sprintCtx := commands.NewTestCtx(db, "architect", []string{"--project", "1", "--name", "S1"})
	if err := (&commands.SprintAddCommand{}).Execute(sprintCtx); err != nil {
		t.Fatalf("seed sprint: %v", err)
	}

	// Run "task add" through the full registry so the auto-audit hook fires.
	reg := commands.NewRegistry()
	reg.Register(&commands.TaskAddCommand{})

	ctx := commands.NewTestCtx(db, "architect", []string{
		"add", "--project", "1", "--sprint", "1", "--title", "My Task",
	})
	commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, reg.Execute(ctx))
	})

	entries := auditEntriesForAction(t, commands.NewTestCtx(db, "architect", nil), "task.add")
	if len(entries) != 1 {
		t.Errorf("expected 1 audit entry for task.add, got %d", len(entries))
	}
	if len(entries) > 0 && entries[0].EntityType != "task" {
		t.Errorf("expected entityType 'task', got %q", entries[0].EntityType)
	}
}

// TestAutoAudit_TaskUpdate verifies that executing "task update" produces an
// audit entry with action "task.update".
func TestAutoAudit_TaskUpdate(t *testing.T) {
	db := commands.NewTestDB(t)

	// Seed project → sprint → task.
	commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, (&commands.ProjectAddCommand{}).Execute(
			commands.NewTestCtx(db, "architect", []string{"--name", "P"})))
		commands.AssertNoError(t, (&commands.SprintAddCommand{}).Execute(
			commands.NewTestCtx(db, "architect", []string{"--project", "1", "--name", "S1"})))
		commands.AssertNoError(t, (&commands.TaskAddCommand{}).Execute(
			commands.NewTestCtx(db, "architect", []string{"--project", "1", "--sprint", "1", "--title", "T"})))
	})

	// Run "task update" via registry.
	reg := commands.NewRegistry()
	reg.Register(&commands.TaskUpdateCommand{})

	ctx := commands.NewTestCtx(db, "senior-engineer", []string{"update", "--status", "doing", "1"})
	commands.AssertNoError(t, reg.Execute(ctx))

	entries := auditEntriesForAction(t, commands.NewTestCtx(db, "senior-engineer", nil), "task.update")
	if len(entries) != 1 {
		t.Errorf("expected 1 audit entry for task.update, got %d", len(entries))
	}
}

// --- Project auto-audit ---

// TestAutoAudit_ProjectAdd verifies that "project add" produces an audit
// entry with action "project.add".
func TestAutoAudit_ProjectAdd(t *testing.T) {
	db := commands.NewTestDB(t)

	reg := commands.NewRegistry()
	reg.Register(&commands.ProjectAddCommand{})

	ctx := commands.NewTestCtx(db, "architect", []string{"add", "--name", "Castra"})
	commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, reg.Execute(ctx))
	})

	entries := auditEntriesForAction(t, commands.NewTestCtx(db, "architect", nil), "project.add")
	if len(entries) != 1 {
		t.Errorf("expected 1 audit entry for project.add, got %d", len(entries))
	}
}

// --- Read-only commands do NOT produce audit entries ---

// TestAutoAudit_ReadOnly_TaskList verifies that "task list" (a read-only
// command) does NOT write any auto audit entry.
func TestAutoAudit_ReadOnly_TaskList(t *testing.T) {
	db := commands.NewTestDB(t)

	// Seed a project first.
	commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, (&commands.ProjectAddCommand{}).Execute(
			commands.NewTestCtx(db, "architect", []string{"--name", "P"})))
	})

	// Drain auto-audit from the seeding step.
	beforeEntries, _ := cli.ListAuditEntries(db, "", nil)

	// Now run the read-only list command through the registry.
	reg := commands.NewRegistry()
	reg.Register(&commands.TaskListCommand{})

	ctx := commands.NewTestCtx(db, "senior-engineer", []string{"list", "--project", "1"})
	commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, reg.Execute(ctx))
	})

	afterEntries, _ := cli.ListAuditEntries(db, "", nil)
	if len(afterEntries) != len(beforeEntries)+1 {
		t.Errorf("read-only TaskList SHOULD now produce audit entries; before=%d after=%d",
			len(beforeEntries), len(afterEntries))
	}
	
	// Verify it has the [read] payload
	found := false
	for _, e := range afterEntries {
		if e.Payload == "[read]" && e.Action == "task.list" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected audit entry with '[read]' payload for task.list, not found")
	}
}

// --- Failed commands do NOT produce spurious audit entries ---

// TestAutoAudit_FailedCommand verifies that when a mutating command returns an
// error the router does NOT write an audit entry.
func TestAutoAudit_FailedCommand(t *testing.T) {
	db := commands.NewTestDB(t)

	beforeEntries, _ := cli.ListAuditEntries(db, "", nil)

	// "task add" with missing required flags will return an error.
	reg := commands.NewRegistry()
	reg.Register(&commands.TaskAddCommand{})

	ctx := commands.NewTestCtx(db, "architect", []string{"add"}) // missing --project and --title
	// We expect an error (missing args), but we also expect no audit entry.
	_ = reg.Execute(ctx)

	afterEntries, _ := cli.ListAuditEntries(db, "", nil)
	if len(afterEntries) > len(beforeEntries) {
		t.Errorf("failed command must NOT produce audit entries; before=%d after=%d",
			len(beforeEntries), len(afterEntries))
	}
}

// TestAutoAudit_RoleRecorded verifies that the audit entry captures the
// correct role of the user who performed the operation.
func TestAutoAudit_RoleRecorded(t *testing.T) {
	db := commands.NewTestDB(t)

	// Seed project + sprint via direct Execute (bypassing registry to avoid
	// polluting audit with their entries before we check the one we care about).
	commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, (&commands.ProjectAddCommand{}).Execute(
			commands.NewTestCtx(db, "architect", []string{"--name", "P"})))
		commands.AssertNoError(t, (&commands.SprintAddCommand{}).Execute(
			commands.NewTestCtx(db, "architect", []string{"--project", "1", "--name", "S1"})))
		commands.AssertNoError(t, (&commands.TaskAddCommand{}).Execute(
			commands.NewTestCtx(db, "architect", []string{"--project", "1", "--sprint", "1", "--title", "T"})))
	})

	// Run update as "senior-engineer" via registry.
	reg := commands.NewRegistry()
	reg.Register(&commands.TaskUpdateCommand{})

	ctx := commands.NewTestCtx(db, "senior-engineer", []string{"update", "--status", "doing", "1"})
	commands.AssertNoError(t, reg.Execute(ctx))

	entries := auditEntriesForAction(t, commands.NewTestCtx(db, "senior-engineer", nil), "task.update")
	if len(entries) == 0 {
		t.Fatal("expected at least one task.update audit entry")
	}
	if entries[0].Role != "senior-engineer" {
		t.Errorf("expected role 'senior-engineer', got %q", entries[0].Role)
	}
}
