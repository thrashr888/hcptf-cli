package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestRunTaskReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RunTaskReadCommand{
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

func TestRunTaskReadHelp(t *testing.T) {
	cmd := &RunTaskReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf runtask read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestRunTaskReadSynopsis(t *testing.T) {
	cmd := &RunTaskReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read run task details by ID" {
		t.Errorf("expected 'Read run task details by ID', got %q", synopsis)
	}
}

func TestRunTaskReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "id only, default format",
			args:           []string{"-id=task-ABC123"},
			expectedID:     "task-ABC123",
			expectedFormat: "table",
		},
		{
			name:           "id with table format",
			args:           []string{"-id=task-XYZ789", "-output=table"},
			expectedID:     "task-XYZ789",
			expectedFormat: "table",
		},
		{
			name:           "id with json format",
			args:           []string{"-id=task-DEF456", "-output=json"},
			expectedID:     "task-DEF456",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &RunTaskReadCommand{}

			flags := cmd.Meta.FlagSet("runtask read")
			flags.StringVar(&cmd.id, "id", "", "Run task ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
