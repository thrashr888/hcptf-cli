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

func newRunShowCommand(ui cli.Ui, svc runReader) *RunShowCommand {
	return &RunShowCommand{
		Meta:   newTestMeta(ui),
		runSvc: svc,
	}
}

func TestRunShowRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunShowCommand(ui, &mockRunReadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestRunShowHandlesError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunReadService{err: errors.New("boom")}
	cmd := newRunShowCommand(ui, svc)

	if code := cmd.Run([]string{"-id=run-1"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastRun != "run-1" {
		t.Fatalf("expected read called")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRunShowOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunReadService{response: &tfe.Run{
		ID:        "run-1",
		Status:    tfe.RunApplied,
		Message:   "hello",
		IsDestroy: false,
		Source:    tfe.RunSourceUI,
		CreatedAt: time.Unix(0, 0),
		Plan: &tfe.Plan{
			ResourceAdditions:    1,
			ResourceChanges:      2,
			ResourceDestructions: 3,
		},
	}}
	cmd := newRunShowCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=run-1", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "run-1" {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestRunShowPassesIncludeOptions(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunReadService{response: &tfe.Run{
		ID:      "run-1",
		Status:  tfe.RunApplied,
		Source:  tfe.RunSourceUI,
		Plan:    &tfe.Plan{},
		Message: "ok",
	}}
	cmd := newRunShowCommand(ui, svc)

	if code := cmd.Run([]string{"-id=run-1", "-include=workspace,plan"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastReadOptions == nil {
		t.Fatalf("expected read-with-options call")
	}
	if len(svc.lastReadOptions.Include) != 2 {
		t.Fatalf("expected include options, got %#v", svc.lastReadOptions.Include)
	}
}
