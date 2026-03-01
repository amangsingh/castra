// Package templates provides the central, shared template filesystem
// embedded from the canonical template sources. All platform-specific
// generators (antigravity, copilot, gemini) read from this single FS
// instead of maintaining their own duplicate copies.
package templates

import "embed"

// FS is the embedded filesystem containing all shared templates.
// Structure:
//
//	rules.md                          — The Universal Constitution
//	roles/<role>/SKILL.md             — Role identity and commands
//	roles/<role>/error_handling.md    — Error handling guidelines (if present)
//	roles/<role>/examples.md          — Usage examples (if present)
//	roles/<role>/scripts/main.go      — Wrapper script source
//	workflows/<name>.md               — All workflow definitions (flat)
//
//go:embed rules.md roles workflows
var FS embed.FS
