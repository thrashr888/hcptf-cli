package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestUserTokenCreateRequiresDescription(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &UserTokenCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-description") {
		t.Fatalf("expected description error, got %q", out)
	}
}

func TestUserTokenCreateHelp(t *testing.T) {
	cmd := &UserTokenCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf usertoken create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-description") {
		t.Error("Help should mention -description flag")
	}
	if !strings.Contains(help, "-expired-at") {
		t.Error("Help should mention -expired-at flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestUserTokenCreateSynopsis(t *testing.T) {
	cmd := &UserTokenCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a user API token" {
		t.Errorf("expected 'Create a user API token', got %q", synopsis)
	}
}

func TestUserTokenCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedDesc   string
		expectedExpiry string
		expectedFmt    string
	}{
		{
			name:           "required flag, default format",
			args:           []string{"-description=Production user token"},
			expectedDesc:   "Production user token",
			expectedExpiry: "",
			expectedFmt:    "table",
		},
		{
			name:           "required flag with table format",
			args:           []string{"-description=Dev user", "-output=table"},
			expectedDesc:   "Dev user",
			expectedExpiry: "",
			expectedFmt:    "table",
		},
		{
			name:           "required flag with json format",
			args:           []string{"-description=CI user token", "-output=json"},
			expectedDesc:   "CI user token",
			expectedExpiry: "",
			expectedFmt:    "json",
		},
		{
			name:           "with expiration date",
			args:           []string{"-description=Temporary token", "-expired-at=2024-12-31T23:59:59Z"},
			expectedDesc:   "Temporary token",
			expectedExpiry: "2024-12-31T23:59:59Z",
			expectedFmt:    "table",
		},
		{
			name:           "all flags with json format",
			args:           []string{"-description=Full test", "-expired-at=2025-06-30T12:00:00Z", "-output=json"},
			expectedDesc:   "Full test",
			expectedExpiry: "2025-06-30T12:00:00Z",
			expectedFmt:    "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &UserTokenCreateCommand{}

			flags := cmd.Meta.FlagSet("usertoken create")
			flags.StringVar(&cmd.description, "description", "", "Token description (required)")
			flags.StringVar(&cmd.expiredAt, "expired-at", "", "Expiration date in ISO 8601 format (e.g., 2024-12-31T23:59:59Z)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the description was set correctly
			if cmd.description != tt.expectedDesc {
				t.Errorf("expected description %q, got %q", tt.expectedDesc, cmd.description)
			}

			// Verify the expired-at was set correctly
			if cmd.expiredAt != tt.expectedExpiry {
				t.Errorf("expected expiredAt %q, got %q", tt.expectedExpiry, cmd.expiredAt)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
