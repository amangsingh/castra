// Package commands — testing.go
// Shared test harness for the commands package.
// This file is compiled only during `go test` runs and must not be imported
// by production code. All exported symbols are prefixed with "Test" per
// Go convention for test helpers in non-_test files.
//
// Provides three building blocks for command tests:
//  1. NewTestDB  — isolated in-memory SQLite database, fully migrated.
//  2. NewTestCtx — assembles a *Context with pre-configured role and args.
//  3. CaptureOutput — redirects os.Stdout so callers can assert on printed output.

package commands

import (
	"bytes"
	"castra/internal/db"
	"database/sql"
	"io"
	"os"
	"testing"
)

// NewTestDB opens an isolated, fully-migrated in-memory SQLite database.
// The database is automatically closed when the test finishes.
//
// Each call returns a brand-new database — no state shared between tests.
func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("NewTestDB: failed to initialize in-memory database: %v", err)
	}

	t.Cleanup(func() { database.Close() })
	return database
}

// NewTestCtx builds a *Context suitable for feeding directly into a command's
// Execute method.
//
//	role  — the Castra role string (e.g. "architect", "senior-engineer")
//	args  — command-level arguments AFTER the subcommand name has been stripped,
//	        e.g. []string{"--name", "Foo"} for `castra project add --name Foo`.
func NewTestCtx(db *sql.DB, role string, args []string) *Context {
	return &Context{
		Role: role,
		DB:   db,
		Args: args,
	}
}

// CaptureOutput redirects os.Stdout for the duration of fn and returns
// everything that was written to it as a string.
//
// Usage:
//
//	out := CaptureOutput(t, func() {
//	    err := cmd.Execute(ctx)
//	    if err != nil { t.Fatal(err) }
//	})
//	if !strings.Contains(out, "expected text") { ... }
func CaptureOutput(t *testing.T, fn func()) string {
	t.Helper()

	// Swap os.Stdout for a pipe.
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("CaptureOutput: os.Pipe() failed: %v", err)
	}
	os.Stdout = w

	// Run the function with the swapped stdout.
	fn()

	// Restore stdout, then drain the pipe.
	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("CaptureOutput: failed to read captured output: %v", err)
	}
	r.Close()

	return buf.String()
}

// AssertNoError is a convenience helper that calls t.Fatalf when err != nil.
// Saves boilerplate in tests that expect a command to succeed.
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// AssertError is a convenience helper that calls t.Fatalf when err is nil.
// Use this to assert that a command correctly rejects invalid input.
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

// AssertOutputContains fails the test if the captured output does not contain
// the expected substring.
func AssertOutputContains(t *testing.T, output, want string) {
	t.Helper()
	if !bytes.Contains([]byte(output), []byte(want)) {
		t.Errorf("output does not contain %q\ngot:\n%s", want, output)
	}
}
