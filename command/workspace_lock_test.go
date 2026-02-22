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

type mockWorkspaceLockService struct {
	response    *tfe.Workspace
	err         error
	lastID      string
	lastOptions tfe.WorkspaceLockOptions
}

func (m *mockWorkspaceLockService) Lock(_ context.Context, workspaceID string, options tfe.WorkspaceLockOptions) (*tfe.Workspace, error) {
	m.lastID = workspaceID
	m.lastOptions = options
	return m.response, m.err
}

func newWorkspaceLockCommand(ui cli.Ui, reader workspaceReader, locker workspaceLocker) *WorkspaceLockCommand {
	return &WorkspaceLockCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: reader,
		lockSvc:      locker,
	}
}

func TestWorkspaceLockRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceLockCommand(ui, &mockWorkspaceReader{}, &mockWorkspaceLockService{})

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

func TestWorkspaceLockHandlesWorkspaceReadError(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceLockCommand(ui, &mockWorkspaceReader{err: errors.New("read failed")}, &mockWorkspaceLockService{})

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "read failed") {
		t.Fatalf("expected read error output")
	}
}

func TestWorkspaceLockHandlesLockError(t *testing.T) {
	ui := cli.NewMockUi()
	locker := &mockWorkspaceLockService{err: errors.New("lock failed")}
	cmd := newWorkspaceLockCommand(ui, &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}, locker)

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod", "-reason=maintenance"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if locker.lastID != "ws-1" {
		t.Fatalf("expected lock call")
	}
	if locker.lastOptions.Reason == nil || *locker.lastOptions.Reason != "maintenance" {
		t.Fatalf("expected reason option")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "lock failed") {
		t.Fatalf("expected lock error output")
	}
}

func TestWorkspaceLockOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	locker := &mockWorkspaceLockService{response: &tfe.Workspace{
		ID:     "ws-1",
		Name:   "prod",
		Locked: true,
	}}
	cmd := newWorkspaceLockCommand(ui, &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}, locker)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=prod", "-reason=maintenance", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "ws-1" || data["Locked"] != true {
		t.Fatalf("unexpected data: %#v", data)
	}
}
