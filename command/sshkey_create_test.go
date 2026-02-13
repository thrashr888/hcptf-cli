package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newSSHKeyCreateCommand(ui cli.Ui, svc sshKeyCreator) *SSHKeyCreateCommand {
	return &SSHKeyCreateCommand{
		Meta:      newTestMeta(ui),
		sshKeySvc: svc,
	}
}

func TestSSHKeyCreateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newSSHKeyCreateCommand(ui, &mockSSHKeyCreateService{})

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
	if code := cmd.Run([]string{"-organization=my-org", "-name=mykey"}); code != 1 {
		t.Fatalf("expected exit 1 missing value, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-value") {
		t.Fatalf("expected value error")
	}
}

func TestSSHKeyCreateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockSSHKeyCreateService{err: errors.New("boom")}
	cmd := newSSHKeyCreateCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=mykey", "-value=ssh-rsa AAAA"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if svc.lastOrg != "my-org" {
		t.Fatalf("expected organization my-org, got %s", svc.lastOrg)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestSSHKeyCreateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockSSHKeyCreateService{response: &tfe.SSHKey{ID: "sshkey-1", Name: "mykey"}}
	cmd := newSSHKeyCreateCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=mykey", "-value=ssh-rsa AAAA"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	if svc.lastOrg != "my-org" {
		t.Fatalf("expected organization my-org, got %s", svc.lastOrg)
	}
	if !strings.Contains(ui.OutputWriter.String(), "mykey") {
		t.Fatalf("expected success output with key name")
	}
}

func TestSSHKeyCreateHelp(t *testing.T) {
	cmd := &SSHKeyCreateCommand{}
	help := cmd.Help()
	if !strings.Contains(help, "sshkey create") {
		t.Fatalf("expected help text, got: %s", help)
	}
}

func TestSSHKeyCreateSynopsis(t *testing.T) {
	cmd := &SSHKeyCreateCommand{}
	syn := cmd.Synopsis()
	if syn == "" {
		t.Fatal("expected non-empty synopsis")
	}
}
