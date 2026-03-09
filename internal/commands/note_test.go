package commands_test

import (
	"strings"
	"testing"

	"castra/internal/commands"
)

// --- Task 53: Note Command Tests ---

// TestNoteAddHappyPath adds a project-scoped note and verifies success.
func TestNoteAddHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)
	projID := seedProject(t, db, "NoteProj")

	cmd := &commands.NoteAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{
		"--project", itoa(projID),
		"--content", "Important note",
		"--tags", "architect",
	})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, cmd.Execute(ctx))
	})
	commands.AssertOutputContains(t, out, "Note added")
}

// TestNoteAddWithTask verifies the optional --task flag is accepted.
func TestNoteAddWithTask(t *testing.T) {
	db := commands.NewTestDB(t)
	projID := seedProject(t, db, "NoteTaskProj")
	milID := seedMilestone(t, db, projID, "M1")
	sprintID := seedSprint(t, db, projID, "S1")
	taskID := seedTask(t, db, projID, milID, sprintID, "Task For Note")

	cmd := &commands.NoteAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{
		"--project", itoa(projID),
		"--task", itoa(taskID),
		"--content", "Task-scoped note",
		"--tags", "engineer",
	})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, cmd.Execute(ctx))
	})
	commands.AssertOutputContains(t, out, "Note added")
}

// TestNoteAddMissingProject verifies --project is required.
func TestNoteAddMissingProject(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.NoteAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{"--content", "Orphan note"})
	commands.AssertError(t, cmd.Execute(ctx))
}

// TestNoteAddMissingContent verifies --content is required.
func TestNoteAddMissingContent(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.NoteAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{"--project", "1"})
	commands.AssertError(t, cmd.Execute(ctx))
}

// TestNoteListHappyPath verifies that list returns notes for a project.
func TestNoteListHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)
	projID := seedProject(t, db, "ListNoteProj")

	// Add a note
	addCtx := commands.NewTestCtx(db, "architect", []string{
		"--project", itoa(projID),
		"--content", "Listed note",
		"--tags", "architect",
	})
	commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, (&commands.NoteAddCommand{}).Execute(addCtx))
	})

	// List
	listCtx := commands.NewTestCtx(db, "architect", []string{"--project", itoa(projID)})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, (&commands.NoteListCommand{}).Execute(listCtx))
	})
	commands.AssertOutputContains(t, out, "Listed note")
}

// TestNoteListMissingProject verifies --project is required for list.
func TestNoteListMissingProject(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.NoteListCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{})
	err := cmd.Execute(ctx)
	if err == nil || !strings.Contains(err.Error(), "project") {
		t.Errorf("expected project-required error, got: %v", err)
	}
}
