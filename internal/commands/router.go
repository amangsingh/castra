package commands

import (
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
	return cmd.Execute(ctx)
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
