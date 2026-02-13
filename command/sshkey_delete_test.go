package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newSSHKeyDeleteCommand(ui cli.Ui, svc sshKeyDeleter) *SSHKeyDeleteCommand {
	return &SSHKeyDeleteCommand{
		Meta:      newTestMeta(ui),
		sshKeySvc: svc,
	}
}

func TestSSHKeyDeleteRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newSSHKeyDeleteCommand(ui, &mockSSHKeyDeleteService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing id, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error, got %q", ui.ErrorWriter.String())
	}
}

func TestSSHKeyDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockSSHKeyDeleteService{err: errors.New("boom")}
	cmd := newSSHKeyDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=sshkey-123", "-force"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastID != "sshkey-123" {
		t.Fatalf("expected lastID sshkey-123, got %q", svc.lastID)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got %q", ui.ErrorWriter.String())
	}
}

func TestSSHKeyDeleteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockSSHKeyDeleteService{}
	cmd := newSSHKeyDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=sshkey-123", "-force"}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if svc.lastID != "sshkey-123" {
		t.Fatalf("expected lastID sshkey-123, got %q", svc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message, got %q", ui.OutputWriter.String())
	}
}
