package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationTokenCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationTokenCreateCommand{
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

func TestOrganizationTokenCreateHelp(t *testing.T) {
	cmd := &OrganizationTokenCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organizationtoken create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
	}
	if !strings.Contains(help, "-expired-at") {
		t.Error("Help should mention -expired-at flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -organization is required")
	}
	if !strings.Contains(help, "ISO 8601") {
		t.Error("Help should mention ISO 8601 date format")
	}
}

func TestOrganizationTokenCreateSynopsis(t *testing.T) {
	cmd := &OrganizationTokenCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create an organization token" {
		t.Errorf("expected 'Create an organization token', got %q", synopsis)
	}
}

func TestOrganizationTokenCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOrg    string
		expectedExp    string
		expectedFormat string
	}{
		{
			name:           "organization only, default format",
			args:           []string{"-organization=test-org"},
			expectedOrg:    "test-org",
			expectedExp:    "",
			expectedFormat: "table",
		},
		{
			name:           "org alias, default format",
			args:           []string{"-org=my-org"},
			expectedOrg:    "my-org",
			expectedExp:    "",
			expectedFormat: "table",
		},
		{
			name:           "organization with expiration",
			args:           []string{"-organization=prod-org", "-expired-at=2024-12-31T23:59:59Z"},
			expectedOrg:    "prod-org",
			expectedExp:    "2024-12-31T23:59:59Z",
			expectedFormat: "table",
		},
		{
			name:           "org alias with expiration and json format",
			args:           []string{"-org=dev-org", "-expired-at=2025-06-30T12:00:00Z", "-output=json"},
			expectedOrg:    "dev-org",
			expectedExp:    "2025-06-30T12:00:00Z",
			expectedFormat: "json",
		},
		{
			name:           "organization with json format",
			args:           []string{"-organization=test-org", "-output=json"},
			expectedOrg:    "test-org",
			expectedExp:    "",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationTokenCreateCommand{}

			flags := cmd.Meta.FlagSet("organizationtoken create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.expiredAt, "expired-at", "", "Expiration date in ISO 8601 format")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the expiredAt was set correctly
			if cmd.expiredAt != tt.expectedExp {
				t.Errorf("expected expired-at %q, got %q", tt.expectedExp, cmd.expiredAt)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
