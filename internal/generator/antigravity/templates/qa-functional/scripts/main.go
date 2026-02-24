package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	role := "qa-functional"
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
