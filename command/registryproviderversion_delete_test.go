package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestRegistryProviderVersionDeleteRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryProviderVersionDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=aws", "-version=1.0.0"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestRegistryProviderVersionDeleteRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryProviderVersionDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-version=1.0.0"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestRegistryProviderVersionDeleteRequiresVersion(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryProviderVersionDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-name=aws"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-version") {
		t.Fatalf("expected version error, got %q", out)
	}
}

func TestRegistryProviderVersionDeleteRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryProviderVersionDeleteCommand{
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

func TestRegistryProviderVersionDeleteHelp(t *testing.T) {
	cmd := &RegistryProviderVersionDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf registryproviderversion delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "-version") {
		t.Error("Help should mention -version flag")
	}
	if !strings.Contains(help, "-namespace") {
		t.Error("Help should mention -namespace flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestRegistryProviderVersionDeleteSynopsis(t *testing.T) {
	cmd := &RegistryProviderVersionDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a version of a private registry provider" {
		t.Errorf("expected 'Delete a version of a private registry provider', got %q", synopsis)
	}
}

func TestRegistryProviderVersionDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedOrg       string
		expectedName      string
		expectedNamespace string
		expectedVersion   string
	}{
		{
			name:              "all required flags",
			args:              []string{"-organization=my-org", "-name=aws", "-version=1.0.0"},
			expectedOrg:       "my-org",
			expectedName:      "aws",
			expectedNamespace: "",
			expectedVersion:   "1.0.0",
		},
		{
			name:              "org alias with required flags",
			args:              []string{"-org=my-org", "-name=custom", "-version=2.5.1"},
			expectedOrg:       "my-org",
			expectedName:      "custom",
			expectedNamespace: "",
			expectedVersion:   "2.5.1",
		},
		{
			name:              "all flags with custom namespace",
			args:              []string{"-org=test-org", "-name=provider", "-namespace=custom-ns", "-version=0.1.0"},
			expectedOrg:       "test-org",
			expectedName:      "provider",
			expectedNamespace: "custom-ns",
			expectedVersion:   "0.1.0",
		},
		{
			name:              "with namespace",
			args:              []string{"-org=prod-org", "-name=infra", "-namespace=prod-ns", "-version=1.2.3"},
			expectedOrg:       "prod-org",
			expectedName:      "infra",
			expectedNamespace: "prod-ns",
			expectedVersion:   "1.2.3",
		},
		{
			name:              "full flags",
			args:              []string{"-org=dev-org", "-name=terraform", "-namespace=hashicorp", "-version=4.1.0"},
			expectedOrg:       "dev-org",
			expectedName:      "terraform",
			expectedNamespace: "hashicorp",
			expectedVersion:   "4.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &RegistryProviderVersionDeleteCommand{}

			flags := cmd.Meta.FlagSet("registryproviderversion delete")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.name, "name", "", "Provider name (required)")
			flags.StringVar(&cmd.namespace, "namespace", "", "Namespace (defaults to organization)")
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

			// Verify the namespace was set correctly
			if cmd.namespace != tt.expectedNamespace {
				t.Errorf("expected namespace %q, got %q", tt.expectedNamespace, cmd.namespace)
			}

			// Verify the version was set correctly
			if cmd.version != tt.expectedVersion {
				t.Errorf("expected version %q, got %q", tt.expectedVersion, cmd.version)
			}
		})
	}
}
