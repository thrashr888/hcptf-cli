package command

import (
	"strings"
	"testing"
)

func TestOrganizationListHelp(t *testing.T) {
	cmd := &OrganizationListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organization list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "List organizations") {
		t.Error("Help should describe what the command does")
	}
}

func TestOrganizationListSynopsis(t *testing.T) {
	cmd := &OrganizationListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List organizations" {
		t.Errorf("expected 'List organizations', got %q", synopsis)
	}
}

func TestOrganizationListFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedFormat string
	}{
		{"default format", []string{}, "table"},
		{"table format", []string{"-output=table"}, "table"},
		{"json format", []string{"-output=json"}, "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationListCommand{}

			flags := cmd.Meta.FlagSet("organization list")
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
