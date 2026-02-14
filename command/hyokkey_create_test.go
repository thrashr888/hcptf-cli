package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestHYOKKeyCreateRequiresHyokConfigID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKKeyCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-hyok-config-id") {
		t.Fatalf("expected hyok-config-id error, got %q", out)
	}
}

func TestHYOKKeyCreateHelp(t *testing.T) {
	cmd := &HYOKKeyCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf hyokkey create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-hyok-config-id") {
		t.Error("Help should mention -hyok-config-id flag")
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

func TestHYOKKeyCreateSynopsis(t *testing.T) {
	cmd := &HYOKKeyCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Check for and register new HYOK customer key versions" {
		t.Errorf("expected 'Check for and register new HYOK customer key versions', got %q", synopsis)
	}
}

func TestHYOKKeyCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		expectedHyokConfig string
		expectedFormat     string
	}{
		{
			name:               "required flags with default format",
			args:               []string{"-hyok-config-id=hyokc-123456"},
			expectedHyokConfig: "hyokc-123456",
			expectedFormat:     "table",
		},
		{
			name:               "required flags with json output",
			args:               []string{"-hyok-config-id=hyokc-abc123", "-output=json"},
			expectedHyokConfig: "hyokc-abc123",
			expectedFormat:     "json",
		},
		{
			name:               "required flags with table output",
			args:               []string{"-hyok-config-id=hyokc-xyz789", "-output=table"},
			expectedHyokConfig: "hyokc-xyz789",
			expectedFormat:     "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &HYOKKeyCreateCommand{}

			flags := cmd.Meta.FlagSet("hyokkey create")
			flags.StringVar(&cmd.hyokConfigID, "hyok-config-id", "", "HYOK configuration ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the hyok-config-id was set correctly
			if cmd.hyokConfigID != tt.expectedHyokConfig {
				t.Errorf("expected hyok-config-id %q, got %q", tt.expectedHyokConfig, cmd.hyokConfigID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
