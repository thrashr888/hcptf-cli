package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestChangeRequestReadRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ChangeRequestReadCommand{
		Meta: newTestMeta(ui),
	}

	// Test missing id
	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing id, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestReadHelp(t *testing.T) {
	cmd := &ChangeRequestReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf changerequest read") {
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
	if !strings.Contains(help, "HCP Terraform Plus or Enterprise") {
		t.Error("Help should mention plan requirements")
	}
}

func TestChangeRequestReadSynopsis(t *testing.T) {
	cmd := &ChangeRequestReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show details of a specific change request" {
		t.Errorf("expected 'Show details of a specific change request', got %q", synopsis)
	}
}

func TestChangeRequestReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "id with default format",
			args:           []string{"-id=wscr-abc123"},
			expectedID:     "wscr-abc123",
			expectedFormat: "table",
		},
		{
			name:           "id with table format",
			args:           []string{"-id=wscr-xyz789", "-output=table"},
			expectedID:     "wscr-xyz789",
			expectedFormat: "table",
		},
		{
			name:           "id with json format",
			args:           []string{"-id=wscr-test456", "-output=json"},
			expectedID:     "wscr-test456",
			expectedFormat: "json",
		},
		{
			name:           "different id format",
			args:           []string{"-id=wscr-prod999"},
			expectedID:     "wscr-prod999",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ChangeRequestReadCommand{}

			flags := cmd.Meta.FlagSet("changerequest read")
			flags.StringVar(&cmd.id, "id", "", "Change request ID (required)")
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
