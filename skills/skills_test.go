package skills

import (
	"strings"
	"testing"
)

func TestAdapterRegistry(t *testing.T) {
	names := AdapterNames()

	// Should have Claude and Codex adapters
	if len(names) < 2 {
		t.Errorf("expected at least 2 adapters, got %d", len(names))
	}

	// Check all adapters exist
	for _, name := range []string{"claude", "codex"} {
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
	skill := NewSkill("version-analysis", "Analyze git history for semantic versioning")
	skill.Instructions = "Analyze commits since the last tag and suggest version."
	skill.AddTrigger("version")
	skill.AddTrigger("semver")
	skill.AddDependency("git")
	skill.AddScript("scripts/analyze.sh")

	data, err := adapter.Marshal(skill)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	content := string(data)

	// Check frontmatter
	if !strings.HasPrefix(content, "---") {
		t.Error("expected Markdown to start with frontmatter")
	}
	if !strings.Contains(content, "name: version-analysis") {
		t.Error("expected name in frontmatter")
	}
	if !strings.Contains(content, "triggers:") {
		t.Error("expected triggers in frontmatter")
	}
	if !strings.Contains(content, "dependencies:") {
		t.Error("expected dependencies in frontmatter")
	}

	// Check title
	if !strings.Contains(content, "# Version Analysis") {
		t.Error("expected title in content")
	}

	// Check sections
	if !strings.Contains(content, "## Instructions") {
		t.Error("expected Instructions section")
	}
	if !strings.Contains(content, "## Scripts") {
		t.Error("expected Scripts section")
	}

	// Test round-trip
	parsed, err := adapter.Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if parsed.Name != skill.Name {
		t.Errorf("round-trip: expected Name '%s', got '%s'", skill.Name, parsed.Name)
	}
	if parsed.Description != skill.Description {
		t.Errorf("round-trip: expected Description '%s', got '%s'", skill.Description, parsed.Description)
	}
}

func TestCodexAdapter(t *testing.T) {
	adapter, ok := GetAdapter("codex")
	if !ok {
		t.Fatal("Codex adapter not found")
	}

	// Test marshal
	skill := NewSkill("version-analysis", "Analyze git history for semantic versioning")
	skill.Instructions = "Analyze commits and suggest version."

	data, err := adapter.Marshal(skill)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	content := string(data)

	// Check frontmatter (Codex uses simpler format)
	if !strings.HasPrefix(content, "---") {
		t.Error("expected Markdown to start with frontmatter")
	}
	if !strings.Contains(content, "name: version-analysis") {
		t.Error("expected name in frontmatter")
	}
	if !strings.Contains(content, "description: Analyze git history") {
		t.Error("expected description in frontmatter")
	}

	// Check instructions follow frontmatter
	if !strings.Contains(content, "Analyze commits and suggest version.") {
		t.Error("expected instructions in content")
	}

	// Test round-trip
	parsed, err := adapter.Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if parsed.Name != skill.Name {
		t.Errorf("round-trip: expected Name '%s', got '%s'", skill.Name, parsed.Name)
	}
}

func TestConvert(t *testing.T) {
	// Create a Claude skill
	claudeMD := `---
name: test-skill
description: A test skill
---

# Test Skill

Test instructions here.
`

	// Convert to Codex format
	codexData, err := Convert([]byte(claudeMD), "claude", "codex")
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	content := string(codexData)
	if !strings.Contains(content, "name: test-skill") {
		t.Error("expected name in converted output")
	}
	if !strings.Contains(content, "description: A test skill") {
		t.Error("expected description in converted output")
	}
}
