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

// Adapter converts between canonical Skill and tool-specific formats.
type Adapter interface {
	// Name returns the adapter identifier (e.g., "claude", "codex").
	Name() string

	// SkillFileName returns the skill definition filename (e.g., "SKILL.md").
	SkillFileName() string

	// DefaultDir returns the default directory name for skills.
	DefaultDir() string

	// Parse converts tool-specific bytes to canonical Skill.
	Parse(data []byte) (*Skill, error)

	// Marshal converts canonical Skill to tool-specific bytes.
	Marshal(skill *Skill) ([]byte, error)

	// ReadFile reads from path and returns canonical Skill.
	ReadFile(path string) (*Skill, error)

	// WriteFile writes canonical Skill to path.
	WriteFile(skill *Skill, path string) error

	// WriteSkillDir writes the complete skill directory structure.
	WriteSkillDir(skill *Skill, baseDir string) error
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

// Convert converts skill data from one format to another.
func (r *Registry) Convert(data []byte, from, to string) ([]byte, error) {
	fromAdapter, ok := r.GetAdapter(from)
	if !ok {
		return nil, fmt.Errorf("unknown source adapter: %s", from)
	}

	toAdapter, ok := r.GetAdapter(to)
	if !ok {
		return nil, fmt.Errorf("unknown target adapter: %s", to)
	}

	skill, err := fromAdapter.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", from, err)
	}

	return toAdapter.Marshal(skill)
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

// ReadCanonicalFile reads a canonical skill.json file.
func ReadCanonicalFile(path string) (*Skill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ReadError{Path: path, Err: err}
	}

	var skill Skill
	if err := json.Unmarshal(data, &skill); err != nil {
		return nil, &ParseError{Format: "canonical", Path: path, Err: err}
	}

	return &skill, nil
}

// WriteCanonicalFile writes a canonical skill.json file.
func WriteCanonicalFile(skill *Skill, path string) error {
	data, err := json.MarshalIndent(skill, "", "  ")
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

// ReadCanonicalDir reads all skill.json files from subdirectories.
func ReadCanonicalDir(dir string) ([]*Skill, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, &ReadError{Path: dir, Err: err}
	}

	var skills []*Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillPath := filepath.Join(dir, entry.Name(), "skill.json")
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			continue
		}

		skill, err := ReadCanonicalFile(skillPath)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

// WriteSkillsToDir writes multiple skills to a directory using the specified adapter.
func WriteSkillsToDir(skills []*Skill, dir string, adapterName string) error {
	adapter, ok := GetAdapter(adapterName)
	if !ok {
		return fmt.Errorf("unknown adapter: %s", adapterName)
	}

	if err := os.MkdirAll(dir, DefaultDirMode); err != nil {
		return &WriteError{Path: dir, Err: err}
	}

	for _, skill := range skills {
		if err := adapter.WriteSkillDir(skill, dir); err != nil {
			return err
		}
	}

	return nil
}
