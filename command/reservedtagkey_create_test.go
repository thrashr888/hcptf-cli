package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newReservedTagKeyCreateCommand(ui cli.Ui, svc reservedTagKeyCreator) *ReservedTagKeyCreateCommand {
	return &ReservedTagKeyCreateCommand{
		Meta:              newTestMeta(ui),
		reservedTagKeySvc: svc,
	}
}

func TestReservedTagKeyCreateRequiresOrg(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newReservedTagKeyCreateCommand(ui, &mockReservedTagKeyCreateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing organization, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected organization error, got %q", ui.ErrorWriter.String())
	}
}

func TestReservedTagKeyCreateRequiresKey(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newReservedTagKeyCreateCommand(ui, &mockReservedTagKeyCreateService{})

	if code := cmd.Run([]string{"-org=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing key, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-key") {
		t.Fatalf("expected key error, got %q", ui.ErrorWriter.String())
	}
}

func TestReservedTagKeyCreateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockReservedTagKeyCreateService{err: errors.New("boom")}
	cmd := newReservedTagKeyCreateCommand(ui, svc)

	if code := cmd.Run([]string{"-org=my-org", "-key=environment"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastOrg != "my-org" {
		t.Fatalf("expected lastOrg my-org, got %q", svc.lastOrg)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got %q", ui.ErrorWriter.String())
	}
}

func TestReservedTagKeyCreateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockReservedTagKeyCreateService{
		response: &tfe.ReservedTagKey{
			ID:  "rtk-123",
			Key: "environment",
		},
	}
	cmd := newReservedTagKeyCreateCommand(ui, svc)

	if code := cmd.Run([]string{"-org=my-org", "-key=environment"}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if svc.lastOrg != "my-org" {
		t.Fatalf("expected lastOrg my-org, got %q", svc.lastOrg)
	}
	if !strings.Contains(ui.OutputWriter.String(), "created successfully") {
		t.Fatalf("expected success message, got %q", ui.OutputWriter.String())
	}
}
