package core

import "testing"

func TestNewCommand(t *testing.T) {
	cmd := NewCommand("release", "Execute release workflow")

	if cmd.Name != "release" {
		t.Errorf("expected Name 'release', got '%s'", cmd.Name)
	}
	if cmd.Description != "Execute release workflow" {
		t.Errorf("expected Description 'Execute release workflow', got '%s'", cmd.Description)
	}
}

func TestCommandAddRequiredArgument(t *testing.T) {
	cmd := NewCommand("test", "test")

	cmd.AddRequiredArgument("version", "Semantic version", "v1.2.3")

	if len(cmd.Arguments) != 1 {
		t.Errorf("expected 1 argument, got %d", len(cmd.Arguments))
	}

	arg := cmd.Arguments[0]
	if arg.Name != "version" {
		t.Errorf("expected Name 'version', got '%s'", arg.Name)
	}
	if !arg.Required {
		t.Error("expected argument to be required")
	}
	if arg.Hint != "v1.2.3" {
		t.Errorf("expected Hint 'v1.2.3', got '%s'", arg.Hint)
	}
}

func TestCommandAddOptionalArgument(t *testing.T) {
	cmd := NewCommand("test", "test")

	cmd.AddOptionalArgument("format", "Output format", "json")

	if len(cmd.Arguments) != 1 {
		t.Errorf("expected 1 argument, got %d", len(cmd.Arguments))
	}

	arg := cmd.Arguments[0]
	if arg.Required {
		t.Error("expected argument to be optional")
	}
	if arg.Default != "json" {
		t.Errorf("expected Default 'json', got '%s'", arg.Default)
	}
}

func TestCommandAddProcessStep(t *testing.T) {
	cmd := NewCommand("test", "test")

	cmd.AddProcessStep("Run validation")
	cmd.AddProcessStep("Generate changelog")
	cmd.AddProcessStep("Create tag")

	if len(cmd.Process) != 3 {
		t.Errorf("expected 3 process steps, got %d", len(cmd.Process))
	}

	if cmd.Process[0] != "Run validation" {
		t.Errorf("expected first step 'Run validation', got '%s'", cmd.Process[0])
	}
}

func TestCommandAddDependency(t *testing.T) {
	cmd := NewCommand("test", "test")

	cmd.AddDependency("git")
	cmd.AddDependency("gh")

	if len(cmd.Dependencies) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(cmd.Dependencies))
	}
}

func TestCommandAddExample(t *testing.T) {
	cmd := NewCommand("test", "test")

	cmd.AddExample("Basic usage", "/release v1.0.0", "Release created")

	if len(cmd.Examples) != 1 {
		t.Errorf("expected 1 example, got %d", len(cmd.Examples))
	}

	ex := cmd.Examples[0]
	if ex.Input != "/release v1.0.0" {
		t.Errorf("expected Input '/release v1.0.0', got '%s'", ex.Input)
	}
}
