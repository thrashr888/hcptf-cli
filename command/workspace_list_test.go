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

type mockWorkspaceService struct {
	response    *tfe.WorkspaceList
	err         error
	lastOrg     string
	lastOptions *tfe.WorkspaceListOptions
}

func (m *mockWorkspaceService) List(_ context.Context, organization string, options *tfe.WorkspaceListOptions) (*tfe.WorkspaceList, error) {
	m.lastOrg = organization
	m.lastOptions = options
	return m.response, m.err
}

func newWorkspaceListCommand(ui cli.Ui, svc workspaceLister) *WorkspaceListCommand {
	return &WorkspaceListCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: svc,
	}
}

func TestWorkspaceListCommandRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceListCommand(ui, &mockWorkspaceService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestWorkspaceListCommandHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceService{err: errors.New("boom")}
	cmd := newWorkspaceListCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if svc.lastOrg != "my-org" {
		t.Fatalf("expected organization my-org, got %s", svc.lastOrg)
	}

	if svc.lastOptions == nil || svc.lastOptions.ListOptions.PageSize != 100 {
		t.Fatalf("expected list options with page size 100, got %#v", svc.lastOptions)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "boom") {
		t.Fatalf("expected error output, got %q", out)
	}
}

func TestWorkspaceListCommandOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceService{
		response: &tfe.WorkspaceList{Items: []*tfe.Workspace{{
			ID:               "ws-123",
			Name:             "prod",
			TerraformVersion: "1.7.0",
			AutoApply:        true,
			Locked:           false,
		}}},
	}
	cmd := newWorkspaceListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-output=json"})
	})

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	var rows []map[string]string
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if len(rows) != 1 || rows[0]["Name"] != "prod" {
		t.Fatalf("unexpected rows: %#v", rows)
	}
}

func TestWorkspaceListCommandPassesFilterFlags(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceService{
		response: &tfe.WorkspaceList{Items: []*tfe.Workspace{}},
	}
	cmd := newWorkspaceListCommand(ui, svc)

	code := cmd.Run([]string{
		"-organization=my-org",
		"-search=prod",
		"-tags=env:prod",
		"-exclude-tags=archived",
		"-wildcard-name=prod-*",
		"-project-id=prj-123",
		"-current-run-status=planned",
		"-include=project,current_run",
		"-sort=-name",
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if svc.lastOptions == nil {
		t.Fatalf("expected list options")
	}
	if svc.lastOptions.Search != "prod" || svc.lastOptions.ProjectID != "prj-123" {
		t.Fatalf("expected filters in options, got %#v", svc.lastOptions)
	}
	if len(svc.lastOptions.Include) != 2 {
		t.Fatalf("expected include options, got %#v", svc.lastOptions.Include)
	}
	if svc.lastOptions.Sort != "-name" {
		t.Fatalf("expected sort value, got %q", svc.lastOptions.Sort)
	}
}
