package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestWorkspaceTagAddRequiresWorkspaceID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &WorkspaceTagAddCommand{
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

func TestWorkspaceTagAddRequiresTags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &WorkspaceTagAddCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-workspace-id=ws-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-tags") {
		t.Fatalf("expected tags error, got %q", out)
	}
}

func TestWorkspaceTagAddRequiresEmptyWorkspaceID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &WorkspaceTagAddCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-workspace-id=", "-tags=prod"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace-id") {
		t.Fatalf("expected workspace-id error, got %q", out)
	}
}

func TestWorkspaceTagAddRequiresEmptyTags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &WorkspaceTagAddCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-workspace-id=ws-123", "-tags="})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-tags") {
		t.Fatalf("expected tags error, got %q", out)
	}
}

func TestWorkspaceTagAddHelp(t *testing.T) {
	cmd := &WorkspaceTagAddCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf workspacetag add") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-workspace-id") {
		t.Error("Help should mention -workspace-id flag")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag alias")
	}
	if !strings.Contains(help, "-tags") {
		t.Error("Help should mention -tags flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
}

func TestWorkspaceTagAddSynopsis(t *testing.T) {
	cmd := &WorkspaceTagAddCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Add tags to a workspace" {
		t.Errorf("expected 'Add tags to a workspace', got %q", synopsis)
	}
}

func TestWorkspaceTagAddFlagParsing(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedWsID string
		expectedTags string
	}{
		{
			name:         "workspace-id flag",
			args:         []string{"-workspace-id=ws-123", "-tags=prod"},
			expectedWsID: "ws-123",
			expectedTags: "prod",
		},
		{
			name:         "id alias flag",
			args:         []string{"-id=ws-456", "-tags=staging"},
			expectedWsID: "ws-456",
			expectedTags: "staging",
		},
		{
			name:         "multiple tags comma-separated",
			args:         []string{"-workspace-id=ws-789", "-tags=prod,us-west-2,team-a"},
			expectedWsID: "ws-789",
			expectedTags: "prod,us-west-2,team-a",
		},
		{
			name:         "tags with spaces",
			args:         []string{"-id=ws-abc", "-tags=prod, staging, dev"},
			expectedWsID: "ws-abc",
			expectedTags: "prod, staging, dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &WorkspaceTagAddCommand{}

			flags := cmd.Meta.FlagSet("workspacetag add")
			flags.StringVar(&cmd.workspaceID, "workspace-id", "", "Workspace ID (required)")
			flags.StringVar(&cmd.workspaceID, "id", "", "Workspace ID (alias)")
			flags.StringVar(&cmd.tags, "tags", "", "Comma-separated list of tag names to add (required)")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the workspace ID was set correctly
			if cmd.workspaceID != tt.expectedWsID {
				t.Errorf("expected workspace ID %q, got %q", tt.expectedWsID, cmd.workspaceID)
			}

			// Verify the tags were set correctly
			if cmd.tags != tt.expectedTags {
				t.Errorf("expected tags %q, got %q", tt.expectedTags, cmd.tags)
			}
		})
	}
}
