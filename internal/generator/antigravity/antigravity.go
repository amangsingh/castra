package antigravity

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"castra/internal/generator/common"
	"castra/internal/generator/templates"
)

func InitWorkspace(baseDir string) error {
	// 1. Create base directories
	dirs := []string{
		".agent/rules",
		".agent/skills",
		".agent/workflows",
		".github/agents",
	}
	if err := common.EnsureDirs(baseDir, dirs); err != nil {
		return err
	}

	// 2. Write rules.md → .agent/rules/rules.md
	if err := common.WriteTemplateFile(templates.FS, "rules.md", filepath.Join(baseDir, ".agent/rules/rules.md")); err != nil {
		return err
	}

	// 3. Walk roles/ → .agent/skills/<role>/
	err := common.WalkAndMap(templates.FS, "roles", filepath.Join(baseDir, ".agent/skills"), func(path string, d fs.DirEntry) (string, bool, error) {
		// Calculate relative path from "roles/"
		relPath, err := filepath.Rel("roles", path)
		if err != nil {
			return "", false, err
		}

		if relPath == "." {
			return "", false, nil
		}

		// Skip scripts directories — handled separately via compilation
		if d.IsDir() && d.Name() == "scripts" {
			return "", true, nil
		}

		// Skip scripts source files (reached via non-dir walk)
		parts := strings.Split(relPath, string(filepath.Separator))
		for _, p := range parts {
			if p == "scripts" {
				return "", true, nil
			}
		}

		return relPath, false, nil
	})
	if err != nil {
		return err
	}

	// 4. Walk workflows/ → .agent/workflows/ (flat)
	err = common.WalkAndMap(templates.FS, "workflows", filepath.Join(baseDir, ".agent/workflows"), func(path string, d fs.DirEntry) (string, bool, error) {
		if d.IsDir() {
			return "", false, nil
		}
		return d.Name(), false, nil
	})
	if err != nil {
		return err
	}

	// 5. Compile role wrapper scripts into binaries
	if err := compileRoleScripts(baseDir); err != nil {
		return err
	}

	return nil
}

// compileRoleScripts finds each role's scripts/main.go in the shared templates,
// compiles it to a binary, and places it at .agent/skills/<role>/scripts/castra.
// Falls back to a shell script if Go is not available.
func compileRoleScripts(baseDir string) error {
	entries, err := templates.FS.ReadDir("roles")
	if err != nil {
		return fmt.Errorf("failed to read roles dir: %w", err)
	}

	goAvailable := isGoAvailable()

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		role := entry.Name()
		srcPath := fmt.Sprintf("roles/%s/scripts/main.go", role)

		// Check if this role has a scripts/main.go
		content, err := templates.FS.ReadFile(srcPath)
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
