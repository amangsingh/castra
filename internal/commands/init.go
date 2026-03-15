package commands

import (
	"castra/internal/config"
	"castra/internal/generator"
	"fmt"
	"os"
)

type InitCommand struct{}

func (c *InitCommand) Name() string        { return "init" }
func (c *InitCommand) Description() string { return "Initialize workspace from castra.yaml" }
func (c *InitCommand) Usage() string       { return "castra init" }

func (c *InitCommand) Execute(ctx *Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get CWD: %v", err)
	}

	yamlPath := cwd + "/castra.yaml"

	// If castra.yaml does not exist, generate template and exit.
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		if err := config.GenerateTemplate(cwd); err != nil {
			return fmt.Errorf("failed to generate castra.yaml template: %v", err)
		}
		fmt.Println("castra.yaml not found. A template has been generated.")
		fmt.Println("Edit castra.yaml to declare your agent configuration, then re-run `castra init`.")
		return nil
	}

	// castra.yaml exists — parse and generate workspace.
	cfg, err := config.Parse(yamlPath)
	if err != nil {
		return fmt.Errorf("failed to parse castra.yaml: %v", err)
	}

	initialized := false
	for vendor, vCfg := range cfg.Vendors {
		if err := generator.InitWorkspaceFromConfig(cwd, vendor, vCfg); err != nil {
			return fmt.Errorf("failed to initialize %s workspace: %v", vendor, err)
		}
		initialized = true
	}

	if !initialized {
		fmt.Println("No supported vendors found in castra.yaml.")
	} else {
		fmt.Println("Castra workspace initialized.")
	}
	return nil
}

func (c *InitCommand) AuditInfo() (string, int64, string) {
	return "project", 0, "init"
}
