package mcp

import (
	"testing"
)

func TestGetAdapter(t *testing.T) {
	adapters := []string{"claude", "cursor", "windsurf", "vscode", "codex", "cline", "roo", "kiro"}

	for _, name := range adapters {
		t.Run(name, func(t *testing.T) {
			adapter, ok := GetAdapter(name)
			if !ok {
				t.Errorf("Adapter %q not found", name)
				return
			}
			if adapter.Name() != name {
				t.Errorf("Adapter name mismatch: expected %q, got %q", name, adapter.Name())
			}
		})
	}
}

func TestAdapterNames(t *testing.T) {
	names := AdapterNames()
	if len(names) < 8 {
		t.Errorf("Expected at least 8 adapters, got %d", len(names))
	}
}

func TestSupportedTools(t *testing.T) {
	tools := SupportedTools()
	expected := []string{"claude", "cursor", "windsurf", "vscode", "codex", "cline", "roo", "kiro"}

	if len(tools) != len(expected) {
		t.Errorf("Expected %d tools, got %d", len(expected), len(tools))
	}

	for i, tool := range expected {
		if tools[i] != tool {
			t.Errorf("Tool mismatch at index %d: expected %q, got %q", i, tool, tools[i])
		}
	}
}

func TestConvert(t *testing.T) {
	claudeJSON := []byte(`{
		"mcpServers": {
			"test": {
				"command": "npx",
				"args": ["-y", "test-server"]
			}
		}
	}`)

	// Convert Claude to VS Code
	vscodeData, err := Convert(claudeJSON, "claude", "vscode")
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify the output contains VS Code's "servers" key
	if len(vscodeData) == 0 {
		t.Error("Convert returned empty data")
	}

	// Parse back to verify
	vscodeAdapter, _ := GetAdapter("vscode")
	cfg, err := vscodeAdapter.Parse(vscodeData)
	if err != nil {
		t.Fatalf("Failed to parse converted data: %v", err)
	}

	server, ok := cfg.GetServer("test")
	if !ok {
		t.Fatal("test server not found after conversion")
	}
	if server.Command != "npx" {
		t.Errorf("Expected command 'npx', got %q", server.Command)
	}
}

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	if cfg == nil {
		t.Fatal("NewConfig returned nil")
	}
	if cfg.Servers == nil {
		t.Error("Servers map should be initialized")
	}
}

func TestTransportConstants(t *testing.T) {
	if TransportStdio != "stdio" {
		t.Errorf("TransportStdio mismatch")
	}
	if TransportHTTP != "http" {
		t.Errorf("TransportHTTP mismatch")
	}
	if TransportSSE != "sse" {
		t.Errorf("TransportSSE mismatch")
	}
}
