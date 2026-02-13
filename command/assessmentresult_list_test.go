package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAssessmentResultListRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AssessmentResultListCommand{
		Meta: newTestMeta(ui),
	}

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 org")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 workspace")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-workspace") {
		t.Fatalf("expected workspace error")
	}
}

func TestAssessmentResultListHelp(t *testing.T) {
	cmd := &AssessmentResultListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf assessmentresult list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
	}
	if !strings.Contains(help, "-workspace") {
		t.Error("Help should mention -workspace flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "Health assessments") {
		t.Error("Help should describe health assessments")
	}
}

func TestAssessmentResultListSynopsis(t *testing.T) {
	cmd := &AssessmentResultListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List health assessment results for a workspace" {
		t.Errorf("expected 'List health assessment results for a workspace', got %q", synopsis)
	}
}

func TestAssessmentResultListFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOrg string
		expectedWS  string
		expectedFmt string
	}{
		{
			name:        "organization and workspace, default format",
			args:        []string{"-organization=my-org", "-workspace=my-ws"},
			expectedOrg: "my-org",
			expectedWS:  "my-ws",
			expectedFmt: "table",
		},
		{
			name:        "org alias",
			args:        []string{"-org=test-org", "-workspace=test-ws"},
			expectedOrg: "test-org",
			expectedWS:  "test-ws",
			expectedFmt: "table",
		},
		{
			name:        "organization and workspace, table format",
			args:        []string{"-org=my-org", "-workspace=prod", "-output=table"},
			expectedOrg: "my-org",
			expectedWS:  "prod",
			expectedFmt: "table",
		},
		{
			name:        "organization and workspace, json format",
			args:        []string{"-org=prod-org", "-workspace=staging", "-output=json"},
			expectedOrg: "prod-org",
			expectedWS:  "staging",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AssessmentResultListCommand{}

			flags := cmd.Meta.FlagSet("assessmentresult list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.workspace, "workspace", "", "Workspace name (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the workspace was set correctly
			if cmd.workspace != tt.expectedWS {
				t.Errorf("expected workspace %q, got %q", tt.expectedWS, cmd.workspace)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
