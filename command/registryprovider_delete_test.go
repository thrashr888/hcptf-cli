package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newRegistryProviderDeleteCommand(ui cli.Ui, svc registryProviderDeleter) *RegistryProviderDeleteCommand {
	return &RegistryProviderDeleteCommand{
		Meta:                newTestMeta(ui),
		registryProviderSvc: svc,
	}
}

func TestRegistryProviderDeleteRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRegistryProviderDeleteCommand(ui, &mockRegistryProviderDeleteService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 org")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 name")
	}
}

func TestRegistryProviderDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryProviderDeleteService{err: errors.New("boom")}
	cmd := newRegistryProviderDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org", "-name=aws"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID.Name != "aws" {
		t.Fatalf("expected provider name recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRegistryProviderDeleteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryProviderDeleteService{}
	cmd := newRegistryProviderDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org", "-name=aws"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastID.RegistryName != tfe.RegistryName("private") {
		t.Fatalf("expected private registry")
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message")
	}
}
