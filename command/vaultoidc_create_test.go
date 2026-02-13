package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVaultOIDCCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VaultoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestVaultOIDCCreateRequiresAddress(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VaultoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-address") {
		t.Fatalf("expected address error, got %q", out)
	}
}

func TestVaultOIDCCreateRequiresRole(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VaultoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-address=https://vault.example.com:8200"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-role") {
		t.Fatalf("expected role error, got %q", out)
	}
}

func TestVaultOIDCCreateRequiresNamespace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VaultoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-address=https://vault.example.com:8200", "-role=terraform-role"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-namespace") {
		t.Fatalf("expected namespace error, got %q", out)
	}
}

func TestVaultOIDCCreateRequiresEmptyAddress(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VaultoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-address=", "-role=terraform-role", "-namespace=admin"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-address") {
		t.Fatalf("expected address error, got %q", out)
	}
}

func TestVaultOIDCCreateRequiresEmptyRole(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VaultoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-address=https://vault.example.com:8200", "-role=", "-namespace=admin"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-role") {
		t.Fatalf("expected role error, got %q", out)
	}
}

func TestVaultOIDCCreateRequiresEmptyNamespace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VaultoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-address=https://vault.example.com:8200", "-role=terraform-role", "-namespace="})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-namespace") {
		t.Fatalf("expected namespace error, got %q", out)
	}
}

func TestVaultOIDCCreateHelp(t *testing.T) {
	cmd := &VaultoidcCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf vaultoidc create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
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

func TestVaultOIDCCreateSynopsis(t *testing.T) {
	cmd := &VaultoidcCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a Vault OIDC configuration for dynamic credentials" {
		t.Errorf("expected 'Create a Vault OIDC configuration for dynamic credentials', got %q", synopsis)
	}
}

func TestVaultOIDCCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedOrg      string
		expectedAddress  string
		expectedRole     string
		expectedNS       string
		expectedAuthPath string
		expectedFmt      string
	}{
		{
			name:             "all required flags",
			args:             []string{"-organization=my-org", "-address=https://vault.example.com:8200", "-role=terraform-role", "-namespace=admin"},
			expectedOrg:      "my-org",
			expectedAddress:  "https://vault.example.com:8200",
			expectedRole:     "terraform-role",
			expectedNS:       "admin",
			expectedAuthPath: "jwt",
			expectedFmt:      "table",
		},
		{
			name:             "org alias flag",
			args:             []string{"-org=test-org", "-address=https://vault.example.com:8200", "-role=test-role", "-namespace=root"},
			expectedOrg:      "test-org",
			expectedAddress:  "https://vault.example.com:8200",
			expectedRole:     "test-role",
			expectedNS:       "root",
			expectedAuthPath: "jwt",
			expectedFmt:      "table",
		},
		{
			name:             "with custom auth path",
			args:             []string{"-organization=my-org", "-address=https://vault.example.com:8200", "-role=terraform-role", "-namespace=admin", "-auth-path=jwt-auth"},
			expectedOrg:      "my-org",
			expectedAddress:  "https://vault.example.com:8200",
			expectedRole:     "terraform-role",
			expectedNS:       "admin",
			expectedAuthPath: "jwt-auth",
			expectedFmt:      "table",
		},
		{
			name:             "with json output",
			args:             []string{"-organization=my-org", "-address=https://vault.example.com:8200", "-role=terraform-role", "-namespace=admin", "-output=json"},
			expectedOrg:      "my-org",
			expectedAddress:  "https://vault.example.com:8200",
			expectedRole:     "terraform-role",
			expectedNS:       "admin",
			expectedAuthPath: "jwt",
			expectedFmt:      "json",
		},
		{
			name:             "with table output",
			args:             []string{"-org=test-org", "-address=https://my-vault.vault.cloud:8200", "-role=my-role", "-namespace=admin", "-output=table"},
			expectedOrg:      "test-org",
			expectedAddress:  "https://my-vault.vault.cloud:8200",
			expectedRole:     "my-role",
			expectedNS:       "admin",
			expectedAuthPath: "jwt",
			expectedFmt:      "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VaultoidcCreateCommand{}

			flags := cmd.Meta.FlagSet("vaultoidc create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.address, "address", "", "Vault instance address (required)")
			flags.StringVar(&cmd.roleName, "role", "", "Vault JWT auth role name (required)")
			flags.StringVar(&cmd.namespace, "namespace", "", "Vault namespace (required)")
			flags.StringVar(&cmd.jwtAuthPath, "auth-path", "jwt", "Vault JWT auth mount path (default: jwt)")
			flags.StringVar(&cmd.tlsCACertificate, "encoded-cacert", "", "Base64-encoded CA certificate (optional)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
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
