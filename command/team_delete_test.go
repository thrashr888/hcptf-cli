package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestTeamDeleteRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=test-team"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestTeamDeleteRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestTeamDeleteRequiresBothFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamDeleteCommand{
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

func TestTeamDeleteHelp(t *testing.T) {
	cmd := &TeamDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf team delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestTeamDeleteSynopsis(t *testing.T) {
	cmd := &TeamDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a team" {
		t.Errorf("expected 'Delete a team', got %q", synopsis)
	}
}

func TestTeamDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedOrg   string
		expectedName  string
		expectedForce bool
	}{
		{
			name:          "org and name, no force",
			args:          []string{"-organization=my-org", "-name=old-team"},
			expectedOrg:   "my-org",
			expectedName:  "old-team",
			expectedForce: false,
		},
		{
			name:          "org alias and name",
			args:          []string{"-org=my-org", "-name=deprecated"},
			expectedOrg:   "my-org",
			expectedName:  "deprecated",
			expectedForce: false,
		},
		{
			name:          "org, name with force",
			args:          []string{"-org=prod-org", "-name=legacy", "-force"},
			expectedOrg:   "prod-org",
			expectedName:  "legacy",
			expectedForce: true,
		},
		{
			name:          "org, name with explicit force=true",
			args:          []string{"-org=test-org", "-name=temporary", "-force=true"},
			expectedOrg:   "test-org",
			expectedName:  "temporary",
			expectedForce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &TeamDeleteCommand{}

			flags := cmd.Meta.FlagSet("team delete")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.name, "name", "", "Team name (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
