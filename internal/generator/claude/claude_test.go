package claude_test

import (
	"os"
	"path/filepath"
	"testing"

	"castra/internal/config"
	"castra/internal/generator/claude"
	"castra/internal/generator/testgolden"
)

func generate(baseDir string) error {
	gen := &claude.ClaudeGenerator{}
	return gen.InitWorkspace(baseDir, config.VendorConfig{})
}

func TestClaudeGenerator_Golden(t *testing.T) {
	testgolden.Run(t, "testdata/golden", generate)
}

func TestClaudeGenerator_StructureChecks(t *testing.T) {
	dir := t.TempDir()
	if err := generate(dir); err != nil {
		t.Fatalf("InitWorkspace returned error: %v", err)
	}

	expectedFiles := []string{
		"CLAUDE.md",
		".claude/agents/architect.md",
		".claude/agents/senior-engineer.md",
		".claude/workflows/build_cycle.md",
	}

	for _, rel := range expectedFiles {
		full := filepath.Join(dir, rel)
		if _, err := os.Stat(full); os.IsNotExist(err) {
			t.Errorf("expected file not found: %s", rel)
		}
	}
}
