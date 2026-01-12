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

// Adapter converts between canonical Command and tool-specific formats.
type Adapter interface {
	// Name returns the adapter identifier (e.g., "claude", "gemini", "codex").
	Name() string

	// FileExtension returns the file extension for this format (e.g., ".md", ".toml").
	FileExtension() string

	// DefaultDir returns the default directory name for commands.
	DefaultDir() string

	// Parse converts tool-specific bytes to canonical Command.
	Parse(data []byte) (*Command, error)

	// Marshal converts canonical Command to tool-specific bytes.
	Marshal(cmd *Command) ([]byte, error)

	// ReadFile reads from path and returns canonical Command.
	ReadFile(path string) (*Command, error)

	// WriteFile writes canonical Command to path.
	WriteFile(cmd *Command, path string) error
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

// Convert converts command data from one format to another.
func (r *Registry) Convert(data []byte, from, to string) ([]byte, error) {
	fromAdapter, ok := r.GetAdapter(from)
	if !ok {
		return nil, fmt.Errorf("unknown source adapter: %s", from)
	}

	toAdapter, ok := r.GetAdapter(to)
	if !ok {
		return nil, fmt.Errorf("unknown target adapter: %s", to)
	}

	cmd, err := fromAdapter.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", from, err)
	}

	return toAdapter.Marshal(cmd)
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

// ReadCanonicalFile reads a canonical command.json file.
func ReadCanonicalFile(path string) (*Command, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ReadError{Path: path, Err: err}
	}

	var cmd Command
	if err := json.Unmarshal(data, &cmd); err != nil {
		return nil, &ParseError{Format: "canonical", Path: path, Err: err}
	}

	return &cmd, nil
}

// WriteCanonicalFile writes a canonical command.json file.
func WriteCanonicalFile(cmd *Command, path string) error {
	data, err := json.MarshalIndent(cmd, "", "  ")
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

// ReadCanonicalDir reads all command.json files from a directory.
func ReadCanonicalDir(dir string) ([]*Command, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, &ReadError{Path: dir, Err: err}
	}

	var commands []*Command
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		cmd, err := ReadCanonicalFile(path)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}

// WriteCommandsToDir writes multiple commands to a directory using the specified adapter.
func WriteCommandsToDir(commands []*Command, dir string, adapterName string) error {
	adapter, ok := GetAdapter(adapterName)
	if !ok {
		return fmt.Errorf("unknown adapter: %s", adapterName)
	}

	if err := os.MkdirAll(dir, DefaultDirMode); err != nil {
		return &WriteError{Path: dir, Err: err}
	}

	for _, cmd := range commands {
		filename := cmd.Name + adapter.FileExtension()
		path := filepath.Join(dir, filename)
		if err := adapter.WriteFile(cmd, path); err != nil {
			return err
		}
	}

	return nil
}
