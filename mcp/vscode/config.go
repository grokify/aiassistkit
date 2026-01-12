package vscode

// Config represents the VS Code MCP configuration.
// Note: VS Code uses "servers" not "mcpServers".
type Config struct {
	// Inputs defines input variables for sensitive data.
	Inputs []InputVariable `json:"inputs,omitempty"`

	// Servers maps server names to their configurations.
	Servers map[string]ServerConfig `json:"servers"`
}

// InputVariable represents a placeholder for sensitive values.
type InputVariable struct {
	// Type is the input type, typically "promptString".
	Type string `json:"type"`

	// ID is the unique identifier referenced as ${input:id}.
	ID string `json:"id"`

	// Description is the user-friendly prompt text.
	Description string `json:"description"`

	// Password hides input when true.
	Password bool `json:"password,omitempty"`
}

// ServerConfig represents a VS Code MCP server configuration.
type ServerConfig struct {
	// Type is required: "stdio", "http", or "sse".
	Type string `json:"type"`

	// --- STDIO Server Fields ---
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	EnvFile string            `json:"envFile,omitempty"`

	// --- HTTP/SSE Server Fields ---
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

// NewConfig creates a new VS Code config.
func NewConfig() *Config {
	return &Config{
		Servers: make(map[string]ServerConfig),
	}
}
