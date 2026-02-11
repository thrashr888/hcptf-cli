package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newApplyLogsCommand(ui cli.Ui, svc applyLogReader) *ApplyLogsCommand {
	return &ApplyLogsCommand{
		Meta:        newTestMeta(ui),
		applyLogSvc: svc,
	}
}

func TestApplyLogsRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newApplyLogsCommand(ui, &mockApplyLogService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestApplyLogsHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockApplyLogService{err: errors.New("boom")}
	cmd := newApplyLogsCommand(ui, svc)

	if code := cmd.Run([]string{"-id=apply-1"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID != "apply-1" {
		t.Fatalf("expected apply id recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestApplyLogsOutputsRaw(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockApplyLogService{reader: strings.NewReader("logdata")}
	cmd := newApplyLogsCommand(ui, svc)

	if code := cmd.Run([]string{"-id=apply-1"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if !strings.Contains(ui.OutputWriter.String(), "logdata") {
		t.Fatalf("expected raw output")
	}
}

func TestApplyLogsOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockApplyLogService{reader: strings.NewReader("logdata")}
	cmd := newApplyLogsCommand(ui, svc)

	if code := cmd.Run([]string{"-id=apply-1", "-output=json"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if !strings.Contains(ui.OutputWriter.String(), "apply-1") {
		t.Fatalf("expected apply id in json output")
	}
}
