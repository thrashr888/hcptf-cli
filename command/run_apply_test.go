package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newRunApplyCommand(ui cli.Ui, runs runApplier) *RunApplyCommand {
	return &RunApplyCommand{
		Meta:   newTestMeta(ui),
		runSvc: runs,
	}
}

func TestRunApplyRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunApplyCommand(ui, &mockRunApplyService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestRunApplyHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	runs := &mockRunApplyService{err: errors.New("boom")}
	cmd := newRunApplyCommand(ui, runs)

	if code := cmd.Run([]string{"-id=run-1", "-comment=ok"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if runs.lastRun != "run-1" {
		t.Fatalf("expected run id recorded")
	}
	if runs.lastOptions.Comment == nil || *runs.lastOptions.Comment != "ok" {
		t.Fatalf("expected comment option")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRunApplySuccess(t *testing.T) {
	ui := cli.NewMockUi()
	runs := &mockRunApplyService{}
	cmd := newRunApplyCommand(ui, runs)

	if code := cmd.Run([]string{"-id=run-1"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if runs.lastRun != "run-1" {
		t.Fatalf("expected run id recorded")
	}
	if !strings.Contains(ui.OutputWriter.String(), "Run run-1") {
		t.Fatalf("expected success message")
	}
}
