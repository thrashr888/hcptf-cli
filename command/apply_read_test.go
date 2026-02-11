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
