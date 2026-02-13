package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackConfigurationListRequiresStackID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackConfigurationListCommand{
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

func TestStackConfigurationListHelp(t *testing.T) {
	cmd := &StackConfigurationListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stackconfiguration list") {
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

func TestStackConfigurationListSynopsis(t *testing.T) {
	cmd := &StackConfigurationListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List stack configurations for a stack" {
		t.Errorf("expected 'List stack configurations for a stack', got %q", synopsis)
	}
}

func TestStackConfigurationListFlagParsing(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedStackID string
		expectedFormat  string
	}{
		{
			name:            "stack-id, default format",
			args:            []string{"-stack-id=st-123"},
			expectedStackID: "st-123",
			expectedFormat:  "table",
		},
		{
			name:            "stack-id with table format",
			args:            []string{"-stack-id=st-456", "-output=table"},
			expectedStackID: "st-456",
			expectedFormat:  "table",
		},
		{
			name:            "stack-id with json format",
			args:            []string{"-stack-id=st-789", "-output=json"},
			expectedStackID: "st-789",
			expectedFormat:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackConfigurationListCommand{}

			flags := cmd.Meta.FlagSet("stackconfiguration list")
			flags.StringVar(&cmd.stackID, "stack-id", "", "Stack ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the stack ID was set correctly
			if cmd.stackID != tt.expectedStackID {
				t.Errorf("expected stack ID %q, got %q", tt.expectedStackID, cmd.stackID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
