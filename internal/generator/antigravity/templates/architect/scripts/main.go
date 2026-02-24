package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// The Architect always uses the --role architect flag
	role := "architect"

	// Construct the command: castra <args> --role architect
	// We append --role at the end, as our main.go logic handles it globally if present.
	// However, for clarity and strictness, we'll just append it.
	args := append(os.Args[1:], "--role", role)

	cmd := exec.Command("castra", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing castra as %s: %v\n", role, err)
		os.Exit(1)
	}
}
