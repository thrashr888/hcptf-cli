package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestTeamShowRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamShowCommand{
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

func TestTeamShowRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamShowCommand{
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

func TestTeamShowRequiresBothFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamShowCommand{
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

func TestTeamShowHelp(t *testing.T) {
	cmd := &TeamShowCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf team show") {
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

func TestTeamShowSynopsis(t *testing.T) {
	cmd := &TeamShowCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show team details" {
		t.Errorf("expected 'Show team details', got %q", synopsis)
	}
}

func TestTeamShowFlagParsing(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedOrg  string
		expectedName string
		expectedFmt  string
	}{
		{
			name:         "org and name, default format",
			args:         []string{"-organization=my-org", "-name=developers"},
			expectedOrg:  "my-org",
			expectedName: "developers",
			expectedFmt:  "table",
		},
		{
			name:         "org alias and name",
			args:         []string{"-org=my-org", "-name=admins"},
			expectedOrg:  "my-org",
			expectedName: "admins",
			expectedFmt:  "table",
		},
		{
			name:         "org, name, table format",
			args:         []string{"-org=prod-org", "-name=ops", "-output=table"},
			expectedOrg:  "prod-org",
			expectedName: "ops",
			expectedFmt:  "table",
		},
		{
			name:         "org, name, json format",
			args:         []string{"-org=test-org", "-name=security", "-output=json"},
			expectedOrg:  "test-org",
			expectedName: "security",
			expectedFmt:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &TeamShowCommand{}

			flags := cmd.Meta.FlagSet("team show")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.name, "name", "", "Team name (required)")
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

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
