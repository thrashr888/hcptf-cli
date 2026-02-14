package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOAuthClientUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthClientUpdateCommand{
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

func TestOAuthClientUpdateHelp(t *testing.T) {
	cmd := &OAuthClientUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf oauthclient update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestOAuthClientUpdateSynopsis(t *testing.T) {
	cmd := &OAuthClientUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update OAuth client settings" {
		t.Errorf("expected 'Update OAuth client settings', got %q", synopsis)
	}
}

func TestOAuthClientUpdateValidatesOrganizationScoped(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthClientUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=oc-123", "-organization-scoped=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "organization-scoped") {
		t.Fatalf("expected organization-scoped validation error, got %q", out)
	}
}

func TestOAuthClientUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedName   string
		expectedToken  string
		expectedScoped string
		expectedFmt    string
	}{
		{
			name:           "id only, no updates",
			args:           []string{"-id=oc-XKFwG6ggfA9n7t1K"},
			expectedID:     "oc-XKFwG6ggfA9n7t1K",
			expectedName:   "",
			expectedToken:  "",
			expectedScoped: "",
			expectedFmt:    "table",
		},
		{
			name:           "update name",
			args:           []string{"-id=oc-ABC123XYZ456", "-name=GitHub Production"},
			expectedID:     "oc-ABC123XYZ456",
			expectedName:   "GitHub Production",
			expectedToken:  "",
			expectedScoped: "",
			expectedFmt:    "table",
		},
		{
			name:           "rotate oauth token",
			args:           []string{"-id=oc-test12345678", "-oauth-token-string=ghp_newtoken"},
			expectedID:     "oc-test12345678",
			expectedName:   "",
			expectedToken:  "ghp_newtoken",
			expectedScoped: "",
			expectedFmt:    "table",
		},
		{
			name:           "update organization scoped with json output",
			args:           []string{"-id=oc-prod99999", "-organization-scoped=false", "-output=json"},
			expectedID:     "oc-prod99999",
			expectedName:   "",
			expectedToken:  "",
			expectedScoped: "false",
			expectedFmt:    "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OAuthClientUpdateCommand{}

			flags := cmd.Meta.FlagSet("oauthclient update")
			flags.StringVar(&cmd.id, "id", "", "OAuth client ID (required)")
			flags.StringVar(&cmd.name, "name", "", "Display name for the OAuth client")
			flags.StringVar(&cmd.oauthTokenString, "oauth-token-string", "", "New OAuth token string")
			flags.StringVar(&cmd.organizationScoped, "organization-scoped", "", "Organization scoped")
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

			// Verify the oauth token was set correctly
			if cmd.oauthTokenString != tt.expectedToken {
				t.Errorf("expected oauthTokenString %q, got %q", tt.expectedToken, cmd.oauthTokenString)
			}

			// Verify the organization scoped was set correctly
			if cmd.organizationScoped != tt.expectedScoped {
				t.Errorf("expected organizationScoped %q, got %q", tt.expectedScoped, cmd.organizationScoped)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
