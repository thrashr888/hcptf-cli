package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
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

	code := cmd.Run([]string{"-organization=my-org", "-name=aws", "-version=1.0.0", "-force"})
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

	code := cmd.Run([]string{"-organization=my-org", "-name=aws", "-version=1.0.0", "-force"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
	if svc.lastID != (tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			RegistryName: tfe.PrivateRegistry,
			Namespace:    "my-org",
			Name:         "aws",
		},
		Version: "1.0.0",
	}) {
		t.Fatalf("expected version id %#v, got %#v", tfe.RegistryProviderVersionID{
			RegistryProviderID: tfe.RegistryProviderID{
				RegistryName: tfe.PrivateRegistry,
				Namespace:    "my-org",
				Name:         "aws",
			},
			Version: "1.0.0",
		}, svc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted") {
		t.Fatalf("expected success output")
	}
}

func TestRegistryProviderVersionDeletePromptsWithoutForce(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("no\n")
	cmd := newRegistryProviderVersionDeleteCommand(ui, &mockRegistryProviderVersionDeleteService{})

	code := cmd.Run([]string{"-organization=my-org", "-name=aws", "-version=1.0.0"})
	if code != 0 {
		t.Fatalf("expected exit 0 on cancel, got %d", code)
	}
	if !strings.Contains(ui.OutputWriter.String(), "Deletion cancelled") {
		t.Fatalf("expected cancellation message")
	}
}

func TestRegistryProviderVersionDeleteHelp(t *testing.T) {
	cmd := &RegistryProviderVersionDeleteCommand{}
	help := cmd.Help()
	if !strings.Contains(help, "registryproviderversion delete") {
		t.Fatalf("expected help text, got: %s", help)
	}
	if !strings.Contains(help, "-force") {
		t.Fatalf("expected help text to mention force flag")
	}
	if !strings.Contains(help, "-y") {
		t.Fatalf("expected help text to mention y flag")
	}
}

func TestRegistryProviderVersionDeleteSuccessWithYesFlag(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryProviderVersionDeleteService{}
	cmd := newRegistryProviderVersionDeleteCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=aws", "-version=1.0.0", "-y"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted") {
		t.Fatalf("expected success output")
	}
}

func TestRegistryProviderVersionDeleteSynopsis(t *testing.T) {
	cmd := &RegistryProviderVersionDeleteCommand{}
	syn := cmd.Synopsis()
	if syn == "" {
		t.Fatal("expected non-empty synopsis")
	}
}
