package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOAuthClientCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthClientCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{
		"-service-provider=github",
		"-http-url=https://github.com",
		"-api-url=https://api.github.com",
		"-oauth-token-string=ghp_test",
	})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestOAuthClientCreateRequiresServiceProvider(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthClientCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{
		"-organization=test-org",
		"-http-url=https://github.com",
		"-api-url=https://api.github.com",
		"-oauth-token-string=ghp_test",
	})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-service-provider") {
		t.Fatalf("expected service-provider error, got %q", out)
	}
}

func TestOAuthClientCreateRequiresHTTPURL(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthClientCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{
		"-organization=test-org",
		"-service-provider=github",
		"-api-url=https://api.github.com",
		"-oauth-token-string=ghp_test",
	})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-http-url") {
		t.Fatalf("expected http-url error, got %q", out)
	}
}

func TestOAuthClientCreateRequiresAPIURL(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthClientCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{
		"-organization=test-org",
		"-service-provider=github",
		"-http-url=https://github.com",
		"-oauth-token-string=ghp_test",
	})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-api-url") {
		t.Fatalf("expected api-url error, got %q", out)
	}
}

func TestOAuthClientCreateRequiresOAuthToken(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthClientCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{
		"-organization=test-org",
		"-service-provider=github",
		"-http-url=https://github.com",
		"-api-url=https://api.github.com",
	})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-oauth-token-string") {
		t.Fatalf("expected oauth-token-string error, got %q", out)
	}
}

func TestOAuthClientCreateHelp(t *testing.T) {
	cmd := &OAuthClientCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf oauthclient create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-service-provider") {
		t.Error("Help should mention -service-provider flag")
	}
	if !strings.Contains(help, "-http-url") {
		t.Error("Help should mention -http-url flag")
	}
	if !strings.Contains(help, "-api-url") {
		t.Error("Help should mention -api-url flag")
	}
	if !strings.Contains(help, "-oauth-token-string") {
		t.Error("Help should mention -oauth-token-string flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
}

func TestOAuthClientCreateSynopsis(t *testing.T) {
	cmd := &OAuthClientCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new OAuth client for VCS integration" {
		t.Errorf("expected 'Create a new OAuth client for VCS integration', got %q", synopsis)
	}
}

func TestOAuthClientCreateValidatesOrganizationScoped(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthClientCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{
		"-organization=test-org",
		"-service-provider=github",
		"-http-url=https://github.com",
		"-api-url=https://api.github.com",
		"-oauth-token-string=ghp_test",
		"-organization-scoped=invalid",
	})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "organization-scoped") {
		t.Fatalf("expected organization-scoped validation error, got %q", out)
	}
}

func TestOAuthClientCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedOrg      string
		expectedProvider string
		expectedHTTPURL  string
		expectedAPIURL   string
		expectedToken    string
		expectedName     string
		expectedScoped   string
		expectedFmt      string
	}{
		{
			name: "all required flags, default values",
			args: []string{
				"-organization=my-org",
				"-service-provider=github",
				"-http-url=https://github.com",
				"-api-url=https://api.github.com",
				"-oauth-token-string=ghp_test",
			},
			expectedOrg:      "my-org",
			expectedProvider: "github",
			expectedHTTPURL:  "https://github.com",
			expectedAPIURL:   "https://api.github.com",
			expectedToken:    "ghp_test",
			expectedName:     "",
			expectedScoped:   "true",
			expectedFmt:      "table",
		},
		{
			name: "org alias with name",
			args: []string{
				"-org=test-org",
				"-service-provider=gitlab_hosted",
				"-http-url=https://gitlab.com",
				"-api-url=https://gitlab.com/api/v4",
				"-oauth-token-string=glpat_test",
				"-name=GitLab Production",
			},
			expectedOrg:      "test-org",
			expectedProvider: "gitlab_hosted",
			expectedHTTPURL:  "https://gitlab.com",
			expectedAPIURL:   "https://gitlab.com/api/v4",
			expectedToken:    "glpat_test",
			expectedName:     "GitLab Production",
			expectedScoped:   "true",
			expectedFmt:      "table",
		},
		{
			name: "github enterprise with organization scoped false",
			args: []string{
				"-org=prod-org",
				"-service-provider=github_enterprise",
				"-http-url=https://github.example.com",
				"-api-url=https://github.example.com/api/v3",
				"-oauth-token-string=ghp_enterprise",
				"-organization-scoped=false",
			},
			expectedOrg:      "prod-org",
			expectedProvider: "github_enterprise",
			expectedHTTPURL:  "https://github.example.com",
			expectedAPIURL:   "https://github.example.com/api/v3",
			expectedToken:    "ghp_enterprise",
			expectedName:     "",
			expectedScoped:   "false",
			expectedFmt:      "table",
		},
		{
			name: "all options with json output",
			args: []string{
				"-organization=dev-org",
				"-service-provider=ado_server",
				"-http-url=https://ado.example.com",
				"-api-url=https://ado.example.com/_api",
				"-oauth-token-string=ado_token",
				"-name=Azure DevOps",
				"-organization-scoped=true",
				"-output=json",
			},
			expectedOrg:      "dev-org",
			expectedProvider: "ado_server",
			expectedHTTPURL:  "https://ado.example.com",
			expectedAPIURL:   "https://ado.example.com/_api",
			expectedToken:    "ado_token",
			expectedName:     "Azure DevOps",
			expectedScoped:   "true",
			expectedFmt:      "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OAuthClientCreateCommand{}

			flags := cmd.Meta.FlagSet("oauthclient create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.serviceProvider, "service-provider", "", "VCS provider (required)")
			flags.StringVar(&cmd.name, "name", "", "Display name for the OAuth client")
			flags.StringVar(&cmd.httpURL, "http-url", "", "VCS provider HTTP URL (required)")
			flags.StringVar(&cmd.apiURL, "api-url", "", "VCS provider API URL (required)")
			flags.StringVar(&cmd.oauthTokenString, "oauth-token-string", "", "OAuth token string (required)")
			flags.StringVar(&cmd.organizationScoped, "organization-scoped", "true", "Organization scoped")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the service provider was set correctly
			if cmd.serviceProvider != tt.expectedProvider {
				t.Errorf("expected serviceProvider %q, got %q", tt.expectedProvider, cmd.serviceProvider)
			}

			// Verify the HTTP URL was set correctly
			if cmd.httpURL != tt.expectedHTTPURL {
				t.Errorf("expected httpURL %q, got %q", tt.expectedHTTPURL, cmd.httpURL)
			}

			// Verify the API URL was set correctly
			if cmd.apiURL != tt.expectedAPIURL {
				t.Errorf("expected apiURL %q, got %q", tt.expectedAPIURL, cmd.apiURL)
			}

			// Verify the OAuth token was set correctly
			if cmd.oauthTokenString != tt.expectedToken {
				t.Errorf("expected oauthTokenString %q, got %q", tt.expectedToken, cmd.oauthTokenString)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
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
