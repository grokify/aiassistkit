package core

// Task represents a unit of work assigned to an agent within a team.
type Task struct {
	// Name is the task identifier (e.g., "qa-validation", "docs-validation").
	Name string `json:"name" yaml:"name"`

	// Description explains what this task accomplishes.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Agent is the name of the agent assigned to this task.
	Agent string `json:"agent" yaml:"agent"`

	// DependsOn lists task names that must complete before this task.
	DependsOn []string `json:"depends_on,omitempty" yaml:"depends_on,omitempty"`

	// Subtasks are the individual checks or actions within this task.
	Subtasks []Subtask `json:"subtasks,omitempty" yaml:"subtasks,omitempty"`

	// Inputs are data or context passed to the task.
	Inputs []string `json:"inputs,omitempty" yaml:"inputs,omitempty"`

	// Outputs are data or artifacts produced by the task.
	Outputs []string `json:"outputs,omitempty" yaml:"outputs,omitempty"`
}

// NewTask creates a new Task with the given name and agent.
func NewTask(name, agent string) *Task {
	return &Task{
		Name:  name,
		Agent: agent,
	}
}

// WithDescription sets the description and returns the task for chaining.
func (t *Task) WithDescription(desc string) *Task {
	t.Description = desc
	return t
}

// AddDependency adds a task dependency.
func (t *Task) AddDependency(taskName string) *Task {
	t.DependsOn = append(t.DependsOn, taskName)
	return t
}

// AddSubtask adds a subtask to the task.
func (t *Task) AddSubtask(subtask Subtask) *Task {
	t.Subtasks = append(t.Subtasks, subtask)
	return t
}

// AddSubtasks adds multiple subtasks to the task.
func (t *Task) AddSubtasks(subtasks ...Subtask) *Task {
	t.Subtasks = append(t.Subtasks, subtasks...)
	return t
}

// HasDependencies returns true if this task depends on other tasks.
func (t *Task) HasDependencies() bool {
	return len(t.DependsOn) > 0
}

// HasSubtasks returns true if this task has subtasks.
func (t *Task) HasSubtasks() bool {
	return len(t.Subtasks) > 0
}

// RequiredSubtaskCount returns the number of required subtasks.
func (t *Task) RequiredSubtaskCount() int {
	count := 0
	for _, st := range t.Subtasks {
		if st.Required {
			count++
		}
	}
	return count
}

// SubtaskNames returns the names of all subtasks.
func (t *Task) SubtaskNames() []string {
	names := make([]string, len(t.Subtasks))
	for i, st := range t.Subtasks {
		names[i] = st.Name
	}
	return names
}
