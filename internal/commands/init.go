package commands

import (
	"castra/internal/generator"
	"flag"
	"fmt"
	"log"
	"os"
)

func HandleInit() {
	// Parse flags for init
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	antigravity := fs.Bool("antigravity", false, "Initialize for Antigravity platform")

	// Simple approach: Use FilterArgs on everything after os.Args[1]
	argsToParse := FilterArgs(os.Args[2:])
	fs.Parse(argsToParse)

	if !*antigravity {
		fmt.Println("Error: initialization requires a target platform.")
		fmt.Println("Usage: castra init --antigravity")
		fmt.Println("(Support for other platforms coming soon)")
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get CWD: %v", err)
	}

	database := GetDB()
	database.Close()

	if err := generator.InitWorkspace(cwd); err != nil {
		log.Fatalf("Failed to init workspace: %v", err)
	}
	fmt.Println("Castra initialized for Antigravity.")
}
