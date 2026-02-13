package command

import (
	"strings"
	"testing"
)

func TestAccountShowHelp(t *testing.T) {
	cmd := &AccountShowCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf account show") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "table") {
		t.Error("Help should mention table format")
	}
	if !strings.Contains(help, "json") {
		t.Error("Help should mention json format")
	}
	if !strings.Contains(help, "authentication") {
		t.Error("Help should mention authentication requirement")
	}
}

func TestAccountShowSynopsis(t *testing.T) {
	cmd := &AccountShowCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show current account details" {
		t.Errorf("expected 'Show current account details', got %q", synopsis)
	}
}

func TestAccountShowFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedFormat string
	}{
		{
			name:           "default format",
			args:           []string{},
			expectedFormat: "table",
		},
		{
			name:           "table format",
			args:           []string{"-output=table"},
			expectedFormat: "table",
		},
		{
			name:           "json format",
			args:           []string{"-output=json"},
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AccountShowCommand{}

			flags := cmd.Meta.FlagSet("account show")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
