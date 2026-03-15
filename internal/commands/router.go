package commands

import (
	"castra/internal/cli"
	"castra/internal/persona"
	"database/sql"
	"fmt"
	"os"
)

// Context provides the execution context for a command.
type Context struct {
	Role string
	DB   *sql.DB
	Args []string
}

// Command is the interface that all commands must implement.
type Command interface {
	Name() string
	Description() string
	Execute(ctx *Context) error
	Usage() string
}

// MutatingCommand is an optional interface for commands that modify database
// state. When a command implements this interface the router automatically
// records an audit entry after a successful Execute.
//
// AuditInfo returns:
//   - entityType: e.g. "task", "project", "sprint" (empty string = generic)
//   - entityID:   the ID of the affected entity (0 if not applicable, e.g. project.add before ID known)
//   - action:     human-readable description, e.g. "task.add", "task.update"
type MutatingCommand interface {
	AuditInfo() (entityType string, entityID int64, action string)
}

// ReadCommand is an optional interface for commands that perform read-only
// queries. The router will automatically log an audit entry for reads.
type ReadCommand interface {
	ReadInfo() (entityType string, action string)
}

// PersonaCompliant is an optional interface that commands implement to declare
// which roles are permitted to execute them. The router enforces this restriction
// as a first-pass Persona Linter before the command's Execute is called.
//
// If the invoking role is not in AllowedRoles(), the router:
//  1. Logs a persona_non_compliance audit entry (best-effort).
//  2. Returns a rejection error — command execution is aborted immediately.
type PersonaCompliant interface {
	AllowedRoles() []string
}

// Registry manages a set of commands and subcommands.
type Registry struct {
	commands map[string]Command
}

// NewRegistry creates a new command registry.
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

// Register adds a command to the registry.
func (r *Registry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

// Execute dispatches the command to the appropriate handler.
// After a successful Execute, if the command implements MutatingCommand or ReadCommand
// the router automatically records a best-effort audit entry.
func (r *Registry) Execute(ctx *Context) error {
	if len(ctx.Args) == 0 {
		return fmt.Errorf("no command provided")
	}

	name := ctx.Args[0]
	cmd, ok := r.commands[name]
	if !ok {
		return fmt.Errorf("unknown command: %s", name)
	}

	// Remove the command name from args for the handler
	ctx.Args = ctx.Args[1:]

	// --- Persona Linter (Gate 3: The Self) ---
	// Validate role compliance before any execution occurs.
	if pc, ok := cmd.(PersonaCompliant); ok {
		linter := persona.NewLinter()
		if err := linter.Lint(ctx.Role, name, pc.AllowedRoles()); err != nil {
			auditMsg := fmt.Sprintf("[persona_non_compliance] role '%s' attempted '%s'; allowed roles: %v", ctx.Role, name, pc.AllowedRoles())
			_ = cli.AddAuditEntry(ctx.DB, "system", 0, "persona_non_compliance", ctx.Role, auditMsg)
			return err
		}
	}

	if err := cmd.Execute(ctx); err != nil {
		// Centralized error logging: record non-fatal command errors.
		cli.LogError(ctx.DB, ctx.Role, name, err.Error(), cli.SeverityError)
		return err
	}

	// Auto-audit: record a write operation for any MutatingCommand.
	if mc, ok := cmd.(MutatingCommand); ok {
		entityType, entityID, action := mc.AuditInfo()
		// Best-effort — never fail the operation if audit logging fails.
		_ = cli.AddAuditEntry(ctx.DB, entityType, entityID, action, ctx.Role, "[auto]")
	}

	// Auto-audit: record a read operation for any ReadCommand.
	if rc, ok := cmd.(ReadCommand); ok {
		entityType, action := rc.ReadInfo()
		// Entity ID is 0 for generic reads
		_ = cli.AddAuditEntry(ctx.DB, entityType, 0, action, ctx.Role, "[read]")
	}

	return nil
}

// PrintUsage prints the usage of all registered commands.
func (r *Registry) PrintUsage() {
	for _, cmd := range r.commands {
		fmt.Printf("  %-10s %s\n", cmd.Name(), cmd.Description())
	}
}

// SubCommand wraps multiple commands as subcommands of a parent command.
type SubCommand struct {
	name        string
	description string
	registry    *Registry
}

func NewSubCommand(name, description string) *SubCommand {
	return &SubCommand{
		name:        name,
		description: description,
		registry:    NewRegistry(),
	}
}

func (s *SubCommand) Register(cmd Command) {
	s.registry.Register(cmd)
}

func (s *SubCommand) Name() string        { return s.name }
func (s *SubCommand) Description() string { return s.description }
func (s *SubCommand) Usage() string       { return fmt.Sprintf("castra %s [subcommand]", s.name) }

func (s *SubCommand) Execute(ctx *Context) error {
	if len(ctx.Args) == 0 {
		fmt.Printf("Usage: %s\n\nSubcommands:\n", s.Usage())
		s.registry.PrintUsage()
		os.Exit(1)
	}
	return s.registry.Execute(ctx)
}
