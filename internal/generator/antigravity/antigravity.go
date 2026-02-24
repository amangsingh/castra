package antigravity

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed templates/*
var templatesFS embed.FS

func InitWorkspace(baseDir string) error {
	// 1. Create base directories
	dirs := []string{
		".agent/rules",
		".agent/skills",
		".agent/workflows",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(baseDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 2. Walk the embedded templates directory
	return fs.WalkDir(templatesFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path from "templates/"
		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		// Determine destination path
		var destPath string
		if relPath == "rules.md" {
			destPath = filepath.Join(baseDir, ".agent/rules/rules.md")
		} else if filepath.Base(filepath.Dir(relPath)) == "workflows" {
			// Files in "workflows" subdirectory go to .agent/workflows/
			destPath = filepath.Join(baseDir, ".agent/workflows", filepath.Base(relPath))
		} else {
			// Everything else goes to .agent/skills/
			destPath = filepath.Join(baseDir, ".agent/skills", relPath)
		}

		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
		} else {
			content, err := templatesFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read template %s: %w", path, err)
			}
			if err := os.WriteFile(destPath, content, 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", destPath, err)
			}
		}

		return nil
	})
}
