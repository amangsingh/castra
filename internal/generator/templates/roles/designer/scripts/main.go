package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// The castra CLI requires the role flag. We inject it automatically.
	args := []string{"castra"}
	args = append(args, os.Args[1:]...)
	args = append(args, "--role", "designer")

	cmd := exec.Command("castra", args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.Sys().(syscall.WaitStatus).ExitStatus())
		}
		log.Fatalf("failed to execute castra: %v", err)
	}
}
