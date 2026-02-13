package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestGCPOIDCReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcReadCommand{
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

func TestGCPOIDCReadRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcReadCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "required") {
		t.Fatalf("expected required error, got %q", out)
	}
}

func TestGCPOIDCReadHelp(t *testing.T) {
	cmd := &GCPoidcReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf gcpoidc read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
	if !strings.Contains(help, "gcpoidc-") {
		t.Error("Help should show ID format example")
	}
}

func TestGCPOIDCReadSynopsis(t *testing.T) {
	cmd := &GCPoidcReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read GCP OIDC configuration details" {
		t.Errorf("expected 'Read GCP OIDC configuration details', got %q", synopsis)
	}
}

func TestGCPOIDCReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "id flag with default format",
			args:           []string{"-id=gcpoidc-ABC123"},
			expectedID:     "gcpoidc-ABC123",
			expectedFormat: "table",
		},
		{
			name:           "id flag with json format",
			args:           []string{"-id=gcpoidc-XYZ789", "-output=json"},
			expectedID:     "gcpoidc-XYZ789",
			expectedFormat: "json",
		},
		{
			name:           "id flag with table format",
			args:           []string{"-id=gcpoidc-DEF456", "-output=table"},
			expectedID:     "gcpoidc-DEF456",
			expectedFormat: "table",
		},
		{
			name:           "long id format",
			args:           []string{"-id=gcpoidc-1234567890ABCDEF"},
			expectedID:     "gcpoidc-1234567890ABCDEF",
			expectedFormat: "table",
		},
		{
			name:           "id with underscores",
			args:           []string{"-id=gcpoidc-test_config_123", "-output=json"},
			expectedID:     "gcpoidc-test_config_123",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &GCPoidcReadCommand{}

			flags := cmd.Meta.FlagSet("gcpoidc read")
			flags.StringVar(&cmd.id, "id", "", "GCP OIDC configuration ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
