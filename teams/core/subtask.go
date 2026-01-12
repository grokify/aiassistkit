package core

// Subtask represents an individual check or action within a task.
type Subtask struct {
	// Name is the subtask identifier (e.g., "build", "tests", "lint").
	Name string `json:"name" yaml:"name"`

	// Description explains what this subtask validates or does.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Command is the CLI command to execute (e.g., "go build ./...").
	Command string `json:"command,omitempty" yaml:"command,omitempty"`

	// Pattern is a regex pattern to search for (presence indicates failure).
	Pattern string `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// File is a specific file that must exist.
	File string `json:"file,omitempty" yaml:"file,omitempty"`

	// Files is a glob pattern for files to check (used with Pattern).
	Files string `json:"files,omitempty" yaml:"files,omitempty"`

	// Required indicates if failure blocks the workflow (NO-GO vs WARN).
	Required bool `json:"required" yaml:"required"`

	// ExpectedOutput describes what successful execution looks like.
	ExpectedOutput string `json:"expected_output,omitempty" yaml:"expected_output,omitempty"`

	// Timeout in seconds for command execution.
	Timeout int `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

// NewSubtask creates a new Subtask with the given name.
func NewSubtask(name string) *Subtask {
	return &Subtask{
		Name:     name,
		Required: true, // Default to required
	}
}

// WithCommand sets the command and returns the subtask for chaining.
func (s *Subtask) WithCommand(cmd string) *Subtask {
	s.Command = cmd
	return s
}

// WithPattern sets the pattern and returns the subtask for chaining.
func (s *Subtask) WithPattern(pattern string) *Subtask {
	s.Pattern = pattern
	return s
}

// WithFile sets the file and returns the subtask for chaining.
func (s *Subtask) WithFile(file string) *Subtask {
	s.File = file
	return s
}

// WithFiles sets the files glob and returns the subtask for chaining.
func (s *Subtask) WithFiles(files string) *Subtask {
	s.Files = files
	return s
}

// Optional marks the subtask as optional (WARN instead of NO-GO on failure).
func (s *Subtask) Optional() *Subtask {
	s.Required = false
	return s
}

// IsCommandBased returns true if this subtask executes a command.
func (s *Subtask) IsCommandBased() bool {
	return s.Command != ""
}

// IsPatternBased returns true if this subtask searches for a pattern.
func (s *Subtask) IsPatternBased() bool {
	return s.Pattern != ""
}

// IsFileBased returns true if this subtask checks for file existence.
func (s *Subtask) IsFileBased() bool {
	return s.File != ""
}

// Type returns the subtask type based on its configuration.
func (s *Subtask) Type() string {
	if s.Command != "" {
		return "command"
	}
	if s.Pattern != "" {
		return "pattern"
	}
	if s.File != "" {
		return "file"
	}
	return "unknown"
}
