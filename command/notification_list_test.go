package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestNotificationListRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-workspace=test-workspace"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestNotificationListRequiresWorkspace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace") {
		t.Fatalf("expected workspace error, got %q", out)
	}
}

func TestNotificationListHelp(t *testing.T) {
	cmd := &NotificationListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf notification list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-workspace") {
		t.Error("Help should mention -workspace flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestNotificationListSynopsis(t *testing.T) {
	cmd := &NotificationListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List notification configurations for a workspace" {
		t.Errorf("expected 'List notification configurations for a workspace', got %q", synopsis)
	}
}

func TestNotificationListFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOrg    string
		expectedWS     string
		expectedFormat string
	}{
		{
			name:           "required flags only, default format",
			args:           []string{"-organization=my-org", "-workspace=my-workspace"},
			expectedOrg:    "my-org",
			expectedWS:     "my-workspace",
			expectedFormat: "table",
		},
		{
			name:           "org alias with workspace",
			args:           []string{"-org=test-org", "-workspace=prod-workspace"},
			expectedOrg:    "test-org",
			expectedWS:     "prod-workspace",
			expectedFormat: "table",
		},
		{
			name:           "with json output",
			args:           []string{"-org=prod-org", "-workspace=prod-ws", "-output=json"},
			expectedOrg:    "prod-org",
			expectedWS:     "prod-ws",
			expectedFormat: "json",
		},
		{
			name:           "with table output",
			args:           []string{"-organization=dev-org", "-workspace=dev-ws", "-output=table"},
			expectedOrg:    "dev-org",
			expectedWS:     "dev-ws",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &NotificationListCommand{}

			flags := cmd.Meta.FlagSet("notification list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.workspace, "workspace", "", "Workspace name (required)")
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

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
