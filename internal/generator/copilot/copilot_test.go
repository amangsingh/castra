package copilot_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"castra/internal/generator/copilot"
	"castra/internal/generator/testgolden"
)

// TestCopilotGenerator_Golden verifies the full file structure produced by
// InitWorkspace matches the approved golden snapshot.
//
// Run with -update to regenerate golden files after intentional template changes:
//
//	go test ./internal/generator/copilot/... -update
func TestCopilotGenerator_Golden(t *testing.T) {
	testgolden.Run(t, "testdata/golden", copilot.InitWorkspace)
}

// TestCopilotGenerator_StructureChecks provides fast hermetic assertions
// about the directory structure without relying on golden file presence.
func TestCopilotGenerator_StructureChecks(t *testing.T) {
	dir := t.TempDir()
	if err := copilot.InitWorkspace(dir); err != nil {
		t.Fatalf("InitWorkspace returned error: %v", err)
	}

	expectedFiles := []string{
		".github/copilot-instructions.md",
		// Agent files — one per role
		".github/agents/architect.md",
		".github/agents/designer.md",
		".github/agents/doc-writer.md",
		".github/agents/junior-engineer.md",
		".github/agents/qa-functional.md",
		".github/agents/security-ops.md",
		".github/agents/senior-engineer.md",
		// Workflow files
		".github/castra-workflows/build_cycle.md",
		".github/castra-workflows/audit_cycle.md",
		".github/castra-workflows/plan_project.md",
	}

	for _, rel := range expectedFiles {
		full := filepath.Join(dir, rel)
		if _, err := os.Stat(full); os.IsNotExist(err) {
			t.Errorf("expected file not found: %s", rel)
		}
	}
}

// TestCopilotGenerator_AgentFilesAreFlatMd verifies that every file under
// .github/agents/ is a flat .md file (no subdirectories).
func TestCopilotGenerator_AgentFilesAreFlatMd(t *testing.T) {
	dir := t.TempDir()
	if err := copilot.InitWorkspace(dir); err != nil {
		t.Fatalf("InitWorkspace returned error: %v", err)
	}

	agentsDir := filepath.Join(dir, ".github/agents")
	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		t.Fatalf("cannot read agents dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal(".github/agents is empty — expected one .md file per role")
	}

	for _, e := range entries {
		if e.IsDir() {
			t.Errorf(".github/agents must be flat but found directory: %s", e.Name())
			continue
		}
		if !strings.HasSuffix(e.Name(), ".md") {
			t.Errorf("unexpected file extension in .github/agents: %s", e.Name())
		}
	}
}

// TestCopilotGenerator_CopilotInstructionsNotEmpty verifies that
// copilot-instructions.md is written and non-empty.
func TestCopilotGenerator_CopilotInstructionsNotEmpty(t *testing.T) {
	dir := t.TempDir()
	if err := copilot.InitWorkspace(dir); err != nil {
		t.Fatalf("InitWorkspace returned error: %v", err)
	}

	path := filepath.Join(dir, ".github/copilot-instructions.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("copilot-instructions.md not found: %v", err)
	}
	if len(content) == 0 {
		t.Error("copilot-instructions.md is empty")
	}
}

// TestCopilotGenerator_WorkflowFilesAreMd verifies all files under
// .github/castra-workflows/ have the .md extension.
func TestCopilotGenerator_WorkflowFilesAreMd(t *testing.T) {
	dir := t.TempDir()
	if err := copilot.InitWorkspace(dir); err != nil {
		t.Fatalf("InitWorkspace returned error: %v", err)
	}

	wfDir := filepath.Join(dir, ".github/castra-workflows")
	entries, err := os.ReadDir(wfDir)
	if err != nil {
		t.Fatalf("cannot read castra-workflows dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal(".github/castra-workflows is empty")
	}

	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".md") {
			t.Errorf("unexpected file in castra-workflows (expected .md): %s", e.Name())
		}
	}
}
