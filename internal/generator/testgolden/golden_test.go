package testgolden_test

import (
	"os"
	"path/filepath"
	"testing"

	"castra/internal/generator/testgolden"
)

// simpleGen writes two known files into baseDir. These files are the fixtures
// for testdata/simple golden files.
func simpleGen(baseDir string) error {
	files := map[string]string{
		"hello.txt":     "hello, world\n",
		"sub/nested.md": "# nested\nsome content\n",
	}
	for rel, content := range files {
		full := filepath.Join(baseDir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

// TestGoldenRun_HappyPath verifies that Run passes when generated output
// exactly matches the golden files.
//
// Run with -update to regenerate all golden files:
//
//	go test ./internal/generator/testgolden/... -update
func TestGoldenRun_HappyPath(t *testing.T) {
	testgolden.Run(t, "testdata/simple", simpleGen)
}

// TestGoldenRun_GeneratorError verifies that Run immediately fails the test
// when the generator function itself returns an error.
func TestGoldenRun_GeneratorError(t *testing.T) {
	errGen := func(baseDir string) error {
		return os.ErrPermission // any non-nil error
	}

	result := testgolden.Capture(t, "testdata/simple", errGen)
	if !result.Failed {
		t.Error("expected Run to fail when generator returns an error")
	}
	if result.FatalMsg == "" {
		t.Error("expected a fatal message when generator returns an error")
	}
}

// TestGoldenRun_ContentMismatch verifies that Run calls t.Errorf with a
// diff when generated content differs from the golden file.
func TestGoldenRun_ContentMismatch(t *testing.T) {
	mismatchGen := func(baseDir string) error {
		// Write hello.txt with wrong content.
		if err := os.WriteFile(filepath.Join(baseDir, "hello.txt"), []byte("wrong content\n"), 0644); err != nil {
			return err
		}
		// Write the other expected file correctly.
		sub := filepath.Join(baseDir, "sub")
		if err := os.MkdirAll(sub, 0755); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(sub, "nested.md"), []byte("# nested\nsome content\n"), 0644)
	}

	result := testgolden.Capture(t, "testdata/simple", mismatchGen)
	if !result.Failed {
		t.Error("expected Run to fail for content mismatch in hello.txt")
	}
	if result.ErrorMsg == "" {
		t.Error("expected an error message describing the mismatch")
	}
}

// TestGoldenRun_MissingGolden verifies that Run calls t.Errorf when the
// generator produces a file that has no corresponding golden.
func TestGoldenRun_MissingGolden(t *testing.T) {
	extraGen := func(baseDir string) error {
		if err := simpleGen(baseDir); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(baseDir, "unexpected.txt"), []byte("new\n"), 0644)
	}

	result := testgolden.Capture(t, "testdata/simple", extraGen)
	if !result.Failed {
		t.Error("expected Run to fail for file with no golden counterpart")
	}
}

// TestGoldenRun_OrphanDetection verifies that Run calls t.Errorf when a
// golden file exists but the generator no longer produces it.
func TestGoldenRun_OrphanDetection(t *testing.T) {
	partialGen := func(baseDir string) error {
		// Produces only hello.txt; testdata/simple also has sub/nested.md.
		return os.WriteFile(filepath.Join(baseDir, "hello.txt"), []byte("hello, world\n"), 0644)
	}

	result := testgolden.Capture(t, "testdata/simple", partialGen)
	if !result.Failed {
		t.Error("expected Run to fail for orphaned golden sub/nested.md")
	}
}
