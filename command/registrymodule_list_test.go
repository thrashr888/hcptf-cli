package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestRegistryModuleListRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RegistryModuleListCommand{
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

func TestRegistryModuleListHelp(t *testing.T) {
	cmd := &RegistryModuleListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	if !strings.Contains(help, "hcptf registrymodule list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
}

func TestRegistryModuleListSynopsis(t *testing.T) {
	cmd := &RegistryModuleListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List private registry modules in an organization" {
		t.Errorf("expected 'List private registry modules in an organization', got %q", synopsis)
	}
}

func TestRegistryModuleListFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOrg    string
		expectedFormat string
	}{
		{"organization with default", []string{"-organization=my-org"}, "my-org", "table"},
		{"org alias", []string{"-org=my-org"}, "my-org", "table"},
		{"json output", []string{"-organization=my-org", "-output=json"}, "my-org", "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &RegistryModuleListCommand{}

			flags := cmd.Meta.FlagSet("registrymodule list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.format, "output", "table", "Output format")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
