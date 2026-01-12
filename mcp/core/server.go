package core

// Server represents a canonical MCP server configuration that can be
// converted to/from various AI assistant formats (Claude, Cursor, VS Code, etc.).
type Server struct {
	// Transport specifies the communication protocol (stdio, http, sse).
	// If empty, it will be inferred from Command (stdio) or URL (http).
	Transport TransportType `json:"transport,omitempty"`

	// --- STDIO Server Fields ---

	// Command is the executable to launch for stdio servers.
	Command string `json:"command,omitempty"`

	// Args are the command-line arguments passed to the command.
	Args []string `json:"args,omitempty"`

	// Env contains environment variables for the server process.
	Env map[string]string `json:"env,omitempty"`

	// EnvFile is the path to an environment file to load (VS Code feature).
	EnvFile string `json:"envFile,omitempty"`

	// Cwd is the working directory for the server process (Codex feature).
	Cwd string `json:"cwd,omitempty"`

	// --- HTTP/SSE Server Fields ---

	// URL is the endpoint for HTTP/SSE servers.
	URL string `json:"url,omitempty"`

	// Headers contains HTTP headers for authentication or configuration.
	Headers map[string]string `json:"headers,omitempty"`

	// BearerTokenEnvVar is the name of an env var containing a bearer token (Codex feature).
	BearerTokenEnvVar string `json:"bearerTokenEnvVar,omitempty"`

	// --- Tool Control Fields ---

	// EnabledTools is an allow-list of tools to expose (Codex/Cline feature).
	EnabledTools []string `json:"enabledTools,omitempty"`

	// DisabledTools is a deny-list of tools to hide (Codex/Cline feature).
	DisabledTools []string `json:"disabledTools,omitempty"`

	// AlwaysAllow lists tools that don't require user approval (Cline feature).
	AlwaysAllow []string `json:"alwaysAllow,omitempty"`

	// --- Server State ---

	// Enabled indicates whether the server is active. Defaults to true.
	Enabled *bool `json:"enabled,omitempty"`

	// --- Timeout Configuration ---

	// StartupTimeoutSec is the timeout for server startup in seconds (Codex feature).
	StartupTimeoutSec int `json:"startupTimeoutSec,omitempty"`

	// ToolTimeoutSec is the timeout for tool execution in seconds (Codex feature).
	ToolTimeoutSec int `json:"toolTimeoutSec,omitempty"`

	// NetworkTimeoutSec is the timeout for network operations (Cline feature).
	NetworkTimeoutSec int `json:"networkTimeoutSec,omitempty"`
}

// IsStdio returns true if the server is configured for stdio transport.
func (s *Server) IsStdio() bool {
	if s.Transport != "" {
		return s.Transport == TransportStdio
	}
	return s.Command != ""
}

// IsHTTP returns true if the server is configured for HTTP transport.
func (s *Server) IsHTTP() bool {
	if s.Transport != "" {
		return s.Transport == TransportHTTP
	}
	return s.URL != "" && s.Transport != TransportSSE
}

// IsSSE returns true if the server is configured for SSE transport.
func (s *Server) IsSSE() bool {
	return s.Transport == TransportSSE
}

// IsRemote returns true if the server is accessed over the network.
func (s *Server) IsRemote() bool {
	return s.IsHTTP() || s.IsSSE()
}

// IsEnabled returns whether the server is enabled. Defaults to true if not set.
func (s *Server) IsEnabled() bool {
	if s.Enabled == nil {
		return true
	}
	return *s.Enabled
}

// SetEnabled sets the enabled state of the server.
func (s *Server) SetEnabled(enabled bool) {
	s.Enabled = &enabled
}

// InferTransport returns the inferred transport type based on configuration.
func (s *Server) InferTransport() TransportType {
	if s.Transport != "" {
		return s.Transport
	}
	if s.Command != "" {
		return TransportStdio
	}
	if s.URL != "" {
		return TransportHTTP
	}
	return ""
}

// Validate checks if the server configuration is valid.
func (s *Server) Validate() error {
	if s.Command == "" && s.URL == "" {
		return ErrNoCommandOrURL
	}
	if s.Command != "" && s.URL != "" {
		return ErrBothCommandAndURL
	}
	return nil
}
