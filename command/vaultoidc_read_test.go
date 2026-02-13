package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVaultOIDCReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VaultoidcReadCommand{
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

func TestVaultOIDCReadHelp(t *testing.T) {
	cmd := &VaultoidcReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf vaultoidc read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
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

func TestVaultOIDCReadSynopsis(t *testing.T) {
	cmd := &VaultoidcReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read Vault OIDC configuration details" {
		t.Errorf("expected 'Read Vault OIDC configuration details', got %q", synopsis)
	}
}

func TestVaultOIDCReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id flag",
			args:        []string{"-id=voidc-ABC123"},
			expectedID:  "voidc-ABC123",
			expectedFmt: "table",
		},
		{
			name:        "id with json output",
			args:        []string{"-id=voidc-XYZ789", "-output=json"},
			expectedID:  "voidc-XYZ789",
			expectedFmt: "json",
		},
		{
			name:        "id with table output",
			args:        []string{"-id=voidc-TEST123", "-output=table"},
			expectedID:  "voidc-TEST123",
			expectedFmt: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VaultoidcReadCommand{}

			flags := cmd.Meta.FlagSet("vaultoidc read")
			flags.StringVar(&cmd.id, "id", "", "Vault OIDC configuration ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the ID was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
