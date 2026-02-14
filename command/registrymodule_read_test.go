package command

import (
	"strings"
	"testing"
)

func TestRegistryModuleReadHelp(t *testing.T) {
	cmd := &RegistryModuleReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf registrymodule read") {
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
	if !strings.Contains(help, "-namespace") {
		t.Error("Help should mention -namespace flag")
	}
	if !strings.Contains(help, "-registry-name") {
		t.Error("Help should mention -registry-name flag")
	}
	if !strings.Contains(help, "private registry module") {
		t.Error("Help should mention private registry module")
	}
}

func TestRegistryModuleReadSynopsis(t *testing.T) {
	cmd := &RegistryModuleReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show details of a private registry module" {
		t.Errorf("expected 'Show details of a private registry module', got %q", synopsis)
	}
}

func TestRegistryModuleReadFlagParsing(t *testing.T) {
	tests := []struct {
		name                 string
		args                 []string
		expectedOrganization string
		expectedName         string
		expectedProvider     string
		expectedNamespace    string
		expectedRegistryName string
		expectedFormat       string
	}{
		{
			name:                 "basic flags with defaults",
			args:                 []string{"-organization=my-org", "-name=vpc", "-provider=aws"},
			expectedOrganization: "my-org",
			expectedName:         "vpc",
			expectedProvider:     "aws",
			expectedNamespace:    "",
			expectedRegistryName: "private",
			expectedFormat:       "table",
		},
		{
			name:                 "using org alias",
			args:                 []string{"-org=test-org", "-name=network", "-provider=azure"},
			expectedOrganization: "test-org",
			expectedName:         "network",
			expectedProvider:     "azure",
			expectedNamespace:    "",
			expectedRegistryName: "private",
			expectedFormat:       "table",
		},
		{
			name:                 "with custom namespace",
			args:                 []string{"-org=my-org", "-name=storage", "-provider=gcp", "-namespace=custom-ns"},
			expectedOrganization: "my-org",
			expectedName:         "storage",
			expectedProvider:     "gcp",
			expectedNamespace:    "custom-ns",
			expectedRegistryName: "private",
			expectedFormat:       "table",
		},
		{
			name:                 "with json output",
			args:                 []string{"-organization=my-org", "-name=compute", "-provider=aws", "-output=json"},
			expectedOrganization: "my-org",
			expectedName:         "compute",
			expectedProvider:     "aws",
			expectedNamespace:    "",
			expectedRegistryName: "private",
			expectedFormat:       "json",
		},
		{
			name:                 "with public registry",
			args:                 []string{"-org=my-org", "-name=database", "-provider=aws", "-registry-name=public"},
			expectedOrganization: "my-org",
			expectedName:         "database",
			expectedProvider:     "aws",
			expectedNamespace:    "",
			expectedRegistryName: "public",
			expectedFormat:       "table",
		},
		{
			name:                 "all flags specified",
			args:                 []string{"-org=my-org", "-name=vpc", "-provider=aws", "-namespace=shared", "-registry-name=public", "-output=json"},
			expectedOrganization: "my-org",
			expectedName:         "vpc",
			expectedProvider:     "aws",
			expectedNamespace:    "shared",
			expectedRegistryName: "public",
			expectedFormat:       "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &RegistryModuleReadCommand{}

			flags := cmd.Meta.FlagSet("registrymodule read")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.name, "name", "", "Module name (required)")
			flags.StringVar(&cmd.provider, "provider", "", "Provider name (required)")
			flags.StringVar(&cmd.namespace, "namespace", "", "Namespace (defaults to organization)")
			flags.StringVar(&cmd.registryName, "registry-name", "private", "Registry name: public or private (default: private)")
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

			// Verify the namespace was set correctly
			if cmd.namespace != tt.expectedNamespace {
				t.Errorf("expected namespace %q, got %q", tt.expectedNamespace, cmd.namespace)
			}

			// Verify the registryName was set correctly
			if cmd.registryName != tt.expectedRegistryName {
				t.Errorf("expected registryName %q, got %q", tt.expectedRegistryName, cmd.registryName)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
