package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationTokenReadRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationTokenReadCommand{
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

func TestOrganizationTokenReadHelp(t *testing.T) {
	cmd := &OrganizationTokenReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organizationtoken read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -organization is required")
	}
}

func TestOrganizationTokenReadSynopsis(t *testing.T) {
	cmd := &OrganizationTokenReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show organization token details" {
		t.Errorf("expected 'Show organization token details', got %q", synopsis)
	}
}

func TestOrganizationTokenReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOrg    string
		expectedFormat string
	}{
		{
			name:           "organization only, default format",
			args:           []string{"-organization=test-org"},
			expectedOrg:    "test-org",
			expectedFormat: "table",
		},
		{
			name:           "org alias, default format",
			args:           []string{"-org=my-org"},
			expectedOrg:    "my-org",
			expectedFormat: "table",
		},
		{
			name:           "organization with table format",
			args:           []string{"-organization=prod-org", "-output=table"},
			expectedOrg:    "prod-org",
			expectedFormat: "table",
		},
		{
			name:           "org alias with json format",
			args:           []string{"-org=dev-org", "-output=json"},
			expectedOrg:    "dev-org",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationTokenReadCommand{}

			flags := cmd.Meta.FlagSet("organizationtoken read")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
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
