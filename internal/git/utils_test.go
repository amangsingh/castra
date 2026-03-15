package git

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestDiscoverCommonDir(t *testing.T) {
	// Since we are running in the castra repo, we should be able to find the common dir.
	commonDir, err := DiscoverCommonDir()
	if err != nil {
		t.Fatalf("DiscoverCommonDir failed: %v", err)
	}

	if !filepath.IsAbs(commonDir) {
		t.Errorf("Expected absolute path, got: %s", commonDir)
	}

	if !strings.HasSuffix(commonDir, ".git") {
		t.Errorf("Expected path to end with .git, got: %s", commonDir)
	}
}

// We'll skip complex worktree tests in CI/automated environments for now, 
// but local manual verification was done in the research phase.
