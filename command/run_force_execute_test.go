package command

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

type mockRunForceExecuteService struct {
	err     error
	lastRun string
}

func (m *mockRunForceExecuteService) ForceExecute(_ context.Context, runID string) error {
	m.lastRun = runID
	return m.err
}

func newRunForceExecuteCommand(ui cli.Ui, svc runForceExecutor) *RunForceExecuteCommand {
	return &RunForceExecuteCommand{
		Meta:   newTestMeta(ui),
		runSvc: svc,
	}
}

func TestRunForceExecuteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunForceExecuteCommand(ui, &mockRunForceExecuteService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestRunForceExecuteHandlesError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunForceExecuteService{err: errors.New("boom")}
	cmd := newRunForceExecuteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=run-1"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastRun != "run-1" {
		t.Fatalf("expected force execute call")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRunForceExecuteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunForceExecuteService{}
	cmd := newRunForceExecuteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=run-1"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastRun != "run-1" {
		t.Fatalf("expected force execute call")
	}
	if !strings.Contains(ui.OutputWriter.String(), "force-executed") {
		t.Fatalf("expected success output")
	}
}
