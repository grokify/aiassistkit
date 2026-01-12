package cline

// Config represents the Cline MCP configuration.
type Config struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

// ServerConfig represents a Cline MCP server configuration.
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

	// --- Cline-specific Fields ---

	// AlwaysAllow lists tools that don't require user approval.
	AlwaysAllow []string `json:"alwaysAllow,omitempty"`

	// Disabled indicates whether the server is disabled.
	Disabled bool `json:"disabled,omitempty"`
}

// NewConfig creates a new Cline config.
func NewConfig() *Config {
	return &Config{
		MCPServers: make(map[string]ServerConfig),
	}
}
