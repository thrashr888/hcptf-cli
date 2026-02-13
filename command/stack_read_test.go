package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackReadCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestStackReadHelp(t *testing.T) {
	cmd := &StackReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stack read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "Example:") {
		t.Error("Help should contain examples")
	}
}

func TestStackReadSynopsis(t *testing.T) {
	cmd := &StackReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read stack details" {
		t.Errorf("expected 'Read stack details', got %q", synopsis)
	}
}

func TestStackReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id, default format",
			args:        []string{"-id=st-abc123"},
			expectedID:  "st-abc123",
			expectedFmt: "table",
		},
		{
			name:        "id, table format",
			args:        []string{"-id=st-xyz789", "-output=table"},
			expectedID:  "st-xyz789",
			expectedFmt: "table",
		},
		{
			name:        "id, json format",
			args:        []string{"-id=st-test456", "-output=json"},
			expectedID:  "st-test456",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackReadCommand{}

			flags := cmd.Meta.FlagSet("stack read")
			flags.StringVar(&cmd.stackID, "id", "", "Stack ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.stackID != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.stackID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
