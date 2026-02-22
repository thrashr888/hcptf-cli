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

type mockWorkspaceForceUnlockService struct {
	response *tfe.Workspace
	err      error
	lastID   string
}

func (m *mockWorkspaceForceUnlockService) ForceUnlock(_ context.Context, workspaceID string) (*tfe.Workspace, error) {
	m.lastID = workspaceID
	return m.response, m.err
}

func newWorkspaceForceUnlockCommand(ui cli.Ui, reader workspaceReader, forceUnlocker workspaceForceUnlocker) *WorkspaceForceUnlockCommand {
	return &WorkspaceForceUnlockCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: reader,
		forceSvc:     forceUnlocker,
	}
}

func TestWorkspaceForceUnlockRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceForceUnlockCommand(ui, &mockWorkspaceReader{}, &mockWorkspaceForceUnlockService{})

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

func TestWorkspaceForceUnlockHandlesError(t *testing.T) {
	ui := cli.NewMockUi()
	forceSvc := &mockWorkspaceForceUnlockService{err: errors.New("force unlock failed")}
	cmd := newWorkspaceForceUnlockCommand(ui, &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}, forceSvc)

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if forceSvc.lastID != "ws-1" {
		t.Fatalf("expected force unlock call")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "force unlock failed") {
		t.Fatalf("expected error output")
	}
}

func TestWorkspaceForceUnlockOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	forceSvc := &mockWorkspaceForceUnlockService{response: &tfe.Workspace{
		ID:     "ws-1",
		Name:   "prod",
		Locked: false,
	}}
	cmd := newWorkspaceForceUnlockCommand(ui, &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}, forceSvc)

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
