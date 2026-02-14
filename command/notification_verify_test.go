package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestNotificationVerifyRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationVerifyCommand{
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

func TestNotificationVerifyHelp(t *testing.T) {
	cmd := &NotificationVerifyCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf notification verify") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestNotificationVerifySynopsis(t *testing.T) {
	cmd := &NotificationVerifyCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Verify a notification configuration" {
		t.Errorf("expected 'Verify a notification configuration', got %q", synopsis)
	}
}

func TestNotificationVerifyFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "id only, default format",
			args:           []string{"-id=nc-ABC123"},
			expectedID:     "nc-ABC123",
			expectedFormat: "table",
		},
		{
			name:           "id with table output",
			args:           []string{"-id=nc-XYZ789", "-output=table"},
			expectedID:     "nc-XYZ789",
			expectedFormat: "table",
		},
		{
			name:           "id with json output",
			args:           []string{"-id=nc-TEST456", "-output=json"},
			expectedID:     "nc-TEST456",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &NotificationVerifyCommand{}

			flags := cmd.Meta.FlagSet("notification verify")
			flags.StringVar(&cmd.id, "id", "", "Notification configuration ID (required)")
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
