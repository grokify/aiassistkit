package core

import (
	"errors"
	"fmt"
)

// Common errors for MCP configuration.
var (
	// ErrNoCommandOrURL is returned when a server has neither command nor URL.
	ErrNoCommandOrURL = errors.New("server must have either command (stdio) or url (http/sse)")

	// ErrBothCommandAndURL is returned when a server has both command and URL.
	ErrBothCommandAndURL = errors.New("server cannot have both command and url")

	// ErrInvalidTransport is returned when a transport type is invalid.
	ErrInvalidTransport = errors.New("invalid transport type")

	// ErrServerNotFound is returned when a server is not found.
	ErrServerNotFound = errors.New("server not found")

	// ErrEmptyConfig is returned when configuration is empty.
	ErrEmptyConfig = errors.New("configuration is empty")
)

// ServerValidationError wraps a validation error with the server name.
type ServerValidationError struct {
	Name string
	Err  error
}

func (e *ServerValidationError) Error() string {
	return fmt.Sprintf("server %q: %v", e.Name, e.Err)
}

func (e *ServerValidationError) Unwrap() error {
	return e.Err
}

// ParseError represents an error parsing a configuration file.
type ParseError struct {
	Format string
	Path   string
	Err    error
}

func (e *ParseError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("failed to parse %s config from %s: %v", e.Format, e.Path, e.Err)
	}
	return fmt.Sprintf("failed to parse %s config: %v", e.Format, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// WriteError represents an error writing a configuration file.
type WriteError struct {
	Format string
	Path   string
	Err    error
}

func (e *WriteError) Error() string {
	return fmt.Sprintf("failed to write %s config to %s: %v", e.Format, e.Path, e.Err)
}

func (e *WriteError) Unwrap() error {
	return e.Err
}
