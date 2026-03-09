package commands

import (
	antigravitygen "castra/internal/generator/antigravity"
	copilotgen "castra/internal/generator/copilot"
	geminigen "castra/internal/generator/gemini"
	"flag"
	"fmt"
	"os"
)

type InitCommand struct{}

func (c *InitCommand) Name() string        { return "init" }
func (c *InitCommand) Description() string { return "Initialize workspace" }
func (c *InitCommand) Usage() string       { return "castra init --antigravity|--copilot|--gemini" }

func (c *InitCommand) Execute(ctx *Context) error {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	useAntigravity := fs.Bool("antigravity", false, "Initialize for Antigravity platform")
	useCopilot := fs.Bool("copilot", false, "Initialize GitHub Copilot agent templates")
	useGemini := fs.Bool("gemini", false, "Initialize Gemini Code Assist agent templates")

	// FilterArgs is a helper to skip --role and its value
	argsToParse := ctx.Args
	fs.Parse(argsToParse)

	if !*useAntigravity && !*useCopilot && !*useGemini {
		return fmt.Errorf("initialization requires a target platform.\nUsage: %s", c.Usage())
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get CWD: %v", err)
	}

	if *useAntigravity {
		if err := antigravitygen.InitWorkspace(cwd); err != nil {
			return fmt.Errorf("failed to init Antigravity workspace: %v", err)
		}
		fmt.Println("Castra initialized for Antigravity.")
	}

	if *useCopilot {
		if err := copilotgen.InitWorkspace(cwd); err != nil {
			return fmt.Errorf("failed to init Copilot workspace: %v", err)
		}
		fmt.Println("Castra initialized for GitHub Copilot.")
	}

	if *useGemini {
		if err := geminigen.InitWorkspace(cwd); err != nil {
			return fmt.Errorf("failed to init Gemini workspace: %v", err)
		}
		fmt.Println("Castra initialized for Gemini Code Assist.")
	}

	return nil
}
