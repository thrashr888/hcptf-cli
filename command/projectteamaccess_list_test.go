package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestProjectTeamAccessListHelp(t *testing.T) {
	cmd := &ProjectTeamAccessListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf projectteamaccess list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-project-id") {
		t.Error("Help should mention -project-id flag")
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
}

func TestProjectTeamAccessListSynopsis(t *testing.T) {
	cmd := &ProjectTeamAccessListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List team access for a project" {
		t.Errorf("expected 'List team access for a project', got %q", synopsis)
	}
}

func TestProjectTeamAccessListFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedProjID string
		expectedOutput string
	}{
		{
			name:           "required flags only",
			args:           []string{"-project-id=prj-123"},
			expectedProjID: "prj-123",
			expectedOutput: "table",
		},
		{
			name:           "with json output",
			args:           []string{"-project-id=prj-456", "-output=json"},
			expectedProjID: "prj-456",
			expectedOutput: "json",
		},
		{
			name:           "with table output",
			args:           []string{"-project-id=prj-abc", "-output=table"},
			expectedProjID: "prj-abc",
			expectedOutput: "table",
		},
		{
			name:           "different project id",
			args:           []string{"-project-id=prj-xyz123"},
			expectedProjID: "prj-xyz123",
			expectedOutput: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ProjectTeamAccessListCommand{}

			flags := cmd.Meta.FlagSet("projectteamaccess list")
			flags.StringVar(&cmd.projectID, "project-id", "", "Project ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify flags were set correctly
			if cmd.projectID != tt.expectedProjID {
				t.Errorf("expected project-id %q, got %q", tt.expectedProjID, cmd.projectID)
			}
			if cmd.format != tt.expectedOutput {
				t.Errorf("expected output %q, got %q", tt.expectedOutput, cmd.format)
			}
		})
	}
}

func TestProjectTeamAccessListRequiresProjectID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectTeamAccessListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "project-id") {
		t.Fatalf("expected project-id error, got %q", out)
	}
}
