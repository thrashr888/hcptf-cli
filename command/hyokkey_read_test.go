package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestHYOKKeyReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKKeyReadCommand{
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

func TestHYOKKeyReadHelp(t *testing.T) {
	cmd := &HYOKKeyReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf hyokkey read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
	if !strings.Contains(help, "HYOK") {
		t.Error("Help should explain HYOK feature")
	}
}

func TestHYOKKeyReadSynopsis(t *testing.T) {
	cmd := &HYOKKeyReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show HYOK customer key version details" {
		t.Errorf("expected 'Show HYOK customer key version details', got %q", synopsis)
	}
}

func TestHYOKKeyReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "required flags with default format",
			args:           []string{"-id=keyv-123456"},
			expectedID:     "keyv-123456",
			expectedFormat: "table",
		},
		{
			name:           "required flags with json output",
			args:           []string{"-id=keyv-abc123", "-output=json"},
			expectedID:     "keyv-abc123",
			expectedFormat: "json",
		},
		{
			name:           "required flags with table output",
			args:           []string{"-id=keyv-xyz789", "-output=table"},
			expectedID:     "keyv-xyz789",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &HYOKKeyReadCommand{}

			flags := cmd.Meta.FlagSet("hyokkey read")
			flags.StringVar(&cmd.id, "id", "", "HYOK customer key version ID (required)")
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
