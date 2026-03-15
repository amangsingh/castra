package gemini

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"castra/internal/config"
	"castra/internal/generator"
	"castra/internal/generator/common"
	"castra/internal/generator/templates"
)

type GeminiGenerator struct{}

func init() {
	generator.Register("gemini", &GeminiGenerator{})
}

// InitWorkspace implements the generator.Generator interface.
func (g *GeminiGenerator) InitWorkspace(baseDir string, _ config.VendorConfig) error {
	return InitWorkspace(baseDir)
}

// InitWorkspace generates the Gemini Code Assist configuration files
// from the shared templates package.
func InitWorkspace(baseDir string) error {
	dirs := []string{
		".gemini/agents",
		".gemini/workflows",
	}
	if err := common.EnsureDirs(baseDir, dirs); err != nil {
		return err
	}

	// 2. Map rules.md → GEMINI.md
	if err := common.WriteTemplateFile(templates.FS, "rules.md", filepath.Join(baseDir, "GEMINI.md")); err != nil {
		return err
	}

	// 2.5 Create .gemini/settings.json
	settingsJSON := []byte(`{
  "project.features.enableAssistant": true,
  "project.features.customInstructions": true
}`)
	settingsPath := filepath.Join(baseDir, ".gemini/settings.json")
	if err := os.WriteFile(settingsPath, settingsJSON, 0644); err != nil {
		return fmt.Errorf("failed to write settings.json: %w", err)
	}

	// 3. Map roles/*/SKILL.md → .gemini/agents/<role>.md
	err := common.WalkAndMap(templates.FS, "roles", filepath.Join(baseDir, ".gemini/agents"), func(path string, d fs.DirEntry) (string, bool, error) {
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

	// 4. Map workflows/*.md → .gemini/workflows/*.md
	err = common.WalkAndMap(templates.FS, "workflows", filepath.Join(baseDir, ".gemini/workflows"), func(path string, d fs.DirEntry) (string, bool, error) {
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
