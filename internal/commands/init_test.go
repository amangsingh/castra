package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"castra/internal/commands"
)

// --- Task 55: Init Command Tests ---

func TestInitHappyPaths(t *testing.T) {
	// Setup temp directory
	tempDir, err := os.MkdirTemp("", "castra-init-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change working directory to temp dir for the test
	origCwd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origCwd)

	db := commands.NewTestDB(t)
	cmd := &commands.InitCommand{}

	// Test Antigravity Platform
	t.Run("Antigravity", func(t *testing.T) {
		ctx := commands.NewTestCtx(db, "architect", []string{"--antigravity"})
		out := commands.CaptureOutput(t, func() {
			commands.AssertNoError(t, cmd.Execute(ctx))
		})
		commands.AssertOutputContains(t, out, "initialized for Antigravity")

		// Verify directory structure (based on what antigravitygen.InitWorkspace does)
		// Usually .agent/skills etc.
		if _, err := os.Stat(filepath.Join(tempDir, ".agent")); os.IsNotExist(err) {
			t.Error("Expected .agent directory to be created")
		}
	})

	// Test Gemini Platform
	t.Run("Gemini", func(t *testing.T) {
		ctx := commands.NewTestCtx(db, "architect", []string{"--gemini"})
		out := commands.CaptureOutput(t, func() {
			commands.AssertNoError(t, cmd.Execute(ctx))
		})
		commands.AssertOutputContains(t, out, "initialized for Gemini")

		if _, err := os.Stat(filepath.Join(tempDir, ".gemini")); os.IsNotExist(err) {
			t.Error("Expected .gemini directory to be created")
		}
	})
}

func TestInitSadPaths(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &commands.InitCommand{}

	// No flags
	ctx1 := commands.NewTestCtx(db, "architect", []string{})
	commands.AssertError(t, cmd.Execute(ctx1))

	// Invalid flag (handled by flag.ExitOnError usually, but Execute might return error)
	// Note: flag.FlagSet with ExitOnError will exit the process if it's not captured.
	// But in our command implementation we don't handle error from fs.Parse if ExitOnError is set.
	// Actually init.go uses ExitOnError. This might be hard to test for "invalid flag"
	// without redirecting stderr and catching the panic/exit if we don't change ExitOnError.
}
