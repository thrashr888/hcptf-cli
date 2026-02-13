package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestChangeRequestListRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ChangeRequestListCommand{
		Meta: newTestMeta(ui),
	}

	// Test missing organization
	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error, got %q", ui.ErrorWriter.String())
	}

	// Test missing workspace
	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing workspace, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-workspace") {
		t.Fatalf("expected workspace error, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestListHelp(t *testing.T) {
	cmd := &ChangeRequestListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf changerequest list") {
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
	if !strings.Contains(help, "HCP Terraform Plus or Enterprise") {
		t.Error("Help should mention plan requirements")
	}
}

func TestChangeRequestListSynopsis(t *testing.T) {
	cmd := &ChangeRequestListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List change requests for a workspace" {
		t.Errorf("expected 'List change requests for a workspace', got %q", synopsis)
	}
}

func TestChangeRequestListFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedOrg      string
		expectedWorkspace string
		expectedFormat   string
	}{
		{
			name:             "organization and workspace with default format",
			args:             []string{"-organization=my-org", "-workspace=prod"},
			expectedOrg:      "my-org",
			expectedWorkspace: "prod",
			expectedFormat:   "table",
		},
		{
			name:             "org alias flag",
			args:             []string{"-org=test-org", "-workspace=staging"},
			expectedOrg:      "test-org",
			expectedWorkspace: "staging",
			expectedFormat:   "table",
		},
		{
			name:             "organization and workspace with table format",
			args:             []string{"-organization=my-org", "-workspace=dev", "-output=table"},
			expectedOrg:      "my-org",
			expectedWorkspace: "dev",
			expectedFormat:   "table",
		},
		{
			name:             "organization and workspace with json format",
			args:             []string{"-organization=acme", "-workspace=prod", "-output=json"},
			expectedOrg:      "acme",
			expectedWorkspace: "prod",
			expectedFormat:   "json",
		},
		{
			name:             "org alias with json format",
			args:             []string{"-org=test-org", "-workspace=qa", "-output=json"},
			expectedOrg:      "test-org",
			expectedWorkspace: "qa",
			expectedFormat:   "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ChangeRequestListCommand{}

			flags := cmd.Meta.FlagSet("changerequest list")
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
			if cmd.workspace != tt.expectedWorkspace {
				t.Errorf("expected workspace %q, got %q", tt.expectedWorkspace, cmd.workspace)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
