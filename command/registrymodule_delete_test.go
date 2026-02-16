package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newRegistryModuleDeleteCommand(ui cli.Ui, svc registryModuleDeleter) *RegistryModuleDeleteCommand {
	return &RegistryModuleDeleteCommand{
		Meta:              newTestMeta(ui),
		registryModuleSvc: svc,
	}
}

func TestRegistryModuleDeleteRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRegistryModuleDeleteCommand(ui, &mockRegistryModuleDeleteService{})

	code := cmd.Run([]string{"-name=vpc"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestRegistryModuleDeleteRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRegistryModuleDeleteCommand(ui, &mockRegistryModuleDeleteService{})

	code := cmd.Run([]string{"-organization=test-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestRegistryModuleDeleteRequiresBothFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRegistryModuleDeleteCommand(ui, &mockRegistryModuleDeleteService{})

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "required") {
		t.Fatalf("expected required error, got %q", out)
	}
}

func TestRegistryModuleDeleteHelp(t *testing.T) {
	cmd := &RegistryModuleDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf registry module delete") {
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
	if !strings.Contains(help, "private registry module") {
		t.Error("Help should mention private registry module")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "-y") {
		t.Error("Help should mention -y flag")
	}
}

func TestRegistryModuleDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryModuleDeleteService{err: errors.New("boom")}
	cmd := newRegistryModuleDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org", "-name=vpc", "-force"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastOrg != "my-org" || svc.lastName != "vpc" {
		t.Fatalf("unexpected delete args %q %q", svc.lastOrg, svc.lastName)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRegistryModuleDeleteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryModuleDeleteService{}
	cmd := newRegistryModuleDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org", "-name=vpc", "-force"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastOrg != "my-org" || svc.lastName != "vpc" {
		t.Fatalf("unexpected delete args %q %q", svc.lastOrg, svc.lastName)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message")
	}
}

func TestRegistryModuleDeleteSuccessWithYesFlag(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryModuleDeleteService{}
	cmd := newRegistryModuleDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org", "-name=vpc", "-y"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastOrg != "my-org" || svc.lastName != "vpc" {
		t.Fatalf("unexpected delete args %q %q", svc.lastOrg, svc.lastName)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message")
	}
}

func TestRegistryModuleDeletePromptsWithoutForce(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("no\n")
	cmd := newRegistryModuleDeleteCommand(ui, &mockRegistryModuleDeleteService{})

	if code := cmd.Run([]string{"-organization=my-org", "-name=vpc"}); code != 0 {
		t.Fatalf("expected exit 0 on cancel, got %d", code)
	}
	if !strings.Contains(ui.OutputWriter.String(), "Deletion cancelled") {
		t.Fatalf("expected cancellation message")
	}
}

func TestRegistryModuleDeleteSynopsis(t *testing.T) {
	cmd := &RegistryModuleDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a private registry module" {
		t.Errorf("expected 'Delete a private registry module', got %q", synopsis)
	}
}

func TestRegistryModuleDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name                 string
		args                 []string
		expectedOrganization string
		expectedName         string
	}{
		{
			name:                 "basic flags",
			args:                 []string{"-organization=my-org", "-name=vpc"},
			expectedOrganization: "my-org",
			expectedName:         "vpc",
		},
		{
			name:                 "using org alias",
			args:                 []string{"-org=test-org", "-name=network"},
			expectedOrganization: "test-org",
			expectedName:         "network",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &RegistryModuleDeleteCommand{}

			flags := cmd.Meta.FlagSet("registrymodule delete")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.name, "name", "", "Module name (required)")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			if cmd.organization != tt.expectedOrganization {
				t.Errorf("expected organization %q, got %q", tt.expectedOrganization, cmd.organization)
			}

			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}
		})
	}
}
