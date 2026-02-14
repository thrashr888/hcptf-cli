package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackDeploymentListRequiresStackID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackDeploymentListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-stack-id") {
		t.Fatalf("expected stack-id error, got %q", out)
	}
}

func TestStackDeploymentListHelp(t *testing.T) {
	cmd := &StackDeploymentListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stackdeployment list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-stack-id") {
		t.Error("Help should mention -stack-id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -stack-id is required")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "Example:") {
		t.Error("Help should contain examples")
	}
}

func TestStackDeploymentListSynopsis(t *testing.T) {
	cmd := &StackDeploymentListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List deployments for a stack" {
		t.Errorf("expected 'List deployments for a stack', got %q", synopsis)
	}
}

func TestStackDeploymentListFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedStack string
		expectedFmt   string
	}{
		{
			name:          "stack-id, default format",
			args:          []string{"-stack-id=st-abc123"},
			expectedStack: "st-abc123",
			expectedFmt:   "table",
		},
		{
			name:          "stack-id, table format",
			args:          []string{"-stack-id=st-xyz789", "-output=table"},
			expectedStack: "st-xyz789",
			expectedFmt:   "table",
		},
		{
			name:          "stack-id, json format",
			args:          []string{"-stack-id=st-test456", "-output=json"},
			expectedStack: "st-test456",
			expectedFmt:   "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackDeploymentListCommand{}

			flags := cmd.Meta.FlagSet("stackdeployment list")
			flags.StringVar(&cmd.stackID, "stack-id", "", "Stack ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the stack-id was set correctly
			if cmd.stackID != tt.expectedStack {
				t.Errorf("expected stack-id %q, got %q", tt.expectedStack, cmd.stackID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
