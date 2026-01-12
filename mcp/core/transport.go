// Package core provides the canonical types for MCP server configuration
// that can be converted to/from various AI assistant formats.
package core

// TransportType represents the communication protocol for an MCP server.
type TransportType string

const (
	// TransportStdio represents local servers using standard input/output streams.
	TransportStdio TransportType = "stdio"

	// TransportHTTP represents remote servers using HTTP/Streamable HTTP transport.
	TransportHTTP TransportType = "http"

	// TransportSSE represents remote servers using Server-Sent Events (legacy).
	TransportSSE TransportType = "sse"
)

// String returns the string representation of the transport type.
func (t TransportType) String() string {
	return string(t)
}

// IsLocal returns true if the transport type runs locally (stdio).
func (t TransportType) IsLocal() bool {
	return t == TransportStdio
}

// IsRemote returns true if the transport type is network-based.
func (t TransportType) IsRemote() bool {
	return t == TransportHTTP || t == TransportSSE
}

// Valid returns true if the transport type is a known valid type.
func (t TransportType) Valid() bool {
	switch t {
	case TransportStdio, TransportHTTP, TransportSSE:
		return true
	default:
		return false
	}
}
