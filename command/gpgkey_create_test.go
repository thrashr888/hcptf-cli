package command

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newGPGKeyCreateCommand(ui cli.Ui, svc gpgKeyCreator) *GPGKeyCreateCommand {
	return &GPGKeyCreateCommand{
		Meta:      newTestMeta(ui),
		gpgKeySvc: svc,
	}
}

func TestGPGKeyCreateRequiresNamespace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newGPGKeyCreateCommand(ui, &mockGPGKeyCreateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-namespace") {
		t.Fatalf("expected namespace error")
	}
}

func TestGPGKeyCreateRequiresInput(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newGPGKeyCreateCommand(ui, &mockGPGKeyCreateService{})

	if code := cmd.Run([]string{"-namespace=org"}); code != 1 {
		t.Fatalf("expected exit 1 when no key provided")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "ascii-armor") {
		t.Fatalf("expected ascii armor error")
	}
}

func TestGPGKeyCreateRejectsBothInputs(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newGPGKeyCreateCommand(ui, &mockGPGKeyCreateService{})

	if code := cmd.Run([]string{"-namespace=org", "-ascii-armor=data", "-file=key.asc"}); code != 1 {
		t.Fatalf("expected exit 1 when both inputs provided")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "cannot specify both") {
		t.Fatalf("expected dual-input error")
	}
}

func TestGPGKeyCreateReadsFile(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockGPGKeyCreateService{response: &tfe.GPGKey{ID: "key-1", KeyID: "abc"}}
	cmd := newGPGKeyCreateCommand(ui, svc)

	dir := t.TempDir()
	path := filepath.Join(dir, "key.asc")
	if err := os.WriteFile(path, []byte("file-data"), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	if code := cmd.Run([]string{"-namespace=org", "-file=" + path}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if svc.lastOptions.AsciiArmor != "file-data" {
		t.Fatalf("expected ascii armor from file")
	}
}

func TestGPGKeyCreateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockGPGKeyCreateService{err: errors.New("boom")}
	cmd := newGPGKeyCreateCommand(ui, svc)

	if code := cmd.Run([]string{"-namespace=org", "-ascii-armor=data"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestGPGKeyCreateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockGPGKeyCreateService{response: &tfe.GPGKey{ID: "key-1", KeyID: "abc", Namespace: "org"}}
	cmd := newGPGKeyCreateCommand(ui, svc)

	if code := cmd.Run([]string{"-namespace=org", "-ascii-armor=data"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastRegistry != tfe.PrivateRegistry {
		t.Fatalf("expected private registry")
	}
	if !strings.Contains(ui.OutputWriter.String(), "created successfully") {
		t.Fatalf("expected success message")
	}
}
