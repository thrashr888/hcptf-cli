package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOAuthTokenUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthTokenUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-ssh-key=test-key"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestOAuthTokenUpdateHelp(t *testing.T) {
	cmd := &OAuthTokenUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf oauthtoken update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestOAuthTokenUpdateSynopsis(t *testing.T) {
	cmd := &OAuthTokenUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update OAuth token settings" {
		t.Errorf("expected 'Update OAuth token settings', got %q", synopsis)
	}
}

func TestOAuthTokenUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedSSHKey string
		expectedFmt    string
	}{
		{
			name:           "id and ssh-key",
			args:           []string{"-id=ot-hmAyP66qk2AMVdbJ", "-ssh-key=-----BEGIN RSA PRIVATE KEY-----"},
			expectedID:     "ot-hmAyP66qk2AMVdbJ",
			expectedSSHKey: "-----BEGIN RSA PRIVATE KEY-----",
			expectedFmt:    "table",
		},
		{
			name:           "id and ssh-key with json output",
			args:           []string{"-id=ot-ABC123XYZ456", "-ssh-key=ssh-rsa AAAAB3NzaC1yc2E", "-output=json"},
			expectedID:     "ot-ABC123XYZ456",
			expectedSSHKey: "ssh-rsa AAAAB3NzaC1yc2E",
			expectedFmt:    "json",
		},
		{
			name:           "id and multiline ssh-key",
			args:           []string{"-id=ot-test12345678", "-ssh-key=-----BEGIN OPENSSH PRIVATE KEY-----\nkey content here\n-----END OPENSSH PRIVATE KEY-----"},
			expectedID:     "ot-test12345678",
			expectedSSHKey: "-----BEGIN OPENSSH PRIVATE KEY-----\nkey content here\n-----END OPENSSH PRIVATE KEY-----",
			expectedFmt:    "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OAuthTokenUpdateCommand{}

			flags := cmd.Meta.FlagSet("oauthtoken update")
			flags.StringVar(&cmd.id, "id", "", "OAuth token ID (required)")
			flags.StringVar(&cmd.sshKey, "ssh-key", "", "SSH private key content")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the ssh-key was set correctly
			if cmd.sshKey != tt.expectedSSHKey {
				t.Errorf("expected sshKey %q, got %q", tt.expectedSSHKey, cmd.sshKey)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
