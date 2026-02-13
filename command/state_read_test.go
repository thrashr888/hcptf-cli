package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStateReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StateReadCommand{
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

func TestStateReadHelp(t *testing.T) {
	cmd := &StateReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf state read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestStateReadSynopsis(t *testing.T) {
	cmd := &StateReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if !strings.Contains(synopsis, "state") {
		t.Errorf("expected synopsis to mention 'state', got %q", synopsis)
	}
}

func TestStateReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id with default format",
			args:        []string{"-id=sv-123"},
			expectedID:  "sv-123",
			expectedFmt: "table",
		},
		{
			name:        "id with table format",
			args:        []string{"-id=sv-456", "-output=table"},
			expectedID:  "sv-456",
			expectedFmt: "table",
		},
		{
			name:        "id with json format",
			args:        []string{"-id=sv-789", "-output=json"},
			expectedID:  "sv-789",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StateReadCommand{}

			flags := cmd.Meta.FlagSet("state read")
			flags.StringVar(&cmd.stateVersionID, "id", "", "State version ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the ID was set correctly
			if cmd.stateVersionID != tt.expectedID {
				t.Errorf("expected ID %q, got %q", tt.expectedID, cmd.stateVersionID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
