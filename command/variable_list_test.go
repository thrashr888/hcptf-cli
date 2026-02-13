package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVariableListRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableListCommand{
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

func TestVariableListRequiresWorkspace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace") {
		t.Fatalf("expected workspace error, got %q", out)
	}
}

func TestVariableListHelp(t *testing.T) {
	cmd := &VariableListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf variable list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-workspace") {
		t.Error("Help should mention -workspace flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "table") {
		t.Error("Help should mention table format")
	}
	if !strings.Contains(help, "json") {
		t.Error("Help should mention json format")
	}
	if !strings.Contains(help, "Example") {
		t.Error("Help should contain examples")
	}
}

func TestVariableListSynopsis(t *testing.T) {
	cmd := &VariableListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List variables for a workspace" {
		t.Errorf("expected 'List variables for a workspace', got %q", synopsis)
	}
}

func TestVariableListFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOrg string
		expectedWs  string
		expectedFmt string
	}{
		{
			name:        "organization and workspace, default format",
			args:        []string{"-organization=my-org", "-workspace=my-workspace"},
			expectedOrg: "my-org",
			expectedWs:  "my-workspace",
			expectedFmt: "table",
		},
		{
			name:        "org alias",
			args:        []string{"-org=test-org", "-workspace=prod"},
			expectedOrg: "test-org",
			expectedWs:  "prod",
			expectedFmt: "table",
		},
		{
			name:        "organization and workspace, table format",
			args:        []string{"-org=my-org", "-workspace=dev", "-output=table"},
			expectedOrg: "my-org",
			expectedWs:  "dev",
			expectedFmt: "table",
		},
		{
			name:        "organization and workspace, json format",
			args:        []string{"-org=prod-org", "-workspace=staging", "-output=json"},
			expectedOrg: "prod-org",
			expectedWs:  "staging",
			expectedFmt: "json",
		},
		{
			name:        "full flags with json output",
			args:        []string{"-organization=enterprise-org", "-workspace=production", "-output=json"},
			expectedOrg: "enterprise-org",
			expectedWs:  "production",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VariableListCommand{}

			flags := cmd.Meta.FlagSet("variable list")
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
			if cmd.workspace != tt.expectedWs {
				t.Errorf("expected workspace %q, got %q", tt.expectedWs, cmd.workspace)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
