package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newRegistryProviderVersionCreateCommand(ui cli.Ui, svc registryProviderVersionCreator) *RegistryProviderVersionCreateCommand {
	return &RegistryProviderVersionCreateCommand{
		Meta:       newTestMeta(ui),
		versionSvc: svc,
	}
}

func TestRegistryProviderVersionCreateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRegistryProviderVersionCreateCommand(ui, &mockRegistryProviderVersionCreateService{})

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

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org", "-name=aws", "-version=1.0.0"}); code != 1 {
		t.Fatalf("expected exit 1 missing key-id, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-key-id") {
		t.Fatalf("expected key-id error")
	}
}

func TestRegistryProviderVersionCreateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryProviderVersionCreateService{err: errors.New("boom")}
	cmd := newRegistryProviderVersionCreateCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=aws", "-version=1.0.0", "-key-id=ABC123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRegistryProviderVersionCreateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRegistryProviderVersionCreateService{response: &tfe.RegistryProviderVersion{
		ID:        "provver-1",
		Version:   "1.0.0",
		KeyID:     "ABC123",
		Protocols: []string{"5.0", "6.0"},
		Links: map[string]interface{}{
			"shasums-upload":     "https://example.com/shasums",
			"shasums-sig-upload": "https://example.com/shasums-sig",
		},
	}}
	cmd := newRegistryProviderVersionCreateCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=aws", "-version=1.0.0", "-key-id=ABC123"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
	if !strings.Contains(ui.OutputWriter.String(), "1.0.0") {
		t.Fatalf("expected success output with version")
	}
}

func TestRegistryProviderVersionCreateHelp(t *testing.T) {
	cmd := &RegistryProviderVersionCreateCommand{}
	help := cmd.Help()
	if !strings.Contains(help, "registryproviderversion create") {
		t.Fatalf("expected help text, got: %s", help)
	}
}

func TestRegistryProviderVersionCreateSynopsis(t *testing.T) {
	cmd := &RegistryProviderVersionCreateCommand{}
	syn := cmd.Synopsis()
	if syn == "" {
		t.Fatal("expected non-empty synopsis")
	}
}
