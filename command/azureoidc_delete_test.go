package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAzureOIDCDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AzureoidcDeleteCommand{
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

func TestAzureOIDCDeleteHelp(t *testing.T) {
	cmd := &AzureoidcDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf azureoidc delete") {
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
	if !strings.Contains(help, "Azure OIDC") {
		t.Error("Help should describe Azure OIDC configuration")
	}
	if !strings.Contains(help, "WARNING") {
		t.Error("Help should contain a warning about deletion")
	}
}

func TestAzureOIDCDeleteSynopsis(t *testing.T) {
	cmd := &AzureoidcDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete an Azure OIDC configuration" {
		t.Errorf("expected 'Delete an Azure OIDC configuration', got %q", synopsis)
	}
}

func TestAzureOIDCDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id only",
			args:          []string{"-id=azoidc-ABC123"},
			expectedID:    "azoidc-ABC123",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=azoidc-XYZ789", "-force"},
			expectedID:    "azoidc-XYZ789",
			expectedForce: true,
		},
		{
			name:          "id with force=true",
			args:          []string{"-id=azoidc-DEF456", "-force=true"},
			expectedID:    "azoidc-DEF456",
			expectedForce: true,
		},
		{
			name:          "id with force=false",
			args:          []string{"-id=azoidc-GHI789", "-force=false"},
			expectedID:    "azoidc-GHI789",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AzureoidcDeleteCommand{}

			flags := cmd.Meta.FlagSet("azureoidc delete")
			flags.StringVar(&cmd.id, "id", "", "Azure OIDC configuration ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the ID was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected ID %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
