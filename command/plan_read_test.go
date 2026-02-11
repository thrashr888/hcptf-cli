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

func newPlanReadCommand(ui cli.Ui, svc planReader) *PlanReadCommand {
	return &PlanReadCommand{
		Meta:    newTestMeta(ui),
		planSvc: svc,
	}
}

func TestPlanReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newPlanReadCommand(ui, &mockPlanService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestPlanReadHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPlanService{err: errors.New("boom")}
	cmd := newPlanReadCommand(ui, svc)

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

func TestPlanReadOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPlanService{response: &tfe.Plan{
		ID:                "plan-1",
		Status:            tfe.PlanFinished,
		HasChanges:        true,
		ResourceAdditions: 1,
		StatusTimestamps: &tfe.PlanStatusTimestamps{
			QueuedAt:   time.Unix(0, 0),
			StartedAt:  time.Unix(100, 0),
			FinishedAt: time.Unix(200, 0),
		},
	}}
	cmd := newPlanReadCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=plan-1", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "plan-1" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
