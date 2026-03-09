package commands_test

import (
	"testing"

	"castra/internal/commands"
)

// --- Task 54: Log Command Tests ---

func TestLogAddHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.LogAddCommand{}

	ctx := commands.NewTestCtx(db, "architect", []string{
		"--msg", "Something happened",
		"--type", "project",
		"--entity", "1",
	})

	commands.AssertNoError(t, cmd.Execute(ctx))
}

func TestLogAddMissingMsg(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.LogAddCommand{}

	ctx := commands.NewTestCtx(db, "architect", []string{"--type", "project"})
	commands.AssertError(t, cmd.Execute(ctx))
}

func TestLogListHappyPath(t *testing.T) {
	db := commands.NewTestDB(t)

	// Add a log entry
	addCmd := &commands.LogAddCommand{}
	addCtx := commands.NewTestCtx(db, "architect", []string{"--msg", "Audit this"})
	commands.AssertNoError(t, addCmd.Execute(addCtx))

	// List log entries
	listCmd := &commands.LogListCommand{}
	listCtx := commands.NewTestCtx(db, "architect", []string{})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, listCmd.Execute(listCtx))
	})

	commands.AssertOutputContains(t, out, "Audit this")
}

func TestLogListFiltering(t *testing.T) {
	db := commands.NewTestDB(t)

	addCmd := &commands.LogAddCommand{}
	addCtx := commands.NewTestCtx(db, "architect", []string{"--msg", "Secret", "--type", "task", "--entity", "42"})
	commands.AssertNoError(t, addCmd.Execute(addCtx))

	listCmd := &commands.LogListCommand{}
	listCtx := commands.NewTestCtx(db, "architect", []string{"--type", "task", "--entity", "42"})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, listCmd.Execute(listCtx))
	})

	commands.AssertOutputContains(t, out, "Secret")
	commands.AssertOutputContains(t, out, "task/42")
}
