package common_test

import (
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"castra/internal/generator/common"
)

// Define -update flag so 'go test ./... -update' doesn't fail for this package.
var _ = flag.Bool("update", false, "unused")

// ---------------------------------------------------------------------------
// EnsureDirs
// ---------------------------------------------------------------------------

// TestEnsureDirs_CreatesNestedDirectories verifies that EnsureDirs creates
// all specified directories including intermediate parents.
func TestEnsureDirs_CreatesNestedDirectories(t *testing.T) {
	base := t.TempDir()

	dirs := []string{
		"a/b/c",
		"x/y",
		"z",
	}

	if err := common.EnsureDirs(base, dirs); err != nil {
		t.Fatalf("EnsureDirs returned error: %v", err)
	}

	for _, d := range dirs {
		full := filepath.Join(base, d)
		info, err := os.Stat(full)
		if os.IsNotExist(err) {
			t.Errorf("directory not created: %s", d)
			continue
		}
		if err != nil {
			t.Errorf("error stating %s: %v", d, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected directory but got file at: %s", d)
		}
	}
}

// TestEnsureDirs_IdempotentOnExistingDirs verifies that calling EnsureDirs
// twice for the same paths does not return an error.
func TestEnsureDirs_IdempotentOnExistingDirs(t *testing.T) {
	base := t.TempDir()
	dirs := []string{"a/b", "c"}

	if err := common.EnsureDirs(base, dirs); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if err := common.EnsureDirs(base, dirs); err != nil {
		t.Errorf("second call (idempotency) returned error: %v", err)
	}
}

// TestEnsureDirs_EmptyList verifies that an empty list is a no-op.
func TestEnsureDirs_EmptyList(t *testing.T) {
	base := t.TempDir()
	if err := common.EnsureDirs(base, nil); err != nil {
		t.Errorf("EnsureDirs with nil dirs returned error: %v", err)
	}
	if err := common.EnsureDirs(base, []string{}); err != nil {
		t.Errorf("EnsureDirs with empty slice returned error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// WriteTemplateFile
// ---------------------------------------------------------------------------

// makeFS returns a minimal in-memory FS for testing.
func makeFS(files map[string]string) fs.FS {
	m := make(fstest.MapFS)
	for name, content := range files {
		m[name] = &fstest.MapFile{Data: []byte(content)}
	}
	return m
}

// TestWriteTemplateFile_WritesContent verifies that the file is created with
// the exact content from the embedded FS.
func TestWriteTemplateFile_WritesContent(t *testing.T) {
	tfs := makeFS(map[string]string{
		"rules.md": "# Rules\nsome content\n",
	})

	dest := filepath.Join(t.TempDir(), "output/rules.md")
	if err := common.WriteTemplateFile(tfs, "rules.md", dest); err != nil {
		t.Fatalf("WriteTemplateFile returned error: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	want := "# Rules\nsome content\n"
	if string(got) != want {
		t.Errorf("content mismatch:\nwant: %q\n got: %q", want, string(got))
	}
}

// TestWriteTemplateFile_CreatesParentDirs verifies that intermediate directories
// are created automatically.
func TestWriteTemplateFile_CreatesParentDirs(t *testing.T) {
	tfs := makeFS(map[string]string{"file.txt": "data"})

	dest := filepath.Join(t.TempDir(), "deeply/nested/dir/file.txt")
	if err := common.WriteTemplateFile(tfs, "file.txt", dest); err != nil {
		t.Fatalf("WriteTemplateFile returned error: %v", err)
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Error("destination file was not created")
	}
}

// TestWriteTemplateFile_MissingSourceReturnsError verifies that a missing
// source file in the template FS returns a non-nil error.
func TestWriteTemplateFile_MissingSourceReturnsError(t *testing.T) {
	tfs := makeFS(map[string]string{})
	dest := filepath.Join(t.TempDir(), "out.md")

	if err := common.WriteTemplateFile(tfs, "nonexistent.md", dest); err == nil {
		t.Error("expected error for missing source file, got nil")
	}
}

// ---------------------------------------------------------------------------
// WalkAndMap
// ---------------------------------------------------------------------------

// TestWalkAndMap_MapsFilesToDestination verifies that files are correctly
// mapped and written to the destination base directory.
func TestWalkAndMap_MapsFilesToDestination(t *testing.T) {
	tfs := makeFS(map[string]string{
		"src/a.txt":     "a content",
		"src/sub/b.txt": "b content",
	})

	destBase := t.TempDir()

	err := common.WalkAndMap(tfs, "src", destBase, func(path string, d fs.DirEntry) (string, bool, error) {
		relPath, _ := filepath.Rel("src", path)
		if relPath == "." {
			return "", false, nil // don't skip root
		}
		return relPath, false, nil
	})
	if err != nil {
		t.Fatalf("WalkAndMap returned error: %v", err)
	}

	wantFiles := map[string]string{
		"a.txt":     "a content",
		"sub/b.txt": "b content",
	}

	for rel, wantContent := range wantFiles {
		full := filepath.Join(destBase, rel)
		got, err := os.ReadFile(full)
		if os.IsNotExist(err) {
			t.Errorf("expected output file not found: %s", rel)
			continue
		}
		if err != nil {
			t.Errorf("error reading %s: %v", rel, err)
			continue
		}
		if string(got) != wantContent {
			t.Errorf("content mismatch for %s:\nwant: %q\n got: %q", rel, wantContent, string(got))
		}
	}
}

// TestWalkAndMap_SkipLogicHonoured verifies that entries with skip=true are
// not written to the destination.
func TestWalkAndMap_SkipLogicHonoured(t *testing.T) {
	tfs := makeFS(map[string]string{
		"src/keep.txt": "keep",
		"src/skip.txt": "skip",
	})

	destBase := t.TempDir()

	err := common.WalkAndMap(tfs, "src", destBase, func(path string, d fs.DirEntry) (string, bool, error) {
		if d.IsDir() {
			return "", false, nil // don't skip root dir
		}
		if d.Name() == "skip.txt" {
			return "", true, nil // skip this file
		}
		return d.Name(), false, nil
	})
	if err != nil {
		t.Fatalf("WalkAndMap returned error: %v", err)
	}

	// keep.txt must exist.
	if _, err := os.Stat(filepath.Join(destBase, "keep.txt")); os.IsNotExist(err) {
		t.Error("keep.txt was not written")
	}
	// skip.txt must NOT exist.
	if _, err := os.Stat(filepath.Join(destBase, "skip.txt")); !os.IsNotExist(err) {
		t.Error("skip.txt was written but should have been skipped")
	}
}

// TestWalkAndMap_CreatesSubDirectories verifies that directory entries trigger
// subdirectory creation in the destination.
func TestWalkAndMap_CreatesSubDirectories(t *testing.T) {
	tfs := makeFS(map[string]string{
		"src/sub/file.txt": "content",
	})

	destBase := t.TempDir()

	err := common.WalkAndMap(tfs, "src", destBase, func(path string, d fs.DirEntry) (string, bool, error) {
		relPath, _ := filepath.Rel("src", path)
		if relPath == "." {
			return "", false, nil
		}
		return relPath, false, nil
	})
	if err != nil {
		t.Fatalf("WalkAndMap returned error: %v", err)
	}

	subDir := filepath.Join(destBase, "sub")
	if info, err := os.Stat(subDir); os.IsNotExist(err) {
		t.Error("subdirectory 'sub' was not created")
	} else if err == nil && !info.IsDir() {
		t.Error("'sub' exists but is not a directory")
	}
}

// TestWalkAndMap_MapperErrorAborts verifies that a non-nil mapper error
// causes WalkAndMap to return that error immediately.
func TestWalkAndMap_MapperErrorAborts(t *testing.T) {
	tfs := makeFS(map[string]string{
		"src/a.txt": "a",
	})

	destBase := t.TempDir()
	sentinelErr := os.ErrInvalid

	err := common.WalkAndMap(tfs, "src", destBase, func(path string, d fs.DirEntry) (string, bool, error) {
		if !d.IsDir() {
			return "", false, sentinelErr
		}
		return "", false, nil // don't skip dir
	})

	if err == nil {
		t.Error("expected error from mapper, got nil")
	}
}
