package core

import (
	"encoding/json"
	"io/fs"
	"os"
)

// DefaultFileMode is the default permission mode for configuration files.
// This can be used by adapters or overridden with WriteFileWithMode.
const DefaultFileMode fs.FileMode = 0600

// Config represents the canonical MCP configuration that can be
// converted to/from various AI assistant formats.
type Config struct {
	// Servers is a map of server names to their configurations.
	Servers map[string]Server `json:"servers"`

	// Inputs defines input variables for sensitive data like API keys (VS Code feature).
	Inputs []InputVariable `json:"inputs,omitempty"`
}

// InputVariable represents a placeholder for sensitive configuration values.
// When referenced using ${input:variable-id}, the user is prompted for the value.
type InputVariable struct {
	// Type is the input type, typically "promptString".
	Type string `json:"type"`

	// ID is the unique identifier referenced in server config as ${input:id}.
	ID string `json:"id"`

	// Description is the user-friendly prompt text.
	Description string `json:"description"`

	// Password hides typed input when true (for API keys and passwords).
	Password bool `json:"password,omitempty"`
}

// NewConfig creates a new empty Config.
func NewConfig() *Config {
	return &Config{
		Servers: make(map[string]Server),
	}
}

// AddServer adds a server to the configuration.
func (c *Config) AddServer(name string, server Server) {
	if c.Servers == nil {
		c.Servers = make(map[string]Server)
	}
	c.Servers[name] = server
}

// RemoveServer removes a server from the configuration.
func (c *Config) RemoveServer(name string) {
	delete(c.Servers, name)
}

// GetServer returns a server by name and whether it exists.
func (c *Config) GetServer(name string) (Server, bool) {
	server, ok := c.Servers[name]
	return server, ok
}

// ServerNames returns a slice of all server names.
func (c *Config) ServerNames() []string {
	names := make([]string, 0, len(c.Servers))
	for name := range c.Servers {
		names = append(names, name)
	}
	return names
}

// AddInput adds an input variable to the configuration.
func (c *Config) AddInput(input InputVariable) {
	c.Inputs = append(c.Inputs, input)
}

// GetInput returns an input variable by ID and whether it exists.
func (c *Config) GetInput(id string) (InputVariable, bool) {
	for _, input := range c.Inputs {
		if input.ID == id {
			return input, true
		}
	}
	return InputVariable{}, false
}

// StdioServers returns only the servers configured for stdio transport.
func (c *Config) StdioServers() map[string]Server {
	result := make(map[string]Server)
	for name, server := range c.Servers {
		if server.IsStdio() {
			result[name] = server
		}
	}
	return result
}

// RemoteServers returns only the servers configured for HTTP/SSE transport.
func (c *Config) RemoteServers() map[string]Server {
	result := make(map[string]Server)
	for name, server := range c.Servers {
		if server.IsRemote() {
			result[name] = server
		}
	}
	return result
}

// EnabledServers returns only the servers that are enabled.
func (c *Config) EnabledServers() map[string]Server {
	result := make(map[string]Server)
	for name, server := range c.Servers {
		if server.IsEnabled() {
			result[name] = server
		}
	}
	return result
}

// Merge combines another config into this one. Servers and inputs from
// the other config override those in this config with the same name/ID.
func (c *Config) Merge(other *Config) {
	if other == nil {
		return
	}
	for name, server := range other.Servers {
		c.Servers[name] = server
	}
	// Merge inputs, replacing by ID
	inputMap := make(map[string]InputVariable)
	for _, input := range c.Inputs {
		inputMap[input.ID] = input
	}
	for _, input := range other.Inputs {
		inputMap[input.ID] = input
	}
	c.Inputs = make([]InputVariable, 0, len(inputMap))
	for _, input := range inputMap {
		c.Inputs = append(c.Inputs, input)
	}
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	for name, server := range c.Servers {
		if err := server.Validate(); err != nil {
			return &ServerValidationError{Name: name, Err: err}
		}
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (c *Config) MarshalJSON() ([]byte, error) {
	type Alias Config
	return json.Marshal((*Alias)(c))
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := (*Alias)(c)
	return json.Unmarshal(data, aux)
}

// WriteFile writes the config to a file in JSON format using DefaultFileMode.
func (c *Config) WriteFile(path string) error {
	return c.WriteFileWithMode(path, DefaultFileMode)
}

// WriteFileWithMode writes the config to a file in JSON format with the specified permission mode.
func (c *Config) WriteFileWithMode(path string, mode fs.FileMode) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, mode)
}

// ReadFile reads a config from a JSON file.
func ReadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
