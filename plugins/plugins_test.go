package plugins

import (
	"encoding/json"
	"testing"
)

func TestAdapterRegistry(t *testing.T) {
	names := AdapterNames()

	// Should have at least Claude and Gemini adapters
	if len(names) < 2 {
		t.Errorf("expected at least 2 adapters, got %d", len(names))
	}

	// Check Claude adapter exists
	claude, ok := GetAdapter("claude")
	if !ok {
		t.Error("expected Claude adapter to be registered")
	}
	if claude.Name() != "claude" {
		t.Errorf("expected Claude adapter name 'claude', got '%s'", claude.Name())
	}

	// Check Gemini adapter exists
	gemini, ok := GetAdapter("gemini")
	if !ok {
		t.Error("expected Gemini adapter to be registered")
	}
	if gemini.Name() != "gemini" {
		t.Errorf("expected Gemini adapter name 'gemini', got '%s'", gemini.Name())
	}
}

func TestClaudeAdapter(t *testing.T) {
	adapter, ok := GetAdapter("claude")
	if !ok {
		t.Fatal("Claude adapter not found")
	}

	// Test marshal
	plugin := NewPlugin("test-plugin", "1.0.0", "A test plugin")
	plugin.Commands = "commands"
	plugin.Skills = "skills"

	data, err := adapter.Marshal(plugin)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify JSON structure
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse marshaled JSON: %v", err)
	}

	if result["name"] != "test-plugin" {
		t.Errorf("expected name 'test-plugin', got '%v'", result["name"])
	}
	if result["commands"] != "./commands/" {
		t.Errorf("expected commands './commands/', got '%v'", result["commands"])
	}

	// Test round-trip
	parsed, err := adapter.Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if parsed.Name != plugin.Name {
		t.Errorf("round-trip: expected Name '%s', got '%s'", plugin.Name, parsed.Name)
	}
}

func TestGeminiAdapter(t *testing.T) {
	adapter, ok := GetAdapter("gemini")
	if !ok {
		t.Fatal("Gemini adapter not found")
	}

	// Test marshal
	plugin := NewPlugin("test-plugin", "1.0.0", "A test plugin")
	plugin.Context = "This is the plugin context"
	plugin.AddMCPServer("github", MCPServer{
		Command: "npx",
		Args:    []string{"-y", "@modelcontextprotocol/server-github"},
	})

	data, err := adapter.Marshal(plugin)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify JSON structure
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse marshaled JSON: %v", err)
	}

	if result["name"] != "test-plugin" {
		t.Errorf("expected name 'test-plugin', got '%v'", result["name"])
	}
	if result["contextFileName"] != "GEMINI.md" {
		t.Errorf("expected contextFileName 'GEMINI.md', got '%v'", result["contextFileName"])
	}

	// Test round-trip
	parsed, err := adapter.Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if parsed.Name != plugin.Name {
		t.Errorf("round-trip: expected Name '%s', got '%s'", plugin.Name, parsed.Name)
	}
	if len(parsed.MCPServers) != 1 {
		t.Errorf("round-trip: expected 1 MCP server, got %d", len(parsed.MCPServers))
	}
}

func TestConvert(t *testing.T) {
	// Create a Claude plugin JSON
	claudeJSON := `{
		"name": "convert-test",
		"version": "1.0.0",
		"description": "Test conversion",
		"commands": "./commands/"
	}`

	// Convert to Gemini format
	geminiData, err := Convert([]byte(claudeJSON), "claude", "gemini")
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(geminiData, &result); err != nil {
		t.Fatalf("Failed to parse converted JSON: %v", err)
	}

	if result["name"] != "convert-test" {
		t.Errorf("expected name 'convert-test', got '%v'", result["name"])
	}
}
