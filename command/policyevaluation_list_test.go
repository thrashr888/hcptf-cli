package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicyEvaluationListRequiresTaskStageID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicyEvaluationListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-task-stage-id") {
		t.Fatalf("expected task-stage-id error, got %q", out)
	}
}

func TestPolicyEvaluationListHelp(t *testing.T) {
	cmd := &PolicyEvaluationListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policyevaluation list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-task-stage-id") {
		t.Error("Help should mention -task-stage-id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "policy evaluation") {
		t.Error("Help should describe policy evaluations")
	}
}

func TestPolicyEvaluationListSynopsis(t *testing.T) {
	cmd := &PolicyEvaluationListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List policy evaluations for a task stage" {
		t.Errorf("expected 'List policy evaluations for a task stage', got %q", synopsis)
	}
}

func TestPolicyEvaluationListFlagParsing(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedStageID   string
		expectedFmt       string
	}{
		{
			name:            "task-stage-id flag",
			args:            []string{"-task-stage-id=ts-abc123"},
			expectedStageID: "ts-abc123",
			expectedFmt:     "table",
		},
		{
			name:            "with json output",
			args:            []string{"-task-stage-id=ts-xyz789", "-output=json"},
			expectedStageID: "ts-xyz789",
			expectedFmt:     "json",
		},
		{
			name:            "with table output",
			args:            []string{"-task-stage-id=ts-test456", "-output=table"},
			expectedStageID: "ts-test456",
			expectedFmt:     "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicyEvaluationListCommand{}

			flags := cmd.Meta.FlagSet("policyevaluation list")
			flags.StringVar(&cmd.taskStageID, "task-stage-id", "", "Task Stage ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the task stage ID was set correctly
			if cmd.taskStageID != tt.expectedStageID {
				t.Errorf("expected task-stage-id %q, got %q", tt.expectedStageID, cmd.taskStageID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
