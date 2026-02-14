package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAccountCreateRequiresEmail(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AccountCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-username=testuser", "-password=password123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-email") {
		t.Fatalf("expected email error, got %q", out)
	}
}

func TestAccountCreateRequiresUsername(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AccountCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-email=test@example.com", "-password=password123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-username") {
		t.Fatalf("expected username error, got %q", out)
	}
}

func TestAccountCreateRequiresPassword(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AccountCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-email=test@example.com", "-username=testuser"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-password") {
		t.Fatalf("expected password error, got %q", out)
	}
}

func TestAccountCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AccountCreateCommand{
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

func TestAccountCreateValidatesPasswordLength(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AccountCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-email=test@example.com", "-username=testuser", "-password=short"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "at least 8 characters") {
		t.Fatalf("expected password length error, got %q", out)
	}
}

func TestAccountCreateHelp(t *testing.T) {
	cmd := &AccountCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf account create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-email") {
		t.Error("Help should mention -email flag")
	}
	if !strings.Contains(help, "-username") {
		t.Error("Help should mention -username flag")
	}
	if !strings.Contains(help, "-password") {
		t.Error("Help should mention -password flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
	if !strings.Contains(help, "min 8 characters") {
		t.Error("Help should mention password minimum length")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
}

func TestAccountCreateSynopsis(t *testing.T) {
	cmd := &AccountCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new user account" {
		t.Errorf("expected 'Create a new user account', got %q", synopsis)
	}
}

func TestAccountCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedEmail    string
		expectedUsername string
		expectedPassword string
		expectedFormat   string
	}{
		{
			name:             "all required flags with default format",
			args:             []string{"-email=test@example.com", "-username=testuser", "-password=password123"},
			expectedEmail:    "test@example.com",
			expectedUsername: "testuser",
			expectedPassword: "password123",
			expectedFormat:   "table",
		},
		{
			name:             "all flags with table format",
			args:             []string{"-email=admin@example.com", "-username=adminuser", "-password=securepass123", "-output=table"},
			expectedEmail:    "admin@example.com",
			expectedUsername: "adminuser",
			expectedPassword: "securepass123",
			expectedFormat:   "table",
		},
		{
			name:             "all flags with json format",
			args:             []string{"-email=user@test.com", "-username=myuser", "-password=mypassword123", "-output=json"},
			expectedEmail:    "user@test.com",
			expectedUsername: "myuser",
			expectedPassword: "mypassword123",
			expectedFormat:   "json",
		},
		{
			name:             "flags with special characters",
			args:             []string{"-email=john.doe+test@example.com", "-username=john_doe123", "-password=P@ssw0rd!"},
			expectedEmail:    "john.doe+test@example.com",
			expectedUsername: "john_doe123",
			expectedPassword: "P@ssw0rd!",
			expectedFormat:   "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AccountCreateCommand{}

			flags := cmd.Meta.FlagSet("account create")
			flags.StringVar(&cmd.email, "email", "", "Email address (required)")
			flags.StringVar(&cmd.username, "username", "", "Username (required)")
			flags.StringVar(&cmd.password, "password", "", "Password (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the email was set correctly
			if cmd.email != tt.expectedEmail {
				t.Errorf("expected email %q, got %q", tt.expectedEmail, cmd.email)
			}

			// Verify the username was set correctly
			if cmd.username != tt.expectedUsername {
				t.Errorf("expected username %q, got %q", tt.expectedUsername, cmd.username)
			}

			// Verify the password was set correctly
			if cmd.password != tt.expectedPassword {
				t.Errorf("expected password %q, got %q", tt.expectedPassword, cmd.password)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
