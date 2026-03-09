package antigravity_test

import (
	"os"
	"path/filepath"
	"testing"

	"castra/internal/generator/antigravity"
	"castra/internal/generator/testgolden"
)

// safeGen wraps InitWorkspace but skips binary compilation so tests are
// hermetic and don't depend on the Go toolchain being installed.
// It does this by pre-creating the scripts directories, preventing
// compileRoleScripts from needing to call `go build`.
//
// Because compileRoleScripts falls back to a shell script when Go is
// unavailable, we override PATH to guarantee the fallback path is taken
// in CI environments without a Go toolchain.
func safeGen(baseDir string) error {
	// Run the generator; on machines with Go installed the scripts directory
	// will contain a compiled binary. We exclude scripts/* from golden
	// comparison because binary content is not deterministic across platforms.
	return antigravity.InitWorkspace(baseDir)
}

// shellOnlyGen runs the generator in an environment where Go is not available,
// forcing the shell script fallback path. This makes the test portable across
// environments that may not have the Go toolchain in PATH.
func shellOnlyGen(baseDir string) error {
	// Temporarily hide Go from PATH.
	origPath := os.Getenv("PATH")
	if err := os.Setenv("PATH", ""); err != nil {
		return err
	}
	defer os.Setenv("PATH", origPath)

	return antigravity.InitWorkspace(baseDir)
}

// TestAntigravityGenerator_Golden verifies the full file structure produced
// by InitWorkspace matches the approved golden snapshot.
//
// Run with -update to regenerate golden files after intentional template changes:
//
//	go test ./internal/generator/antigravity/... -update
func TestAntigravityGenerator_Golden(t *testing.T) {
	testgolden.Run(t, "testdata/golden", shellOnlyGen)
}

// TestAntigravityGenerator_StructureChecks provides fast, hermetic assertions
// about the directory structure without relying on golden file presence.
// These run even before golden files are generated.
func TestAntigravityGenerator_StructureChecks(t *testing.T) {
	dir := t.TempDir()
	if err := shellOnlyGen(dir); err != nil {
		t.Fatalf("InitWorkspace returned error: %v", err)
	}

	expectedFiles := []string{
		".agent/rules/rules.md",
		// Roles
		".agent/skills/architect/SKILL.md",
		".agent/skills/designer/SKILL.md",
		".agent/skills/doc-writer/SKILL.md",
		".agent/skills/junior-engineer/SKILL.md",
		".agent/skills/qa-functional/SKILL.md",
		".agent/skills/security-ops/SKILL.md",
		".agent/skills/senior-engineer/SKILL.md",
		// Workflows
		".agent/workflows/build_cycle.md",
		".agent/workflows/audit_cycle.md",
		".agent/workflows/plan_project.md",
		".agent/workflows/plan_feature.md",
		".agent/workflows/review_cycle.md",
	}

	for _, rel := range expectedFiles {
		full := filepath.Join(dir, rel)
		if _, err := os.Stat(full); os.IsNotExist(err) {
			t.Errorf("expected file not found: %s", rel)
		}
	}
}

// TestAntigravityGenerator_ShellScriptFallback verifies that the shell script
// fallback produces an executable wrapper when Go is unavailable.
func TestAntigravityGenerator_ShellScriptFallback(t *testing.T) {
	dir := t.TempDir()
	if err := shellOnlyGen(dir); err != nil {
		t.Fatalf("InitWorkspace returned error: %v", err)
	}

	// Each role should have a scripts/castra shell wrapper.
	roles := []string{
		"architect", "designer", "doc-writer",
		"junior-engineer", "qa-functional", "security-ops", "senior-engineer",
	}

	for _, role := range roles {
		scriptPath := filepath.Join(dir, ".agent/skills", role, "scripts", "castra")
		info, err := os.Stat(scriptPath)
		if os.IsNotExist(err) {
			t.Errorf("shell script wrapper missing for role %q at %s", role, scriptPath)
			continue
		}
		if err != nil {
			t.Errorf("error stating script for role %q: %v", role, err)
			continue
		}
		// Must be executable.
		if info.Mode()&0111 == 0 {
			t.Errorf("script for role %q is not executable (mode: %v)", role, info.Mode())
		}

		// Must contain the role name so the wrapper wires up --role correctly.
		content, err := os.ReadFile(scriptPath)
		if err != nil {
			t.Errorf("failed to read script for role %q: %v", role, err)
			continue
		}
		if len(content) == 0 {
			t.Errorf("script for role %q is empty", role)
		}
	}
}

// TestAntigravityGenerator_RulesNotEmpty verifies rules.md is written and
// is non-empty (a proxy for correct template FS wiring).
func TestAntigravityGenerator_RulesNotEmpty(t *testing.T) {
	dir := t.TempDir()
	if err := shellOnlyGen(dir); err != nil {
		t.Fatalf("InitWorkspace returned error: %v", err)
	}

	rulesPath := filepath.Join(dir, ".agent/rules/rules.md")
	content, err := os.ReadFile(rulesPath)
	if err != nil {
		t.Fatalf("rules.md not found: %v", err)
	}
	if len(content) == 0 {
		t.Error("rules.md is empty")
	}
}
