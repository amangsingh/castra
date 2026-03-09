package copilot

import (
	"io/fs"
	"path/filepath"
	"strings"

	"castra/internal/generator/common"
	"castra/internal/generator/templates"
)

// InitWorkspace generates the .github/ Copilot configuration files
// from the shared templates package.
func InitWorkspace(baseDir string) error {
	// 1. Create base directories
	dirs := []string{
		".github/agents",
		".github/castra-workflows",
	}
	if err := common.EnsureDirs(baseDir, dirs); err != nil {
		return err
	}

	// 2. Map rules.md → .github/copilot-instructions.md
	if err := common.WriteTemplateFile(templates.FS, "rules.md", filepath.Join(baseDir, ".github/copilot-instructions.md")); err != nil {
		return err
	}

	// 3. Map roles/*/SKILL.md → .github/agents/<role>.md
	err := common.WalkAndMap(templates.FS, "roles", filepath.Join(baseDir, ".github/agents"), func(path string, d fs.DirEntry) (string, bool, error) {
		if d.IsDir() {
			return "", false, nil
		}
		if d.Name() != "SKILL.md" {
			return "", true, nil
		}

		// Calculate role name: roles/<role>/SKILL.md
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			return "", true, nil
		}
		role := parts[1]
		return role + ".md", false, nil
	})
	if err != nil {
		return err
	}

	// 4. Map workflows/*.md → .github/castra-workflows/*.md
	err = common.WalkAndMap(templates.FS, "workflows", filepath.Join(baseDir, ".github/castra-workflows"), func(path string, d fs.DirEntry) (string, bool, error) {
		if d.IsDir() {
			return "", false, nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return "", true, nil
		}
		return d.Name(), false, nil
	})

	return err
}
