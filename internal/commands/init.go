package commands

import (
	antigravitygen "castra/internal/generator/antigravity"
	copilotgen "castra/internal/generator/copilot"
	geminigen "castra/internal/generator/gemini"
	"flag"
	"fmt"
	"log"
	"os"
)

func HandleInit() {
	// Parse flags for init
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	useAntigravity := fs.Bool("antigravity", false, "Initialize for Antigravity platform")
	useCopilot := fs.Bool("copilot", false, "Initialize GitHub Copilot agent templates")
	useGemini := fs.Bool("gemini", false, "Initialize Gemini Code Assist agent templates")

	// Simple approach: Use FilterArgs on everything after os.Args[1]
	argsToParse := FilterArgs(os.Args[2:])
	fs.Parse(argsToParse)

	if !*useAntigravity && !*useCopilot && !*useGemini {
		fmt.Println("Error: initialization requires a target platform.")
		fmt.Println("Usage: castra init --antigravity")
		fmt.Println("       castra init --copilot")
		fmt.Println("       castra init --gemini")
		fmt.Println("       castra init --antigravity --copilot --gemini")
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get CWD: %v", err)
	}

	database := GetDB()
	database.Close()

	if *useAntigravity {
		if err := antigravitygen.InitWorkspace(cwd); err != nil {
			log.Fatalf("Failed to init Antigravity workspace: %v", err)
		}
		fmt.Println("Castra initialized for Antigravity.")
	}

	if *useCopilot {
		if err := copilotgen.InitWorkspace(cwd); err != nil {
			log.Fatalf("Failed to init Copilot workspace: %v", err)
		}
		fmt.Println("Castra initialized for GitHub Copilot.")
	}

	if *useGemini {
		if err := geminigen.InitWorkspace(cwd); err != nil {
			log.Fatalf("Failed to init Gemini workspace: %v", err)
		}
		fmt.Println("Castra initialized for Gemini Code Assist.")
	}
}
