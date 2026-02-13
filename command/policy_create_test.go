package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicyCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicyCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=test-policy", "-policy-file=test.sentinel"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestPolicyCreateRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicyCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-policy-file=test.sentinel"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestPolicyCreateRequiresPolicyFile(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicyCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-name=test-policy"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-policy-file") {
		t.Fatalf("expected policy-file error, got %q", out)
	}
}

func TestPolicyCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicyCreateCommand{
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

func TestPolicyCreateHelp(t *testing.T) {
	cmd := &PolicyCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policy create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "-policy-file") {
		t.Error("Help should mention -policy-file flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestPolicyCreateSynopsis(t *testing.T) {
	cmd := &PolicyCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new policy" {
		t.Errorf("expected 'Create a new policy', got %q", synopsis)
	}
}

func TestPolicyCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOrg    string
		expectedName   string
		expectedDesc   string
		expectedEnf    string
		expectedFile   string
		expectedFormat string
	}{
		{
			name:           "required flags only, default values",
			args:           []string{"-organization=my-org", "-name=test-policy", "-policy-file=policy.sentinel"},
			expectedOrg:    "my-org",
			expectedName:   "test-policy",
			expectedDesc:   "",
			expectedEnf:    "advisory",
			expectedFile:   "policy.sentinel",
			expectedFormat: "table",
		},
		{
			name:           "org alias with all optional flags",
			args:           []string{"-org=prod-org", "-name=prod-policy", "-description=Production policy", "-enforce=hard-mandatory", "-policy-file=prod.sentinel"},
			expectedOrg:    "prod-org",
			expectedName:   "prod-policy",
			expectedDesc:   "Production policy",
			expectedEnf:    "hard-mandatory",
			expectedFile:   "prod.sentinel",
			expectedFormat: "table",
		},
		{
			name:           "soft-mandatory enforcement with json output",
			args:           []string{"-org=test-org", "-name=test-policy", "-enforce=soft-mandatory", "-policy-file=test.sentinel", "-output=json"},
			expectedOrg:    "test-org",
			expectedName:   "test-policy",
			expectedDesc:   "",
			expectedEnf:    "soft-mandatory",
			expectedFile:   "test.sentinel",
			expectedFormat: "json",
		},
		{
			name:           "all flags with table output",
			args:           []string{"-organization=dev-org", "-name=dev-policy", "-description=Dev policy", "-enforce=advisory", "-policy-file=dev.sentinel", "-output=table"},
			expectedOrg:    "dev-org",
			expectedName:   "dev-policy",
			expectedDesc:   "Dev policy",
			expectedEnf:    "advisory",
			expectedFile:   "dev.sentinel",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicyCreateCommand{}

			flags := cmd.Meta.FlagSet("policy create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.name, "name", "", "Policy name (required)")
			flags.StringVar(&cmd.description, "description", "", "Policy description")
			flags.StringVar(&cmd.enforce, "enforce", "advisory", "Enforcement level: advisory, soft-mandatory, or hard-mandatory")
			flags.StringVar(&cmd.policyFile, "policy-file", "", "Path to policy file (required)")
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

			// Verify the enforce was set correctly
			if cmd.enforce != tt.expectedEnf {
				t.Errorf("expected enforce %q, got %q", tt.expectedEnf, cmd.enforce)
			}

			// Verify the policy-file was set correctly
			if cmd.policyFile != tt.expectedFile {
				t.Errorf("expected policyFile %q, got %q", tt.expectedFile, cmd.policyFile)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
