// Package claude provides an adapter for Claude Code / Claude Desktop
// MCP configuration files (.mcp.json).
//
// Claude's format is the de-facto standard adopted by most AI assistants
// including Cursor, Windsurf, and Cline.
//
// File locations (per https://code.claude.com/docs/en/mcp):
//   - Project scope: .mcp.json (in project root, checked into source control)
//   - User/Local scope: ~/.claude.json (mcpServers field, private to user)
//   - Enterprise managed:
//   - macOS: /Library/Application Support/ClaudeCode/managed-mcp.json
//   - Linux/WSL: /etc/claude-code/managed-mcp.json
//   - Windows: C:\Program Files\ClaudeCode\managed-mcp.json
package claude

// Config represents the Claude MCP configuration file format.
// This is the top-level structure for .mcp.json files.
type Config struct {
	// MCPServers maps server names to their configurations.
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

// ServerConfig represents a single MCP server in Claude's format.
type ServerConfig struct {
	// Type specifies the transport type: "stdio", "http", or "sse".
	// If omitted, inferred from Command (stdio) or URL (http).
	Type string `json:"type,omitempty"`

	// --- STDIO Server Fields ---

	// Command is the executable to run for stdio servers.
	Command string `json:"command,omitempty"`

	// Args are command-line arguments for the executable.
	Args []string `json:"args,omitempty"`

	// Env contains environment variables for the server process.
	// Supports variable expansion: ${VAR} and ${VAR:-default}
	Env map[string]string `json:"env,omitempty"`

	// --- HTTP/SSE Server Fields ---

	// URL is the endpoint for remote servers.
	URL string `json:"url,omitempty"`

	// Headers contains HTTP headers for authentication.
	Headers map[string]string `json:"headers,omitempty"`
}

// NewConfig creates a new empty Claude config.
func NewConfig() *Config {
	return &Config{
		MCPServers: make(map[string]ServerConfig),
	}
}

// AddServer adds a server to the configuration.
func (c *Config) AddServer(name string, server ServerConfig) {
	if c.MCPServers == nil {
		c.MCPServers = make(map[string]ServerConfig)
	}
	c.MCPServers[name] = server
}

// RemoveServer removes a server from the configuration.
func (c *Config) RemoveServer(name string) {
	delete(c.MCPServers, name)
}

// GetServer returns a server by name.
func (c *Config) GetServer(name string) (ServerConfig, bool) {
	server, ok := c.MCPServers[name]
	return server, ok
}

// ServerNames returns a list of all server names.
func (c *Config) ServerNames() []string {
	names := make([]string, 0, len(c.MCPServers))
	for name := range c.MCPServers {
		names = append(names, name)
	}
	return names
}

// IsStdio returns true if the server uses stdio transport.
func (s *ServerConfig) IsStdio() bool {
	if s.Type != "" {
		return s.Type == "stdio"
	}
	return s.Command != ""
}

// IsHTTP returns true if the server uses HTTP transport.
func (s *ServerConfig) IsHTTP() bool {
	if s.Type != "" {
		return s.Type == "http"
	}
	return s.URL != "" && s.Type != "sse"
}

// IsSSE returns true if the server uses SSE transport.
func (s *ServerConfig) IsSSE() bool {
	return s.Type == "sse"
}
