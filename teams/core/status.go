package core

// Status represents the result of a subtask or task.
type Status string

const (
	// StatusGo indicates the check passed.
	StatusGo Status = "GO"

	// StatusNoGo indicates the check failed (blocking).
	StatusNoGo Status = "NO-GO"

	// StatusWarn indicates the check failed but is non-blocking.
	StatusWarn Status = "WARN"

	// StatusSkip indicates the check was skipped.
	StatusSkip Status = "SKIP"

	// StatusPending indicates the check has not run yet.
	StatusPending Status = "PENDING"

	// StatusRunning indicates the check is currently executing.
	StatusRunning Status = "RUNNING"
)

// String returns the string representation of the status.
func (s Status) String() string {
	return string(s)
}

// Emoji returns the emoji representation of the status.
func (s Status) Emoji() string {
	switch s {
	case StatusGo:
		return "🟢"
	case StatusNoGo:
		return "🔴"
	case StatusWarn:
		return "🟡"
	case StatusSkip:
		return "⚪"
	case StatusPending:
		return "⏳"
	case StatusRunning:
		return "🔄"
	default:
		return "❓"
	}
}

// IsBlocking returns true if this status should block the workflow.
func (s Status) IsBlocking() bool {
	return s == StatusNoGo
}

// IsPassing returns true if this status indicates success.
func (s Status) IsPassing() bool {
	return s == StatusGo || s == StatusWarn || s == StatusSkip
}

// SubtaskResult holds the result of a subtask execution.
type SubtaskResult struct {
	Name    string `json:"name"`
	Status  Status `json:"status"`
	Message string `json:"message,omitempty"`
	Output  string `json:"output,omitempty"`
}

// TaskResult holds the result of a task execution.
type TaskResult struct {
	Name     string          `json:"name"`
	Agent    string          `json:"agent"`
	Status   Status          `json:"status"`
	Subtasks []SubtaskResult `json:"subtasks,omitempty"`
}

// TeamResult holds the result of a team execution.
type TeamResult struct {
	Name    string       `json:"name"`
	Status  Status       `json:"status"`
	Tasks   []TaskResult `json:"tasks"`
	Version string       `json:"version,omitempty"` // Target version if applicable
}

// ComputeTaskStatus computes the overall status from subtask results.
func ComputeTaskStatus(subtasks []SubtaskResult) Status {
	hasNoGo := false
	hasWarn := false
	allSkip := true

	for _, st := range subtasks {
		if st.Status != StatusSkip {
			allSkip = false
		}
		if st.Status == StatusNoGo {
			hasNoGo = true
		}
		if st.Status == StatusWarn {
			hasWarn = true
		}
	}

	if hasNoGo {
		return StatusNoGo
	}
	if allSkip {
		return StatusSkip
	}
	if hasWarn {
		return StatusWarn
	}
	return StatusGo
}

// ComputeTeamStatus computes the overall status from task results.
func ComputeTeamStatus(tasks []TaskResult) Status {
	hasNoGo := false
	hasWarn := false
	allSkip := true

	for _, t := range tasks {
		if t.Status != StatusSkip {
			allSkip = false
		}
		if t.Status == StatusNoGo {
			hasNoGo = true
		}
		if t.Status == StatusWarn {
			hasWarn = true
		}
	}

	if hasNoGo {
		return StatusNoGo
	}
	if allSkip {
		return StatusSkip
	}
	if hasWarn {
		return StatusWarn
	}
	return StatusGo
}
