package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationShowRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationShowCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestOrganizationShowHelp(t *testing.T) {
	cmd := &OrganizationShowCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organization show") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -name is required")
	}
}

func TestOrganizationShowSynopsis(t *testing.T) {
	cmd := &OrganizationShowCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show organization details" {
		t.Errorf("expected 'Show organization details', got %q", synopsis)
	}
}

func TestOrganizationShowFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedName   string
		expectedFormat string
	}{
		{"name and default format", []string{"-name=test-org"}, "test-org", "table"},
		{"name and table format", []string{"-name=my-org", "-output=table"}, "my-org", "table"},
		{"name and json format", []string{"-name=prod-org", "-output=json"}, "prod-org", "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationShowCommand{}

			flags := cmd.Meta.FlagSet("organization show")
			flags.StringVar(&cmd.name, "name", "", "Organization name (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
