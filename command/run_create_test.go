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

func newRunCreateCommand(ui cli.Ui, reader workspaceReader, runs runCreator) *RunCreateCommand {
	return &RunCreateCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: reader,
		runSvc:       runs,
	}
}

func TestRunCreateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunCreateCommand(ui, &mockWorkspaceReader{}, &mockRunCreateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing workspace")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected workspace error")
	}
}

func TestRunCreateHandlesWorkspaceError(t *testing.T) {
	ui := cli.NewMockUi()
	reader := &mockWorkspaceReader{err: errors.New("boom")}
	runs := &mockRunCreateService{}
	cmd := newRunCreateCommand(ui, reader, runs)

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if reader.lastOrg != "my-org" || reader.lastName != "prod" {
		t.Fatalf("unexpected read args")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRunCreateHandlesRunError(t *testing.T) {
	ui := cli.NewMockUi()
	reader := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	runs := &mockRunCreateService{err: errors.New("run failed")}
	cmd := newRunCreateCommand(ui, reader, runs)

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if runs.lastOptions.Workspace == nil || runs.lastOptions.Workspace.ID != "ws-1" {
		t.Fatalf("expected workspace in options")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "run failed") {
		t.Fatalf("expected run error output")
	}
}

func TestRunCreateHandlesNameAlias(t *testing.T) {
	ui := cli.NewMockUi()
	reader := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	runs := &mockRunCreateService{response: &tfe.Run{
		ID:        "run-1",
		Status:    tfe.RunApplied,
		Message:   "hello",
		CreatedAt: time.Unix(0, 0),
	}}
	cmd := newRunCreateCommand(ui, reader, runs)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-message=hi", "-output=json"})
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

func TestRunCreateWithRefreshOnly(t *testing.T) {
	ui := cli.NewMockUi()
	reader := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	runs := &mockRunCreateService{response: &tfe.Run{
		ID:        "run-1",
		Status:    tfe.RunApplied,
		Message:   "refresh only",
		CreatedAt: time.Unix(0, 0),
	}}
	cmd := newRunCreateCommand(ui, reader, runs)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=prod", "-refresh-only", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if runs.lastOptions.RefreshOnly == nil || !*runs.lastOptions.RefreshOnly {
		t.Fatalf("expected refresh-only option to be set, got %#v", runs.lastOptions)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "run-1" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
