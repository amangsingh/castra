package generator

import (
	"castra/internal/config"
	"fmt"
	"sync"
)

// Generator defines the interface for workspace providers.
type Generator interface {
	InitWorkspace(baseDir string, cfg config.VendorConfig) error
}

var (
	registry = make(map[string]Generator)
	mu       sync.RWMutex
)

// Register adds a generator to the global registry.
func Register(name string, gen Generator) {
	mu.Lock()
	defer mu.Unlock()
	registry[name] = gen
}

// Get retrieves a generator from the registry.
func Get(name string) (Generator, bool) {
	mu.RLock()
	defer mu.RUnlock()
	gen, ok := registry[name]
	return gen, ok
}

// List returns all registered generator names.
func List() []string {
	mu.RLock()
	defer mu.RUnlock()
	var names []string
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// InitWorkspaceFromConfig dispatches the workspace generation to the appropriate vendor generator.
func InitWorkspaceFromConfig(baseDir string, vendor string, cfg config.VendorConfig) error {
	gen, ok := Get(vendor)
	if !ok {
		return fmt.Errorf("generator for vendor %q not found", vendor)
	}
	return gen.InitWorkspace(baseDir, cfg)
}
