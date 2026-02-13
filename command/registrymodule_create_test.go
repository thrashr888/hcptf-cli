package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestRegistryModuleCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryModuleCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=vpc", "-provider=aws"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestRegistryModuleCreateRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryModuleCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-provider=aws"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestRegistryModuleCreateRequiresProvider(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryModuleCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-name=vpc"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-provider") {
		t.Fatalf("expected provider error, got %q", out)
	}
}

func TestRegistryModuleCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryModuleCreateCommand{
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

func TestRegistryModuleCreateHelp(t *testing.T) {
	cmd := &RegistryModuleCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf registrymodule create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "-provider") {
		t.Error("Help should mention -provider flag")
	}
	if !strings.Contains(help, "-registry-name") {
		t.Error("Help should mention -registry-name flag")
	}
	if !strings.Contains(help, "-no-code") {
		t.Error("Help should mention -no-code flag")
	}
	if !strings.Contains(help, "private registry module") {
		t.Error("Help should mention private registry module")
	}
}

func TestRegistryModuleCreateSynopsis(t *testing.T) {
	cmd := &RegistryModuleCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new private registry module" {
		t.Errorf("expected 'Create a new private registry module', got %q", synopsis)
	}
}

func TestRegistryModuleCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name                 string
		args                 []string
		expectedOrganization string
		expectedName         string
		expectedProvider     string
		expectedRegistryName string
		expectedNoCode       bool
		expectedFormat       string
	}{
		{
			name:                 "basic flags with defaults",
			args:                 []string{"-organization=my-org", "-name=vpc", "-provider=aws"},
			expectedOrganization: "my-org",
			expectedName:         "vpc",
			expectedProvider:     "aws",
			expectedRegistryName: "private",
			expectedNoCode:       false,
			expectedFormat:       "table",
		},
		{
			name:                 "using org alias",
			args:                 []string{"-org=test-org", "-name=network", "-provider=azure"},
			expectedOrganization: "test-org",
			expectedName:         "network",
			expectedProvider:     "azure",
			expectedRegistryName: "private",
			expectedNoCode:       false,
			expectedFormat:       "table",
		},
		{
			name:                 "with no-code enabled",
			args:                 []string{"-org=my-org", "-name=storage", "-provider=gcp", "-no-code"},
			expectedOrganization: "my-org",
			expectedName:         "storage",
			expectedProvider:     "gcp",
			expectedRegistryName: "private",
			expectedNoCode:       true,
			expectedFormat:       "table",
		},
		{
			name:                 "with json output",
			args:                 []string{"-organization=my-org", "-name=compute", "-provider=aws", "-output=json"},
			expectedOrganization: "my-org",
			expectedName:         "compute",
			expectedProvider:     "aws",
			expectedRegistryName: "private",
			expectedNoCode:       false,
			expectedFormat:       "json",
		},
		{
			name:                 "with custom registry name",
			args:                 []string{"-org=my-org", "-name=database", "-provider=aws", "-registry-name=public"},
			expectedOrganization: "my-org",
			expectedName:         "database",
			expectedProvider:     "aws",
			expectedRegistryName: "public",
			expectedNoCode:       false,
			expectedFormat:       "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &RegistryModuleCreateCommand{}

			flags := cmd.Meta.FlagSet("registrymodule create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.name, "name", "", "Module name (required)")
			flags.StringVar(&cmd.provider, "provider", "", "Provider name (required)")
			flags.StringVar(&cmd.registryName, "registry-name", "private", "Registry name: public or private (default: private)")
			flags.BoolVar(&cmd.noCode, "no-code", false, "Enable no-code publishing workflow")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrganization {
				t.Errorf("expected organization %q, got %q", tt.expectedOrganization, cmd.organization)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the provider was set correctly
			if cmd.provider != tt.expectedProvider {
				t.Errorf("expected provider %q, got %q", tt.expectedProvider, cmd.provider)
			}

			// Verify the registryName was set correctly
			if cmd.registryName != tt.expectedRegistryName {
				t.Errorf("expected registryName %q, got %q", tt.expectedRegistryName, cmd.registryName)
			}

			// Verify the noCode was set correctly
			if cmd.noCode != tt.expectedNoCode {
				t.Errorf("expected noCode %v, got %v", tt.expectedNoCode, cmd.noCode)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
