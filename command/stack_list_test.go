package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackListRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackListCommand{
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

func TestStackListHelp(t *testing.T) {
	cmd := &StackListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stack list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -organization is required")
	}
	if !strings.Contains(help, "-project") {
		t.Error("Help should mention -project flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
}

func TestStackListSynopsis(t *testing.T) {
	cmd := &StackListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List stacks in an organization or project" {
		t.Errorf("expected 'List stacks in an organization or project', got %q", synopsis)
	}
}

func TestStackListFlagParsing(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedOrg     string
		expectedProject string
		expectedFmt     string
	}{
		{
			name:            "organization, default format",
			args:            []string{"-organization=my-org"},
			expectedOrg:     "my-org",
			expectedProject: "",
			expectedFmt:     "table",
		},
		{
			name:            "org alias",
			args:            []string{"-org=test-org"},
			expectedOrg:     "test-org",
			expectedProject: "",
			expectedFmt:     "table",
		},
		{
			name:            "organization with project filter",
			args:            []string{"-org=my-org", "-project=prj-123"},
			expectedOrg:     "my-org",
			expectedProject: "prj-123",
			expectedFmt:     "table",
		},
		{
			name:            "organization, table format",
			args:            []string{"-org=my-org", "-output=table"},
			expectedOrg:     "my-org",
			expectedProject: "",
			expectedFmt:     "table",
		},
		{
			name:            "organization, json format",
			args:            []string{"-org=prod-org", "-output=json"},
			expectedOrg:     "prod-org",
			expectedProject: "",
			expectedFmt:     "json",
		},
		{
			name:            "all flags together",
			args:            []string{"-org=my-org", "-project=prj-abc", "-output=json"},
			expectedOrg:     "my-org",
			expectedProject: "prj-abc",
			expectedFmt:     "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackListCommand{}

			flags := cmd.Meta.FlagSet("stack list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.project, "project", "", "Filter by project ID")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the project was set correctly
			if cmd.project != tt.expectedProject {
				t.Errorf("expected project %q, got %q", tt.expectedProject, cmd.project)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
