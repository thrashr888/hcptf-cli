package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newRunLogsCommand(ui cli.Ui, runSvc runReader, planLogSvc planLogReader, applyLogSvc applyLogReader) *RunLogsCommand {
	return &RunLogsCommand{
		Meta:        newTestMeta(ui),
		runSvc:      runSvc,
		planLogSvc:  planLogSvc,
		applyLogSvc: applyLogSvc,
	}
}

func TestRunLogsRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunLogsCommand(ui, &mockRunReadService{}, &mockPlanLogService{}, &mockApplyLogService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error, got: %s", ui.ErrorWriter.String())
	}
}

func TestRunLogsAutoSelectsApply(t *testing.T) {
	ui := cli.NewMockUi()
	runSvc := &mockRunReadService{
		response: &tfe.Run{
			ID: "run-1",
			Plan: &tfe.Plan{
				ID: "plan-1",
			},
			Apply: &tfe.Apply{
				ID: "apply-1",
			},
		},
	}
	applyLogSvc := &mockApplyLogService{reader: strings.NewReader("apply log output")}
	planLogSvc := &mockPlanLogService{reader: strings.NewReader("plan log output")}
	cmd := newRunLogsCommand(ui, runSvc, planLogSvc, applyLogSvc)

	if code := cmd.Run([]string{"-id=run-1"}); code != 0 {
		t.Fatalf("expected exit 0, got %d; errors: %s", code, ui.ErrorWriter.String())
	}
	if applyLogSvc.lastID != "apply-1" {
		t.Fatalf("expected apply log service called with apply-1, got %q", applyLogSvc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "apply log output") {
		t.Fatalf("expected apply logs in output, got: %s", ui.OutputWriter.String())
	}
}

func TestRunLogsAutoSelectsPlan(t *testing.T) {
	ui := cli.NewMockUi()
	runSvc := &mockRunReadService{
		response: &tfe.Run{
			ID: "run-2",
			Plan: &tfe.Plan{
				ID: "plan-2",
			},
			// Apply is nil — auto should select plan
		},
	}
	planLogSvc := &mockPlanLogService{reader: strings.NewReader("plan log output")}
	applyLogSvc := &mockApplyLogService{}
	cmd := newRunLogsCommand(ui, runSvc, planLogSvc, applyLogSvc)

	if code := cmd.Run([]string{"-id=run-2"}); code != 0 {
		t.Fatalf("expected exit 0, got %d; errors: %s", code, ui.ErrorWriter.String())
	}
	if planLogSvc.lastID != "plan-2" {
		t.Fatalf("expected plan log service called with plan-2, got %q", planLogSvc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "plan log output") {
		t.Fatalf("expected plan logs in output, got: %s", ui.OutputWriter.String())
	}
}

func TestRunLogsPlanPhaseFlag(t *testing.T) {
	ui := cli.NewMockUi()
	runSvc := &mockRunReadService{
		response: &tfe.Run{
			ID: "run-3",
			Plan: &tfe.Plan{
				ID: "plan-3",
			},
			Apply: &tfe.Apply{
				ID: "apply-3",
			},
		},
	}
	planLogSvc := &mockPlanLogService{reader: strings.NewReader("plan forced")}
	applyLogSvc := &mockApplyLogService{reader: strings.NewReader("apply forced")}
	cmd := newRunLogsCommand(ui, runSvc, planLogSvc, applyLogSvc)

	// Even though apply exists, -phase=plan should show plan logs
	if code := cmd.Run([]string{"-id=run-3", "-phase=plan"}); code != 0 {
		t.Fatalf("expected exit 0, got %d; errors: %s", code, ui.ErrorWriter.String())
	}
	if planLogSvc.lastID != "plan-3" {
		t.Fatalf("expected plan log service called with plan-3, got %q", planLogSvc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "plan forced") {
		t.Fatalf("expected plan logs in output, got: %s", ui.OutputWriter.String())
	}
}

func TestRunLogsApplyPhaseFlag(t *testing.T) {
	ui := cli.NewMockUi()
	runSvc := &mockRunReadService{
		response: &tfe.Run{
			ID: "run-4",
			Plan: &tfe.Plan{
				ID: "plan-4",
			},
			Apply: &tfe.Apply{
				ID: "apply-4",
			},
		},
	}
	planLogSvc := &mockPlanLogService{reader: strings.NewReader("plan data")}
	applyLogSvc := &mockApplyLogService{reader: strings.NewReader("apply data")}
	cmd := newRunLogsCommand(ui, runSvc, planLogSvc, applyLogSvc)

	if code := cmd.Run([]string{"-id=run-4", "-phase=apply"}); code != 0 {
		t.Fatalf("expected exit 0, got %d; errors: %s", code, ui.ErrorWriter.String())
	}
	if applyLogSvc.lastID != "apply-4" {
		t.Fatalf("expected apply log service called with apply-4, got %q", applyLogSvc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "apply data") {
		t.Fatalf("expected apply logs in output, got: %s", ui.OutputWriter.String())
	}
}

func TestRunLogsHandlesError(t *testing.T) {
	ui := cli.NewMockUi()
	runSvc := &mockRunReadService{err: errors.New("run not found")}
	cmd := newRunLogsCommand(ui, runSvc, &mockPlanLogService{}, &mockApplyLogService{})

	if code := cmd.Run([]string{"-id=run-bad"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if runSvc.lastRun != "run-bad" {
		t.Fatalf("expected run service called with run-bad, got %q", runSvc.lastRun)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "run not found") {
		t.Fatalf("expected error in output, got: %s", ui.ErrorWriter.String())
	}
}
