package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestRunTaskDetachRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RunTaskDetachCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-workspace=test-ws", "-workspace-runtask-id=wr-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestRunTaskDetachRequiresWorkspace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RunTaskDetachCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-workspace-runtask-id=wr-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace") {
		t.Fatalf("expected workspace error, got %q", out)
	}
}

func TestRunTaskDetachRequiresWorkspaceRunTaskID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RunTaskDetachCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-workspace=test-ws"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace-runtask-id") {
		t.Fatalf("expected workspace-runtask-id error, got %q", out)
	}
}

func TestRunTaskDetachHelp(t *testing.T) {
	cmd := &RunTaskDetachCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf runtask detach") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-workspace") {
		t.Error("Help should mention -workspace flag")
	}
	if !strings.Contains(help, "-workspace-runtask-id") {
		t.Error("Help should mention -workspace-runtask-id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestRunTaskDetachSynopsis(t *testing.T) {
	cmd := &RunTaskDetachCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Detach a run task from a workspace" {
		t.Errorf("expected 'Detach a run task from a workspace', got %q", synopsis)
	}
}

func TestRunTaskDetachFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedOrg   string
		expectedWS    string
		expectedWRTID string
		expectedForce bool
	}{
		{
			name:          "required flags only, no force",
			args:          []string{"-organization=my-org", "-workspace=my-ws", "-workspace-runtask-id=wr-123"},
			expectedOrg:   "my-org",
			expectedWS:    "my-ws",
			expectedWRTID: "wr-123",
			expectedForce: false,
		},
		{
			name:          "org alias with force flag",
			args:          []string{"-org=test-org", "-workspace=test-ws", "-workspace-runtask-id=wr-456", "-force"},
			expectedOrg:   "test-org",
			expectedWS:    "test-ws",
			expectedWRTID: "wr-456",
			expectedForce: true,
		},
		{
			name:          "all options with force=true",
			args:          []string{"-organization=prod-org", "-workspace=prod-ws", "-workspace-runtask-id=wr-789", "-force=true"},
			expectedOrg:   "prod-org",
			expectedWS:    "prod-ws",
			expectedWRTID: "wr-789",
			expectedForce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &RunTaskDetachCommand{}

			flags := cmd.Meta.FlagSet("runtask detach")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.workspace, "workspace", "", "Workspace name (required)")
			flags.StringVar(&cmd.workspaceRunTaskID, "workspace-runtask-id", "", "Workspace run task ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force detach without confirmation")

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

			// Verify the workspace-runtask-id was set correctly
			if cmd.workspaceRunTaskID != tt.expectedWRTID {
				t.Errorf("expected workspaceRunTaskID %q, got %q", tt.expectedWRTID, cmd.workspaceRunTaskID)
			}

			// Verify the force was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
