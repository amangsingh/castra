package antigravity

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"castra/internal/config"
	"castra/internal/generator/common"
	"castra/internal/generator/templates"
)

// GenerateConfigTemplate is a wrapper around the agnostic config template generator.
func GenerateConfigTemplate(baseDir string) error {
	return config.GenerateTemplate(baseDir)
}

// InitWorkspaceFromConfig is the entry point for generating the .agent workspace for Antigravity.
func InitWorkspaceFromConfig(baseDir string, cfg config.VendorConfig) error {
	if len(cfg.Roles) == 0 {
		return fmt.Errorf("antigravity: no roles specified. Please add at least one role to 'antigravity.roles'")
	}

	// Validate roles against available templates.
	availableRoles, err := listAvailableRoles()
	if err != nil {
		return err
	}
	for _, role := range cfg.Roles {
		if !availableRoles[role] {
			return fmt.Errorf("castra.yaml: unknown role %q. Available roles: %s", role, strings.Join(keys(availableRoles), ", "))
		}
	}

	// Create base directories.
	dirs := []string{
		".agent/rules",
		".agent/skills",
		".agent/workflows",
		".github/agents",
	}
	if err := common.EnsureDirs(baseDir, dirs); err != nil {
		return err
	}

	// Write rules.md.
	if err := common.WriteTemplateFile(templates.FS, "rules.md", filepath.Join(baseDir, ".agent/rules/rules.md")); err != nil {
		return err
	}

	// Walk workflows/ → .agent/workflows/ (flat).
	err = common.WalkAndMap(templates.FS, "workflows", filepath.Join(baseDir, ".agent/workflows"), func(path string, d fs.DirEntry) (string, bool, error) {
		if d.IsDir() {
			return "", false, nil
		}
		return d.Name(), false, nil
	})
	if err != nil {
		return err
	}

	// Generate only the requested roles.
	for _, role := range cfg.Roles {
		roleTemplateDir := "roles/" + role
		roleDestDir := filepath.Join(baseDir, ".agent/skills", role)

		err := common.WalkAndMap(templates.FS, roleTemplateDir, roleDestDir, func(path string, d fs.DirEntry) (string, bool, error) {
			relPath, err := filepath.Rel(roleTemplateDir, path)
			if err != nil {
				return "", false, err
			}
			if relPath == "." {
				return "", false, nil
			}
			// Skip scripts source files — handled separately via compilation.
			parts := strings.Split(relPath, string(filepath.Separator))
			for _, p := range parts {
				if p == "scripts" {
					return "", true, nil
				}
			}
			return relPath, false, nil
		})
		if err != nil {
			return fmt.Errorf("failed to generate role %s: %w", role, err)
		}
	}

	// Compile role wrapper scripts only for selected roles.
	if err := compileRoleScriptsFor(baseDir, cfg.Roles); err != nil {
		return err
	}

	return nil
}

// listAvailableRoles returns a set of role names available in the templates FS.
func listAvailableRoles() (map[string]bool, error) {
	entries, err := templates.FS.ReadDir("roles")
	if err != nil {
		return nil, fmt.Errorf("failed to list available roles: %w", err)
	}
	m := make(map[string]bool, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			m[e.Name()] = true
		}
	}
	return m, nil
}

// compileRoleScriptsFor compiles scripts only for the specified roles.
func compileRoleScriptsFor(baseDir string, roles []string) error {
	goAvailable := isGoAvailable()

	for _, role := range roles {
		srcPath := fmt.Sprintf("roles/%s/scripts/main.go", role)
		content, err := templates.FS.ReadFile(srcPath)
		if err != nil {
			continue // No script for this role.
		}

		scriptsDir := filepath.Join(baseDir, ".agent/skills", role, "scripts")
		if err := os.MkdirAll(scriptsDir, 0755); err != nil {
			return fmt.Errorf("failed to create scripts dir for %s: %w", role, err)
		}

		if goAvailable {
			if err := compileToBinary(content, scriptsDir); err != nil {
				if err := generateShellScript(role, scriptsDir); err != nil {
					return fmt.Errorf("failed to generate shell script for %s: %w", role, err)
				}
			}
		} else {
			if err := generateShellScript(role, scriptsDir); err != nil {
				return fmt.Errorf("failed to generate shell script for %s: %w", role, err)
			}
		}
	}
	return nil
}

// keys returns the sorted keys of a string bool map.
func keys(m map[string]bool) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
