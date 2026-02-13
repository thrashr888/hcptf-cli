package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestRegistryModuleDeleteVersionRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryModuleDeleteVersionCommand{
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

func TestRegistryModuleDeleteVersionRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryModuleDeleteVersionCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestRegistryModuleDeleteVersionRequiresProvider(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryModuleDeleteVersionCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-name=vpc"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-provider") {
		t.Fatalf("expected provider error, got %q", out)
	}
}

func TestRegistryModuleDeleteVersionRequiresVersion(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryModuleDeleteVersionCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-name=vpc", "-provider=aws"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-version") {
		t.Fatalf("expected version error, got %q", out)
	}
}

func TestRegistryModuleDeleteVersionHelp(t *testing.T) {
	cmd := &RegistryModuleDeleteVersionCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf registrymodule delete-version") {
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
	if !strings.Contains(help, "-version") {
		t.Error("Help should mention -version flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
}

func TestRegistryModuleDeleteVersionSynopsis(t *testing.T) {
	cmd := &RegistryModuleDeleteVersionCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a version of a private registry module" {
		t.Errorf("expected 'Delete a version of a private registry module', got %q", synopsis)
	}
}

func TestRegistryModuleDeleteVersionFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedOrg      string
		expectedName     string
		expectedProvider string
		expectedNS       string
		expectedRegName  string
		expectedVersion  string
	}{
		{
			name:             "required flags only",
			args:             []string{"-organization=my-org", "-name=vpc", "-provider=aws", "-version=1.0.0"},
			expectedOrg:      "my-org",
			expectedName:     "vpc",
			expectedProvider: "aws",
			expectedRegName:  "private",
			expectedVersion:  "1.0.0",
		},
		{
			name:             "org alias flag",
			args:             []string{"-org=test-org", "-name=vpc", "-provider=aws", "-version=2.0.0"},
			expectedOrg:      "test-org",
			expectedName:     "vpc",
			expectedProvider: "aws",
			expectedRegName:  "private",
			expectedVersion:  "2.0.0",
		},
		{
			name:             "with namespace",
			args:             []string{"-org=my-org", "-name=vpc", "-provider=aws", "-namespace=custom-ns", "-version=1.0.0"},
			expectedOrg:      "my-org",
			expectedName:     "vpc",
			expectedProvider: "aws",
			expectedNS:       "custom-ns",
			expectedRegName:  "private",
			expectedVersion:  "1.0.0",
		},
		{
			name:             "with registry name",
			args:             []string{"-org=my-org", "-name=vpc", "-provider=aws", "-registry-name=private", "-version=1.0.0"},
			expectedOrg:      "my-org",
			expectedName:     "vpc",
			expectedProvider: "aws",
			expectedRegName:  "private",
			expectedVersion:  "1.0.0",
		},
		{
			name:             "all flags",
			args:             []string{"-org=my-org", "-name=s3-bucket", "-provider=aws", "-namespace=custom", "-registry-name=private", "-version=3.2.1"},
			expectedOrg:      "my-org",
			expectedName:     "s3-bucket",
			expectedProvider: "aws",
			expectedNS:       "custom",
			expectedRegName:  "private",
			expectedVersion:  "3.2.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &RegistryModuleDeleteVersionCommand{}

			flags := cmd.Meta.FlagSet("registrymodule delete-version")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.name, "name", "", "Module name (required)")
			flags.StringVar(&cmd.provider, "provider", "", "Provider name (required)")
			flags.StringVar(&cmd.namespace, "namespace", "", "Namespace")
			flags.StringVar(&cmd.registryName, "registry-name", "private", "Registry name")
			flags.StringVar(&cmd.version, "version", "", "Version string (required)")

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

			// Verify the provider was set correctly
			if cmd.provider != tt.expectedProvider {
				t.Errorf("expected provider %q, got %q", tt.expectedProvider, cmd.provider)
			}

			// Verify the namespace was set correctly
			if cmd.namespace != tt.expectedNS {
				t.Errorf("expected namespace %q, got %q", tt.expectedNS, cmd.namespace)
			}

			// Verify the registry name was set correctly
			if cmd.registryName != tt.expectedRegName {
				t.Errorf("expected registry name %q, got %q", tt.expectedRegName, cmd.registryName)
			}

			// Verify the version was set correctly
			if cmd.version != tt.expectedVersion {
				t.Errorf("expected version %q, got %q", tt.expectedVersion, cmd.version)
			}
		})
	}
}
