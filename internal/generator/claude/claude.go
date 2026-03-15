package claude

import (
	"io/fs"
	"path/filepath"
	"strings"

	"castra/internal/config"
	"castra/internal/generator"
	"castra/internal/generator/common"
	"castra/internal/generator/templates"
)

type ClaudeGenerator struct{}

func init() {
	generator.Register("claude", &ClaudeGenerator{})
}

// InitWorkspace implements the generator.Generator interface.
func (g *ClaudeGenerator) InitWorkspace(baseDir string, _ config.VendorConfig) error {
	// 1. Create base directories
	dirs := []string{
		".claude/agents",
		".claude/workflows",
	}
	if err := common.EnsureDirs(baseDir, dirs); err != nil {
		return err
	}

	// 2. Map rules.md → CLAUDE.md
	if err := common.WriteTemplateFile(templates.FS, "rules.md", filepath.Join(baseDir, "CLAUDE.md")); err != nil {
		return err
	}

	// 3. Map roles/*/SKILL.md → .claude/agents/<role>.md
	err := common.WalkAndMap(templates.FS, "roles", filepath.Join(baseDir, ".claude/agents"), func(path string, d fs.DirEntry) (string, bool, error) {
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

	// 4. Map workflows/*.md → .claude/workflows/*.md
	err = common.WalkAndMap(templates.FS, "workflows", filepath.Join(baseDir, ".claude/workflows"), func(path string, d fs.DirEntry) (string, bool, error) {
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
