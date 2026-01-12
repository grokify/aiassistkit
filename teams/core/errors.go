package core

import "fmt"

// ParseError represents an error during parsing.
type ParseError struct {
	Format string
	Path   string
	Err    error
}

func (e *ParseError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("failed to parse %s file %s: %v", e.Format, e.Path, e.Err)
	}
	return fmt.Sprintf("failed to parse %s: %v", e.Format, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// MarshalError represents an error during marshaling.
type MarshalError struct {
	Format string
	Err    error
}

func (e *MarshalError) Error() string {
	return fmt.Sprintf("failed to marshal %s: %v", e.Format, e.Err)
}

func (e *MarshalError) Unwrap() error {
	return e.Err
}

// ReadError represents an error during file reading.
type ReadError struct {
	Path string
	Err  error
}

func (e *ReadError) Error() string {
	return fmt.Sprintf("failed to read %s: %v", e.Path, e.Err)
}

func (e *ReadError) Unwrap() error {
	return e.Err
}

// WriteError represents an error during file writing.
type WriteError struct {
	Path string
	Err  error
}

func (e *WriteError) Error() string {
	return fmt.Sprintf("failed to write %s: %v", e.Path, e.Err)
}

func (e *WriteError) Unwrap() error {
	return e.Err
}

// AdapterError represents an error with an adapter.
type AdapterError struct {
	Name string
}

func (e *AdapterError) Error() string {
	return fmt.Sprintf("adapter not found: %s", e.Name)
}
