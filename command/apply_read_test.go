package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newApplyReadCommand(ui cli.Ui, svc applyReader) *ApplyReadCommand {
	return &ApplyReadCommand{
		Meta:     newTestMeta(ui),
		applySvc: svc,
	}
}

func TestApplyReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newApplyReadCommand(ui, &mockApplyService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestApplyReadHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockApplyService{err: errors.New("boom")}
	cmd := newApplyReadCommand(ui, svc)

	if code := cmd.Run([]string{"-id=apply-1"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID != "apply-1" {
		t.Fatalf("expected id recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestApplyReadOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockApplyService{response: &tfe.Apply{
		ID:                "apply-1",
		Status:            tfe.ApplyFinished,
		ResourceAdditions: 1,
		StatusTimestamps: &tfe.ApplyStatusTimestamps{
			QueuedAt:   time.Unix(0, 0),
			StartedAt:  time.Unix(100, 0),
			FinishedAt: time.Unix(200, 0),
		},
	}}
	cmd := newApplyReadCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=apply-1", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "apply-1" {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestApplyReadRunID(t *testing.T) {
	ui := cli.NewMockUi()
	applySvc := &mockApplyService{response: &tfe.Apply{
		ID:     "apply-123",
		Status: tfe.ApplyFinished,
	}}
	runSvc := &mockRunReadService{response: &tfe.Run{
		ID: "run-123",
		Apply: &tfe.Apply{
			ID: "apply-123",
		},
	}}
	cmd := &ApplyReadCommand{
		Meta:     newTestMeta(ui),
		applySvc: applySvc,
		runSvc:   runSvc,
	}

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-run-id=run-123", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "apply-123" {
		t.Fatalf("unexpected data: %#v", data)
	}
	if runSvc.lastRun != "run-123" {
		t.Fatalf("expected run id recorded, got: %s", runSvc.lastRun)
	}
	if applySvc.lastID != "apply-123" {
		t.Fatalf("expected apply read id recorded, got: %s", applySvc.lastID)
	}
}

func TestApplyReadRunReadFailure(t *testing.T) {
	ui := cli.NewMockUi()
	runSvc := &mockRunReadService{err: errors.New("run failed")}
	cmd := &ApplyReadCommand{
		Meta:     newTestMeta(ui),
		runSvc:   runSvc,
		applySvc: &mockApplyService{},
	}

	if code := cmd.Run([]string{"-run-id=run-123"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if runSvc.lastRun != "run-123" {
		t.Fatalf("expected run id recorded, got: %s", runSvc.lastRun)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "Error reading run: run failed") {
		t.Fatalf("expected run read error output, got: %s", ui.ErrorWriter.String())
	}
}

func TestApplyReadRunWithNoApply(t *testing.T) {
	ui := cli.NewMockUi()
	runSvc := &mockRunReadService{response: &tfe.Run{
		ID: "run-123",
	}}
	cmd := &ApplyReadCommand{
		Meta:     newTestMeta(ui),
		runSvc:   runSvc,
		applySvc: &mockApplyService{},
	}

	if code := cmd.Run([]string{"-id=run-123"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "Error: run has no apply") {
		t.Fatalf("expected no apply error, got: %s", ui.ErrorWriter.String())
	}
}
