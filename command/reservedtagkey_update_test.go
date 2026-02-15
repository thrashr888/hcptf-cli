package command

import (
	"context"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

type mockReservedTagKeyUpdater struct {
	lastID      string
	lastOptions tfe.ReservedTagKeyUpdateOptions
	response    *tfe.ReservedTagKey
	err         error
}

func (m *mockReservedTagKeyUpdater) Update(ctx context.Context, reservedTagKeyID string, options tfe.ReservedTagKeyUpdateOptions) (*tfe.ReservedTagKey, error) {
	m.lastID = reservedTagKeyID
	m.lastOptions = options
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func newReservedTagKeyUpdateCommand(ui cli.Ui, svc reservedTagKeyUpdater) *ReservedTagKeyUpdateCommand {
	return &ReservedTagKeyUpdateCommand{
		Meta:              newTestMeta(ui),
		reservedTagKeySvc: svc,
	}
}

func TestReservedTagKeyUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newReservedTagKeyUpdateCommand(ui, &mockReservedTagKeyUpdater{})

	if code := cmd.Run([]string{"-key=environment"}); code != 1 {
		t.Fatalf("expected exit 1 missing id, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error, got %q", ui.ErrorWriter.String())
	}
}

func TestReservedTagKeyUpdateRequiresFields(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newReservedTagKeyUpdateCommand(ui, &mockReservedTagKeyUpdater{})

	if code := cmd.Run([]string{"-id=rtk-1"}); code != 1 {
		t.Fatalf("expected exit 1 missing update fields, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "at least one of -key or -disable-overrides") {
		t.Fatalf("expected update fields error, got %q", ui.ErrorWriter.String())
	}
}

func TestReservedTagKeyUpdateInvalidDisableOverrides(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newReservedTagKeyUpdateCommand(ui, &mockReservedTagKeyUpdater{})

	if code := cmd.Run([]string{"-id=rtk-1", "-disable-overrides=not-bool"}); code != 1 {
		t.Fatalf("expected exit 1 invalid bool, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "true or false") {
		t.Fatalf("expected bool validation error, got %q", ui.ErrorWriter.String())
	}
}

func TestReservedTagKeyUpdateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockReservedTagKeyUpdater{err: errors.New("boom")}
	cmd := newReservedTagKeyUpdateCommand(ui, svc)

	if code := cmd.Run([]string{"-id=rtk-1", "-key=environment"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected api error output, got %q", ui.ErrorWriter.String())
	}
}

func TestReservedTagKeyUpdateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockReservedTagKeyUpdater{
		response: &tfe.ReservedTagKey{
			ID:               "rtk-1",
			Key:              "environment",
			DisableOverrides: true,
		},
	}
	cmd := newReservedTagKeyUpdateCommand(ui, svc)

	if code := cmd.Run([]string{"-id=rtk-1", "-key=environment", "-disable-overrides=true"}); code != 0 {
		t.Fatalf("expected exit 0, got %d; err=%q", code, ui.ErrorWriter.String())
	}
	if svc.lastID != "rtk-1" {
		t.Fatalf("expected id rtk-1, got %q", svc.lastID)
	}
	if svc.lastOptions.Key == nil || *svc.lastOptions.Key != "environment" {
		t.Fatalf("expected key option to be set")
	}
	if svc.lastOptions.DisableOverrides == nil || *svc.lastOptions.DisableOverrides != true {
		t.Fatalf("expected disable-overrides option to be true")
	}
	if !strings.Contains(ui.OutputWriter.String(), "updated successfully") {
		t.Fatalf("expected success output, got %q", ui.OutputWriter.String())
	}
}
