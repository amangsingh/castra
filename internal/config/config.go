package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// VendorConfig holds settings for a specific workspace provider (antigravity, claude, etc.)
type VendorConfig struct {
	Roles []string
}

// CastraConfig maps vendor names to their respective configurations.
type CastraConfig struct {
	Vendors map[string]VendorConfig
}

const configTemplate = `# castra.yaml — Workspace Configuration
#
# This file declares your agent workspace. Run 'castra init' after configuring.
#
# Available roles:
#   architect, designer, doc-writer, junior-engineer,
#   qa-functional, security-ops, senior-engineer

# Antigravity (Advanced Agentic Coding Environment)
antigravity:
  roles:
    - architect
    - senior-engineer
    - qa-functional

# Claude (Anthropic's Project Environment)
# claude:
#   roles:
#     - architect
#     - senior-engineer

# Copilot (GitHub Copilot Agent Environment)
# copilot:
#   roles:
#     - architect
#     - senior-engineer

# Gemini (Google Code Assist Environment)
# gemini:
#   roles:
#     - architect
#     - senior-engineer
`

// GenerateTemplate writes a default castra.yaml if it doesn't already exist.
func GenerateTemplate(destDir string) error {
	dest := filepath.Join(destDir, "castra.yaml")
	if _, err := os.Stat(dest); err == nil {
		return nil // Already exists
	}
	return os.WriteFile(dest, []byte(configTemplate), 0644)
}

// Parse reads castra.yaml and returns a vendor-agnostic configuration object.
func Parse(yamlPath string) (*CastraConfig, error) {
	f, err := os.Open(yamlPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &CastraConfig{
		Vendors: make(map[string]VendorConfig),
	}

	var currentVendor string
	var inRoles bool

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Skip comments and empty lines
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Detect top-level vendor key (no leading indentation)
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			if strings.HasSuffix(trimmed, ":") {
				currentVendor = strings.TrimSuffix(trimmed, ":")
				cfg.Vendors[currentVendor] = VendorConfig{Roles: []string{}}
				inRoles = false
			}
			continue
		}

		if currentVendor != "" {
			indent := len(line) - len(strings.TrimLeft(line, " \t"))
			if indent >= 1 && indent <= 3 {
				// Second-level key (roles:)
				if strings.HasPrefix(trimmed, "roles:") {
					inRoles = true
					continue
				}
			}

			if inRoles && strings.HasPrefix(trimmed, "- ") {
				role := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
				if role != "" {
					vendorCfg := cfg.Vendors[currentVendor]
					vendorCfg.Roles = append(vendorCfg.Roles, role)
					cfg.Vendors[currentVendor] = vendorCfg
				}
			}
		}
	}

	return cfg, scanner.Err()
}
