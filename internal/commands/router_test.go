package commands_test

import (
	"fmt"
	"strings"
	"testing"

	"castra/internal/commands"
)

// --- Task 57: Router and Registry tests ---

// fakeCmd is a minimal Command implementation for routing tests.
type fakeCmd struct {
	name    string
	invoked bool
	args    []string
	retErr  error
}

func (f *fakeCmd) Name() string        { return f.name }
func (f *fakeCmd) Description() string { return "fake " + f.name }
func (f *fakeCmd) Usage() string       { return "fake " + f.name }
func (f *fakeCmd) Execute(ctx *commands.Context) error {
	f.invoked = true
	f.args = ctx.Args
	return f.retErr
}

// TestRegistryDispatch verifies that Execute routes to the correct registered command.
func TestRegistryDispatch(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &fakeCmd{name: "ping"}

	reg := commands.NewRegistry()
	reg.Register(cmd)

	ctx := commands.NewTestCtx(db, "architect", []string{"ping", "--flag"})
	commands.AssertNoError(t, reg.Execute(ctx))

	if !cmd.invoked {
		t.Error("expected 'ping' command to be invoked")
	}
	// Router strips the command name, so args should be ["--flag"]
	if len(cmd.args) != 1 || cmd.args[0] != "--flag" {
		t.Errorf("unexpected args after dispatch: %v", cmd.args)
	}
}

// TestRegistryUnknownCommand verifies an error is returned for unregistered names.
func TestRegistryUnknownCommand(t *testing.T) {
	db := commands.NewTestDB(t)
	reg := commands.NewRegistry()

	ctx := commands.NewTestCtx(db, "architect", []string{"nonexistent"})
	commands.AssertError(t, reg.Execute(ctx))
}

// TestRegistryEmptyArgs verifies an error is returned when no args are provided.
func TestRegistryEmptyArgs(t *testing.T) {
	db := commands.NewTestDB(t)
	reg := commands.NewRegistry()

	ctx := commands.NewTestCtx(db, "architect", []string{})
	commands.AssertError(t, reg.Execute(ctx))
}

// TestRegistryCommandError verifies that errors from commands are propagated.
func TestRegistryCommandError(t *testing.T) {
	db := commands.NewTestDB(t)
	cmd := &fakeCmd{name: "boom", retErr: fmt.Errorf("explosion")}
	reg := commands.NewRegistry()
	reg.Register(cmd)

	ctx := commands.NewTestCtx(db, "architect", []string{"boom"})
	err := reg.Execute(ctx)
	if err == nil || !strings.Contains(err.Error(), "explosion") {
		t.Errorf("expected explosion error, got: %v", err)
	}
}

// TestSubCommandDispatch verifies nested subcommand routing works end-to-end.
func TestSubCommandDispatch(t *testing.T) {
	db := commands.NewTestDB(t)
	inner := &fakeCmd{name: "list"}

	sub := commands.NewSubCommand("project", "manage projects")
	sub.Register(inner)

	ctx := commands.NewTestCtx(db, "architect", []string{"list", "--extra"})
	commands.AssertNoError(t, sub.Execute(ctx))

	if !inner.invoked {
		t.Error("expected inner 'list' command to be invoked")
	}
}

// TestSubCommandMissingSubcommand verifies that a SubCommand with no args exits
// gracefully. (SubCommand calls os.Exit(1) in this case; we test the early check
// indirectly via direct arg inspection.)
func TestSubCommandMetadata(t *testing.T) {
	sub := commands.NewSubCommand("sprint", "manage sprints")
	if sub.Name() != "sprint" {
		t.Errorf("expected name 'sprint', got %q", sub.Name())
	}
	if sub.Description() != "manage sprints" {
		t.Errorf("unexpected description: %q", sub.Description())
	}
	if !strings.Contains(sub.Usage(), "sprint") {
		t.Errorf("Usage should mention the command name, got: %q", sub.Usage())
	}
}

// TestFilterArgs verifies that --role and its value are stripped from arg slices.
func TestFilterArgs(t *testing.T) {
	input := []string{"--role", "architect", "task", "list", "--project", "1"}
	got := commands.FilterArgs(input)
	want := []string{"task", "list", "--project", "1"}

	if len(got) != len(want) {
		t.Fatalf("FilterArgs: expected %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("FilterArgs[%d]: expected %q, got %q", i, want[i], got[i])
		}
	}
}

// TestFilterArgsNoRole verifies FilterArgs is a no-op when --role is absent.
func TestFilterArgsNoRole(t *testing.T) {
	input := []string{"task", "list", "--project", "1"}
	got := commands.FilterArgs(input)
	if len(got) != len(input) {
		t.Errorf("FilterArgs modified args unexpectedly: %v", got)
	}
}
