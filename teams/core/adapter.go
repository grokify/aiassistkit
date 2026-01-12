package core

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"gopkg.in/yaml.v3"
)

// DefaultFileMode is the default permission for generated files.
const DefaultFileMode fs.FileMode = 0600

// DefaultDirMode is the default permission for generated directories.
const DefaultDirMode fs.FileMode = 0700

// Adapter converts between canonical Team definitions and tool-specific formats.
type Adapter interface {
	// Name returns the adapter identifier (e.g., "agentspec", "crewai").
	Name() string

	// FileExtension returns the file extension for team files.
	FileExtension() string

	// Parse converts tool-specific bytes to canonical Team.
	Parse(data []byte) (*Team, error)

	// Marshal converts canonical Team to tool-specific bytes.
	Marshal(team *Team) ([]byte, error)

	// ReadFile reads from path and returns canonical Team.
	ReadFile(path string) (*Team, error)

	// WriteFile writes canonical Team to path.
	WriteFile(team *Team, path string) error
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

// ReadTeamFile reads a team file (YAML or JSON) and returns the Team.
func ReadTeamFile(path string) (*Team, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ReadError{Path: path, Err: err}
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		return ParseYAML(data, path)
	case ".json":
		return ParseJSON(data, path)
	default:
		// Try YAML first (more permissive), then JSON
		team, err := ParseYAML(data, path)
		if err == nil {
			return team, nil
		}
		return ParseJSON(data, path)
	}
}

// WriteTeamFile writes a Team to a file in YAML format.
func WriteTeamFile(team *Team, path string) error {
	data, err := yaml.Marshal(team)
	if err != nil {
		return &MarshalError{Format: "yaml", Err: err}
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, DefaultDirMode); err != nil {
		return &WriteError{Path: path, Err: err}
	}

	if err := os.WriteFile(path, data, DefaultFileMode); err != nil {
		return &WriteError{Path: path, Err: err}
	}

	return nil
}

// WriteTeamJSON writes a Team to a file in JSON format.
func WriteTeamJSON(team *Team, path string) error {
	data, err := json.MarshalIndent(team, "", "  ")
	if err != nil {
		return &MarshalError{Format: "json", Err: err}
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

// ParseYAML parses YAML bytes into a Team.
func ParseYAML(data []byte, path string) (*Team, error) {
	var team Team
	if err := yaml.Unmarshal(data, &team); err != nil {
		return nil, &ParseError{Format: "yaml", Path: path, Err: err}
	}
	return &team, nil
}

// ParseJSON parses JSON bytes into a Team.
func ParseJSON(data []byte, path string) (*Team, error) {
	var team Team
	if err := json.Unmarshal(data, &team); err != nil {
		return nil, &ParseError{Format: "json", Path: path, Err: err}
	}
	return &team, nil
}

// ReadTeamDir reads all team files from a directory.
func ReadTeamDir(dir string) ([]*Team, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, &ReadError{Path: dir, Err: err}
	}

	var teams []*Team
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		team, err := ReadTeamFile(path)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}
