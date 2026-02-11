package command

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

type mockRunTriggerCreateService struct {
	response    *tfe.RunTrigger
	err         error
	lastWSID    string
	lastOptions tfe.RunTriggerCreateOptions
}

func (m *mockRunTriggerCreateService) Create(_ context.Context, workspaceID string, options tfe.RunTriggerCreateOptions) (*tfe.RunTrigger, error) {
	m.lastWSID = workspaceID
	m.lastOptions = options
	return m.response, m.err
}

type mockWorkspaceMultiReadService struct {
	workspaces map[string]*tfe.Workspace
	err        error
}

func (m *mockWorkspaceMultiReadService) Read(_ context.Context, organization, workspace string) (*tfe.Workspace, error) {
	if m.err != nil {
		return nil, m.err
	}
	if ws, ok := m.workspaces[workspace]; ok {
		return ws, nil
	}
	return nil, errors.New("workspace not found")
}

func newRunTriggerCreateCommand(ui cli.Ui, wsSvc workspaceReader, rtSvc runTriggerCreator) *RunTriggerCreateCommand {
	return &RunTriggerCreateCommand{
		Meta:          newTestMeta(ui),
		workspaceSvc:  wsSvc,
		runTriggerSvc: rtSvc,
	}
}

func TestRunTriggerCreateCommandRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTriggerCreateCommand(ui, nil, nil)

	if code := cmd.Run([]string{"-workspace=test", "-source-workspace=source"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestRunTriggerCreateCommandRequiresWorkspace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTriggerCreateCommand(ui, nil, nil)

	if code := cmd.Run([]string{"-organization=my-org", "-source-workspace=source"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace") {
		t.Fatalf("expected workspace error, got %q", out)
	}
}

func TestRunTriggerCreateCommandRequiresSourceWorkspace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTriggerCreateCommand(ui, nil, nil)

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=target"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-source-workspace") {
		t.Fatalf("expected source-workspace error, got %q", out)
	}
}

func TestRunTriggerCreateCommandHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	wsSvc := &mockWorkspaceMultiReadService{
		workspaces: map[string]*tfe.Workspace{
			"target": {ID: "ws-target"},
			"source": {ID: "ws-source"},
		},
	}
	rtSvc := &mockRunTriggerCreateService{err: errors.New("api error")}
	cmd := newRunTriggerCreateCommand(ui, wsSvc, rtSvc)

	code := cmd.Run([]string{
		"-organization=my-org",
		"-workspace=target",
		"-source-workspace=source",
	})

	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if rtSvc.lastWSID != "ws-target" {
		t.Fatalf("expected target workspace ID ws-target, got %s", rtSvc.lastWSID)
	}

	if rtSvc.lastOptions.Sourceable == nil || rtSvc.lastOptions.Sourceable.ID != "ws-source" {
		t.Fatalf("expected source workspace ID ws-source, got %#v", rtSvc.lastOptions)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "api error") {
		t.Fatalf("expected error output, got %q", out)
	}
}

func TestRunTriggerCreateCommandSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	wsSvc := &mockWorkspaceMultiReadService{
		workspaces: map[string]*tfe.Workspace{
			"target": {ID: "ws-target"},
			"source": {ID: "ws-source"},
		},
	}
	rtSvc := &mockRunTriggerCreateService{
		response: &tfe.RunTrigger{
			ID:             "rt-new123",
			WorkspaceName:  "target",
			SourceableName: "source",
			CreatedAt:      time.Now(),
		},
	}
	cmd := newRunTriggerCreateCommand(ui, wsSvc, rtSvc)

	code := cmd.Run([]string{
		"-organization=my-org",
		"-workspace=target",
		"-source-workspace=source",
	})

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if out := ui.OutputWriter.String(); !strings.Contains(out, "successfully") {
		t.Fatalf("expected success message, got %q", out)
	}
}

func TestRunTriggerCreateCommandHandlesWorkspaceNotFound(t *testing.T) {
	ui := cli.NewMockUi()
	wsSvc := &mockWorkspaceMultiReadService{
		workspaces: map[string]*tfe.Workspace{},
	}
	cmd := newRunTriggerCreateCommand(ui, wsSvc, nil)

	code := cmd.Run([]string{
		"-organization=my-org",
		"-workspace=target",
		"-source-workspace=source",
	})

	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "not found") {
		t.Fatalf("expected workspace not found error, got %q", out)
	}
}
