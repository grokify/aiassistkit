// Package core provides canonical types for multi-agent team orchestration.
package core

// Process defines how tasks are executed within a team.
type Process string

const (
	// ProcessSequential executes tasks one after another in order.
	ProcessSequential Process = "sequential"

	// ProcessParallel executes independent tasks concurrently.
	ProcessParallel Process = "parallel"

	// ProcessHierarchical uses a manager agent to delegate to specialists.
	ProcessHierarchical Process = "hierarchical"
)

// String returns the string representation of the process.
func (p Process) String() string {
	return string(p)
}

// IsValid checks if the process type is valid.
func (p Process) IsValid() bool {
	switch p {
	case ProcessSequential, ProcessParallel, ProcessHierarchical:
		return true
	default:
		return false
	}
}
