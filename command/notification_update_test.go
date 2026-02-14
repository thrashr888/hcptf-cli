package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestNotificationUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=updated-name"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestNotificationUpdateInvalidEnabledValue(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=nc-123", "-enabled=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "must be 'true' or 'false'") {
		t.Fatalf("expected enabled validation error, got %q", out)
	}
}

func TestNotificationUpdateHelp(t *testing.T) {
	cmd := &NotificationUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf notification update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestNotificationUpdateSynopsis(t *testing.T) {
	cmd := &NotificationUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update notification configuration settings" {
		t.Errorf("expected 'Update notification configuration settings', got %q", synopsis)
	}
}

func TestNotificationUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedName   string
		expectedURL    string
		expectedFormat string
	}{
		{
			name:           "id only",
			args:           []string{"-id=nc-ABC123"},
			expectedID:     "nc-ABC123",
			expectedName:   "",
			expectedURL:    "",
			expectedFormat: "table",
		},
		{
			name:           "id with name",
			args:           []string{"-id=nc-XYZ789", "-name=updated-notification"},
			expectedID:     "nc-XYZ789",
			expectedName:   "updated-notification",
			expectedURL:    "",
			expectedFormat: "table",
		},
		{
			name:           "id with url and json output",
			args:           []string{"-id=nc-TEST456", "-url=https://example.com/webhook", "-output=json"},
			expectedID:     "nc-TEST456",
			expectedName:   "",
			expectedURL:    "https://example.com/webhook",
			expectedFormat: "json",
		},
		{
			name:           "id with name and url",
			args:           []string{"-id=nc-PROD999", "-name=prod-notif", "-url=https://slack.webhook.url"},
			expectedID:     "nc-PROD999",
			expectedName:   "prod-notif",
			expectedURL:    "https://slack.webhook.url",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &NotificationUpdateCommand{}

			flags := cmd.Meta.FlagSet("notification update")
			flags.StringVar(&cmd.id, "id", "", "Notification configuration ID (required)")
			flags.StringVar(&cmd.name, "name", "", "Notification configuration name")
			flags.StringVar(&cmd.enabled, "enabled", "", "Enable notification configuration (true/false)")
			flags.StringVar(&cmd.url, "url", "", "Webhook URL")
			flags.StringVar(&cmd.token, "token", "", "Token for authentication")
			flags.StringVar(&cmd.triggers, "triggers", "", "Comma-separated list of trigger types")
			flags.StringVar(&cmd.emailAddresses, "email-addresses", "", "Comma-separated list of email addresses (TFE only)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the url was set correctly
			if cmd.url != tt.expectedURL {
				t.Errorf("expected url %q, got %q", tt.expectedURL, cmd.url)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
