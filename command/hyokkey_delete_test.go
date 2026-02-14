package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestHYOKKeyDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKKeyDeleteCommand{
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

func TestHYOKKeyDeleteHelp(t *testing.T) {
	cmd := &HYOKKeyDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf hyokkey delete") {
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
	if !strings.Contains(help, "HYOK") {
		t.Error("Help should explain HYOK feature")
	}
	if !strings.Contains(help, "Revoke") {
		t.Error("Help should mention revocation")
	}
}

func TestHYOKKeyDeleteSynopsis(t *testing.T) {
	cmd := &HYOKKeyDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Revoke a HYOK customer key version" {
		t.Errorf("expected 'Revoke a HYOK customer key version', got %q", synopsis)
	}
}

func TestHYOKKeyDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "required flags without force",
			args:          []string{"-id=keyv-123456"},
			expectedID:    "keyv-123456",
			expectedForce: false,
		},
		{
			name:          "required flags with force",
			args:          []string{"-id=keyv-abc123", "-force"},
			expectedID:    "keyv-abc123",
			expectedForce: true,
		},
		{
			name:          "different id without force",
			args:          []string{"-id=keyv-xyz789"},
			expectedID:    "keyv-xyz789",
			expectedForce: false,
		},
		{
			name:          "different id with force",
			args:          []string{"-id=keyv-test001", "-force"},
			expectedID:    "keyv-test001",
			expectedForce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &HYOKKeyDeleteCommand{}

			flags := cmd.Meta.FlagSet("hyokkey delete")
			flags.StringVar(&cmd.id, "id", "", "HYOK customer key version ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force revocation without confirmation")

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
