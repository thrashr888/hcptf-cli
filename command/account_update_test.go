package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAccountUpdateRequiresAtLeastOneField(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AccountUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-password=mypassword"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	output := ui.ErrorWriter.String()
	if !strings.Contains(output, "at least one of") {
		t.Fatalf("expected 'at least one of' error, got %q", output)
	}
	if !strings.Contains(output, "-email") {
		t.Fatalf("expected error to mention -email flag, got %q", output)
	}
	if !strings.Contains(output, "-username") {
		t.Fatalf("expected error to mention -username flag, got %q", output)
	}
	if !strings.Contains(output, "-new-password") {
		t.Fatalf("expected error to mention -new-password flag, got %q", output)
	}
}

func TestAccountUpdateRequiresPassword(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AccountUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-email=newemail@example.com"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	output := ui.ErrorWriter.String()
	if !strings.Contains(output, "-password") {
		t.Fatalf("expected password error, got %q", output)
	}
	if !strings.Contains(output, "required") {
		t.Fatalf("expected 'required' in error message, got %q", output)
	}
}

func TestAccountUpdateValidatesPasswordLength(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AccountUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{
		"-email=newemail@example.com",
		"-password=current",
		"-new-password=short",
	})
	if code != 1 {
		t.Fatalf("expected exit 1 for short password, got %d", code)
	}

	output := ui.ErrorWriter.String()
	if !strings.Contains(output, "at least 8 characters") {
		t.Fatalf("expected password length error, got %q", output)
	}
}

func TestAccountUpdatePasswordChangeNotSupported(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AccountUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{
		"-email=newemail@example.com",
		"-password=currentpassword",
		"-new-password=newsecurepassword123",
	})
	if code != 1 {
		t.Fatalf("expected exit 1 for password change, got %d", code)
	}

	output := ui.ErrorWriter.String()
	if !strings.Contains(output, "not yet supported") {
		t.Fatalf("expected 'not yet supported' error, got %q", output)
	}
	if !strings.Contains(output, "web UI") {
		t.Fatalf("expected web UI suggestion, got %q", output)
	}
}

func TestAccountUpdateHelp(t *testing.T) {
	cmd := &AccountUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf account update") {
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
	if !strings.Contains(help, "-new-password") {
		t.Error("Help should mention -new-password flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -password is required")
	}
	if !strings.Contains(help, "Example:") {
		t.Error("Help should contain examples")
	}
}

func TestAccountUpdateSynopsis(t *testing.T) {
	cmd := &AccountUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update account details" {
		t.Errorf("expected 'Update account details', got %q", synopsis)
	}
}

func TestAccountUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedEmail   string
		expectedUser    string
		expectedPass    string
		expectedNewPass string
		expectedFormat  string
	}{
		{
			name:           "email update",
			args:           []string{"-email=newemail@example.com", "-password=current"},
			expectedEmail:  "newemail@example.com",
			expectedPass:   "current",
			expectedFormat: "table",
		},
		{
			name:           "username update",
			args:           []string{"-username=newusername", "-password=current"},
			expectedUser:   "newusername",
			expectedPass:   "current",
			expectedFormat: "table",
		},
		{
			name:            "password update",
			args:            []string{"-password=current", "-new-password=newsecure123"},
			expectedPass:    "current",
			expectedNewPass: "newsecure123",
			expectedFormat:  "table",
		},
		{
			name:           "email and username",
			args:           []string{"-email=new@example.com", "-username=newuser", "-password=current"},
			expectedEmail:  "new@example.com",
			expectedUser:   "newuser",
			expectedPass:   "current",
			expectedFormat: "table",
		},
		{
			name:            "all fields",
			args:            []string{"-email=new@example.com", "-username=newuser", "-password=current", "-new-password=newsecure123"},
			expectedEmail:   "new@example.com",
			expectedUser:    "newuser",
			expectedPass:    "current",
			expectedNewPass: "newsecure123",
			expectedFormat:  "table",
		},
		{
			name:           "with json output",
			args:           []string{"-email=new@example.com", "-password=current", "-output=json"},
			expectedEmail:  "new@example.com",
			expectedPass:   "current",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AccountUpdateCommand{}

			flags := cmd.Meta.FlagSet("account update")
			flags.StringVar(&cmd.email, "email", "", "New email address")
			flags.StringVar(&cmd.username, "username", "", "New username")
			flags.StringVar(&cmd.password, "password", "", "Current password (required for changes)")
			flags.StringVar(&cmd.newPassword, "new-password", "", "New password")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the email was set correctly
			if cmd.email != tt.expectedEmail {
				t.Errorf("expected email %q, got %q", tt.expectedEmail, cmd.email)
			}

			// Verify the username was set correctly
			if cmd.username != tt.expectedUser {
				t.Errorf("expected username %q, got %q", tt.expectedUser, cmd.username)
			}

			// Verify the password was set correctly
			if cmd.password != tt.expectedPass {
				t.Errorf("expected password %q, got %q", tt.expectedPass, cmd.password)
			}

			// Verify the new password was set correctly
			if cmd.newPassword != tt.expectedNewPass {
				t.Errorf("expected new password %q, got %q", tt.expectedNewPass, cmd.newPassword)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
