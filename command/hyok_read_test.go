package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestHYOKReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKReadCommand{
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

func TestHYOKReadHelp(t *testing.T) {
	cmd := &HYOKReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf hyok read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "HYOK") {
		t.Error("Help should mention HYOK")
	}
}

func TestHYOKReadSynopsis(t *testing.T) {
	cmd := &HYOKReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show HYOK configuration details" {
		t.Errorf("expected 'Show HYOK configuration details', got %q", synopsis)
	}
}

func TestHYOKReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "id, default format",
			args:           []string{"-id=hyokc-123456"},
			expectedID:     "hyokc-123456",
			expectedFormat: "table",
		},
		{
			name:           "id, table format",
			args:           []string{"-id=hyokc-abcdef", "-output=table"},
			expectedID:     "hyokc-abcdef",
			expectedFormat: "table",
		},
		{
			name:           "id, json format",
			args:           []string{"-id=hyokc-xyz789", "-output=json"},
			expectedID:     "hyokc-xyz789",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &HYOKReadCommand{}

			flags := cmd.Meta.FlagSet("hyok read")
			flags.StringVar(&cmd.id, "id", "", "HYOK configuration ID (required)")
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
