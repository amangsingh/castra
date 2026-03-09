package commands_test

import (
	"strings"
	"testing"

	"castra/internal/commands"
)

// --- Task 49: Project Command Tests ---

// TestProjectAddHappyPath adds a project and verifies the success output.
func TestProjectAddHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.ProjectAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{"--name", "Alpha", "--desc", "First project"})

	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, cmd.Execute(ctx))
	})
	commands.AssertOutputContains(t, out, "Project added with ID")
}

// TestProjectAddMissingName verifies that --name is required.
func TestProjectAddMissingName(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.ProjectAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{"--desc", "No name given"})
	commands.AssertError(t, cmd.Execute(ctx))
}

// TestProjectAddNonArchitect verifies RBAC — only architect can add projects.
func TestProjectAddNonArchitect(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.ProjectAddCommand{}
	ctx := commands.NewTestCtx(db, "senior-engineer", []string{"--name", "Sneaky"})
	err := cmd.Execute(ctx)
	if err == nil || !strings.Contains(err.Error(), "architect") {
		t.Errorf("expected architect-only error, got: %v", err)
	}
}

// TestProjectListHappyPath inserts a project via add and verifies list returns it.
func TestProjectListHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)

	// Add a project first
	addCmd := &commands.ProjectAddCommand{}
	addCtx := commands.NewTestCtx(db, "architect", []string{"--name", "ListMe"})
	commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, addCmd.Execute(addCtx))
	})

	// Now list
	listCmd := &commands.ProjectListCommand{}
	listCtx := commands.NewTestCtx(db, "senior-engineer", []string{})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, listCmd.Execute(listCtx))
	})
	commands.AssertOutputContains(t, out, "ListMe")
}

// TestProjectDeleteHappyPath soft-deletes an existing project.
func TestProjectDeleteHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)

	// Create a project to delete
	addCmd := &commands.ProjectAddCommand{}
	addCtx := commands.NewTestCtx(db, "architect", []string{"--name", "ToDelete"})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, addCmd.Execute(addCtx))
	})

	// Extract ID from output like "Project added with ID: 1"
	var id int64
	if _, err := extractID(out, "Project added with ID: %d", &id); err != nil {
		// Fall back: just try ID 1
		id = 1
	}

	delCmd := &commands.ProjectDeleteCommand{}
	delCtx := commands.NewTestCtx(db, "architect", []string{itoa(id)})
	out2 := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, delCmd.Execute(delCtx))
	})
	commands.AssertOutputContains(t, out2, "soft deleted")
}

// TestProjectDeleteNonArchitect verifies RBAC on delete.
func TestProjectDeleteNonArchitect(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.ProjectDeleteCommand{}
	ctx := commands.NewTestCtx(db, "qa-functional", []string{"1"})
	err := cmd.Execute(ctx)
	if err == nil || !strings.Contains(err.Error(), "architect") {
		t.Errorf("expected architect-only error, got: %v", err)
	}
}
