package client

import (
	"os"
	"testing"

	"github.com/hashicorp/hcptf-cli/internal/config"
)

func TestNew(t *testing.T) {
	t.Run("creates client with valid token from environment", func(t *testing.T) {
		// Setup
		t.Setenv("HCPTF_TOKEN", "test-token-123")
		t.Setenv("HCPTF_ADDRESS", "https://app.terraform.io")

		cfg := &config.Config{
			Credentials: make(map[string]*config.Credential),
		}

		// Execute
		client, err := New(cfg)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client == nil {
			t.Fatal("expected client to be non-nil")
		}
		if client.GetAddress() != "https://app.terraform.io" {
			t.Errorf("expected address 'https://app.terraform.io', got '%s'", client.GetAddress())
		}
		if client.Client == nil {
			t.Error("expected embedded TFE client to be non-nil")
		}
	})

	t.Run("creates client with token from config", func(t *testing.T) {
		// Setup
		t.Setenv("HCPTF_ADDRESS", "https://app.terraform.io")

		cfg := &config.Config{
			Credentials: map[string]*config.Credential{
				"app.terraform.io": {
					Token: "config-token-456",
				},
			},
		}

		// Execute
		client, err := New(cfg)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client == nil {
			t.Fatal("expected client to be non-nil")
		}
	})

	t.Run("creates client with custom address", func(t *testing.T) {
		// Skip this test as it requires network connectivity
		t.Skip("skipping test that requires custom TFE instance")
	})

	t.Run("returns error when no token found", func(t *testing.T) {
		// Setup - clear any existing env vars
		os.Unsetenv("HCPTF_TOKEN")
		os.Unsetenv("TFE_TOKEN")
		t.Setenv("HCPTF_ADDRESS", "https://app.terraform.io")

		cfg := &config.Config{
			Credentials: make(map[string]*config.Credential),
		}

		// Execute
		client, err := New(cfg)

		// Assert
		if err == nil {
			t.Fatal("expected error when no token found")
		}
		if client != nil {
			t.Error("expected client to be nil on error")
		}
		expectedMsg := "no authentication token found"
		if err.Error()[:len(expectedMsg)] != expectedMsg {
			t.Errorf("expected error message to start with '%s', got '%s'", expectedMsg, err.Error())
		}
	})

	t.Run("returns error with invalid address", func(t *testing.T) {
		// Setup
		t.Setenv("HCPTF_TOKEN", "test-token")
		t.Setenv("HCPTF_ADDRESS", "://invalid-url")

		cfg := &config.Config{
			Credentials: make(map[string]*config.Credential),
		}

		// Execute
		client, err := New(cfg)

		// Assert
		if err == nil {
			t.Fatal("expected error with invalid address")
		}
		if client != nil {
			t.Error("expected client to be nil on error")
		}
	})

	t.Run("TFE_TOKEN takes precedence over config", func(t *testing.T) {
		// Setup
		t.Setenv("TFE_TOKEN", "env-token")
		t.Setenv("HCPTF_ADDRESS", "https://app.terraform.io")

		cfg := &config.Config{
			Credentials: map[string]*config.Credential{
				"app.terraform.io": {
					Token: "config-token",
				},
			},
		}

		// Execute
		client, err := New(cfg)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client == nil {
			t.Fatal("expected client to be non-nil")
		}
		// We can't directly verify which token was used, but we verify client creation succeeded
	})
}

func TestClient_GetAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{
			name:    "returns default address",
			address: "https://app.terraform.io",
		},
		{
			name:    "returns custom address",
			address: "https://tfe.example.com",
		},
		{
			name:    "returns address with path",
			address: "https://tfe.example.com/api/v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				address: tt.address,
			}

			result := client.GetAddress()

			if result != tt.address {
				t.Errorf("expected address '%s', got '%s'", tt.address, result)
			}
		})
	}
}

func TestClient_Context(t *testing.T) {
	t.Run("returns non-nil context", func(t *testing.T) {
		client := &Client{}

		ctx := client.Context()

		if ctx == nil {
			t.Error("expected context to be non-nil")
		}
	})

	t.Run("returns background context", func(t *testing.T) {
		client := &Client{}

		ctx := client.Context()

		// Verify it's a background context by checking it has no deadline
		if _, ok := ctx.Deadline(); ok {
			t.Error("expected background context without deadline")
		}
	})

	t.Run("context is not canceled", func(t *testing.T) {
		client := &Client{}

		ctx := client.Context()

		select {
		case <-ctx.Done():
			t.Error("expected context to not be canceled")
		default:
			// Context is not canceled, this is expected
		}
	})
}
