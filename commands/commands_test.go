package commands

import (
	"strings"
	"testing"
)

func TestAdapterRegistry(t *testing.T) {
	names := AdapterNames()

	// Should have Claude, Gemini, and Codex adapters
	if len(names) < 3 {
		t.Errorf("expected at least 3 adapters, got %d", len(names))
	}

	// Check all adapters exist
	for _, name := range []string{"claude", "gemini", "codex"} {
		adapter, ok := GetAdapter(name)
		if !ok {
			t.Errorf("expected %s adapter to be registered", name)
			continue
		}
		if adapter.Name() != name {
			t.Errorf("expected adapter name '%s', got '%s'", name, adapter.Name())
		}
	}
}

func TestClaudeAdapter(t *testing.T) {
	adapter, ok := GetAdapter("claude")
	if !ok {
		t.Fatal("Claude adapter not found")
	}

	// Test marshal
	cmd := NewCommand("release", "Execute full release workflow")
	cmd.AddRequiredArgument("version", "Semantic version", "v1.2.3")
	cmd.AddProcessStep("Run validation checks")
	cmd.AddProcessStep("Generate changelog")
	cmd.Instructions = "Execute a full release workflow with validation."

	data, err := adapter.Marshal(cmd)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	content := string(data)

	// Check frontmatter
	if !strings.HasPrefix(content, "---") {
		t.Error("expected Markdown to start with frontmatter")
	}
	if !strings.Contains(content, "description: Execute full release workflow") {
		t.Error("expected description in frontmatter")
	}

	// Check title
	if !strings.Contains(content, "# Release") {
		t.Error("expected title in content")
	}

	// Check arguments section
	if !strings.Contains(content, "## Arguments") {
		t.Error("expected Arguments section")
	}
	if !strings.Contains(content, "**version** (required)") {
		t.Error("expected version argument")
	}

	// Check process section
	if !strings.Contains(content, "## Process") {
		t.Error("expected Process section")
	}

	// Test round-trip
	parsed, err := adapter.Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if parsed.Description != cmd.Description {
		t.Errorf("round-trip: expected Description '%s', got '%s'", cmd.Description, parsed.Description)
	}
}

func TestGeminiAdapter(t *testing.T) {
	adapter, ok := GetAdapter("gemini")
	if !ok {
		t.Fatal("Gemini adapter not found")
	}

	// Test marshal
	cmd := NewCommand("release", "Execute full release workflow")
	cmd.AddRequiredArgument("version", "Semantic version", "v1.2.3")
	cmd.AddProcessStep("Run validation")
	cmd.Instructions = "Execute release workflow"

	data, err := adapter.Marshal(cmd)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	content := string(data)

	// Check TOML structure
	if !strings.Contains(content, "[command]") {
		t.Error("expected [command] section")
	}
	// TOML library may use single or double quotes
	if !strings.Contains(content, "name = 'release'") && !strings.Contains(content, "name = \"release\"") {
		t.Errorf("expected name in command section, got:\n%s", content)
	}
	if !strings.Contains(content, "[[arguments]]") {
		t.Error("expected [[arguments]] section")
	}
	if !strings.Contains(content, "[content]") {
		t.Error("expected [content] section")
	}

	// Test round-trip
	parsed, err := adapter.Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if parsed.Name != cmd.Name {
		t.Errorf("round-trip: expected Name '%s', got '%s'", cmd.Name, parsed.Name)
	}
	if len(parsed.Arguments) != 1 {
		t.Errorf("round-trip: expected 1 argument, got %d", len(parsed.Arguments))
	}
}

func TestCodexAdapter(t *testing.T) {
	adapter, ok := GetAdapter("codex")
	if !ok {
		t.Fatal("Codex adapter not found")
	}

	// Test marshal
	cmd := NewCommand("release", "Execute full release workflow")
	cmd.AddRequiredArgument("version", "Semantic version", "v1.2.3")
	cmd.AddProcessStep("Run validation")
	cmd.Instructions = "Execute release workflow"

	data, err := adapter.Marshal(cmd)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	content := string(data)

	// Check frontmatter
	if !strings.HasPrefix(content, "---") {
		t.Error("expected Markdown to start with frontmatter")
	}
	if !strings.Contains(content, "description: Execute full release workflow") {
		t.Error("expected description in frontmatter")
	}
	if !strings.Contains(content, "argument-hint:") {
		t.Error("expected argument-hint in frontmatter")
	}

	// Check arguments use $UPPERCASE format
	if !strings.Contains(content, "$VERSION") {
		t.Error("expected $VERSION in arguments section")
	}

	// Test round-trip
	parsed, err := adapter.Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if parsed.Description != cmd.Description {
		t.Errorf("round-trip: expected Description '%s', got '%s'", cmd.Description, parsed.Description)
	}
}

func TestConvert(t *testing.T) {
	// Create a Claude command
	claudeMD := `---
description: Test command
---

# Test

Test command description
`

	// Convert to Gemini format
	geminiData, err := Convert([]byte(claudeMD), "claude", "gemini")
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	content := string(geminiData)
	if !strings.Contains(content, "[command]") {
		t.Error("expected TOML [command] section in converted output")
	}
	// TOML library may use single or double quotes
	if !strings.Contains(content, "description = 'Test command'") && !strings.Contains(content, "description = \"Test command\"") {
		t.Errorf("expected description in converted output, got:\n%s", content)
	}
}
