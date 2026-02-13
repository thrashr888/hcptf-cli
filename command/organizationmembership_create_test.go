package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationMembershipCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationMembershipCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-email=user@example.com", "-team-ids=team-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestOrganizationMembershipCreateRequiresEmail(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationMembershipCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-team-ids=team-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-email") {
		t.Fatalf("expected email error, got %q", out)
	}
}

func TestOrganizationMembershipCreateRequiresTeamIDs(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationMembershipCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-email=user@example.com"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-team-ids") {
		t.Fatalf("expected team-ids error, got %q", out)
	}
}

func TestOrganizationMembershipCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationMembershipCreateCommand{
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

func TestOrganizationMembershipCreateHelp(t *testing.T) {
	cmd := &OrganizationMembershipCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organizationmembership create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-email") {
		t.Error("Help should mention -email flag")
	}
	if !strings.Contains(help, "-team-ids") {
		t.Error("Help should mention -team-ids flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestOrganizationMembershipCreateSynopsis(t *testing.T) {
	cmd := &OrganizationMembershipCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Invite a user to join an organization" {
		t.Errorf("expected 'Invite a user to join an organization', got %q", synopsis)
	}
}

func TestOrganizationMembershipCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedOrg     string
		expectedEmail   string
		expectedTeamIDs string
		expectedFormat  string
	}{
		{
			name:            "org, email, and single team, default format",
			args:            []string{"-organization=my-org", "-email=user@example.com", "-team-ids=team-123"},
			expectedOrg:     "my-org",
			expectedEmail:   "user@example.com",
			expectedTeamIDs: "team-123",
			expectedFormat:  "table",
		},
		{
			name:            "org alias, email, and multiple teams",
			args:            []string{"-org=my-org", "-email=admin@example.com", "-team-ids=team-123,team-456"},
			expectedOrg:     "my-org",
			expectedEmail:   "admin@example.com",
			expectedTeamIDs: "team-123,team-456",
			expectedFormat:  "table",
		},
		{
			name:            "org, email, team, json format",
			args:            []string{"-org=prod-org", "-email=dev@example.com", "-team-ids=team-789", "-output=json"},
			expectedOrg:     "prod-org",
			expectedEmail:   "dev@example.com",
			expectedTeamIDs: "team-789",
			expectedFormat:  "json",
		},
		{
			name:            "org, email with multiple teams, json format",
			args:            []string{"-org=test-org", "-email=security@example.com", "-team-ids=team-abc,team-def,team-ghi", "-output=json"},
			expectedOrg:     "test-org",
			expectedEmail:   "security@example.com",
			expectedTeamIDs: "team-abc,team-def,team-ghi",
			expectedFormat:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationMembershipCreateCommand{}

			flags := cmd.Meta.FlagSet("organizationmembership create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.email, "email", "", "Email address of user to invite (required)")
			flags.StringVar(&cmd.teamIDs, "team-ids", "", "Comma-separated list of team IDs (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the email was set correctly
			if cmd.email != tt.expectedEmail {
				t.Errorf("expected email %q, got %q", tt.expectedEmail, cmd.email)
			}

			// Verify the team IDs were set correctly
			if cmd.teamIDs != tt.expectedTeamIDs {
				t.Errorf("expected team-ids %q, got %q", tt.expectedTeamIDs, cmd.teamIDs)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
