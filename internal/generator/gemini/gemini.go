package gemini

import (
"fmt"
"io/fs"
"os"
"path/filepath"
"strings"

"castra/internal/generator/templates"
)

// InitWorkspace generates the Gemini Code Assist configuration files
// from the shared templates package.
func InitWorkspace(baseDir string) error {
	dirs := []string{
		".gemini/agents",
		".gemini/workflows",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(baseDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 2. Map rules.md → GEMINI.md
	rulesContent, err := templates.FS.ReadFile("rules.md")
	if err != nil {
		return fmt.Errorf("failed to read rules.md: %w", err)
	}
	destPath := filepath.Join(baseDir, "GEMINI.md")
	if err := os.WriteFile(destPath, rulesContent, 0644); err != nil {
		return fmt.Errorf("failed to write GEMINI.md: %w", err)
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
	roleEntries, err := templates.FS.ReadDir("roles")
	if err != nil {
		return fmt.Errorf("failed to read roles dir: %w", err)
	}
	for _, entry := range roleEntries {
		if !entry.IsDir() {
			continue
		}
		role := entry.Name()
		skillPath := fmt.Sprintf("roles/%s/SKILL.md", role)
		content, err := templates.FS.ReadFile(skillPath)
		if err != nil {
			continue
		}
		agentDest := filepath.Join(baseDir, ".gemini/agents", role+".md")
		if err := os.WriteFile(agentDest, content, 0644); err != nil {
			return fmt.Errorf("failed to write agent file %s: %w", agentDest, err)
		}
	}

	// 4. Map workflows/*.md → .gemini/workflows/*.md
	err = fs.WalkDir(templates.FS, "workflows", func(path string, d fs.DirEntry, err error) error {
if err != nil {
return err
}
if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		content, err := templates.FS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read workflow %s: %w", path, err)
		}
		wfDest := filepath.Join(baseDir, ".gemini/workflows", d.Name())
		if err := os.WriteFile(wfDest, content, 0644); err != nil {
			return fmt.Errorf("failed to write workflow %s: %w", wfDest, err)
		}
		return nil
	})

	return err
}
