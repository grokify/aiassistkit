package core

import "fmt"

// ParseError occurs when parsing tool-specific format fails.
type ParseError struct {
	Format string
	Path   string
	Err    error
}

func (e *ParseError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("parse %s (%s): %v", e.Format, e.Path, e.Err)
	}
	return fmt.Sprintf("parse %s: %v", e.Format, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// MarshalError occurs when marshaling to tool-specific format fails.
type MarshalError struct {
	Format string
	Err    error
}

func (e *MarshalError) Error() string {
	return fmt.Sprintf("marshal %s: %v", e.Format, e.Err)
}

func (e *MarshalError) Unwrap() error {
	return e.Err
}

// ReadError occurs when reading a file fails.
type ReadError struct {
	Path string
	Err  error
}

func (e *ReadError) Error() string {
	return fmt.Sprintf("read %s: %v", e.Path, e.Err)
}

func (e *ReadError) Unwrap() error {
	return e.Err
}

// WriteError occurs when writing a file fails.
type WriteError struct {
	Path string
	Err  error
}

func (e *WriteError) Error() string {
	return fmt.Sprintf("write %s: %v", e.Path, e.Err)
}

func (e *WriteError) Unwrap() error {
	return e.Err
}
