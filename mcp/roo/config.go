package roo

// Config represents the Roo Code MCP configuration.
type Config struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

// ServerConfig represents a Roo Code MCP server configuration.
type ServerConfig struct {
	// Type specifies the transport type.
	Type string `json:"type,omitempty"`

	// --- STDIO Server Fields ---
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`

	// --- HTTP/SSE Server Fields ---
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`

	// --- Roo-specific Fields ---
	AlwaysAllow []string `json:"alwaysAllow,omitempty"`
	Disabled    bool     `json:"disabled,omitempty"`
}

// NewConfig creates a new Roo Code config.
func NewConfig() *Config {
	return &Config{
		MCPServers: make(map[string]ServerConfig),
	}
}
