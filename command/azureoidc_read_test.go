package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAzureOIDCReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AzureoidcReadCommand{
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

func TestAzureOIDCReadHelp(t *testing.T) {
	cmd := &AzureoidcReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf azureoidc read") {
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
	if !strings.Contains(help, "Azure OIDC") {
		t.Error("Help should describe Azure OIDC configuration")
	}
}

func TestAzureOIDCReadSynopsis(t *testing.T) {
	cmd := &AzureoidcReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read Azure OIDC configuration details" {
		t.Errorf("expected 'Read Azure OIDC configuration details', got %q", synopsis)
	}
}

func TestAzureOIDCReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id flag",
			args:        []string{"-id=azoidc-ABC123"},
			expectedID:  "azoidc-ABC123",
			expectedFmt: "table",
		},
		{
			name:        "id with json output",
			args:        []string{"-id=azoidc-XYZ789", "-output=json"},
			expectedID:  "azoidc-XYZ789",
			expectedFmt: "json",
		},
		{
			name:        "id with table output",
			args:        []string{"-id=azoidc-DEF456", "-output=table"},
			expectedID:  "azoidc-DEF456",
			expectedFmt: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AzureoidcReadCommand{}

			flags := cmd.Meta.FlagSet("azureoidc read")
			flags.StringVar(&cmd.id, "id", "", "Azure OIDC configuration ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the ID was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected ID %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
