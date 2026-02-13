package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOAuthTokenListRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthTokenListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestOAuthTokenListHelp(t *testing.T) {
	cmd := &OAuthTokenListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf oauthtoken list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -organization is required")
	}
}

func TestOAuthTokenListSynopsis(t *testing.T) {
	cmd := &OAuthTokenListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List OAuth tokens in an organization" {
		t.Errorf("expected 'List OAuth tokens in an organization', got %q", synopsis)
	}
}

func TestOAuthTokenListFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOrg string
		expectedFmt string
	}{
		{
			name:        "organization, default format",
			args:        []string{"-organization=my-org"},
			expectedOrg: "my-org",
			expectedFmt: "table",
		},
		{
			name:        "org alias",
			args:        []string{"-org=test-org"},
			expectedOrg: "test-org",
			expectedFmt: "table",
		},
		{
			name:        "organization, table format",
			args:        []string{"-org=my-org", "-output=table"},
			expectedOrg: "my-org",
			expectedFmt: "table",
		},
		{
			name:        "organization, json format",
			args:        []string{"-org=prod-org", "-output=json"},
			expectedOrg: "prod-org",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OAuthTokenListCommand{}

			flags := cmd.Meta.FlagSet("oauthtoken list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
