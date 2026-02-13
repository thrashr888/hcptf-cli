package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestNotificationCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-workspace=test-ws", "-name=test-notification", "-destination-type=slack"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestNotificationCreateRequiresWorkspace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-name=test-notification", "-destination-type=slack"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace") {
		t.Fatalf("expected workspace error, got %q", out)
	}
}

func TestNotificationCreateRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-workspace=test-ws", "-destination-type=slack"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestNotificationCreateRequiresDestinationType(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-workspace=test-ws", "-name=test-notification"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-destination-type") {
		t.Fatalf("expected destination-type error, got %q", out)
	}
}

func TestNotificationCreateInvalidDestinationType(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-workspace=test-ws", "-name=test-notification", "-destination-type=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "invalid destination-type") {
		t.Fatalf("expected invalid destination-type error, got %q", out)
	}
}

func TestNotificationCreateRequiresURLForNonEmailTypes(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-workspace=test-ws", "-name=test-notification", "-destination-type=slack"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-url") {
		t.Fatalf("expected url error, got %q", out)
	}
}

func TestNotificationCreateRequiresURLForGenericType(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-workspace=test-ws", "-name=test-notification", "-destination-type=generic"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-url") {
		t.Fatalf("expected url error, got %q", out)
	}
}

func TestNotificationCreateRequiresURLForMicrosoftTeamsType(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-workspace=test-ws", "-name=test-notification", "-destination-type=microsoft-teams"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-url") {
		t.Fatalf("expected url error, got %q", out)
	}
}

func TestNotificationCreateWithOrgAlias(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-org=test-org", "-name=test-notification", "-destination-type=slack"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace") {
		t.Fatalf("expected workspace error, got %q", out)
	}
}

func TestNotificationCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "required") {
		t.Fatalf("expected required error, got %q", out)
	}
}

func TestNotificationCreateHelp(t *testing.T) {
	cmd := &NotificationCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf notification create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-workspace") {
		t.Error("Help should mention -workspace flag")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "-destination-type") {
		t.Error("Help should mention -destination-type flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestNotificationCreateSynopsis(t *testing.T) {
	cmd := &NotificationCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new notification configuration" {
		t.Errorf("expected 'Create a new notification configuration', got %q", synopsis)
	}
}

func TestNotificationCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOrg    string
		expectedWS     string
		expectedName   string
		expectedDest   string
		expectedURL    string
		expectedFormat string
	}{
		{
			name:           "required flags only, default values",
			args:           []string{"-organization=my-org", "-workspace=my-ws", "-name=slack-notif", "-destination-type=slack"},
			expectedOrg:    "my-org",
			expectedWS:     "my-ws",
			expectedName:   "slack-notif",
			expectedDest:   "slack",
			expectedURL:    "",
			expectedFormat: "table",
		},
		{
			name:           "org alias with url",
			args:           []string{"-org=test-org", "-workspace=test-ws", "-name=test-notif", "-destination-type=generic", "-url=https://example.com"},
			expectedOrg:    "test-org",
			expectedWS:     "test-ws",
			expectedName:   "test-notif",
			expectedDest:   "generic",
			expectedURL:    "https://example.com",
			expectedFormat: "table",
		},
		{
			name:           "microsoft-teams with json output",
			args:           []string{"-org=prod-org", "-workspace=prod-ws", "-name=teams-notif", "-destination-type=microsoft-teams", "-url=https://teams.webhook.url", "-output=json"},
			expectedOrg:    "prod-org",
			expectedWS:     "prod-ws",
			expectedName:   "teams-notif",
			expectedDest:   "microsoft-teams",
			expectedURL:    "https://teams.webhook.url",
			expectedFormat: "json",
		},
		{
			name:           "email destination type",
			args:           []string{"-organization=dev-org", "-workspace=dev-ws", "-name=email-notif", "-destination-type=email"},
			expectedOrg:    "dev-org",
			expectedWS:     "dev-ws",
			expectedName:   "email-notif",
			expectedDest:   "email",
			expectedURL:    "",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &NotificationCreateCommand{}

			flags := cmd.Meta.FlagSet("notification create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.workspace, "workspace", "", "Workspace name (required)")
			flags.StringVar(&cmd.name, "name", "", "Notification configuration name (required)")
			flags.StringVar(&cmd.destinationType, "destination-type", "", "Destination type: email, slack, generic, microsoft-teams (required)")
			flags.BoolVar(&cmd.enabled, "enabled", true, "Enable notification configuration")
			flags.StringVar(&cmd.url, "url", "", "Webhook URL (required for slack, generic, microsoft-teams)")
			flags.StringVar(&cmd.token, "token", "", "Token for authentication (optional for generic)")
			flags.StringVar(&cmd.triggers, "triggers", "", "Comma-separated list of trigger types")
			flags.StringVar(&cmd.emailAddresses, "email-addresses", "", "Comma-separated list of email addresses (for email type, TFE only)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the workspace was set correctly
			if cmd.workspace != tt.expectedWS {
				t.Errorf("expected workspace %q, got %q", tt.expectedWS, cmd.workspace)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the destination type was set correctly
			if cmd.destinationType != tt.expectedDest {
				t.Errorf("expected destinationType %q, got %q", tt.expectedDest, cmd.destinationType)
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

