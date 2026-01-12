// Package core provides canonical types for AI assistant skill definitions.
package core

// Skill represents a canonical skill definition that can be
// converted to tool-specific formats (Claude, Codex).
type Skill struct {
	// Metadata
	Name        string `json:"name"`
	Description string `json:"description"`

	// Content
	Instructions string `json:"instructions"` // The skill instructions/prompt

	// Resources - paths relative to skill directory
	Scripts    []string `json:"scripts,omitempty"`    // Executable scripts
	References []string `json:"references,omitempty"` // Reference documentation
	Assets     []string `json:"assets,omitempty"`     // Templates, config files

	// Invocation
	Triggers []string `json:"triggers,omitempty"` // Keywords that invoke this skill

	// Dependencies
	Dependencies []string `json:"dependencies,omitempty"` // Required CLI tools
}

// NewSkill creates a new Skill with the given name and description.
func NewSkill(name, description string) *Skill {
	return &Skill{
		Name:        name,
		Description: description,
	}
}

// AddScript adds a script path to the skill.
func (s *Skill) AddScript(path string) {
	s.Scripts = append(s.Scripts, path)
}

// AddReference adds a reference document path to the skill.
func (s *Skill) AddReference(path string) {
	s.References = append(s.References, path)
}

// AddAsset adds an asset path to the skill.
func (s *Skill) AddAsset(path string) {
	s.Assets = append(s.Assets, path)
}

// AddTrigger adds a trigger keyword to the skill.
func (s *Skill) AddTrigger(keyword string) {
	s.Triggers = append(s.Triggers, keyword)
}

// AddDependency adds a dependency to the skill.
func (s *Skill) AddDependency(dep string) {
	s.Dependencies = append(s.Dependencies, dep)
}
