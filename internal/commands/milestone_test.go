package commands_test

import (
	"strings"
	"testing"

	"castra/internal/commands"
)

// --- Task 51: Milestone Command Tests ---

// TestMilestoneAddHappyPath adds a milestone and verifies success output.
func TestMilestoneAddHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)
	projID := seedProject(t, db, "MilProj")

	cmd := &commands.MilestoneAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{
		"--project", itoa(projID), "--name", "v1.0",
	})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, cmd.Execute(ctx))
	})
	commands.AssertOutputContains(t, out, "Milestone added")
}

// TestMilestoneAddMissingProject verifies --project is required.
func TestMilestoneAddMissingProject(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.MilestoneAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{"--name", "Orphan"})
	commands.AssertError(t, cmd.Execute(ctx))
}

// TestMilestoneAddMissingName verifies --name is required.
func TestMilestoneAddMissingName(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.MilestoneAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{"--project", "1"})
	commands.AssertError(t, cmd.Execute(ctx))
}

// TestMilestoneAddNonArchitect verifies Persona Linter rejection on milestone creation.
func TestMilestoneAddNonArchitect(t *testing.T) {
	db := commands.NewTestDB(t)
	reg := commands.NewRegistry()
	reg.Register(&commands.MilestoneAddCommand{})
	ctx := commands.NewTestCtx(db, "junior-engineer", []string{"add", "--project", "1", "--name", "Rogue"})
	err := reg.Execute(ctx)
	if err == nil || !strings.Contains(err.Error(), "Outside my jurisdiction") {
		t.Errorf("expected persona rejection error, got: %v", err)
	}
}

// TestMilestoneListHappyPath verifies that list returns milestones for a project.
func TestMilestoneListHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)
	projID := seedProject(t, db, "ListMilProj")
	seedMilestone(t, db, projID, "Release 1")

	cmd := &commands.MilestoneListCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{"--project", itoa(projID)})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, cmd.Execute(ctx))
	})
	commands.AssertOutputContains(t, out, "Release 1")
}

// TestMilestoneListMissingProject verifies --project is required for list.
func TestMilestoneListMissingProject(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.MilestoneListCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{})
	commands.AssertError(t, cmd.Execute(ctx))
}

// TestMilestoneUpdateHappyPath transitions a milestone from open to completed.
func TestMilestoneUpdateHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)
	projID := seedProject(t, db, "UpdateMilProj")
	milID := seedMilestone(t, db, projID, "To Complete")

	cmd := &commands.MilestoneUpdateCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{
		"--status", "completed", itoa(milID),
	})
	commands.AssertNoError(t, cmd.Execute(ctx))
}

// TestMilestoneUpdateMissingStatus verifies --status is required.
func TestMilestoneUpdateMissingStatus(t *testing.T) {
	db := commands.NewTestDB(t)
	projID := seedProject(t, db, "BadUpdateMilProj")
	milID := seedMilestone(t, db, projID, "Stuck")

	cmd := &commands.MilestoneUpdateCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{itoa(milID)})
	commands.AssertError(t, cmd.Execute(ctx))
}

// TestMilestoneUpdateNonArchitect verifies RBAC on update.
func TestMilestoneUpdateNonArchitect(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.MilestoneUpdateCommand{}
	ctx := commands.NewTestCtx(db, "qa-functional", []string{"--status", "completed", "1"})
	err := cmd.Execute(ctx)
	if err == nil || !strings.Contains(err.Error(), "architect") {
		t.Errorf("expected architect-only error, got: %v", err)
	}
}

// TestMilestoneDeleteHappyPath soft-deletes a milestone.
func TestMilestoneDeleteHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)
	projID := seedProject(t, db, "DelMilProj")
	milID := seedMilestone(t, db, projID, "To Zap")

	cmd := &commands.MilestoneDeleteCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{itoa(milID)})
	commands.AssertNoError(t, cmd.Execute(ctx))
}

// TestMilestoneDeleteMissingID verifies that an ID is required.
func TestMilestoneDeleteMissingID(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.MilestoneDeleteCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{})
	commands.AssertError(t, cmd.Execute(ctx))
}

// TestMilestoneDeleteNonArchitect verifies Persona Linter rejection on delete.
func TestMilestoneDeleteNonArchitect(t *testing.T) {
	db := commands.NewTestDB(t)
	reg := commands.NewRegistry()
	reg.Register(&commands.MilestoneDeleteCommand{})
	ctx := commands.NewTestCtx(db, "security-ops", []string{"delete", "1"})
	err := reg.Execute(ctx)
	if err == nil || !strings.Contains(err.Error(), "Outside my jurisdiction") {
		t.Errorf("expected persona rejection error, got: %v", err)
	}
}
