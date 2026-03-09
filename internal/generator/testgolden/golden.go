// Package testgolden provides a reusable golden file testing framework for
// the castra generator suite.
//
// Usage:
//
//	func TestMyGenerator(t *testing.T) {
//	    testgolden.Run(t, "testdata", mypackage.InitWorkspace)
//	}
//
// Run with -update to regenerate golden files:
//
//	go test ./... -update
package testgolden

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// update is set by the -update flag. When true, Run will write generated
// output back to the testdata directory instead of comparing against it.
var update = flag.Bool("update", false, "regenerate golden files instead of comparing")

// GeneratorFunc is the signature of any castra workspace generator.
// It receives a base directory and writes files into it.
type GeneratorFunc func(baseDir string) error

// reporter is the minimal interface used internally so that we can swap in a
// fake implementation during self-testing without requiring the full
// *testing.T.
type reporter interface {
	Helper()
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	TempDir() string
}

// Run is the single entry point for all generator golden file tests.
//
// It:
//  1. Creates a temporary directory and runs gen against it.
//  2. Walks every file in the generated output.
//  3. In update mode: copies each generated file to testdataDir, creating it
//     if absent. Any golden files that no longer exist in the output are left
//     in place so the developer notices the deletion.
//  4. In comparison mode: reads the corresponding golden file from testdataDir
//     and compares it byte-for-byte with the generated content, reporting a
//     clear diff on mismatch.
func Run(t *testing.T, testdataDir string, gen GeneratorFunc) {
	t.Helper()
	run(t, testdataDir, gen)
}

// CaptureResult holds the outcome of a Capture call for assertion in tests
// that verify the framework's own failure detection.
type CaptureResult struct {
	// Failed is true if the framework would have failed the test.
	Failed bool
	// FatalMsg is the message passed to Fatalf (if any).
	FatalMsg string
	// ErrorMsg is the first message passed to Errorf (if any).
	ErrorMsg string
}

// Capture runs the framework using gen against testdataDir but captures any
// failures instead of propagating them to t. This is intended only for testing
// the framework itself. It never enters -update mode.
func Capture(t *testing.T, testdataDir string, gen GeneratorFunc) CaptureResult {
	t.Helper()

	// Make sure flags are parsed.
	if !flag.Parsed() {
		flag.Parse()
	}

	fake := &fakeReporter{tempDir: t.TempDir()}
	absTestdata, err := filepath.Abs(testdataDir)
	if err != nil {
		t.Fatalf("cannot resolve testdata directory %q: %v", testdataDir, err)
	}

	tmpDir := fake.TempDir()
	if err := gen(tmpDir); err != nil {
		fake.Fatalf("generator returned error: %v", err)
	}

	if !fake.failed {
		compare(fake, tmpDir, absTestdata)
	}

	return CaptureResult{
		Failed:   fake.failed,
		FatalMsg: fake.fatalMsg,
		ErrorMsg: fake.errorMsg,
	}
}

// run executes the full golden file workflow against the provided reporter.
func run(r reporter, testdataDir string, gen GeneratorFunc) {
	r.Helper()

	// Make sure flags are parsed (needed when called outside of TestMain).
	if !flag.Parsed() {
		flag.Parse()
	}

	// Run the generator into a temp directory.
	tmpDir := r.TempDir()
	if err := gen(tmpDir); err != nil {
		r.Fatalf("generator returned error: %v", err)
		return
	}

	// Resolve the testdata directory to an absolute path relative to the
	// calling test's location. If testdataDir is already absolute, this is
	// a no-op.
	absTestdata, err := filepath.Abs(testdataDir)
	if err != nil {
		r.Fatalf("cannot resolve testdata directory %q: %v", testdataDir, err)
		return
	}

	if *update {
		regenerate(r, tmpDir, absTestdata)
	} else {
		compare(r, tmpDir, absTestdata)
	}
}

// regenerate copies every file produced by the generator into testdataDir,
// mirroring the directory structure. It creates testdataDir if it does not
// already exist.
func regenerate(r reporter, generatedDir, testdataDir string) {
	r.Helper()

	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		r.Fatalf("failed to create testdata directory %q: %v", testdataDir, err)
		return
	}

	err := filepath.WalkDir(generatedDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(generatedDir, path)
		if err != nil {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading generated file %q: %w", relPath, err)
		}

		goldenPath := filepath.Join(testdataDir, relPath)
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0755); err != nil {
			return fmt.Errorf("creating parent dir for %q: %w", goldenPath, err)
		}

		if err := os.WriteFile(goldenPath, content, 0644); err != nil {
			return fmt.Errorf("writing golden file %q: %w", goldenPath, err)
		}

		r.Logf("updated golden: %s", relPath)
		return nil
	})

	if err != nil {
		r.Fatalf("error walking generated output: %v", err)
	}
}

