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

type mockRunTriggerListService struct {
	response    *tfe.RunTriggerList
	err         error
	lastWSID    string
	lastOptions *tfe.RunTriggerListOptions
}

func (m *mockRunTriggerListService) List(_ context.Context, workspaceID string, options *tfe.RunTriggerListOptions) (*tfe.RunTriggerList, error) {
	m.lastWSID = workspaceID
	m.lastOptions = options
	return m.response, m.err
}

type mockWorkspaceReadService struct {
	response *tfe.Workspace
	err      error
}

func (m *mockWorkspaceReadService) Read(_ context.Context, organization, workspace string) (*tfe.Workspace, error) {
	return m.response, m.err
}

func newRunTriggerListCommand(ui cli.Ui, wsSvc workspaceReader, rtSvc runTriggerLister) *RunTriggerListCommand {
	return &RunTriggerListCommand{
		Meta:          newTestMeta(ui),
		workspaceSvc:  wsSvc,
		runTriggerSvc: rtSvc,
	}
}

func TestRunTriggerListCommandRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTriggerListCommand(ui, nil, nil)

	if code := cmd.Run([]string{"-workspace=test"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestRunTriggerListCommandRequiresWorkspace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTriggerListCommand(ui, nil, nil)

	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace") {
		t.Fatalf("expected workspace error, got %q", out)
	}
}

func TestRunTriggerListCommandValidatesTriggerType(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTriggerListCommand(ui, nil, nil)

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=test", "-type=invalid"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "inbound") || !strings.Contains(out, "outbound") {
		t.Fatalf("expected type validation error, got %q", out)
	}
}

func TestRunTriggerListCommandHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	wsSvc := &mockWorkspaceReadService{
		response: &tfe.Workspace{ID: "ws-123"},
	}
	rtSvc := &mockRunTriggerListService{err: errors.New("api error")}
	cmd := newRunTriggerListCommand(ui, wsSvc, rtSvc)

	code := cmd.Run([]string{"-organization=my-org", "-workspace=test"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if rtSvc.lastWSID != "ws-123" {
		t.Fatalf("expected workspace ID ws-123, got %s", rtSvc.lastWSID)
	}

	if rtSvc.lastOptions == nil || rtSvc.lastOptions.RunTriggerType != "inbound" {
		t.Fatalf("expected inbound trigger type, got %#v", rtSvc.lastOptions)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "api error") {
		t.Fatalf("expected error output, got %q", out)
	}
}

func TestRunTriggerListCommandSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	wsSvc := &mockWorkspaceReadService{
		response: &tfe.Workspace{ID: "ws-123"},
	}
	rtSvc := &mockRunTriggerListService{
		response: &tfe.RunTriggerList{
			Items: []*tfe.RunTrigger{
				{
					ID:              "rt-abc123",
					WorkspaceName:   "target-workspace",
					SourceableName:  "source-workspace",
					CreatedAt:       time.Now(),
				},
			},
		},
	}
	cmd := newRunTriggerListCommand(ui, wsSvc, rtSvc)

	code := cmd.Run([]string{"-organization=my-org", "-workspace=test"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if rtSvc.lastOptions.RunTriggerType != "inbound" {
		t.Fatalf("expected inbound type by default, got %s", rtSvc.lastOptions.RunTriggerType)
	}
}

func TestRunTriggerListCommandOutboundType(t *testing.T) {
	ui := cli.NewMockUi()
	wsSvc := &mockWorkspaceReadService{
		response: &tfe.Workspace{ID: "ws-123"},
	}
	rtSvc := &mockRunTriggerListService{
		response: &tfe.RunTriggerList{Items: []*tfe.RunTrigger{}},
	}
	cmd := newRunTriggerListCommand(ui, wsSvc, rtSvc)

	code := cmd.Run([]string{"-organization=my-org", "-workspace=test", "-type=outbound"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if rtSvc.lastOptions.RunTriggerType != "outbound" {
		t.Fatalf("expected outbound type, got %s", rtSvc.lastOptions.RunTriggerType)
	}
}
