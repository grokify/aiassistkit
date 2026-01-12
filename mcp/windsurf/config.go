package windsurf

// Config represents the Windsurf MCP configuration.
type Config struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

// ServerConfig represents a Windsurf MCP server configuration.
// Note: Windsurf uses "serverUrl" for HTTP servers instead of "url".
type ServerConfig struct {
	// Type specifies the transport type.
	Type string `json:"type,omitempty"`

	// --- STDIO Server Fields ---
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`

	// --- HTTP Server Fields ---
	// Note: Windsurf uses "serverUrl" instead of "url"
	ServerURL string            `json:"serverUrl,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`

	// --- Tool Control ---
	DisabledTools []string `json:"disabledTools,omitempty"`
}

// NewConfig creates a new Windsurf config.
func NewConfig() *Config {
	return &Config{
		MCPServers: make(map[string]ServerConfig),
	}
}
