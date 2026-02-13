package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestSSHKeyUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &SSHKeyUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=new-name"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestSSHKeyUpdateHelp(t *testing.T) {
	cmd := &SSHKeyUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf sshkey update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestSSHKeyUpdateSynopsis(t *testing.T) {
	cmd := &SSHKeyUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update an SSH key" {
		t.Errorf("expected 'Update an SSH key', got %q", synopsis)
	}
}

func TestSSHKeyUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedID   string
		expectedName string
		expectedFmt  string
	}{
		{
			name:         "id only, default format",
			args:         []string{"-id=sshkey-123abc"},
			expectedID:   "sshkey-123abc",
			expectedName: "",
			expectedFmt:  "table",
		},
		{
			name:         "id and name",
			args:         []string{"-id=sshkey-456def", "-name=updated-key"},
			expectedID:   "sshkey-456def",
			expectedName: "updated-key",
			expectedFmt:  "table",
		},
		{
			name:         "id and name with table format",
			args:         []string{"-id=sshkey-789ghi", "-name=new-key-name", "-output=table"},
			expectedID:   "sshkey-789ghi",
			expectedName: "new-key-name",
			expectedFmt:  "table",
		},
		{
			name:         "id and name with json format",
			args:         []string{"-id=sshkey-abc123", "-name=renamed-key", "-output=json"},
			expectedID:   "sshkey-abc123",
			expectedName: "renamed-key",
			expectedFmt:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &SSHKeyUpdateCommand{}

			flags := cmd.Meta.FlagSet("sshkey update")
			flags.StringVar(&cmd.id, "id", "", "SSH key ID (required)")
			flags.StringVar(&cmd.name, "name", "", "SSH key name")
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

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}

func TestSSHKeyUpdateNameEmpty(t *testing.T) {
	cmd := &SSHKeyUpdateCommand{}

	flags := cmd.Meta.FlagSet("sshkey update")
	flags.StringVar(&cmd.id, "id", "", "SSH key ID (required)")
	flags.StringVar(&cmd.name, "name", "", "SSH key name")
	flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

	if err := flags.Parse([]string{"-id=sshkey-123", "-name="}); err != nil {
		t.Fatalf("flag parsing failed: %v", err)
	}

	if cmd.name != "" {
		t.Errorf("expected empty name, got %q", cmd.name)
	}
}

func TestSSHKeyUpdateNameWithSpecialChars(t *testing.T) {
	tests := []struct {
		name      string
		nameValue string
	}{
		{"with-dashes", "my-ssh-key"},
		{"with-underscores", "my_ssh_key"},
		{"with-numbers", "sshkey123"},
		{"mixed", "my-ssh_key-123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &SSHKeyUpdateCommand{}

			flags := cmd.Meta.FlagSet("sshkey update")
			flags.StringVar(&cmd.id, "id", "", "SSH key ID (required)")
			flags.StringVar(&cmd.name, "name", "", "SSH key name")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse([]string{"-id=sshkey-123", "-name=" + tt.nameValue}); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			if cmd.name != tt.nameValue {
				t.Errorf("expected name %q, got %q", tt.nameValue, cmd.name)
			}
		})
	}
}
