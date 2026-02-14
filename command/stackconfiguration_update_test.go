package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackConfigurationUpdateReturnsError(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackConfigurationUpdateCommand{
		Meta: newTestMeta(ui),
	}

	// Test with no arguments
	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "cannot be updated") {
		t.Fatalf("expected update error message, got %q", out)
	}
}

func TestStackConfigurationUpdateReturnsErrorWithID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackConfigurationUpdateCommand{
		Meta: newTestMeta(ui),
	}

	// Test with ID provided - should still error because updates are not supported
	code := cmd.Run([]string{"-id=stc-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	out := ui.ErrorWriter.String()
	if !strings.Contains(out, "cannot be updated") {
		t.Fatalf("expected update error message, got %q", out)
	}
	if !strings.Contains(out, "immutable") {
		t.Fatalf("expected immutable message, got %q", out)
	}
}

func TestStackConfigurationUpdateHelp(t *testing.T) {
	cmd := &StackConfigurationUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stackconfiguration update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "immutable") {
		t.Error("Help should explain configurations are immutable")
	}
	if !strings.Contains(help, "stackconfiguration create") {
		t.Error("Help should mention alternative commands")
	}
}

func TestStackConfigurationUpdateSynopsis(t *testing.T) {
	cmd := &StackConfigurationUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update stack configuration (not supported - configurations are immutable)" {
		t.Errorf("expected 'Update stack configuration (not supported - configurations are immutable)', got %q", synopsis)
	}
}

func TestStackConfigurationUpdateFlagParsing(t *testing.T) {
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
			name:           "id with json format",
			args:           []string{"-id=stc-456", "-output=json"},
			expectedID:     "stc-456",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackConfigurationUpdateCommand{}

			flags := cmd.Meta.FlagSet("stackconfiguration update")
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
