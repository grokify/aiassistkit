package core

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// DefaultFileMode is the default permission for generated files.
const DefaultFileMode fs.FileMode = 0600

// DefaultDirMode is the default permission for generated directories.
const DefaultDirMode fs.FileMode = 0700

// Adapter converts between canonical Plugin and tool-specific formats.
type Adapter interface {
	// Name returns the adapter identifier (e.g., "claude", "gemini").
	Name() string

	// DefaultPaths returns default file paths for this tool's plugin manifest.
	DefaultPaths() []string

	// Parse converts tool-specific bytes to canonical Plugin.
	Parse(data []byte) (*Plugin, error)

	// Marshal converts canonical Plugin to tool-specific bytes.
	Marshal(plugin *Plugin) ([]byte, error)

	// ReadFile reads from path and returns canonical Plugin.
	ReadFile(path string) (*Plugin, error)

	// WriteFile writes canonical Plugin to path.
	WriteFile(plugin *Plugin, path string) error

	// WritePlugin writes the complete plugin structure to the given directory.
	// This includes the manifest file and any referenced component directories.
	WritePlugin(plugin *Plugin, dir string) error
}

// Registry manages adapter registration and lookup.
type Registry struct {
	mu       sync.RWMutex
	adapters map[string]Adapter
}

// NewRegistry creates a new adapter registry.
func NewRegistry() *Registry {
	return &Registry{
		adapters: make(map[string]Adapter),
	}
}

// Register adds an adapter to the registry.
func (r *Registry) Register(adapter Adapter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters[adapter.Name()] = adapter
}

// GetAdapter returns an adapter by name.
func (r *Registry) GetAdapter(name string) (Adapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	adapter, ok := r.adapters[name]
	return adapter, ok
}

// AdapterNames returns all registered adapter names sorted alphabetically.
func (r *Registry) AdapterNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Convert converts plugin data from one format to another.
func (r *Registry) Convert(data []byte, from, to string) ([]byte, error) {
	fromAdapter, ok := r.GetAdapter(from)
	if !ok {
		return nil, fmt.Errorf("unknown source adapter: %s", from)
	}

	toAdapter, ok := r.GetAdapter(to)
	if !ok {
		return nil, fmt.Errorf("unknown target adapter: %s", to)
	}

	plugin, err := fromAdapter.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", from, err)
	}

	return toAdapter.Marshal(plugin)
}

// DefaultRegistry is the global adapter registry.
var DefaultRegistry = NewRegistry()

// Register adds an adapter to the default registry.
func Register(adapter Adapter) {
	DefaultRegistry.Register(adapter)
}

// GetAdapter returns an adapter from the default registry.
func GetAdapter(name string) (Adapter, bool) {
	return DefaultRegistry.GetAdapter(name)
}

// AdapterNames returns adapter names from the default registry.
func AdapterNames() []string {
	return DefaultRegistry.AdapterNames()
}

// Convert converts using the default registry.
func Convert(data []byte, from, to string) ([]byte, error) {
	return DefaultRegistry.Convert(data, from, to)
}

// ReadCanonicalFile reads a canonical plugin.json file.
func ReadCanonicalFile(path string) (*Plugin, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ReadError{Path: path, Err: err}
	}

	var plugin Plugin
	if err := json.Unmarshal(data, &plugin); err != nil {
		return nil, &ParseError{Format: "canonical", Path: path, Err: err}
	}

	return &plugin, nil
}

// WriteCanonicalFile writes a canonical plugin.json file.
func WriteCanonicalFile(plugin *Plugin, path string) error {
	data, err := json.MarshalIndent(plugin, "", "  ")
	if err != nil {
		return &MarshalError{Format: "canonical", Err: err}
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, DefaultDirMode); err != nil {
		return &WriteError{Path: path, Err: err}
	}

	if err := os.WriteFile(path, append(data, '\n'), DefaultFileMode); err != nil {
		return &WriteError{Path: path, Err: err}
	}

	return nil
}
