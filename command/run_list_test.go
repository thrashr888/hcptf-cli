package command

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

type mockRunService struct {
	response        *tfe.RunList
	err             error
	lastWorkspaceID string
	lastOptions     *tfe.RunListOptions
}

func (m *mockRunService) List(_ context.Context, workspaceID string, options *tfe.RunListOptions) (*tfe.RunList, error) {
	m.lastWorkspaceID = workspaceID
	m.lastOptions = options
	return m.response, m.err
}

func newRunListCommand(ui cli.Ui, ws workspaceReader, runs runLister) *RunListCommand {
	return &RunListCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: ws,
		runSvc:       runs,
	}
}

func TestRunListRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunListCommand(ui, &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws"}}, &mockRunService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org, got %d", code)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing workspace, got %d", code)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected workspace error")
	}
}

func TestRunListHandlesWorkspaceError(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{err: errors.New("boom")}
	runs := &mockRunService{}
	cmd := newRunListCommand(ui, ws, runs)

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if ws.lastOrg != "my-org" || ws.lastName != "prod" {
		t.Fatalf("unexpected workspace request: %#v", ws)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected workspace error output")
	}
}

func TestRunListHandlesRunError(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-123"}}
	runs := &mockRunService{err: errors.New("run failed")}
	cmd := newRunListCommand(ui, ws, runs)

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if runs.lastWorkspaceID != "ws-123" {
		t.Fatalf("expected runs for ws-123, got %s", runs.lastWorkspaceID)
	}

	if runs.lastOptions == nil || runs.lastOptions.ListOptions.PageSize != 50 {
		t.Fatalf("expected run list options, got %#v", runs.lastOptions)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "run failed") {
		t.Fatalf("expected run error output")
	}
}

func TestRunListOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-123"}}
	runs := &mockRunService{response: &tfe.RunList{Items: []*tfe.Run{{
		ID:        "run-1",
		Status:    tfe.RunApplied,
		Source:    tfe.RunSourceUI,
		Message:   "hello world",
		CreatedAt: time.Unix(0, 0),
	}}}}
	cmd := newRunListCommand(ui, ws, runs)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=prod", "-output=json"})
	})

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	var rows []map[string]string
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if len(rows) != 1 || rows[0]["ID"] != "run-1" {
		t.Fatalf("unexpected rows: %#v", rows)
	}
}
