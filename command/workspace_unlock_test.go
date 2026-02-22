package command

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

type mockWorkspaceUnlockService struct {
	response *tfe.Workspace
	err      error
	lastID   string
}

func (m *mockWorkspaceUnlockService) Unlock(_ context.Context, workspaceID string) (*tfe.Workspace, error) {
	m.lastID = workspaceID
	return m.response, m.err
}

func newWorkspaceUnlockCommand(ui cli.Ui, reader workspaceReader, unlocker workspaceUnlocker) *WorkspaceUnlockCommand {
	return &WorkspaceUnlockCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: reader,
		unlockSvc:    unlocker,
	}
}

func TestWorkspaceUnlockRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceUnlockCommand(ui, &mockWorkspaceReader{}, &mockWorkspaceUnlockService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected organization error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected name error")
	}
}

func TestWorkspaceUnlockHandlesUnlockError(t *testing.T) {
	ui := cli.NewMockUi()
	unlocker := &mockWorkspaceUnlockService{err: errors.New("unlock failed")}
	cmd := newWorkspaceUnlockCommand(ui, &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}, unlocker)

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if unlocker.lastID != "ws-1" {
		t.Fatalf("expected unlock call")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "unlock failed") {
		t.Fatalf("expected unlock error output")
	}
}

func TestWorkspaceUnlockOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	unlocker := &mockWorkspaceUnlockService{response: &tfe.Workspace{
		ID:     "ws-1",
		Name:   "prod",
		Locked: false,
	}}
	cmd := newWorkspaceUnlockCommand(ui, &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}, unlocker)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=prod", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "ws-1" || data["Locked"] != false {
		t.Fatalf("unexpected data: %#v", data)
	}
}
