package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackConfigurationCreateRequiresStackID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackConfigurationCreateCommand{
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

func TestStackConfigurationCreateReturnsNotSupportedError(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackConfigurationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-stack-id=st-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "not directly supported") {
		t.Fatalf("expected not directly supported error, got %q", out)
	}
}

func TestStackConfigurationCreateHelp(t *testing.T) {
	cmd := &StackConfigurationCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stackconfiguration create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-stack-id") {
		t.Error("Help should mention -stack-id flag")
	}
	if !strings.Contains(help, "-speculative") {
		t.Error("Help should mention -speculative flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
}

func TestStackConfigurationCreateSynopsis(t *testing.T) {
	cmd := &StackConfigurationCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new stack configuration" {
		t.Errorf("expected 'Create a new stack configuration', got %q", synopsis)
	}
}

func TestStackConfigurationCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name                string
		args                []string
		expectedStackID     string
		expectedSpeculative bool
		expectedFormat      string
	}{
		{
			name:                "stack-id, default options",
			args:                []string{"-stack-id=st-123"},
			expectedStackID:     "st-123",
			expectedSpeculative: false,
			expectedFormat:      "table",
		},
		{
			name:                "stack-id with speculative flag",
			args:                []string{"-stack-id=st-456", "-speculative"},
			expectedStackID:     "st-456",
			expectedSpeculative: true,
			expectedFormat:      "table",
		},
		{
			name:                "stack-id with json output",
			args:                []string{"-stack-id=st-789", "-output=json"},
			expectedStackID:     "st-789",
			expectedSpeculative: false,
			expectedFormat:      "json",
		},
		{
			name:                "stack-id with speculative and json output",
			args:                []string{"-stack-id=st-abc", "-speculative", "-output=json"},
			expectedStackID:     "st-abc",
			expectedSpeculative: true,
			expectedFormat:      "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackConfigurationCreateCommand{}

			flags := cmd.Meta.FlagSet("stackconfiguration create")
			flags.StringVar(&cmd.stackID, "stack-id", "", "Stack ID (required)")
			flags.BoolVar(&cmd.speculative, "speculative", false, "Create a speculative configuration (plan-only)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the stack ID was set correctly
			if cmd.stackID != tt.expectedStackID {
				t.Errorf("expected stack ID %q, got %q", tt.expectedStackID, cmd.stackID)
			}

			// Verify the speculative flag was set correctly
			if cmd.speculative != tt.expectedSpeculative {
				t.Errorf("expected speculative %v, got %v", tt.expectedSpeculative, cmd.speculative)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}






