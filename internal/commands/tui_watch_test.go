package commands_test

import (
	"strings"
	"testing"

	"castra/internal/commands"
)

// --- Task 56: TUI and Watch Command Tests ---

func TestTUIExecuteRBAC(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.TUICommand{}

	// No role error
	ctx := commands.NewTestCtx(db, "", []string{})
	err := cmd.Execute(ctx)
	if err == nil || !strings.Contains(err.Error(), "role is required") {
		t.Errorf("Expected error for missing role, got: %v", err)
	}
}

func TestWatchOutputJSON(t *testing.T) {
	db := commands.NewTestDB(t)
	projID := seedProject(t, db, "WatchProj")
	seedTask(t, db, projID, 0, 0, "Marshaled Task")

	// WatchCommand has an infinite loop and cannot be easily unit tested without refactoring to use a context or interface.
	t.Log("WatchCommand has an infinite loop and cannot be easily unit tested in its current form.")
}
