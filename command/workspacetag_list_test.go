package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestWorkspaceTagListRequiresWorkspaceID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &WorkspaceTagListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace-id") {
		t.Fatalf("expected workspace-id error, got %q", out)
	}
}

func TestWorkspaceTagListHelp(t *testing.T) {
	cmd := &WorkspaceTagListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf workspacetag list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-workspace-id") {
		t.Error("Help should mention -workspace-id flag")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag alias")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
}

func TestWorkspaceTagListSynopsis(t *testing.T) {
	cmd := &WorkspaceTagListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List tags for a workspace" {
		t.Errorf("expected 'List tags for a workspace', got %q", synopsis)
	}
}

func TestWorkspaceTagListFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedWsID   string
		expectedFormat string
	}{
		{
			name:           "workspace-id, default format",
			args:           []string{"-workspace-id=ws-123"},
			expectedWsID:   "ws-123",
			expectedFormat: "table",
		},
		{
			name:           "id alias flag",
			args:           []string{"-id=ws-456"},
			expectedWsID:   "ws-456",
			expectedFormat: "table",
		},
		{
			name:           "workspace-id with table format",
			args:           []string{"-workspace-id=ws-789", "-output=table"},
			expectedWsID:   "ws-789",
			expectedFormat: "table",
		},
		{
			name:           "workspace-id with json format",
			args:           []string{"-workspace-id=ws-abc", "-output=json"},
			expectedWsID:   "ws-abc",
			expectedFormat: "json",
		},
		{
			name:           "id alias with json format",
			args:           []string{"-id=ws-xyz", "-output=json"},
			expectedWsID:   "ws-xyz",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &WorkspaceTagListCommand{}

			flags := cmd.Meta.FlagSet("workspacetag list")
			flags.StringVar(&cmd.workspaceID, "workspace-id", "", "Workspace ID (required)")
			flags.StringVar(&cmd.workspaceID, "id", "", "Workspace ID (alias)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the workspace ID was set correctly
			if cmd.workspaceID != tt.expectedWsID {
				t.Errorf("expected workspace ID %q, got %q", tt.expectedWsID, cmd.workspaceID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
