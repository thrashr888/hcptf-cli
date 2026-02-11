package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newGPGKeyUpdateCommand(ui cli.Ui, svc gpgKeyUpdater) *GPGKeyUpdateCommand {
	return &GPGKeyUpdateCommand{
		Meta:      newTestMeta(ui),
		gpgKeySvc: svc,
	}
}

func TestGPGKeyUpdateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newGPGKeyUpdateCommand(ui, &mockGPGKeyUpdateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 namespace")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-namespace") {
		t.Fatalf("expected namespace error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-namespace=old"}); code != 1 {
		t.Fatalf("expected exit 1 key")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-namespace=old", "-key-id=abc"}); code != 1 {
		t.Fatalf("expected exit 1 new ns")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "new-namespace") {
		t.Fatalf("expected new namespace error")
	}
}

func TestGPGKeyUpdateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockGPGKeyUpdateService{err: errors.New("boom")}
	cmd := newGPGKeyUpdateCommand(ui, svc)

	if code := cmd.Run([]string{"-namespace=old", "-key-id=abc", "-new-namespace=new"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID.Namespace != "old" || svc.lastID.KeyID != "abc" {
		t.Fatalf("unexpected key id recorded")
	}
	if svc.lastOptions.Namespace != "new" {
		t.Fatalf("expected namespace option")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestGPGKeyUpdateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockGPGKeyUpdateService{response: &tfe.GPGKey{ID: "key-1", KeyID: "abc", Namespace: "new"}}
	cmd := newGPGKeyUpdateCommand(ui, svc)

	if code := cmd.Run([]string{"-namespace=old", "-key-id=abc", "-new-namespace=new"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if !strings.Contains(ui.OutputWriter.String(), "updated") {
		t.Fatalf("expected success message")
	}
}
