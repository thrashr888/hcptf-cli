package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackCreateRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-project-id=prj-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestStackCreateRequiresProjectID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=test-stack"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-project-id") {
		t.Fatalf("expected project-id error, got %q", out)
	}
}

func TestStackCreateRequiresBothFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackCreateCommand{
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

func TestStackCreateHelp(t *testing.T) {
	cmd := &StackCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stack create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "-project-id") {
		t.Error("Help should mention -project-id flag")
	}
	if !strings.Contains(help, "-description") {
		t.Error("Help should mention -description flag")
	}
	if !strings.Contains(help, "-vcs-identifier") {
		t.Error("Help should mention -vcs-identifier flag")
	}
	if !strings.Contains(help, "-vcs-branch") {
		t.Error("Help should mention -vcs-branch flag")
	}
	if !strings.Contains(help, "-oauth-token-id") {
		t.Error("Help should mention -oauth-token-id flag")
	}
	if !strings.Contains(help, "-service-provider") {
		t.Error("Help should mention -service-provider flag")
	}
	if !strings.Contains(help, "-speculative-enabled") {
		t.Error("Help should mention -speculative-enabled flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
	if !strings.Contains(help, "Example:") {
		t.Error("Help should contain examples")
	}
}

func TestStackCreateSynopsis(t *testing.T) {
	cmd := &StackCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new stack" {
		t.Errorf("expected 'Create a new stack', got %q", synopsis)
	}
}

func TestStackCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name                  string
		args                  []string
		expectedName          string
		expectedProjectID     string
		expectedDescription   string
		expectedVCSIdentifier string
		expectedVCSBranch     string
		expectedOAuthTokenID  string
		expectedServiceProv   string
		expectedSpeculative   bool
		expectedFormat        string
	}{
		{
			name:                "name and project-id, defaults",
			args:                []string{"-name=my-stack", "-project-id=prj-123"},
			expectedName:        "my-stack",
			expectedProjectID:   "prj-123",
			expectedDescription: "",
			expectedServiceProv: "github",
			expectedSpeculative: false,
			expectedFormat:      "table",
		},
		{
			name:                "with description",
			args:                []string{"-name=infra-stack", "-project-id=prj-456", "-description=Infrastructure stack"},
			expectedName:        "infra-stack",
			expectedProjectID:   "prj-456",
			expectedDescription: "Infrastructure stack",
			expectedServiceProv: "github",
			expectedSpeculative: false,
			expectedFormat:      "table",
		},
		{
			name:                  "with VCS configuration",
			args:                  []string{"-name=vcs-stack", "-project-id=prj-789", "-vcs-identifier=myorg/myrepo", "-vcs-branch=main", "-oauth-token-id=ot-abc123"},
			expectedName:          "vcs-stack",
			expectedProjectID:     "prj-789",
			expectedVCSIdentifier: "myorg/myrepo",
			expectedVCSBranch:     "main",
			expectedOAuthTokenID:  "ot-abc123",
			expectedServiceProv:   "github",
			expectedSpeculative:   false,
			expectedFormat:        "table",
		},
		{
			name:                "with speculative enabled",
			args:                []string{"-name=spec-stack", "-project-id=prj-999", "-speculative-enabled"},
			expectedName:        "spec-stack",
			expectedProjectID:   "prj-999",
			expectedServiceProv: "github",
			expectedSpeculative: true,
			expectedFormat:      "table",
		},
		{
			name:                "json output format",
			args:                []string{"-name=json-stack", "-project-id=prj-111", "-output=json"},
			expectedName:        "json-stack",
			expectedProjectID:   "prj-111",
			expectedServiceProv: "github",
			expectedSpeculative: false,
			expectedFormat:      "json",
		},
		{
			name:                  "all flags",
			args:                  []string{"-name=full-stack", "-project-id=prj-222", "-description=Full featured stack", "-vcs-identifier=org/repo", "-vcs-branch=develop", "-oauth-token-id=ot-xyz789", "-service-provider=gitlab", "-speculative-enabled", "-output=json"},
			expectedName:          "full-stack",
			expectedProjectID:     "prj-222",
			expectedDescription:   "Full featured stack",
			expectedVCSIdentifier: "org/repo",
			expectedVCSBranch:     "develop",
			expectedOAuthTokenID:  "ot-xyz789",
			expectedServiceProv:   "gitlab",
			expectedSpeculative:   true,
			expectedFormat:        "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackCreateCommand{}

			flags := cmd.Meta.FlagSet("stack create")
			flags.StringVar(&cmd.name, "name", "", "Stack name (required)")
			flags.StringVar(&cmd.description, "description", "", "Stack description")
			flags.StringVar(&cmd.projectID, "project-id", "", "Project ID (required)")
			flags.StringVar(&cmd.vcsIdentifier, "vcs-identifier", "", "VCS repository identifier (org/repo)")
			flags.StringVar(&cmd.vcsBranch, "vcs-branch", "", "VCS branch (defaults to repo default branch)")
			flags.StringVar(&cmd.oauthTokenID, "oauth-token-id", "", "OAuth token ID for VCS connection")
			flags.StringVar(&cmd.serviceProvider, "service-provider", "github", "VCS service provider")
			flags.BoolVar(&cmd.speculativeEnabled, "speculative-enabled", false, "Enable speculative plans on PRs")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the project-id was set correctly
			if cmd.projectID != tt.expectedProjectID {
				t.Errorf("expected project-id %q, got %q", tt.expectedProjectID, cmd.projectID)
			}

			// Verify the description was set correctly
			if cmd.description != tt.expectedDescription {
				t.Errorf("expected description %q, got %q", tt.expectedDescription, cmd.description)
			}

			// Verify the vcs-identifier was set correctly
			if cmd.vcsIdentifier != tt.expectedVCSIdentifier {
				t.Errorf("expected vcs-identifier %q, got %q", tt.expectedVCSIdentifier, cmd.vcsIdentifier)
			}

			// Verify the vcs-branch was set correctly
			if cmd.vcsBranch != tt.expectedVCSBranch {
				t.Errorf("expected vcs-branch %q, got %q", tt.expectedVCSBranch, cmd.vcsBranch)
			}

			// Verify the oauth-token-id was set correctly
			if cmd.oauthTokenID != tt.expectedOAuthTokenID {
				t.Errorf("expected oauth-token-id %q, got %q", tt.expectedOAuthTokenID, cmd.oauthTokenID)
			}

			// Verify the service-provider was set correctly
			if cmd.serviceProvider != tt.expectedServiceProv {
				t.Errorf("expected service-provider %q, got %q", tt.expectedServiceProv, cmd.serviceProvider)
			}

			// Verify the speculative-enabled was set correctly
			if cmd.speculativeEnabled != tt.expectedSpeculative {
				t.Errorf("expected speculative-enabled %v, got %v", tt.expectedSpeculative, cmd.speculativeEnabled)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}









