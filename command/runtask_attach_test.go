package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestRunTaskAttachRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RunTaskAttachCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-workspace=test-ws", "-runtask-id=task-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestRunTaskAttachRequiresWorkspace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RunTaskAttachCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-runtask-id=task-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace") {
		t.Fatalf("expected workspace error, got %q", out)
	}
}

func TestRunTaskAttachRequiresRunTaskID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RunTaskAttachCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-workspace=test-ws"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-runtask-id") {
		t.Fatalf("expected runtask-id error, got %q", out)
	}
}

func TestRunTaskAttachHelp(t *testing.T) {
	cmd := &RunTaskAttachCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf runtask attach") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-workspace") {
		t.Error("Help should mention -workspace flag")
	}
	if !strings.Contains(help, "-runtask-id") {
		t.Error("Help should mention -runtask-id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestRunTaskAttachSynopsis(t *testing.T) {
	cmd := &RunTaskAttachCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Attach a run task to a workspace" {
		t.Errorf("expected 'Attach a run task to a workspace', got %q", synopsis)
	}
}

func TestRunTaskAttachValidatesEnforcementLevel(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RunTaskAttachCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-workspace=test-ws", "-runtask-id=task-123", "-enforcement-level=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "enforcement-level") {
		t.Fatalf("expected enforcement-level validation error, got %q", out)
	}
}

func TestRunTaskAttachValidatesStage(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RunTaskAttachCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-workspace=test-ws", "-runtask-id=task-123", "-stage=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "stage") {
		t.Fatalf("expected stage validation error, got %q", out)
	}
}

func TestRunTaskAttachFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedOrg      string
		expectedWS       string
		expectedTaskID   string
		expectedEnforce  string
		expectedStage    string
		expectedFormat   string
	}{
		{
			name:             "required flags only, default values",
			args:             []string{"-organization=my-org", "-workspace=my-ws", "-runtask-id=task-123"},
			expectedOrg:      "my-org",
			expectedWS:       "my-ws",
			expectedTaskID:   "task-123",
			expectedEnforce:  "advisory",
			expectedStage:    "post_plan",
			expectedFormat:   "table",
		},
		{
			name:             "org alias with mandatory enforcement",
			args:             []string{"-org=test-org", "-workspace=test-ws", "-runtask-id=task-456", "-enforcement-level=mandatory"},
			expectedOrg:      "test-org",
			expectedWS:       "test-ws",
			expectedTaskID:   "task-456",
			expectedEnforce:  "mandatory",
			expectedStage:    "post_plan",
			expectedFormat:   "table",
		},
		{
			name:             "pre_plan stage with json output",
			args:             []string{"-org=prod-org", "-workspace=prod-ws", "-runtask-id=task-789", "-stage=pre_plan", "-output=json"},
			expectedOrg:      "prod-org",
			expectedWS:       "prod-ws",
			expectedTaskID:   "task-789",
			expectedEnforce:  "advisory",
			expectedStage:    "pre_plan",
			expectedFormat:   "json",
		},
		{
			name:             "all options specified",
			args:             []string{"-organization=dev-org", "-workspace=dev-ws", "-runtask-id=task-012", "-enforcement-level=mandatory", "-stage=pre_apply", "-output=json"},
			expectedOrg:      "dev-org",
			expectedWS:       "dev-ws",
			expectedTaskID:   "task-012",
			expectedEnforce:  "mandatory",
			expectedStage:    "pre_apply",
			expectedFormat:   "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &RunTaskAttachCommand{}

			flags := cmd.Meta.FlagSet("runtask attach")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.workspace, "workspace", "", "Workspace name (required)")
			flags.StringVar(&cmd.runTaskID, "runtask-id", "", "Run task ID (required)")
			flags.StringVar(&cmd.enforcementLevel, "enforcement-level", "advisory", "Enforcement level: advisory or mandatory (default: advisory)")
			flags.StringVar(&cmd.stage, "stage", "post_plan", "Stage: post_plan, pre_plan, or pre_apply (default: post_plan)")
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

			// Verify the runtask-id was set correctly
			if cmd.runTaskID != tt.expectedTaskID {
				t.Errorf("expected runTaskID %q, got %q", tt.expectedTaskID, cmd.runTaskID)
			}

			// Verify the enforcement level was set correctly
			if cmd.enforcementLevel != tt.expectedEnforce {
				t.Errorf("expected enforcementLevel %q, got %q", tt.expectedEnforce, cmd.enforcementLevel)
			}

			// Verify the stage was set correctly
			if cmd.stage != tt.expectedStage {
				t.Errorf("expected stage %q, got %q", tt.expectedStage, cmd.stage)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
