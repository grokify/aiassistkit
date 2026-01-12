package teams

import (
	"testing"
)

func TestNewTeam(t *testing.T) {
	team := NewTeam("release-team", ProcessHierarchical)
	team.WithDescription("Release validation workflow")
	team.WithManager("release-coordinator")

	if team.Name != "release-team" {
		t.Errorf("expected name 'release-team', got '%s'", team.Name)
	}
	if team.Process != ProcessHierarchical {
		t.Errorf("expected process 'hierarchical', got '%s'", team.Process)
	}
	if team.Manager != "release-coordinator" {
		t.Errorf("expected manager 'release-coordinator', got '%s'", team.Manager)
	}
}

func TestTaskWithSubtasks(t *testing.T) {
	task := NewTask("qa-validation", "qa")
	task.WithDescription("Quality Assurance validation")

	task.AddSubtask(Subtask{
		Name:     "build",
		Command:  "go build ./...",
		Required: true,
	})
	task.AddSubtask(Subtask{
		Name:     "tests",
		Command:  "go test -v ./...",
		Required: true,
	})
	task.AddSubtask(Subtask{
		Name:     "lint",
		Command:  "golangci-lint run",
		Required: false,
	})

	if len(task.Subtasks) != 3 {
		t.Errorf("expected 3 subtasks, got %d", len(task.Subtasks))
	}
	if task.RequiredSubtaskCount() != 2 {
		t.Errorf("expected 2 required subtasks, got %d", task.RequiredSubtaskCount())
	}
}

func TestTeamValidation(t *testing.T) {
	// Valid team
	team := NewTeam("release-team", ProcessHierarchical)
	team.WithManager("coordinator")
	team.AddTask(Task{Name: "task1", Agent: "agent1"})

	if err := team.Validate(); err != nil {
		t.Errorf("expected valid team, got error: %v", err)
	}

	// Invalid: hierarchical without manager
	team2 := NewTeam("release-team", ProcessHierarchical)
	team2.AddTask(Task{Name: "task1", Agent: "agent1"})

	if err := team2.Validate(); err == nil {
		t.Error("expected error for hierarchical team without manager")
	}

	// Invalid: no tasks
	team3 := NewTeam("release-team", ProcessSequential)

	if err := team3.Validate(); err == nil {
		t.Error("expected error for team without tasks")
	}
}

func TestTopologicalSort(t *testing.T) {
	team := NewTeam("test-team", ProcessParallel)
	team.AddTask(Task{Name: "task-a", Agent: "agent1"})
	team.AddTask(Task{Name: "task-b", Agent: "agent2", DependsOn: []string{"task-a"}})
	team.AddTask(Task{Name: "task-c", Agent: "agent3", DependsOn: []string{"task-a", "task-b"}})

	sorted, err := team.TopologicalSort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// task-a should come first
	if sorted[0].Name != "task-a" {
		t.Errorf("expected task-a first, got %s", sorted[0].Name)
	}
	// task-b should come second
	if sorted[1].Name != "task-b" {
		t.Errorf("expected task-b second, got %s", sorted[1].Name)
	}
	// task-c should come last
	if sorted[2].Name != "task-c" {
		t.Errorf("expected task-c last, got %s", sorted[2].Name)
	}
}

func TestParallelGroups(t *testing.T) {
	team := NewTeam("test-team", ProcessParallel)
	team.AddTask(Task{Name: "task-a", Agent: "agent1"})
	team.AddTask(Task{Name: "task-b", Agent: "agent2"})
	team.AddTask(Task{Name: "task-c", Agent: "agent3", DependsOn: []string{"task-a", "task-b"}})

	groups, err := team.ParallelGroups()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}

	// First group should have task-a and task-b (can run in parallel)
	if len(groups[0]) != 2 {
		t.Errorf("expected 2 tasks in first group, got %d", len(groups[0]))
	}

	// Second group should have task-c
	if len(groups[1]) != 1 {
		t.Errorf("expected 1 task in second group, got %d", len(groups[1]))
	}
}

func TestStatusComputation(t *testing.T) {
	// All GO
	results := []SubtaskResult{
		{Name: "a", Status: StatusGo},
		{Name: "b", Status: StatusGo},
	}
	if ComputeTaskStatus(results) != StatusGo {
		t.Error("expected GO status for all passing")
	}

	// One NO-GO
	results = []SubtaskResult{
		{Name: "a", Status: StatusGo},
		{Name: "b", Status: StatusNoGo},
	}
	if ComputeTaskStatus(results) != StatusNoGo {
		t.Error("expected NO-GO status when one fails")
	}

	// One WARN (no NO-GO)
	results = []SubtaskResult{
		{Name: "a", Status: StatusGo},
		{Name: "b", Status: StatusWarn},
	}
	if ComputeTaskStatus(results) != StatusWarn {
		t.Error("expected WARN status when one warns")
	}

	// All SKIP
	results = []SubtaskResult{
		{Name: "a", Status: StatusSkip},
		{Name: "b", Status: StatusSkip},
	}
	if ComputeTaskStatus(results) != StatusSkip {
		t.Error("expected SKIP status when all skipped")
	}
}

func TestSubtaskTypes(t *testing.T) {
	st := NewSubtask("build").WithCommand("go build ./...")
	if st.Type() != "command" {
		t.Errorf("expected type 'command', got '%s'", st.Type())
	}

	st = NewSubtask("secrets").WithPattern("password=").WithFiles("**/*.go")
	if st.Type() != "pattern" {
		t.Errorf("expected type 'pattern', got '%s'", st.Type())
	}

	st = NewSubtask("readme").WithFile("README.md")
	if st.Type() != "file" {
		t.Errorf("expected type 'file', got '%s'", st.Type())
	}
}
