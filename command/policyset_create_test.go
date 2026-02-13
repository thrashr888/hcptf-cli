package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicySetCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=test-policyset"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestPolicySetCreateRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestPolicySetCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetCreateCommand{
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

func TestPolicySetCreateHelp(t *testing.T) {
	cmd := &PolicySetCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policyset create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestPolicySetCreateSynopsis(t *testing.T) {
	cmd := &PolicySetCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new policy set" {
		t.Errorf("expected 'Create a new policy set', got %q", synopsis)
	}
}

func TestPolicySetCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOrg    string
		expectedName   string
		expectedDesc   string
		expectedGlobal bool
		expectedFormat string
	}{
		{
			name:           "required flags only, default values",
			args:           []string{"-organization=my-org", "-name=test-policyset"},
			expectedOrg:    "my-org",
			expectedName:   "test-policyset",
			expectedDesc:   "",
			expectedGlobal: false,
			expectedFormat: "table",
		},
		{
			name:           "org alias with all optional flags",
			args:           []string{"-org=prod-org", "-name=prod-policyset", "-description=Production policy set", "-global"},
			expectedOrg:    "prod-org",
			expectedName:   "prod-policyset",
			expectedDesc:   "Production policy set",
			expectedGlobal: true,
			expectedFormat: "table",
		},
		{
			name:           "global flag with json output",
			args:           []string{"-org=test-org", "-name=test-policyset", "-global=true", "-output=json"},
			expectedOrg:    "test-org",
			expectedName:   "test-policyset",
			expectedDesc:   "",
			expectedGlobal: true,
			expectedFormat: "json",
		},
		{
			name:           "all flags with table output",
			args:           []string{"-organization=dev-org", "-name=dev-policyset", "-description=Dev policy set", "-global=false", "-output=table"},
			expectedOrg:    "dev-org",
			expectedName:   "dev-policyset",
			expectedDesc:   "Dev policy set",
			expectedGlobal: false,
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicySetCreateCommand{}

			flags := cmd.Meta.FlagSet("policyset create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.name, "name", "", "Policy set name (required)")
			flags.StringVar(&cmd.description, "description", "", "Policy set description")
			flags.BoolVar(&cmd.global, "global", false, "Apply to all workspaces")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the description was set correctly
			if cmd.description != tt.expectedDesc {
				t.Errorf("expected description %q, got %q", tt.expectedDesc, cmd.description)
			}

			// Verify the global flag was set correctly
			if cmd.global != tt.expectedGlobal {
				t.Errorf("expected global %v, got %v", tt.expectedGlobal, cmd.global)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
