package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationMembershipListRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationMembershipListCommand{
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

func TestOrganizationMembershipListHelp(t *testing.T) {
	cmd := &OrganizationMembershipListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organizationmembership list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate organization flag is required")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org alias")
	}
	if !strings.Contains(help, "-status") {
		t.Error("Help should mention -status flag")
	}
	if !strings.Contains(help, "-email") {
		t.Error("Help should mention -email flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
}

func TestOrganizationMembershipListSynopsis(t *testing.T) {
	cmd := &OrganizationMembershipListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List organization memberships" {
		t.Errorf("expected 'List organization memberships', got %q", synopsis)
	}
}

func TestOrganizationMembershipListFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOrg    string
		expectedStatus string
		expectedEmail  string
		expectedFmt    string
	}{
		{
			name:           "organization, default format",
			args:           []string{"-organization=my-org"},
			expectedOrg:    "my-org",
			expectedStatus: "",
			expectedEmail:  "",
			expectedFmt:    "table",
		},
		{
			name:           "organization alias",
			args:           []string{"-org=test-org"},
			expectedOrg:    "test-org",
			expectedStatus: "",
			expectedEmail:  "",
			expectedFmt:    "table",
		},
		{
			name:           "organization with status filter",
			args:           []string{"-org=my-org", "-status=active"},
			expectedOrg:    "my-org",
			expectedStatus: "active",
			expectedEmail:  "",
			expectedFmt:    "table",
		},
		{
			name:           "organization with invited status",
			args:           []string{"-org=my-org", "-status=invited"},
			expectedOrg:    "my-org",
			expectedStatus: "invited",
			expectedEmail:  "",
			expectedFmt:    "table",
		},
		{
			name:           "organization with email filter",
			args:           []string{"-org=my-org", "-email=user@example.com"},
			expectedOrg:    "my-org",
			expectedStatus: "",
			expectedEmail:  "user@example.com",
			expectedFmt:    "table",
		},
		{
			name:           "organization with table format",
			args:           []string{"-org=prod-org", "-output=table"},
			expectedOrg:    "prod-org",
			expectedStatus: "",
			expectedEmail:  "",
			expectedFmt:    "table",
		},
		{
			name:           "organization with json format",
			args:           []string{"-org=dev-org", "-output=json"},
			expectedOrg:    "dev-org",
			expectedStatus: "",
			expectedEmail:  "",
			expectedFmt:    "json",
		},
		{
			name:           "all filters combined",
			args:           []string{"-org=my-org", "-status=active", "-email=test@example.com", "-output=json"},
			expectedOrg:    "my-org",
			expectedStatus: "active",
			expectedEmail:  "test@example.com",
			expectedFmt:    "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationMembershipListCommand{}

			flags := cmd.Meta.FlagSet("organizationmembership list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.status, "status", "", "Filter by status (invited, active)")
			flags.StringVar(&cmd.email, "email", "", "Filter by email address")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the status was set correctly
			if cmd.status != tt.expectedStatus {
				t.Errorf("expected status %q, got %q", tt.expectedStatus, cmd.status)
			}

			// Verify the email was set correctly
			if cmd.email != tt.expectedEmail {
				t.Errorf("expected email %q, got %q", tt.expectedEmail, cmd.email)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
