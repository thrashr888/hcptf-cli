package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestHYOKListRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKListCommand{
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

func TestHYOKListHelp(t *testing.T) {
	cmd := &HYOKListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf hyok list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -organization is required")
	}
	if !strings.Contains(help, "HYOK") {
		t.Error("Help should mention HYOK")
	}
}

func TestHYOKListSynopsis(t *testing.T) {
	cmd := &HYOKListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List HYOK configurations for an organization" {
		t.Errorf("expected 'List HYOK configurations for an organization', got %q", synopsis)
	}
}

func TestHYOKListFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOrg    string
		expectedFormat string
	}{
		{
			name:           "organization, default format",
			args:           []string{"-organization=my-org"},
			expectedOrg:    "my-org",
			expectedFormat: "table",
		},
		{
			name:           "organization, table format",
			args:           []string{"-organization=test-org", "-output=table"},
			expectedOrg:    "test-org",
			expectedFormat: "table",
		},
		{
			name:           "organization, json format",
			args:           []string{"-organization=prod-org", "-output=json"},
			expectedOrg:    "prod-org",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &HYOKListCommand{}

			flags := cmd.Meta.FlagSet("hyok list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
