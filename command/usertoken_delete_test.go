package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestUserTokenDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &UserTokenDeleteCommand{
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

func TestUserTokenDeleteHelp(t *testing.T) {
	cmd := &UserTokenDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf usertoken delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestUserTokenDeleteSynopsis(t *testing.T) {
	cmd := &UserTokenDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a user API token" {
		t.Errorf("expected 'Delete a user API token', got %q", synopsis)
	}
}

func TestUserTokenDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "required flag, default force",
			args:          []string{"-id=at-abc123xyz"},
			expectedID:    "at-abc123xyz",
			expectedForce: false,
		},
		{
			name:          "required flag with force",
			args:          []string{"-id=at-def456uvw", "-force"},
			expectedID:    "at-def456uvw",
			expectedForce: true,
		},
		{
			name:          "required flag with explicit force true",
			args:          []string{"-id=at-ghi789rst", "-force=true"},
			expectedID:    "at-ghi789rst",
			expectedForce: true,
		},
		{
			name:          "required flag with explicit force false",
			args:          []string{"-id=at-jkl012opq", "-force=false"},
			expectedID:    "at-jkl012opq",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &UserTokenDeleteCommand{
				Meta: newTestMeta(cli.NewMockUi()),
			}

			flags := cmd.Meta.FlagSet("usertoken delete")
			flags.StringVar(&cmd.id, "id", "", "User token ID (required)")
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
