package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicyCheckListRequiresRunID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicyCheckListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-run-id") {
		t.Fatalf("expected run-id error, got %q", out)
	}
}

func TestPolicyCheckListHelp(t *testing.T) {
	cmd := &PolicyCheckListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policycheck list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-run-id") {
		t.Error("Help should mention -run-id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -run-id is required")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
}

func TestPolicyCheckListSynopsis(t *testing.T) {
	cmd := &PolicyCheckListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List policy checks for a run" {
		t.Errorf("expected 'List policy checks for a run', got %q", synopsis)
	}
}

func TestPolicyCheckListFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedRunID  string
		expectedFormat string
	}{
		{
			name:           "run-id with default format",
			args:           []string{"-run-id=run-abc123"},
			expectedRunID:  "run-abc123",
			expectedFormat: "table",
		},
		{
			name:           "run-id with table format",
			args:           []string{"-run-id=run-xyz789", "-output=table"},
			expectedRunID:  "run-xyz789",
			expectedFormat: "table",
		},
		{
			name:           "run-id with json format",
			args:           []string{"-run-id=run-def456", "-output=json"},
			expectedRunID:  "run-def456",
			expectedFormat: "json",
		},
		{
			name:           "different run-id format",
			args:           []string{"-run-id=run-prod123", "-output=table"},
			expectedRunID:  "run-prod123",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicyCheckListCommand{}

			flags := cmd.Meta.FlagSet("policycheck list")
			flags.StringVar(&cmd.runID, "run-id", "", "Run ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the run ID was set correctly
			if cmd.runID != tt.expectedRunID {
				t.Errorf("expected runID %q, got %q", tt.expectedRunID, cmd.runID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
