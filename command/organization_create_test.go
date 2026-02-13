package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationCreateRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-email=test@example.com"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestOrganizationCreateRequiresEmail(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=test-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-email") {
		t.Fatalf("expected email error, got %q", out)
	}
}

func TestOrganizationCreateRequiresBothFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationCreateCommand{
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

func TestOrganizationCreateHelp(t *testing.T) {
	cmd := &OrganizationCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organization create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "-email") {
		t.Error("Help should mention -email flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestOrganizationCreateSynopsis(t *testing.T) {
	cmd := &OrganizationCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new organization" {
		t.Errorf("expected 'Create a new organization', got %q", synopsis)
	}
}

func TestOrganizationCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedName   string
		expectedEmail  string
		expectedFormat string
	}{
		{
			name:           "name, email, default format",
			args:           []string{"-name=test-org", "-email=test@example.com"},
			expectedName:   "test-org",
			expectedEmail:  "test@example.com",
			expectedFormat: "table",
		},
		{
			name:           "name, email, table format",
			args:           []string{"-name=my-org", "-email=admin@example.com", "-output=table"},
			expectedName:   "my-org",
			expectedEmail:  "admin@example.com",
			expectedFormat: "table",
		},
		{
			name:           "name, email, json format",
			args:           []string{"-name=prod-org", "-email=prod@example.com", "-output=json"},
			expectedName:   "prod-org",
			expectedEmail:  "prod@example.com",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationCreateCommand{}

			flags := cmd.Meta.FlagSet("organization create")
			flags.StringVar(&cmd.name, "name", "", "Organization name (required)")
			flags.StringVar(&cmd.email, "email", "", "Admin email address (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the email was set correctly
			if cmd.email != tt.expectedEmail {
				t.Errorf("expected email %q, got %q", tt.expectedEmail, cmd.email)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
