package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackStateListRequiresStackID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackStateListCommand{
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

func TestStackStateListHelp(t *testing.T) {
	cmd := &StackStateListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stackstate list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-stack-id") {
		t.Error("Help should mention -stack-id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
}

func TestStackStateListSynopsis(t *testing.T) {
	cmd := &StackStateListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List state versions for a stack" {
		t.Errorf("expected 'List state versions for a stack', got %q", synopsis)
	}
}

func TestStackStateListFlagParsing(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedStackID string
		expectedFmt     string
	}{
		{
			name:            "stack-id flag only",
			args:            []string{"-stack-id=st-abc123"},
			expectedStackID: "st-abc123",
			expectedFmt:     "table",
		},
		{
			name:            "with json output",
			args:            []string{"-stack-id=st-xyz789", "-output=json"},
			expectedStackID: "st-xyz789",
			expectedFmt:     "json",
		},
		{
			name:            "with table output",
			args:            []string{"-stack-id=st-def456", "-output=table"},
			expectedStackID: "st-def456",
			expectedFmt:     "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackStateListCommand{}

			flags := cmd.Meta.FlagSet("stackstate list")
			flags.StringVar(&cmd.stackID, "stack-id", "", "Stack ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the stack ID was set correctly
			if cmd.stackID != tt.expectedStackID {
				t.Errorf("expected stack ID %q, got %q", tt.expectedStackID, cmd.stackID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
