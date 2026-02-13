package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVaultOIDCUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VaultoidcUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestVaultOIDCUpdateValidatesEmptyID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VaultoidcUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=", "-address=https://vault.example.com:8200"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestVaultOIDCUpdateHelp(t *testing.T) {
	cmd := &VaultoidcUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf vaultoidc update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-address") {
		t.Error("Help should mention -address flag")
	}
	if !strings.Contains(help, "-role") {
		t.Error("Help should mention -role flag")
	}
	if !strings.Contains(help, "-namespace") {
		t.Error("Help should mention -namespace flag")
	}
	if !strings.Contains(help, "-auth-path") {
		t.Error("Help should mention -auth-path flag")
	}
	if !strings.Contains(help, "-encoded-cacert") {
		t.Error("Help should mention -encoded-cacert flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "Vault OIDC") {
		t.Error("Help should describe Vault OIDC configuration")
	}
}

func TestVaultOIDCUpdateSynopsis(t *testing.T) {
	cmd := &VaultoidcUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update Vault OIDC configuration settings" {
		t.Errorf("expected 'Update Vault OIDC configuration settings', got %q", synopsis)
	}
}

func TestVaultOIDCUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedID       string
		expectedAddress  string
		expectedRole     string
		expectedNS       string
		expectedAuthPath string
		expectedFmt      string
	}{
		{
			name:        "id flag only",
			args:        []string{"-id=voidc-ABC123"},
			expectedID:  "voidc-ABC123",
			expectedFmt: "table",
		},
		{
			name:            "id and address",
			args:            []string{"-id=voidc-ABC123", "-address=https://vault.example.com:8200"},
			expectedID:      "voidc-ABC123",
			expectedAddress: "https://vault.example.com:8200",
			expectedFmt:     "table",
		},
		{
			name:         "id and role",
			args:         []string{"-id=voidc-ABC123", "-role=new-terraform-role"},
			expectedID:   "voidc-ABC123",
			expectedRole: "new-terraform-role",
			expectedFmt:  "table",
		},
		{
			name:       "id and namespace",
			args:       []string{"-id=voidc-ABC123", "-namespace=admin"},
			expectedID: "voidc-ABC123",
			expectedNS: "admin",
			expectedFmt: "table",
		},
		{
			name:             "id and auth-path",
			args:             []string{"-id=voidc-ABC123", "-auth-path=jwt-auth"},
			expectedID:       "voidc-ABC123",
			expectedAuthPath: "jwt-auth",
			expectedFmt:      "table",
		},
		{
			name:            "multiple updates",
			args:            []string{"-id=voidc-XYZ789", "-address=https://new-vault.example.com:8200", "-role=updated-role", "-namespace=root"},
			expectedID:      "voidc-XYZ789",
			expectedAddress: "https://new-vault.example.com:8200",
			expectedRole:    "updated-role",
			expectedNS:      "root",
			expectedFmt:     "table",
		},
		{
			name:            "with json output",
			args:            []string{"-id=voidc-ABC123", "-address=https://vault.example.com:8200", "-output=json"},
			expectedID:      "voidc-ABC123",
			expectedAddress: "https://vault.example.com:8200",
			expectedFmt:     "json",
		},
		{
			name:        "with table output",
			args:        []string{"-id=voidc-ABC123", "-role=test-role", "-output=table"},
			expectedID:  "voidc-ABC123",
			expectedRole: "test-role",
			expectedFmt: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VaultoidcUpdateCommand{}

			flags := cmd.Meta.FlagSet("vaultoidc update")
			flags.StringVar(&cmd.id, "id", "", "Vault OIDC configuration ID (required)")
			flags.StringVar(&cmd.address, "address", "", "Vault instance address")
			flags.StringVar(&cmd.roleName, "role", "", "Vault JWT auth role name")
			flags.StringVar(&cmd.namespace, "namespace", "", "Vault namespace")
			flags.StringVar(&cmd.jwtAuthPath, "auth-path", "", "Vault JWT auth mount path")
			flags.StringVar(&cmd.tlsCACertificate, "encoded-cacert", "", "Base64-encoded CA certificate")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the ID was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the address was set correctly
			if cmd.address != tt.expectedAddress {
				t.Errorf("expected address %q, got %q", tt.expectedAddress, cmd.address)
			}

			// Verify the role was set correctly
			if cmd.roleName != tt.expectedRole {
				t.Errorf("expected role %q, got %q", tt.expectedRole, cmd.roleName)
			}

			// Verify the namespace was set correctly
			if cmd.namespace != tt.expectedNS {
				t.Errorf("expected namespace %q, got %q", tt.expectedNS, cmd.namespace)
			}

			// Verify the auth path was set correctly
			if cmd.jwtAuthPath != tt.expectedAuthPath {
				t.Errorf("expected auth path %q, got %q", tt.expectedAuthPath, cmd.jwtAuthPath)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
