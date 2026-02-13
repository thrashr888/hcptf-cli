package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStateListRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StateListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-workspace=test-ws"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestStateListRequiresWorkspace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StateListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace") {
		t.Fatalf("expected workspace error, got %q", out)
	}
}

func TestStateListRequiresBothFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StateListCommand{
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

func TestStateListHelp(t *testing.T) {
	cmd := &StateListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf state list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-workspace") {
		t.Error("Help should mention -workspace flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestStateListSynopsis(t *testing.T) {
	cmd := &StateListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if !strings.Contains(synopsis, "state") {
		t.Errorf("expected synopsis to mention 'state', got %q", synopsis)
	}
}

func TestStateListFlagParsing(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedOrg  string
		expectedWS   string
		expectedFmt  string
	}{
		{
			name:        "org and workspace, default format",
			args:        []string{"-organization=my-org", "-workspace=my-ws"},
			expectedOrg: "my-org",
			expectedWS:  "my-ws",
			expectedFmt: "table",
		},
		{
			name:        "org alias and workspace",
			args:        []string{"-org=my-org", "-workspace=my-ws"},
			expectedOrg: "my-org",
			expectedWS:  "my-ws",
			expectedFmt: "table",
		},
		{
			name:        "org, workspace, table format",
			args:        []string{"-org=test-org", "-workspace=test-ws", "-output=table"},
			expectedOrg: "test-org",
			expectedWS:  "test-ws",
			expectedFmt: "table",
		},
		{
			name:        "org, workspace, json format",
			args:        []string{"-org=prod-org", "-workspace=prod-ws", "-output=json"},
			expectedOrg: "prod-org",
			expectedWS:  "prod-ws",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StateListCommand{}

			flags := cmd.Meta.FlagSet("state list")
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
