package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestSSHKeyDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &SSHKeyDeleteCommand{
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

func TestSSHKeyDeleteHelp(t *testing.T) {
	cmd := &SSHKeyDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf sshkey delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestSSHKeyDeleteSynopsis(t *testing.T) {
	cmd := &SSHKeyDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete an SSH key" {
		t.Errorf("expected 'Delete an SSH key', got %q", synopsis)
	}
}

func TestSSHKeyDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id only, no force",
			args:          []string{"-id=sshkey-123abc"},
			expectedID:    "sshkey-123abc",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=sshkey-456def", "-force"},
			expectedID:    "sshkey-456def",
			expectedForce: true,
		},
		{
			name:          "id with explicit force=true",
			args:          []string{"-id=sshkey-789ghi", "-force=true"},
			expectedID:    "sshkey-789ghi",
			expectedForce: true,
		},
		{
			name:          "id with explicit force=false",
			args:          []string{"-id=sshkey-abc123", "-force=false"},
			expectedID:    "sshkey-abc123",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &SSHKeyDeleteCommand{}

			flags := cmd.Meta.FlagSet("sshkey delete")
			flags.StringVar(&cmd.id, "id", "", "SSH key ID (required)")
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
