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

	// Global flag parsing helpers are tricky with subcommands.
	// We mandate: castra <cmd> --role <role> [subcmd] ...
	// OR castra <cmd> [subcmd] ... --role <role>
	// We use a simple helper to extract the role before dispatching.
	role := getRoleFromArgs()
	if role == "" && os.Args[1] != "init" {
		fmt.Println("Error: --role <architect|senior-engineer|junior-engineer|qa-functional|security-ops|doc-writer> is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		commands.HandleInit()
	case "project":
		commands.HandleProject(role)
	case "milestone":
		commands.HandleMilestone(role)
	case "sprint":
		commands.HandleSprint(role)
	case "task":
		commands.HandleTask(role)
	case "note":
		commands.HandleNote(role)
	case "log":
		commands.HandleLog(role)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: castra <command> --role <role> [subcommand] [flags]")
	fmt.Println("\nRoles: architect, senior-engineer, junior-engineer, qa-functional, security-ops, doc-writer")
	fmt.Println("\nCommands:")
	fmt.Println("  init     Initialize workspace")
	fmt.Println("  project  Manage projects")
	fmt.Println("  sprint   Manage sprints")
	fmt.Println("  task     Manage tasks")
	fmt.Println("  note     Manage project notes")
	fmt.Println("  log      View and add audit log entries")
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
