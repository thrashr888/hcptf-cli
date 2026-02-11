package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newPlanLogsCommand(ui cli.Ui, svc planLogReader) *PlanLogsCommand {
	return &PlanLogsCommand{
		Meta:       newTestMeta(ui),
		planLogSvc: svc,
	}
}

func TestPlanLogsRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newPlanLogsCommand(ui, &mockPlanLogService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestPlanLogsHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPlanLogService{err: errors.New("boom")}
	cmd := newPlanLogsCommand(ui, svc)

	if code := cmd.Run([]string{"-id=plan-1"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID != "plan-1" {
		t.Fatalf("expected plan id recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestPlanLogsOutputsRaw(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPlanLogService{reader: strings.NewReader("hello")}
	cmd := newPlanLogsCommand(ui, svc)

	if code := cmd.Run([]string{"-id=plan-1"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if !strings.Contains(ui.OutputWriter.String(), "hello") {
		t.Fatalf("expected raw logs in output")
	}
}

func TestPlanLogsOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPlanLogService{reader: strings.NewReader("hello")}
	cmd := newPlanLogsCommand(ui, svc)

	if code := cmd.Run([]string{"-id=plan-1", "-output=json"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if !strings.Contains(ui.OutputWriter.String(), "plan-1") {
		t.Fatalf("expected plan id in json output")
	}
}
