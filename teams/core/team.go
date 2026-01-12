package core

// Team represents a multi-agent orchestration definition.
// A team coordinates multiple agents to accomplish a complex workflow.
type Team struct {
	// Name is the team identifier (e.g., "release-team").
	Name string `json:"name" yaml:"name"`

	// Description explains what this team accomplishes.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Process defines how tasks are executed (sequential, parallel, hierarchical).
	Process Process `json:"process" yaml:"process"`

	// Manager is the orchestrating agent (required for hierarchical process).
	Manager string `json:"manager,omitempty" yaml:"manager,omitempty"`

	// Agents lists all agent names participating in this team.
	Agents []string `json:"agents,omitempty" yaml:"agents,omitempty"`

	// Tasks defines the work to be done, with agent assignments and subtasks.
	Tasks []Task `json:"tasks" yaml:"tasks"`

	// Version is the target version for release workflows.
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

// NewTeam creates a new Team with the given name and process type.
func NewTeam(name string, process Process) *Team {
	return &Team{
		Name:    name,
		Process: process,
	}
}

// WithDescription sets the description and returns the team for chaining.
func (t *Team) WithDescription(desc string) *Team {
	t.Description = desc
	return t
}

// WithManager sets the manager agent (for hierarchical process).
func (t *Team) WithManager(manager string) *Team {
	t.Manager = manager
	return t
}

// AddAgent adds an agent to the team.
func (t *Team) AddAgent(agent string) *Team {
	t.Agents = append(t.Agents, agent)
	return t
}

// AddAgents adds multiple agents to the team.
func (t *Team) AddAgents(agents ...string) *Team {
	t.Agents = append(t.Agents, agents...)
	return t
}

// AddTask adds a task to the team.
func (t *Team) AddTask(task Task) *Team {
	t.Tasks = append(t.Tasks, task)
	return t
}

// GetTask returns a task by name, or nil if not found.
func (t *Team) GetTask(name string) *Task {
	for i := range t.Tasks {
		if t.Tasks[i].Name == name {
			return &t.Tasks[i]
		}
	}
	return nil
}

// TaskNames returns the names of all tasks.
func (t *Team) TaskNames() []string {
	names := make([]string, len(t.Tasks))
	for i, task := range t.Tasks {
		names[i] = task.Name
	}
	return names
}

// AgentTasks returns all tasks assigned to a specific agent.
func (t *Team) AgentTasks(agent string) []Task {
	var tasks []Task
	for _, task := range t.Tasks {
		if task.Agent == agent {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// TotalSubtaskCount returns the total number of subtasks across all tasks.
func (t *Team) TotalSubtaskCount() int {
	count := 0
	for _, task := range t.Tasks {
		count += len(task.Subtasks)
	}
	return count
}

// RequiredSubtaskCount returns the total number of required subtasks.
func (t *Team) RequiredSubtaskCount() int {
	count := 0
	for _, task := range t.Tasks {
		count += task.RequiredSubtaskCount()
	}
	return count
}

// Validate checks if the team definition is valid.
func (t *Team) Validate() error {
	if t.Name == "" {
		return &ValidationError{Field: "name", Message: "team name is required"}
	}

	if !t.Process.IsValid() {
		return &ValidationError{Field: "process", Message: "invalid process type"}
	}

	if t.Process == ProcessHierarchical && t.Manager == "" {
		return &ValidationError{Field: "manager", Message: "manager is required for hierarchical process"}
	}

	if len(t.Tasks) == 0 {
		return &ValidationError{Field: "tasks", Message: "at least one task is required"}
	}

	// Validate task dependencies
	taskNames := make(map[string]bool)
	for _, task := range t.Tasks {
		taskNames[task.Name] = true
	}

	for _, task := range t.Tasks {
		for _, dep := range task.DependsOn {
			if !taskNames[dep] {
				return &ValidationError{
					Field:   "tasks." + task.Name + ".depends_on",
					Message: "unknown dependency: " + dep,
				}
			}
		}
	}

	return nil
}

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// TopologicalSort returns tasks in dependency order.
// Tasks with no dependencies come first, followed by tasks whose dependencies are satisfied.
func (t *Team) TopologicalSort() ([]Task, error) {
	// Build dependency graph
	inDegree := make(map[string]int)
	dependents := make(map[string][]string)

	for _, task := range t.Tasks {
		if _, exists := inDegree[task.Name]; !exists {
			inDegree[task.Name] = 0
		}
		for _, dep := range task.DependsOn {
			inDegree[task.Name]++
			dependents[dep] = append(dependents[dep], task.Name)
		}
	}

	// Find tasks with no dependencies
	var queue []string
	for _, task := range t.Tasks {
		if inDegree[task.Name] == 0 {
			queue = append(queue, task.Name)
		}
	}

	// Process queue
	var sorted []Task
	taskMap := make(map[string]Task)
	for _, task := range t.Tasks {
		taskMap[task.Name] = task
	}

	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		sorted = append(sorted, taskMap[name])

		for _, dependent := range dependents[name] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	if len(sorted) != len(t.Tasks) {
		return nil, &ValidationError{Field: "tasks", Message: "circular dependency detected"}
	}

	return sorted, nil
}

// ParallelGroups returns tasks grouped by execution wave.
// Tasks in the same group can run in parallel.
func (t *Team) ParallelGroups() ([][]Task, error) {
	sorted, err := t.TopologicalSort()
	if err != nil {
		return nil, err
	}

	// Calculate the "level" of each task (longest path from a root)
	levels := make(map[string]int)
	for _, task := range sorted {
		maxDepLevel := -1
		for _, dep := range task.DependsOn {
			if levels[dep] > maxDepLevel {
				maxDepLevel = levels[dep]
			}
		}
		levels[task.Name] = maxDepLevel + 1
	}

	// Group tasks by level
	maxLevel := 0
	for _, level := range levels {
		if level > maxLevel {
			maxLevel = level
		}
	}

	groups := make([][]Task, maxLevel+1)
	taskMap := make(map[string]Task)
	for _, task := range t.Tasks {
		taskMap[task.Name] = task
	}

	for name, level := range levels {
		groups[level] = append(groups[level], taskMap[name])
	}

	return groups, nil
}
