package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestTeamAddMemberRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamAddMemberCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-team=developers", "-username=alice"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestTeamAddMemberRequiresTeam(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamAddMemberCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-username=alice"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-team") {
		t.Fatalf("expected team error, got %q", out)
	}
}

func TestTeamAddMemberRequiresUsername(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamAddMemberCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-team=developers"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-username") {
		t.Fatalf("expected username error, got %q", out)
	}
}

func TestTeamAddMemberRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamAddMemberCommand{
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

func TestTeamAddMemberHelp(t *testing.T) {
	cmd := &TeamAddMemberCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf team add-member") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-team") {
		t.Error("Help should mention -team flag")
	}
	if !strings.Contains(help, "-username") {
		t.Error("Help should mention -username flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestTeamAddMemberSynopsis(t *testing.T) {
	cmd := &TeamAddMemberCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Add a member to a team" {
		t.Errorf("expected 'Add a member to a team', got %q", synopsis)
	}
}

func TestTeamAddMemberFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedOrg      string
		expectedTeam     string
		expectedUsername string
	}{
		{
			name:             "org, team, and username",
			args:             []string{"-organization=my-org", "-team=developers", "-username=alice"},
			expectedOrg:      "my-org",
			expectedTeam:     "developers",
			expectedUsername: "alice",
		},
		{
			name:             "org alias, team, and username",
			args:             []string{"-org=my-org", "-team=admins", "-username=bob"},
			expectedOrg:      "my-org",
			expectedTeam:     "admins",
			expectedUsername: "bob",
		},
		{
			name:             "multiple users to different teams",
			args:             []string{"-org=prod-org", "-team=ops", "-username=charlie"},
			expectedOrg:      "prod-org",
			expectedTeam:     "ops",
			expectedUsername: "charlie",
		},
		{
			name:             "test org with security team",
			args:             []string{"-org=test-org", "-team=security", "-username=david"},
			expectedOrg:      "test-org",
			expectedTeam:     "security",
			expectedUsername: "david",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &TeamAddMemberCommand{}

			flags := cmd.Meta.FlagSet("team add-member")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.teamName, "team", "", "Team name (required)")
			flags.StringVar(&cmd.username, "username", "", "Username to add (required)")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the team was set correctly
			if cmd.teamName != tt.expectedTeam {
				t.Errorf("expected team %q, got %q", tt.expectedTeam, cmd.teamName)
			}

			// Verify the username was set correctly
			if cmd.username != tt.expectedUsername {
				t.Errorf("expected username %q, got %q", tt.expectedUsername, cmd.username)
			}
		})
	}
}
