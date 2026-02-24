package antigravity

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	err := fs.WalkDir(templatesFS, "templates", func(path string, d fs.DirEntry, err error) error {
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

		// Skip workflow directories entirely — files are routed to .agent/workflows/,
		// so creating the dir under skills/ would leave it empty.
		if d.IsDir() && d.Name() == "workflows" {
			return nil
		}

		// Skip scripts directories — we handle script compilation separately below.
		if d.IsDir() && d.Name() == "scripts" {
			return nil
		}

		// Determine destination path
		var destPath string
		if relPath == "rules.md" {
			destPath = filepath.Join(baseDir, ".agent/rules/rules.md")
		} else if filepath.Base(filepath.Dir(relPath)) == "workflows" {
			// Files in "workflows" subdirectory go to .agent/workflows/
			destPath = filepath.Join(baseDir, ".agent/workflows", filepath.Base(relPath))
		} else if filepath.Base(filepath.Dir(relPath)) == "scripts" {
			// Skip .go source files — they'll be compiled in the post-walk step.
			return nil
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

	if err != nil {
		return err
	}

	// 3. Compile role wrapper scripts into binaries
	if err := compileRoleScripts(baseDir); err != nil {
		return err
	}

	return nil
}

// compileRoleScripts finds each role's scripts/main.go in the embedded templates,
// compiles it to a binary, and places it at .agent/skills/<role>/scripts/castra.
// Falls back to a shell script if Go is not available.
func compileRoleScripts(baseDir string) error {
	entries, err := templatesFS.ReadDir("templates")
	if err != nil {
		return fmt.Errorf("failed to read templates dir: %w", err)
	}

	goAvailable := isGoAvailable()

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		role := entry.Name()
		srcPath := fmt.Sprintf("templates/%s/scripts/main.go", role)

		// Check if this role has a scripts/main.go
		content, err := templatesFS.ReadFile(srcPath)
		if err != nil {
			continue // No script for this role
		}

		// Create the scripts directory
		scriptsDir := filepath.Join(baseDir, ".agent/skills", role, "scripts")
		if err := os.MkdirAll(scriptsDir, 0755); err != nil {
			return fmt.Errorf("failed to create scripts dir for %s: %w", role, err)
		}

		if goAvailable {
			if err := compileToBinary(content, scriptsDir); err != nil {
				log.Printf("Warning: failed to compile script for %s: %v. Falling back to shell script.", role, err)
				if err := generateShellScript(role, scriptsDir); err != nil {
					return fmt.Errorf("failed to generate shell script for %s: %w", role, err)
				}
			}
		} else {
			log.Printf("Go not found. Generating shell script wrapper for %s.", role)
			if err := generateShellScript(role, scriptsDir); err != nil {
				return fmt.Errorf("failed to generate shell script for %s: %w", role, err)
			}
		}
	}

	return nil
}

// compileToBinary writes the Go source to a temp dir and compiles it.
func compileToBinary(source []byte, destDir string) error {
	// Create a temp directory for the build
	tmpDir, err := os.MkdirTemp("", "castra-build-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Write the source file
	srcFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(srcFile, source, 0644); err != nil {
		return err
	}

	// Determine output binary name
	binaryName := "castra"
	if runtime.GOOS == "windows" {
		binaryName = "castra.exe"
	}
	outPath := filepath.Join(destDir, binaryName)

	// Compile
	cmd := exec.Command("go", "build", "-o", outPath, srcFile)
	cmd.Dir = tmpDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go build failed: %s: %w", string(output), err)
	}

	return nil
}

// generateShellScript creates a shell script wrapper as fallback when Go is unavailable.
func generateShellScript(role, destDir string) error {
	if runtime.GOOS == "windows" {
		// Generate .bat file for Windows
		content := fmt.Sprintf("@echo off\r\ncastra %%* --role %s\r\n", role)
		return os.WriteFile(filepath.Join(destDir, "castra.bat"), []byte(content), 0755)
	}

	// Generate shell script for Unix
	content := fmt.Sprintf("#!/bin/sh\nexec castra \"$@\" --role %s\n", role)
	return os.WriteFile(filepath.Join(destDir, "castra"), []byte(content), 0755)
}

// isGoAvailable checks if the Go toolchain is installed and accessible.
func isGoAvailable() bool {
	_, err := exec.LookPath("go")
	return err == nil
}
