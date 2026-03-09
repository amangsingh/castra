package gemini_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"castra/internal/generator/gemini"
	"castra/internal/generator/testgolden"
)

// TestGeminiGenerator_Golden verifies the full file structure produced by
// InitWorkspace matches the approved golden snapshot.
//
// Run with -update to regenerate golden files after intentional template changes:
//
//	go test ./internal/generator/gemini/... -update
func TestGeminiGenerator_Golden(t *testing.T) {
	testgolden.Run(t, "testdata/golden", gemini.InitWorkspace)
}

// TestGeminiGenerator_StructureChecks provides fast hermetic assertions about
// the directory structure without relying on golden file presence.
func TestGeminiGenerator_StructureChecks(t *testing.T) {
	dir := t.TempDir()
	if err := gemini.InitWorkspace(dir); err != nil {
		t.Fatalf("InitWorkspace returned error: %v", err)
	}

	expectedFiles := []string{
		"GEMINI.md",
		".gemini/settings.json",
		// Agent files — one per role
		".gemini/agents/architect.md",
		".gemini/agents/designer.md",
		".gemini/agents/doc-writer.md",
		".gemini/agents/junior-engineer.md",
		".gemini/agents/qa-functional.md",
		".gemini/agents/security-ops.md",
		".gemini/agents/senior-engineer.md",
		// Workflow files
		".gemini/workflows/build_cycle.md",
		".gemini/workflows/audit_cycle.md",
		".gemini/workflows/plan_project.md",
	}

	for _, rel := range expectedFiles {
		full := filepath.Join(dir, rel)
		if _, err := os.Stat(full); os.IsNotExist(err) {
			t.Errorf("expected file not found: %s", rel)
		}
	}
}

// TestGeminiGenerator_SettingsJsonIsValidAndDeterministic verifies that
// settings.json is valid JSON and produces identical content on repeated runs.
func TestGeminiGenerator_SettingsJsonIsValidAndDeterministic(t *testing.T) {
	runGen := func(t *testing.T) []byte {
		t.Helper()
		dir := t.TempDir()
		if err := gemini.InitWorkspace(dir); err != nil {
			t.Fatalf("InitWorkspace returned error: %v", err)
		}
		content, err := os.ReadFile(filepath.Join(dir, ".gemini/settings.json"))
		if err != nil {
			t.Fatalf("settings.json not found: %v", err)
		}
		return content
	}

	first := runGen(t)
	second := runGen(t)

	// Must be valid JSON.
	var parsed map[string]interface{}
	if err := json.Unmarshal(first, &parsed); err != nil {
		t.Errorf("settings.json is not valid JSON: %v\ncontent: %s", err, first)
	}

	// Must be deterministic across runs.
	if string(first) != string(second) {
		t.Errorf("settings.json content differs between runs:\nfirst:  %s\nsecond: %s", first, second)
	}
}

// TestGeminiGenerator_GeminiMdNotEmpty verifies GEMINI.md is written and
// non-empty (a proxy for correct rules.md → GEMINI.md mapping).
func TestGeminiGenerator_GeminiMdNotEmpty(t *testing.T) {
	dir := t.TempDir()
	if err := gemini.InitWorkspace(dir); err != nil {
		t.Fatalf("InitWorkspace returned error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "GEMINI.md"))
	if err != nil {
		t.Fatalf("GEMINI.md not found: %v", err)
	}
	if len(content) == 0 {
		t.Error("GEMINI.md is empty")
	}
}

// TestGeminiGenerator_AgentFilesAreFlatMd verifies that every file under
// .gemini/agents/ is a flat .md file (no subdirectories).
func TestGeminiGenerator_AgentFilesAreFlatMd(t *testing.T) {
	dir := t.TempDir()
	if err := gemini.InitWorkspace(dir); err != nil {
		t.Fatalf("InitWorkspace returned error: %v", err)
	}

	agentsDir := filepath.Join(dir, ".gemini/agents")
	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		t.Fatalf("cannot read agents dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal(".gemini/agents is empty")
	}

	for _, e := range entries {
		if e.IsDir() {
			t.Errorf(".gemini/agents must be flat but found directory: %s", e.Name())
		}
		if !strings.HasSuffix(e.Name(), ".md") {
			t.Errorf("unexpected file extension in .gemini/agents: %s", e.Name())
		}
	}
}

// TestGeminiGenerator_WorkflowFilesAreMd verifies all files under
// .gemini/workflows/ have the .md extension.
func TestGeminiGenerator_WorkflowFilesAreMd(t *testing.T) {
	dir := t.TempDir()
	if err := gemini.InitWorkspace(dir); err != nil {
		t.Fatalf("InitWorkspace returned error: %v", err)
	}

	wfDir := filepath.Join(dir, ".gemini/workflows")
	entries, err := os.ReadDir(wfDir)
	if err != nil {
		t.Fatalf("cannot read workflows dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal(".gemini/workflows is empty")
	}

	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".md") {
			t.Errorf("unexpected file in .gemini/workflows (expected .md): %s", e.Name())
		}
	}
}
