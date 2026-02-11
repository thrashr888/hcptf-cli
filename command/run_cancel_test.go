package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newRunCancelCommand(ui cli.Ui, svc runCanceler) *RunCancelCommand {
	return &RunCancelCommand{
		Meta:   newTestMeta(ui),
		runSvc: svc,
	}
}

func TestRunCancelRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunCancelCommand(ui, &mockRunCancelService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestRunCancelHandlesCancelError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunCancelService{cancelErr: errors.New("boom")}
	cmd := newRunCancelCommand(ui, svc)

	if code := cmd.Run([]string{"-id=run-1", "-comment=test"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastCancelRun != "run-1" {
		t.Fatalf("expected cancel called")
	}
	if svc.lastCancelOpt.Comment == nil || *svc.lastCancelOpt.Comment != "test" {
		t.Fatalf("expected comment passed")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRunCancelForcePath(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunCancelService{}
	cmd := newRunCancelCommand(ui, svc)

	if code := cmd.Run([]string{"-id=run-1", "-force"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastForceRun != "run-1" {
		t.Fatalf("expected force cancel called")
	}
	if !strings.Contains(ui.OutputWriter.String(), "canceled") {
		t.Fatalf("expected success message")
	}
}
