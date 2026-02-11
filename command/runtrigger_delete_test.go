package command

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

type mockRunTriggerDeleteService struct {
	err    error
	lastID string
}

func (m *mockRunTriggerDeleteService) Delete(_ context.Context, runTriggerID string) error {
	m.lastID = runTriggerID
	return m.err
}

func newRunTriggerDeleteCommand(ui cli.Ui, rtSvc runTriggerDeleter) *RunTriggerDeleteCommand {
	return &RunTriggerDeleteCommand{
		Meta:          newTestMeta(ui),
		runTriggerSvc: rtSvc,
	}
}

func TestRunTriggerDeleteCommandRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTriggerDeleteCommand(ui, nil)

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestRunTriggerDeleteCommandHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	rtSvc := &mockRunTriggerDeleteService{err: errors.New("not found")}
	cmd := newRunTriggerDeleteCommand(ui, rtSvc)

	code := cmd.Run([]string{"-id=rt-123", "-force"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if rtSvc.lastID != "rt-123" {
		t.Fatalf("expected ID rt-123, got %s", rtSvc.lastID)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "not found") {
		t.Fatalf("expected error output, got %q", out)
	}
}

func TestRunTriggerDeleteCommandSuccessWithForce(t *testing.T) {
	ui := cli.NewMockUi()
	rtSvc := &mockRunTriggerDeleteService{}
	cmd := newRunTriggerDeleteCommand(ui, rtSvc)

	code := cmd.Run([]string{"-id=rt-abc123", "-force"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if rtSvc.lastID != "rt-abc123" {
		t.Fatalf("expected ID rt-abc123, got %s", rtSvc.lastID)
	}

	if out := ui.OutputWriter.String(); !strings.Contains(out, "successfully") {
		t.Fatalf("expected success message, got %q", out)
	}
}

func TestRunTriggerDeleteCommandCancelsWithoutConfirmation(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("no\n")
	rtSvc := &mockRunTriggerDeleteService{}
	cmd := newRunTriggerDeleteCommand(ui, rtSvc)

	code := cmd.Run([]string{"-id=rt-abc123"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if rtSvc.lastID != "" {
		t.Fatalf("expected no deletion, but got ID %s", rtSvc.lastID)
	}

	if out := ui.OutputWriter.String(); !strings.Contains(out, "cancelled") {
		t.Fatalf("expected cancellation message, got %q", out)
	}
}

func TestRunTriggerDeleteCommandSuccessWithConfirmation(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("yes\n")
	rtSvc := &mockRunTriggerDeleteService{}
	cmd := newRunTriggerDeleteCommand(ui, rtSvc)

	code := cmd.Run([]string{"-id=rt-abc123"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if rtSvc.lastID != "rt-abc123" {
		t.Fatalf("expected ID rt-abc123, got %s", rtSvc.lastID)
	}

	if out := ui.OutputWriter.String(); !strings.Contains(out, "successfully") {
		t.Fatalf("expected success message, got %q", out)
	}
}
