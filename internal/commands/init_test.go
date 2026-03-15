package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"castra/internal/commands"
	_ "castra/internal/generator/antigravity"
	_ "castra/internal/generator/claude"
	_ "castra/internal/generator/copilot"
	_ "castra/internal/generator/gemini"
)

// Tests for the castra.yaml-driven init command (Tasks 108–110).

func TestInitGeneratesTemplateWhenNoConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "castra-init-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	origCwd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origCwd)

	db := commands.NewTestDB(t)
	cmd := &commands.InitCommand{}

	// No castra.yaml present — should generate template and print instructions.
	ctx := commands.NewTestCtx(db, "architect", []string{})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, cmd.Execute(ctx))
	})

	commands.AssertOutputContains(t, out, "castra.yaml not found")
	commands.AssertOutputContains(t, out, "template has been generated")

	// castra.yaml should now exist.
	if _, err := os.Stat(filepath.Join(tempDir, "castra.yaml")); os.IsNotExist(err) {
		t.Error("Expected castra.yaml to be generated")
	}
}

func TestInitFromConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "castra-init-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	origCwd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origCwd)

	// Write a minimal castra.yaml.
	cfg := "antigravity:\n  roles:\n    - architect\n    - senior-engineer\n"
	if err := os.WriteFile(filepath.Join(tempDir, "castra.yaml"), []byte(cfg), 0644); err != nil {
		t.Fatalf("Failed to write castra.yaml: %v", err)
	}

	db := commands.NewTestDB(t)
	cmd := &commands.InitCommand{}

	ctx := commands.NewTestCtx(db, "architect", []string{})
	out := commands.CaptureOutput(t, func() {
		commands.AssertNoError(t, cmd.Execute(ctx))
	})

	commands.AssertOutputContains(t, out, "Castra workspace initialized")

	// Verify .agent directory was created.
	if _, err := os.Stat(filepath.Join(tempDir, ".agent")); os.IsNotExist(err) {
		t.Error("Expected .agent directory to be created")
	}
	// Verify architect skill was created.
	if _, err := os.Stat(filepath.Join(tempDir, ".agent/skills/architect")); os.IsNotExist(err) {
		t.Error("Expected .agent/skills/architect to be created")
	}
	// Verify senior-engineer skill was created.
	if _, err := os.Stat(filepath.Join(tempDir, ".agent/skills/senior-engineer")); os.IsNotExist(err) {
		t.Error("Expected .agent/skills/senior-engineer to be created")
	}
}

func TestInitFromConfigUnknownRole(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "castra-init-bad-role-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	origCwd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origCwd)

	cfg := "antigravity:\n  roles:\n    - nonexistent-role\n"
	if err := os.WriteFile(filepath.Join(tempDir, "castra.yaml"), []byte(cfg), 0644); err != nil {
		t.Fatalf("Failed to write castra.yaml: %v", err)
	}

	db := commands.NewTestDB(t)
	cmd := &commands.InitCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{})
	commands.AssertError(t, cmd.Execute(ctx))
}

func TestInitFromConfigEmptyRoles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "castra-init-empty-roles-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	origCwd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origCwd)

	// castra.yaml with all roles commented out.
	cfg := "antigravity:\n  roles:\n    # - architect\n"
	if err := os.WriteFile(filepath.Join(tempDir, "castra.yaml"), []byte(cfg), 0644); err != nil {
		t.Fatalf("Failed to write castra.yaml: %v", err)
	}

	db := commands.NewTestDB(t)
	cmd := &commands.InitCommand{}
	ctx := commands.NewTestCtx(db, "architect", []string{})
	commands.AssertError(t, cmd.Execute(ctx))
}
