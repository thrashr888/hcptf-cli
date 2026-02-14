package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestSSHKeyReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &SSHKeyReadCommand{
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

func TestSSHKeyReadHelp(t *testing.T) {
	cmd := &SSHKeyReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf sshkey read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestSSHKeyReadSynopsis(t *testing.T) {
	cmd := &SSHKeyReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read SSH key details" {
		t.Errorf("expected 'Read SSH key details', got %q", synopsis)
	}
}

func TestSSHKeyReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id only, default format",
			args:        []string{"-id=sshkey-123abc"},
			expectedID:  "sshkey-123abc",
			expectedFmt: "table",
		},
		{
			name:        "id with table format",
			args:        []string{"-id=sshkey-456def", "-output=table"},
			expectedID:  "sshkey-456def",
			expectedFmt: "table",
		},
		{
			name:        "id with json format",
			args:        []string{"-id=sshkey-789ghi", "-output=json"},
			expectedID:  "sshkey-789ghi",
			expectedFmt: "json",
		},
		{
			name:        "long id with json format",
			args:        []string{"-id=sshkey-abcdef123456", "-output=json"},
			expectedID:  "sshkey-abcdef123456",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &SSHKeyReadCommand{}

			flags := cmd.Meta.FlagSet("sshkey read")
			flags.StringVar(&cmd.id, "id", "", "SSH key ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
