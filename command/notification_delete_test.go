package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestNotificationDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NotificationDeleteCommand{
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

func TestNotificationDeleteHelp(t *testing.T) {
	cmd := &NotificationDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf notification delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestNotificationDeleteSynopsis(t *testing.T) {
	cmd := &NotificationDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a notification configuration" {
		t.Errorf("expected 'Delete a notification configuration', got %q", synopsis)
	}
}

func TestNotificationDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id only, no force",
			args:          []string{"-id=nc-ABC123"},
			expectedID:    "nc-ABC123",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=nc-XYZ789", "-force"},
			expectedID:    "nc-XYZ789",
			expectedForce: true,
		},
		{
			name:          "different id without force",
			args:          []string{"-id=nc-TEST456"},
			expectedID:    "nc-TEST456",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &NotificationDeleteCommand{}

			flags := cmd.Meta.FlagSet("notification delete")
			flags.StringVar(&cmd.id, "id", "", "Notification configuration ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the force was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
