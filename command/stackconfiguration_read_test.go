package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackConfigurationReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackConfigurationReadCommand{
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

func TestStackConfigurationReadHelp(t *testing.T) {
	cmd := &StackConfigurationReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stackconfiguration read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
}

func TestStackConfigurationReadSynopsis(t *testing.T) {
	cmd := &StackConfigurationReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read stack configuration details" {
		t.Errorf("expected 'Read stack configuration details', got %q", synopsis)
	}
}

func TestStackConfigurationReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "id, default format",
			args:           []string{"-id=stc-123"},
			expectedID:     "stc-123",
			expectedFormat: "table",
		},
		{
			name:           "id with table format",
			args:           []string{"-id=stc-456", "-output=table"},
			expectedID:     "stc-456",
			expectedFormat: "table",
		},
		{
			name:           "id with json format",
			args:           []string{"-id=stc-789", "-output=json"},
			expectedID:     "stc-789",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackConfigurationReadCommand{}

			flags := cmd.Meta.FlagSet("stackconfiguration read")
			flags.StringVar(&cmd.configID, "id", "", "Stack configuration ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the config ID was set correctly
			if cmd.configID != tt.expectedID {
				t.Errorf("expected config ID %q, got %q", tt.expectedID, cmd.configID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
