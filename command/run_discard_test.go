package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newRunDiscardCommand(ui cli.Ui, svc runDiscarder) *RunDiscardCommand {
	return &RunDiscardCommand{
		Meta:   newTestMeta(ui),
		runSvc: svc,
	}
}

func TestRunDiscardRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunDiscardCommand(ui, &mockRunDiscardService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestRunDiscardHandlesError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunDiscardService{err: errors.New("boom")}
	cmd := newRunDiscardCommand(ui, svc)

	if code := cmd.Run([]string{"-id=run-1", "-comment=test"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastRun != "run-1" {
		t.Fatalf("expected discard called")
	}
	if svc.lastOption.Comment == nil || *svc.lastOption.Comment != "test" {
		t.Fatalf("expected comment option")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRunDiscardSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunDiscardService{}
	cmd := newRunDiscardCommand(ui, svc)

	if code := cmd.Run([]string{"-id=run-1"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastRun != "run-1" {
		t.Fatalf("expected discard called")
	}
	if !strings.Contains(ui.OutputWriter.String(), "discarded") {
		t.Fatalf("expected success message")
	}
}
