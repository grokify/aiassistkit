package kiro

// Config represents the Kiro MCP configuration file.
type Config struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

// ServerConfig represents a Kiro MCP server configuration.
type ServerConfig struct {
	// --- STDIO Server Fields ---

	// Command is the executable to launch for stdio servers.
	Command string `json:"command,omitempty"`

	// Args are the command-line arguments passed to the command.
	Args []string `json:"args,omitempty"`

	// Env contains environment variables for the server process.
	// Values can use ${ENV_VAR} syntax for substitution.
	Env map[string]string `json:"env,omitempty"`

	// --- Remote/HTTP Server Fields ---

	// URL is the endpoint for remote HTTP/SSE servers.
	URL string `json:"url,omitempty"`

	// Headers contains HTTP headers for authentication or configuration.
	// Values can use ${ENV_VAR} syntax for secrets.
	Headers map[string]string `json:"headers,omitempty"`

	// --- Server State ---

	// Disabled indicates whether the server is disabled.
	Disabled bool `json:"disabled,omitempty"`
}

// NewConfig creates a new Kiro config.
func NewConfig() *Config {
	return &Config{
		MCPServers: make(map[string]ServerConfig),
	}
}