// compare walks every file in generatedDir and checks it against its golden
// counterpart in testdataDir. It also checks for golden files that the
// generator no longer produces (orphaned goldens).
func compare(r reporter, generatedDir, testdataDir string) {
	r.Helper()

	// Collect all generated files and compare each one.
	var generatedFiles []string

	err := filepath.WalkDir(generatedDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(generatedDir, path)
		if err != nil {
			return err
		}
		generatedFiles = append(generatedFiles, relPath)

		got, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading generated file %q: %w", relPath, err)
		}

		goldenPath := filepath.Join(testdataDir, relPath)
		want, err := os.ReadFile(goldenPath)
		if err != nil {
			if os.IsNotExist(err) {
				r.Errorf("missing golden file for %q\n  run with -update to create it\n  golden path: %s", relPath, goldenPath)
				return nil
			}
			return fmt.Errorf("reading golden file %q: %w", goldenPath, err)
		}

		if !bytes.Equal(want, got) {
			r.Errorf("mismatch for %q:\n%s", relPath, diffOutput(relPath, want, got))
		}

		return nil
	})

	if err != nil {
		r.Fatalf("error walking generated output: %v", err)
		return
	}

	// Check for orphaned golden files (present in testdata but not generated).
	if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
		// No testdata at all — all missing goldens were already reported above.
		return
	}

	goldenSet := make(map[string]struct{}, len(generatedFiles))
	for _, f := range generatedFiles {
		goldenSet[f] = struct{}{}
	}

	_ = filepath.WalkDir(testdataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		relPath, err := filepath.Rel(testdataDir, path)
		if err != nil {
			return err
		}
		if _, ok := goldenSet[relPath]; !ok {
			r.Errorf("orphaned golden file (generator no longer produces this file): %s", relPath)
		}
		return nil
	})
}

// diffOutput produces a human-readable diff between want and got.
// It performs a line-level diff, showing context around changes.
func diffOutput(name string, want, got []byte) string {
	wantLines := strings.Split(string(want), "\n")
	gotLines := strings.Split(string(got), "\n")

	var buf strings.Builder
	fmt.Fprintf(&buf, "--- want/%s\n+++ got/%s\n", name, name)

	maxLines := len(wantLines)
	if len(gotLines) > maxLines {
		maxLines = len(gotLines)
	}

	const context = 3
	type hunk struct{ start, end int }
	var hunks []hunk

	// Find changed line ranges.
	for i := 0; i < maxLines; i++ {
		wl := lineAt(wantLines, i)
		gl := lineAt(gotLines, i)
		if wl != gl {
			start := max(0, i-context)
			end := min(maxLines, i+context+1)
			if len(hunks) > 0 && hunks[len(hunks)-1].end >= start {
				hunks[len(hunks)-1].end = end
			} else {
				hunks = append(hunks, hunk{start, end})
			}
		}
	}

	if len(hunks) == 0 {
		// Should not happen if bytes differ, but handle gracefully.
		buf.WriteString("  (binary difference — content differs but lines appear equal)\n")
		return buf.String()
	}

	for _, h := range hunks {
		fmt.Fprintf(&buf, "@@ -%d +%d @@\n", h.start+1, h.start+1)
		for i := h.start; i < h.end; i++ {
			wl := lineAt(wantLines, i)
			gl := lineAt(gotLines, i)
			if wl == gl {
				fmt.Fprintf(&buf, " %s\n", wl)
			} else {
				if i < len(wantLines) {
					fmt.Fprintf(&buf, "-%s\n", wl)
				}
				if i < len(gotLines) {
					fmt.Fprintf(&buf, "+%s\n", gl)
				}
			}
		}
	}

	return buf.String()
}

// fakeReporter is a test-only reporter that records failures instead of
// propagating them, allowing the framework to test its own error detection.
type fakeReporter struct {
	failed   bool
	fatalMsg string
	errorMsg string
	tempDir  string
}

func (f *fakeReporter) Helper() {}

func (f *fakeReporter) Logf(_ string, _ ...interface{}) {}

func (f *fakeReporter) Errorf(format string, args ...interface{}) {
	f.failed = true
	if f.errorMsg == "" {
		f.errorMsg = fmt.Sprintf(format, args...)
	}
}

func (f *fakeReporter) Fatalf(format string, args ...interface{}) {
	f.failed = true
	f.fatalMsg = fmt.Sprintf(format, args...)
}

func (f *fakeReporter) TempDir() string {
	// Create a unique temp dir each call so multiple calls work correctly.
	dir, err := os.MkdirTemp(f.tempDir, "testgolden-*")
	if err != nil {
		panic(fmt.Sprintf("fakeReporter.TempDir: %v", err))
	}
	return dir
}

func lineAt(lines []string, i int) string {
	if i < len(lines) {
		return lines[i]
	}
	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
