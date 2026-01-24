package core

import (
	"fmt"
	"sync"
)

// Adapter defines the interface for power format adapters.
// Each platform (Kiro, etc.) implements this interface to convert
// between canonical Power and platform-specific formats.
type Adapter interface {
	// Name returns the adapter identifier (e.g., "kiro").
	Name() string

	// GeneratePowerDir creates a complete power directory structure.
	// Returns the paths of all created files.
	GeneratePowerDir(power *Power, outputDir string) ([]string, error)

	// ParsePowerDir reads a power directory and returns a canonical Power.
	ParsePowerDir(dir string) (*Power, error)
}

// registry holds registered adapters.
var (
	registry     = make(map[string]Adapter)
	registryLock sync.RWMutex
)

// Register registers an adapter with the global registry.
// This is typically called from adapter init() functions.
func Register(adapter Adapter) {
	registryLock.Lock()
	defer registryLock.Unlock()
	registry[adapter.Name()] = adapter
}

// Get returns an adapter by name.
func Get(name string) (Adapter, error) {
	registryLock.RLock()
	defer registryLock.RUnlock()

	adapter, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("power adapter not found: %s", name)
	}
	return adapter, nil
}

// List returns all registered adapter names.
func List() []string {
	registryLock.RLock()
	defer registryLock.RUnlock()

	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// GenerateError represents an error during power generation.
type GenerateError struct {
	Format  string
	Path    string
	Message string
	Err     error
}

func (e *GenerateError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("power generation error (%s) at %s: %s: %v", e.Format, e.Path, e.Message, e.Err)
	}
	return fmt.Sprintf("power generation error (%s) at %s: %s", e.Format, e.Path, e.Message)
}

func (e *GenerateError) Unwrap() error {
	return e.Err
}

// ParseError represents an error during power parsing.
type ParseError struct {
	Format string
	Path   string
	Err    error
}

func (e *ParseError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("power parse error (%s) at %s: %v", e.Format, e.Path, e.Err)
	}
	return fmt.Sprintf("power parse error (%s): %v", e.Format, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}
