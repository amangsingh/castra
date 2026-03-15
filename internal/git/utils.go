package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// DiscoverCommonDir returns the absolute path to the common .git directory.
// In a standard repository, this is the .git directory.
// In a worktree, this is the parent .git directory where the common object store and refs are held.
func DiscoverCommonDir() (string, error) {
	// 1. Get the git common directory
	commonDir, err := runGit("rev-parse", "--git-common-dir")
	if err != nil {
		return "", fmt.Errorf("failed to get git common dir: %w", err)
	}

	// 2. Resolve to absolute path (handles if it's already absolute or relative to CWD)
	absPath, err := filepath.Abs(commonDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for %s: %w", commonDir, err)
	}

	return filepath.Clean(absPath), nil
}

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s failed: %v (stderr: %s)", strings.Join(args, " "), err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}
