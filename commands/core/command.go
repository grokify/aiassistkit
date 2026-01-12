// Package core provides canonical types for AI assistant command/prompt definitions.
package core

// Command represents a canonical command/prompt definition that can be
// converted to tool-specific formats (Claude, Gemini, Codex).
type Command struct {
	// Metadata
	Name        string `json:"name"`
	Description string `json:"description"`

	// Arguments
	Arguments []Argument `json:"arguments,omitempty"`

	// Content
	Instructions string `json:"instructions"` // The prompt/instructions content

	// Process steps (for documentation)
	Process []string `json:"process,omitempty"`

	// Dependencies required to run this command
	Dependencies []string `json:"dependencies,omitempty"`

	// Examples of usage
	Examples []Example `json:"examples,omitempty"`
}

// Argument represents a command argument.
type Argument struct {
	Name        string `json:"name"`
	Type        string `json:"type"`                  // string, number, boolean
	Required    bool   `json:"required,omitempty"`    // If true, argument is required
	Default     string `json:"default,omitempty"`     // Default value
	Pattern     string `json:"pattern,omitempty"`     // Regex validation pattern
	Hint        string `json:"hint,omitempty"`        // User-facing hint
	Description string `json:"description,omitempty"` // Detailed description
}

// Example represents a usage example.
type Example struct {
	Description string `json:"description,omitempty"`
	Input       string `json:"input"`
	Output      string `json:"output,omitempty"`
}

// NewCommand creates a new Command with the given name and description.
func NewCommand(name, description string) *Command {
	return &Command{
		Name:        name,
		Description: description,
	}
}

// AddArgument adds an argument to the command.
func (c *Command) AddArgument(arg Argument) {
	c.Arguments = append(c.Arguments, arg)
}

// AddRequiredArgument adds a required string argument.
func (c *Command) AddRequiredArgument(name, description, hint string) {
	c.Arguments = append(c.Arguments, Argument{
		Name:        name,
		Type:        "string",
		Required:    true,
		Description: description,
		Hint:        hint,
	})
}

// AddOptionalArgument adds an optional string argument with a default value.
func (c *Command) AddOptionalArgument(name, description, defaultValue string) {
	c.Arguments = append(c.Arguments, Argument{
		Name:        name,
		Type:        "string",
		Required:    false,
		Default:     defaultValue,
		Description: description,
	})
}

// AddProcessStep adds a step to the process list.
func (c *Command) AddProcessStep(step string) {
	c.Process = append(c.Process, step)
}

// AddDependency adds a dependency to the command.
func (c *Command) AddDependency(dep string) {
	c.Dependencies = append(c.Dependencies, dep)
}

// AddExample adds a usage example.
func (c *Command) AddExample(description, input, output string) {
	c.Examples = append(c.Examples, Example{
		Description: description,
		Input:       input,
		Output:      output,
	})
}
