package core

import (
	"encoding/json"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	if cfg == nil {
		t.Fatal("NewConfig returned nil")
	}
	if cfg.Servers == nil {
		t.Error("Servers map should be initialized")
	}
}

func TestConfigAddServer(t *testing.T) {
	cfg := NewConfig()
	server := Server{
		Transport: TransportStdio,
		Command:   "npx",
		Args:      []string{"-y", "@example/mcp-server"},
	}

	cfg.AddServer("test-server", server)

	got, ok := cfg.GetServer("test-server")
	if !ok {
		t.Fatal("Server not found after adding")
	}
	if got.Command != "npx" {
		t.Errorf("Expected command 'npx', got %q", got.Command)
	}
}

func TestConfigRemoveServer(t *testing.T) {
	cfg := NewConfig()
	cfg.AddServer("test-server", Server{Command: "test"})
	cfg.RemoveServer("test-server")

	_, ok := cfg.GetServer("test-server")
	if ok {
		t.Error("Server should not exist after removal")
	}
}

func TestConfigServerNames(t *testing.T) {
	cfg := NewConfig()
	cfg.AddServer("server-a", Server{Command: "a"})
	cfg.AddServer("server-b", Server{Command: "b"})

	names := cfg.ServerNames()
	if len(names) != 2 {
		t.Errorf("Expected 2 server names, got %d", len(names))
	}
}

func TestConfigStdioServers(t *testing.T) {
	cfg := NewConfig()
	cfg.AddServer("stdio-server", Server{Transport: TransportStdio, Command: "test"})
	cfg.AddServer("http-server", Server{Transport: TransportHTTP, URL: "http://example.com"})

	stdioServers := cfg.StdioServers()
	if len(stdioServers) != 1 {
		t.Errorf("Expected 1 stdio server, got %d", len(stdioServers))
	}
	if _, ok := stdioServers["stdio-server"]; !ok {
		t.Error("stdio-server should be in StdioServers result")
	}
}

func TestConfigRemoteServers(t *testing.T) {
	cfg := NewConfig()
	cfg.AddServer("stdio-server", Server{Transport: TransportStdio, Command: "test"})
	cfg.AddServer("http-server", Server{Transport: TransportHTTP, URL: "http://example.com"})

	remoteServers := cfg.RemoteServers()
	if len(remoteServers) != 1 {
		t.Errorf("Expected 1 remote server, got %d", len(remoteServers))
	}
	if _, ok := remoteServers["http-server"]; !ok {
		t.Error("http-server should be in RemoteServers result")
	}
}

func TestConfigMerge(t *testing.T) {
	cfg1 := NewConfig()
	cfg1.AddServer("server-a", Server{Command: "a"})

	cfg2 := NewConfig()
	cfg2.AddServer("server-b", Server{Command: "b"})

	cfg1.Merge(cfg2)

	if len(cfg1.Servers) != 2 {
		t.Errorf("Expected 2 servers after merge, got %d", len(cfg1.Servers))
	}
}

func TestConfigJSON(t *testing.T) {
	cfg := NewConfig()
	cfg.AddServer("test", Server{
		Transport: TransportStdio,
		Command:   "npx",
		Args:      []string{"-y", "test"},
	})

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded Config
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(decoded.Servers) != 1 {
		t.Errorf("Expected 1 server after round-trip, got %d", len(decoded.Servers))
	}
}

func TestConfigAddInput(t *testing.T) {
	cfg := NewConfig()
	cfg.AddInput(InputVariable{
		Type:        "promptString",
		ID:          "api-key",
		Description: "API Key",
		Password:    true,
	})

	input, ok := cfg.GetInput("api-key")
	if !ok {
		t.Fatal("Input not found after adding")
	}
	if input.Description != "API Key" {
		t.Errorf("Expected description 'API Key', got %q", input.Description)
	}
	if !input.Password {
		t.Error("Expected Password to be true")
	}
}
