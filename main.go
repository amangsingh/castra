package main

import (
	"castra/internal/commands"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Initialize the command registry
	registry := commands.NewDefaultRegistry()

	// Global flag parsing helpers are tricky with subcommands.
	// We mandate: castra <cmd> --role <role> [subcmd] ...
	// OR castra <cmd> [subcmd] ... --role <role>
	role := getRoleFromArgs()
	if role == "" && os.Args[1] != "init" {
		fmt.Println("Error: --role <architect|senior-engineer|junior-engineer|designer|qa-functional|security-ops|doc-writer> is required")
		os.Exit(1)
	}

	// Initialize context with filtered args
	db := commands.GetDB()
	defer db.Close()

	ctx := &commands.Context{
		Role: role,
		DB:   db,
		Args: commands.FilterArgs(os.Args[1:]), // Now using FilterArgs correctly
	}

	if err := registry.Execute(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: castra <command> --role <role> [subcommand] [flags]")
	fmt.Println("\nRoles: architect, senior-engineer, junior-engineer, designer, qa-functional, security-ops, doc-writer")
	fmt.Println("\nCommands:")
	registry := commands.NewDefaultRegistry()
	registry.PrintUsage()
}

func getRoleFromArgs() string {
	for i, arg := range os.Args {
		if arg == "--role" {
			if i+1 < len(os.Args) {
				return os.Args[i+1]
			}
		}
	}
	return ""
}
