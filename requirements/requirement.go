// Package requirements provides dependency checking and HITL installation
// prompts for assistant runtimes. It maps tool names from multi-agent-spec
// `requires` fields to install commands and interactive prompts.
package requirements

// Requirement defines an external tool/binary that an agent may require.
type Requirement struct {
	// Name is the canonical tool name (matches multi-agent-spec requires field).
	Name string `json:"name" yaml:"name"`

	// Purpose describes what this tool is used for.
	Purpose string `json:"purpose" yaml:"purpose"`

	// Check is the command to verify the tool is installed (e.g., "releasekit --version").
	Check string `json:"check" yaml:"check"`

	// InstallMethods lists ways to install this tool, in priority order.
	// The first available method will be suggested to the user.
	InstallMethods []InstallMethod `json:"install_methods" yaml:"install_methods"`

	// Homepage is the tool's documentation URL.
	Homepage string `json:"homepage,omitempty" yaml:"homepage,omitempty"`
}

// InstallMethod defines one way to install a tool.
type InstallMethod struct {
	// Name identifies the method (e.g., "go", "brew", "helm", "apt", "npm").
	Name string `json:"name" yaml:"name"`

	// Command is the install command to run.
	Command string `json:"command" yaml:"command"`

	// Requires lists tools that must be present for this method to work.
	// For example, "go install" requires "go".
	Requires []string `json:"requires,omitempty" yaml:"requires,omitempty"`

	// Platforms limits this method to specific OS (empty = all platforms).
	// Values: "darwin", "linux", "windows"
	Platforms []string `json:"platforms,omitempty" yaml:"platforms,omitempty"`
}

// MissingRequirement represents a tool that is not installed.
type MissingRequirement struct {
	Requirement       Requirement
	AvailableMethods  []InstallMethod // Methods that can be used (prerequisites met)
	SuggestedMethod   *InstallMethod  // First available method (recommended)
}

// CheckResult contains the results of checking requirements.
type CheckResult struct {
	Satisfied []string              // Tools that are installed
	Missing   []MissingRequirement  // Tools that need installation
	Unknown   []string              // Tools not in registry
}

// AllSatisfied returns true if all requirements are met.
func (r CheckResult) AllSatisfied() bool {
	return len(r.Missing) == 0 && len(r.Unknown) == 0
}
