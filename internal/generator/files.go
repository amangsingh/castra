package generator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed templates/*
var templatesFS embed.FS

func InitWorkspace(baseDir string) error {
	// 1. Create directory structure
	dirs := []string{
		".agent/rules",
		".agent/skills",
		".agent/skills",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(baseDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 2. Write Constitutional Rule
	if err := writeTemplate("templates/rules.md", filepath.Join(baseDir, ".agent/rules/rules.md")); err != nil {
		return err
	}
	if err := writeTemplate("templates/castra.md", filepath.Join(baseDir, ".agent/rules/castra.md")); err != nil {
		return err
	}

	// 3. Write Skills
	skills := []string{
		"architect.md",
		"senior-engineer.md",
		"junior-engineer.md",
		"qa-functional.md",
		"security-ops.md",
		"doc-writer.md",
	}

	for _, skill := range skills {
		if err := writeTemplate("templates/"+skill, filepath.Join(baseDir, ".agent/skills/"+skill)); err != nil {
			return err
		}
	}

	return nil
}

func writeTemplate(tmplPath, destPath string) error {
	content, err := templatesFS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
	}

	if err := os.WriteFile(destPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}
	return nil
}
