package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOAuthClientDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthClientDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-force"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestOAuthClientDeleteHelp(t *testing.T) {
	cmd := &OAuthClientDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf oauthclient delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
}

func TestOAuthClientDeleteSynopsis(t *testing.T) {
	cmd := &OAuthClientDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete an OAuth client" {
		t.Errorf("expected 'Delete an OAuth client', got %q", synopsis)
	}
}

func TestOAuthClientDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id only, no force",
			args:          []string{"-id=oc-XKFwG6ggfA9n7t1K"},
			expectedID:    "oc-XKFwG6ggfA9n7t1K",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=oc-ABC123XYZ456", "-force"},
			expectedID:    "oc-ABC123XYZ456",
			expectedForce: true,
		},
		{
			name:          "force flag with different id",
			args:          []string{"-force", "-id=oc-test99999"},
			expectedID:    "oc-test99999",
			expectedForce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OAuthClientDeleteCommand{}

			flags := cmd.Meta.FlagSet("oauthclient delete")
			flags.StringVar(&cmd.id, "id", "", "OAuth client ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
