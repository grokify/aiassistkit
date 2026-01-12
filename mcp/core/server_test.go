package core

import "testing"

func TestServerIsStdio(t *testing.T) {
	tests := []struct {
		name     string
		server   Server
		expected bool
	}{
		{
			name:     "explicit stdio transport",
			server:   Server{Transport: TransportStdio},
			expected: true,
		},
		{
			name:     "inferred from command",
			server:   Server{Command: "npx"},
			expected: true,
		},
		{
			name:     "http transport",
			server:   Server{Transport: TransportHTTP},
			expected: false,
		},
		{
			name:     "url only",
			server:   Server{URL: "http://example.com"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.server.IsStdio(); got != tt.expected {
				t.Errorf("IsStdio() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestServerIsHTTP(t *testing.T) {
	tests := []struct {
		name     string
		server   Server
		expected bool
	}{
		{
			name:     "explicit http transport",
			server:   Server{Transport: TransportHTTP},
			expected: true,
		},
		{
			name:     "inferred from url",
			server:   Server{URL: "http://example.com"},
			expected: true,
		},
		{
			name:     "stdio transport",
			server:   Server{Transport: TransportStdio},
			expected: false,
		},
		{
			name:     "sse transport",
			server:   Server{Transport: TransportSSE, URL: "http://example.com"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.server.IsHTTP(); got != tt.expected {
				t.Errorf("IsHTTP() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestServerIsEnabled(t *testing.T) {
	t.Run("nil enabled defaults to true", func(t *testing.T) {
		server := Server{}
		if !server.IsEnabled() {
			t.Error("Expected IsEnabled to be true by default")
		}
	})

	t.Run("explicit true", func(t *testing.T) {
		enabled := true
		server := Server{Enabled: &enabled}
		if !server.IsEnabled() {
			t.Error("Expected IsEnabled to be true")
		}
	})

	t.Run("explicit false", func(t *testing.T) {
		enabled := false
		server := Server{Enabled: &enabled}
		if server.IsEnabled() {
			t.Error("Expected IsEnabled to be false")
		}
	})
}

func TestServerSetEnabled(t *testing.T) {
	server := Server{}
	server.SetEnabled(false)
	if server.IsEnabled() {
		t.Error("Expected IsEnabled to be false after SetEnabled(false)")
	}

	server.SetEnabled(true)
	if !server.IsEnabled() {
		t.Error("Expected IsEnabled to be true after SetEnabled(true)")
	}
}

func TestServerInferTransport(t *testing.T) {
	tests := []struct {
		name     string
		server   Server
		expected TransportType
	}{
		{
			name:     "explicit transport",
			server:   Server{Transport: TransportSSE},
			expected: TransportSSE,
		},
		{
			name:     "infer stdio from command",
			server:   Server{Command: "npx"},
			expected: TransportStdio,
		},
		{
			name:     "infer http from url",
			server:   Server{URL: "http://example.com"},
			expected: TransportHTTP,
		},
		{
			name:     "empty server",
			server:   Server{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.server.InferTransport(); got != tt.expected {
				t.Errorf("InferTransport() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestServerValidate(t *testing.T) {
	tests := []struct {
		name      string
		server    Server
		wantError bool
	}{
		{
			name:      "valid stdio server",
			server:    Server{Command: "npx"},
			wantError: false,
		},
		{
			name:      "valid http server",
			server:    Server{URL: "http://example.com"},
			wantError: false,
		},
		{
			name:      "no command or url",
			server:    Server{},
			wantError: true,
		},
		{
			name:      "both command and url",
			server:    Server{Command: "npx", URL: "http://example.com"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
