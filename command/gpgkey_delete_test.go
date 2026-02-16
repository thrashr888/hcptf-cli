package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newGPGKeyDeleteCommand(ui cli.Ui, svc gpgKeyDeleter) *GPGKeyDeleteCommand {
	return &GPGKeyDeleteCommand{
		Meta:      newTestMeta(ui),
		gpgKeySvc: svc,
	}
}

func TestGPGKeyDeleteRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newGPGKeyDeleteCommand(ui, &mockGPGKeyDeleteService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 namespace")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-namespace") {
		t.Fatalf("expected namespace error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-namespace=org"}); code != 1 {
		t.Fatalf("expected exit 1 key")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-key-id") {
		t.Fatalf("expected key error")
	}
}

func TestGPGKeyDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockGPGKeyDeleteService{err: errors.New("boom")}
	cmd := newGPGKeyDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-namespace=org", "-key-id=abc", "-force"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID.Namespace != "org" || svc.lastID.KeyID != "abc" {
		t.Fatalf("unexpected key id recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestGPGKeyDeleteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockGPGKeyDeleteService{}
	cmd := newGPGKeyDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-namespace=org", "-key-id=abc", "-force"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastID.RegistryName != tfe.PrivateRegistry {
		t.Fatalf("expected private registry")
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message")
	}
}

func TestGPGKeyDeletePromptsWithoutForce(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("no\n")
	cmd := newGPGKeyDeleteCommand(ui, &mockGPGKeyDeleteService{})

	if code := cmd.Run([]string{"-namespace=org", "-key-id=abc"}); code != 0 {
		t.Fatalf("expected exit 0 on cancel, got %d", code)
	}
	if !strings.Contains(ui.OutputWriter.String(), "Deletion cancelled") {
		t.Fatalf("expected cancellation message")
	}
}

func TestGPGKeyDeleteSuccessWithYesFlag(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockGPGKeyDeleteService{}
	cmd := newGPGKeyDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-namespace=org", "-key-id=abc", "-y"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastID.Namespace != "org" || svc.lastID.KeyID != "abc" {
		t.Fatalf("unexpected key id recorded")
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message")
	}
}
