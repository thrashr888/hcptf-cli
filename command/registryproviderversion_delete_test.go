package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newRegistryProviderVersionDeleteCommand(ui cli.Ui, svc registryProviderVersionDeleter) *RegistryProviderVersionDeleteCommand {
	return &RegistryProviderVersionDeleteCommand{
		Meta:       newTestMeta(ui),
		versionSvc: svc,
	}
}

func TestRegistryProviderVersionDeleteRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRegistryProviderVersionDeleteCommand(ui, &mockRegistryProviderVersionDeleteService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing name, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected name error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org", "-name=aws"}); code != 1 {
		t.Fatalf("expected exit 1 missing version, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-version") {
		t.Fatalf("expected version error")
	}
}

func TestRegistryProviderVersionDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryProviderVersionDeleteService{err: errors.New("not found")}
	cmd := newRegistryProviderVersionDeleteCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=aws", "-version=1.0.0"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "not found") {
		t.Fatalf("expected error output")
	}
}

func TestRegistryProviderVersionDeleteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryProviderVersionDeleteService{}
	cmd := newRegistryProviderVersionDeleteCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=aws", "-version=1.0.0"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted") {
		t.Fatalf("expected success output")
	}
}

func TestRegistryProviderVersionDeleteHelp(t *testing.T) {
	cmd := &RegistryProviderVersionDeleteCommand{}
	help := cmd.Help()
	if !strings.Contains(help, "registryproviderversion delete") {
		t.Fatalf("expected help text, got: %s", help)
	}
}

func TestRegistryProviderVersionDeleteSynopsis(t *testing.T) {
	cmd := &RegistryProviderVersionDeleteCommand{}
	syn := cmd.Synopsis()
	if syn == "" {
		t.Fatal("expected non-empty synopsis")
	}
}
