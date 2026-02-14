package command

import (
	"strings"
	"testing"
)

func TestLogoutHelp(t *testing.T) {
	cmd := &LogoutCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf logout") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-hostname") {
		t.Error("Help should mention -hostname flag")
	}
	if !strings.Contains(help, "app.terraform.io") {
		t.Error("Help should mention default hostname")
	}
	if !strings.Contains(help, "credentials.tfrc.json") {
		t.Error("Help should mention credentials file")
	}
}

func TestLogoutSynopsis(t *testing.T) {
	cmd := &LogoutCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Remove saved credentials" {
		t.Errorf("expected 'Remove saved credentials', got %q", synopsis)
	}
}

func TestLogoutFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedHostname string
	}{
		{
			name:             "default hostname",
			args:             []string{},
			expectedHostname: "app.terraform.io",
		},
		{
			name:             "custom hostname",
			args:             []string{"-hostname=custom.terraform.io"},
			expectedHostname: "custom.terraform.io",
		},
		{
			name:             "enterprise hostname",
			args:             []string{"-hostname=tfe.company.com"},
			expectedHostname: "tfe.company.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &LogoutCommand{}

			flags := cmd.Meta.FlagSet("logout")
			flags.StringVar(&cmd.hostname, "hostname", "app.terraform.io", "HCP Terraform hostname")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the hostname was set correctly
			if cmd.hostname != tt.expectedHostname {
				t.Errorf("expected hostname %q, got %q", tt.expectedHostname, cmd.hostname)
			}
		})
	}
}
