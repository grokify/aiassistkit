package codex

// Config represents the Codex configuration file.
// This is a partial representation focusing on MCP servers.
type Config struct {
	MCPServers map[string]ServerConfig `toml:"mcp_servers"`
}

// ServerConfig represents a Codex MCP server configuration.
type ServerConfig struct {
	// --- STDIO Server Fields ---
	Command string            `toml:"command,omitempty"`
	Args    []string          `toml:"args,omitempty"`
	Env     map[string]string `toml:"env,omitempty"`
	EnvVars []string          `toml:"env_vars,omitempty"` // Additional env vars to whitelist
	Cwd     string            `toml:"cwd,omitempty"`

	// --- HTTP Server Fields ---
	URL               string            `toml:"url,omitempty"`
	BearerTokenEnvVar string            `toml:"bearer_token_env_var,omitempty"`
	HTTPHeaders       map[string]string `toml:"http_headers,omitempty"`
	EnvHTTPHeaders    map[string]string `toml:"env_http_headers,omitempty"` // Header name -> env var name

	// --- Tool Control ---
	EnabledTools  []string `toml:"enabled_tools,omitempty"`
	DisabledTools []string `toml:"disabled_tools,omitempty"`

	// --- Timeouts ---
	StartupTimeoutSec int `toml:"startup_timeout_sec,omitempty"`
	ToolTimeoutSec    int `toml:"tool_timeout_sec,omitempty"`

	// --- State ---
	Enabled *bool `toml:"enabled,omitempty"`
}

// NewConfig creates a new Codex config.
func NewConfig() *Config {
	return &Config{
		MCPServers: make(map[string]ServerConfig),
	}
}
