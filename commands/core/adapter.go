package core

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

// ReadCanonicalFile reads a canonical command file (JSON or Markdown with YAML frontmatter).
// The format is auto-detected based on file extension or content.
func ReadCanonicalFile(path string) (*Command, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ReadError{Path: path, Err: err}
	}

	// Detect format: if it starts with "---" or has .md extension, parse as markdown
	ext := filepath.Ext(path)
	if ext == ".md" || (len(data) >= 3 && string(data[:3]) == "---") {
		cmd, err := ParseCommandMarkdown(data)
		if err != nil {
			return nil, &ParseError{Format: "markdown", Path: path, Err: err}
		}
		// Infer name from filename if not set
		if cmd.Name == "" {
			base := filepath.Base(path)
			cmd.Name = strings.TrimSuffix(base, filepath.Ext(base))
		}
		return cmd, nil
	}

	// Fall back to JSON
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

// ReadCanonicalDir reads all command files (.json or .md) from a directory.
func ReadCanonicalDir(dir string) ([]*Command, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, &ReadError{Path: dir, Err: err}
	}

	var commands []*Command
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".json" && ext != ".md" {
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

// ParseCommandMarkdown parses a Markdown file with YAML frontmatter into a Command.
// The frontmatter should contain: name, description, arguments, dependencies, process.
// The body becomes the instructions.
func ParseCommandMarkdown(data []byte) (*Command, error) {
	content := string(data)

	if !strings.HasPrefix(content, "---") {
		// No frontmatter, treat entire content as instructions
		return &Command{Instructions: strings.TrimSpace(content)}, nil
	}

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return &Command{Instructions: strings.TrimSpace(content)}, nil
	}

	cmd := &Command{}

	// Parse YAML frontmatter
	lines := strings.Split(strings.TrimSpace(parts[1]), "\n")
	var currentKey string
	var listItems []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Check if this is a list item (starts with -)
		if strings.HasPrefix(trimmed, "- ") {
			if currentKey != "" {
				listItems = append(listItems, strings.TrimPrefix(trimmed, "- "))
			}
			continue
		}

		// Process any accumulated list items
		if currentKey != "" && len(listItems) > 0 {
			switch currentKey {
			case "dependencies":
				cmd.Dependencies = listItems
			case "process":
				cmd.Process = listItems
			}
			listItems = nil
		}

		// Parse key: value
		idx := strings.Index(trimmed, ":")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(trimmed[:idx])
		value := strings.TrimSpace(trimmed[idx+1:])
		value = strings.Trim(value, "\"'")

		currentKey = key

		switch key {
		case "name":
			cmd.Name = value
		case "description":
			cmd.Description = value
		case "dependencies":
			if value != "" {
				cmd.Dependencies = parseList(value)
			}
			// Otherwise wait for list items
		case "process":
			if value != "" {
				cmd.Process = parseList(value)
			}
			// Otherwise wait for list items
		case "arguments":
			// Arguments are handled specially - look for inline list or skip
			if value != "" {
				// Could be inline like: [version]
				cmd.Arguments = parseArguments(value)
			}
		}
	}

	// Process any remaining list items
	if currentKey != "" && len(listItems) > 0 {
		switch currentKey {
		case "dependencies":
			cmd.Dependencies = listItems
		case "process":
			cmd.Process = listItems
		}
	}

	// Body becomes instructions
	cmd.Instructions = strings.TrimSpace(parts[2])

	return cmd, nil
}

// parseList parses a comma-separated or bracket-enclosed list.
func parseList(s string) []string {
	s = strings.Trim(s, "[]")
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, "\"'")
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// parseArguments parses an inline arguments list like [version, target].
func parseArguments(s string) []Argument {
	names := parseList(s)
	var args []Argument
	for _, name := range names {
		// Check if required (no ? suffix)
		required := true
		if strings.HasSuffix(name, "?") {
			required = false
			name = strings.TrimSuffix(name, "?")
		}
		args = append(args, Argument{
			Name:     name,
			Type:     "string",
			Required: required,
		})
	}
	return args
}
