package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVaultOIDCDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VaultoidcDeleteCommand{
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

func TestVaultOIDCDeleteHelp(t *testing.T) {
	cmd := &VaultoidcDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf vaultoidc delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "Vault OIDC") {
		t.Error("Help should describe Vault OIDC configuration")
	}
	if !strings.Contains(help, "WARNING") {
		t.Error("Help should contain warning about deletion")
	}
}

func TestVaultOIDCDeleteSynopsis(t *testing.T) {
	cmd := &VaultoidcDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a Vault OIDC configuration" {
		t.Errorf("expected 'Delete a Vault OIDC configuration', got %q", synopsis)
	}
}

func TestVaultOIDCDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id flag only",
			args:          []string{"-id=voidc-ABC123"},
			expectedID:    "voidc-ABC123",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=voidc-XYZ789", "-force"},
			expectedID:    "voidc-XYZ789",
			expectedForce: true,
		},
		{
			name:          "id with force=true",
			args:          []string{"-id=voidc-TEST123", "-force=true"},
			expectedID:    "voidc-TEST123",
			expectedForce: true,
		},
		{
			name:          "id with force=false",
			args:          []string{"-id=voidc-TEST456", "-force=false"},
			expectedID:    "voidc-TEST456",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VaultoidcDeleteCommand{}

			flags := cmd.Meta.FlagSet("vaultoidc delete")
			flags.StringVar(&cmd.id, "id", "", "Vault OIDC configuration ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the ID was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
