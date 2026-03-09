package commands_test

import (
	"strings"
	"testing"

	"castra/internal/commands"
)

// --- Task 50: Sprint Command Tests ---

// TestSprintAddHappyPath adds a sprint and checks success output.
func TestSprintAddHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)
	projID := seedProject(t, db, "SprintProj")

	cmd := &commands.SprintAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{
		"--project", itoa(projID), "--name", "Sprint 1",
	})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, cmd.Execute(ctx))
	})
	commands.AssertOutputContains(t, out, "Sprint added")
}

// TestSprintAddWithDates verifies optional --start and --end flags are accepted.
func TestSprintAddWithDates(t *testing.T) {
	db := commands.NewTestDB(t)
	projID := seedProject(t, db, "DatedSprintProj")

	cmd := &commands.SprintAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{
		"--project", itoa(projID),
		"--name", "Dated Sprint",
		"--start", "2026-01-01",
		"--end", "2026-01-14",
	})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, cmd.Execute(ctx))
	})
	commands.AssertOutputContains(t, out, "Sprint added")
}

// TestSprintAddMissingProject verifies --project is required.
func TestSprintAddMissingProject(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.SprintAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{"--name", "Orphan"})
	commands.AssertError(t, cmd.Execute(ctx))
}

// TestSprintAddNonArchitect verifies RBAC on sprint creation.
func TestSprintAddNonArchitect(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.SprintAddCommand{}
	ctx := commands.NewTestCtx(db, "senior-engineer", []string{"--project", "1", "--name", "Hack"})
	err := cmd.Execute(ctx)
	if err == nil || !strings.Contains(err.Error(), "architect") {
		t.Errorf("expected architect-only error, got: %v", err)
	}
}

// TestSprintListHappyPath verifies that list returns sprints for a project.
func TestSprintListHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)
	projID := seedProject(t, db, "ListSprintProj")
	seedSprint(t, db, projID, "My Sprint")

	cmd := &commands.SprintListCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{"--project", itoa(projID)})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, cmd.Execute(ctx))
	})
	commands.AssertOutputContains(t, out, "My Sprint")
}

// TestSprintListMissingProject verifies --project is required for list.
func TestSprintListMissingProject(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.SprintListCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{})
	commands.AssertError(t, cmd.Execute(ctx))
}
