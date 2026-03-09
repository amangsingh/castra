package commands_test

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	"castra/internal/commands"
)

// extractID parses a single int64 value from a formatted string using fmt.Sscanf.
// Returns the id and any scan error.
func extractID(s, format string, id *int64) (int64, error) {
	_, err := fmt.Sscanf(s, format, id)
	return *id, err
}

// itoa converts an int64 to its string representation (for building Args slices).
func itoa(id int64) string {
	return strconv.FormatInt(id, 10)
}

// seedProject adds a project via the architect role and returns its ID.
func seedProject(t *testing.T, db *sql.DB, name string) int64 {
	t.Helper()
	cmd := &commands.ProjectAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{"--name", name})
	var id int64
	out := commands.CaptureOutput(t, func() {
		if err := cmd.Execute(ctx); err != nil {
			t.Fatalf("seedProject(%q): %v", name, err)
		}
	})
	extractID(out, "Project added with ID: %d", &id) //nolint:errcheck
	return id
}

// seedMilestone adds a milestone to a project and returns its ID.
func seedMilestone(t *testing.T, db *sql.DB, projID int64, name string) int64 {
	t.Helper()
	cmd := &commands.MilestoneAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{
		"--project", itoa(projID), "--name", name,
	})
	var id int64
	out := commands.CaptureOutput(t, func() {
		if err := cmd.Execute(ctx); err != nil {
			t.Fatalf("seedMilestone(%q): %v", name, err)
		}
	})
	extractID(out, "Milestone added: %d", &id) //nolint:errcheck
	return id
}

// seedSprint adds a sprint to a project and returns its ID.
func seedSprint(t *testing.T, db *sql.DB, projID int64, name string) int64 {
	t.Helper()
	cmd := &commands.SprintAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{
		"--project", itoa(projID), "--name", name,
	})
	var id int64
	out := commands.CaptureOutput(t, func() {
		if err := cmd.Execute(ctx); err != nil {
			t.Fatalf("seedSprint(%q): %v", name, err)
		}
	})
	extractID(out, "Sprint added: %d", &id) //nolint:errcheck
	return id
}

// seedTask adds a task to a project+milestone+sprint and returns its ID.
func seedTask(t *testing.T, db *sql.DB, projID, milID, sprintID int64, title string) int64 {
	t.Helper()
	cmd := &commands.TaskAddCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{
		"--project", itoa(projID),
		"--milestone", itoa(milID),
		"--sprint", itoa(sprintID),
		"--title", title,
	})
	var id int64
	out := commands.CaptureOutput(t, func() {
		if err := cmd.Execute(ctx); err != nil {
			t.Fatalf("seedTask(%q): %v", title, err)
		}
	})
	extractID(out, "Task created: %d", &id) //nolint:errcheck
	return id
}
