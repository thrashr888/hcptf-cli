package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestTeamCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamCreateCommand{
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

func TestTeamCreateRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamCreateCommand{
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

func TestTeamCreateRequiresBothFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamCreateCommand{
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

func TestTeamCreateHelp(t *testing.T) {
	cmd := &TeamCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf team create") {
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

func TestTeamCreateSynopsis(t *testing.T) {
	cmd := &TeamCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new team" {
		t.Errorf("expected 'Create a new team', got %q", synopsis)
	}
}

func TestTeamCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		expectedOrg        string
		expectedName       string
		expectedVisibility string
		expectedFmt        string
	}{
		{
			name:               "org and name, default visibility and format",
			args:               []string{"-organization=my-org", "-name=developers"},
			expectedOrg:        "my-org",
			expectedName:       "developers",
			expectedVisibility: "secret",
			expectedFmt:        "table",
		},
		{
			name:               "org alias and name with organization visibility",
			args:               []string{"-org=my-org", "-name=admins", "-visibility=organization"},
			expectedOrg:        "my-org",
			expectedName:       "admins",
			expectedVisibility: "organization",
			expectedFmt:        "table",
		},
		{
			name:               "org, name, json format",
			args:               []string{"-org=prod-org", "-name=ops", "-output=json"},
			expectedOrg:        "prod-org",
			expectedName:       "ops",
			expectedVisibility: "secret",
			expectedFmt:        "json",
		},
		{
			name:               "org, name, visibility, json format",
			args:               []string{"-org=test-org", "-name=security", "-visibility=organization", "-output=json"},
			expectedOrg:        "test-org",
			expectedName:       "security",
			expectedVisibility: "organization",
			expectedFmt:        "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &TeamCreateCommand{}

			flags := cmd.Meta.FlagSet("team create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.name, "name", "", "Team name (required)")
			flags.StringVar(&cmd.visibility, "visibility", "secret", "Team visibility: secret or organization")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

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

			// Verify the visibility was set correctly
			if cmd.visibility != tt.expectedVisibility {
				t.Errorf("expected visibility %q, got %q", tt.expectedVisibility, cmd.visibility)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
