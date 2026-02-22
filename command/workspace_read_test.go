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

type mockWorkspaceReader struct {
	workspace       *tfe.Workspace
	err             error
	lastOrg         string
	lastName        string
	lastReadOptions *tfe.WorkspaceReadOptions
}

func (m *mockWorkspaceReader) Read(_ context.Context, organization, workspace string) (*tfe.Workspace, error) {
	m.lastOrg = organization
	m.lastName = workspace
	return m.workspace, m.err
}

func (m *mockWorkspaceReader) ReadWithOptions(_ context.Context, organization, workspace string, options *tfe.WorkspaceReadOptions) (*tfe.Workspace, error) {
	m.lastOrg = organization
	m.lastName = workspace
	m.lastReadOptions = options
	return m.workspace, m.err
}

func newWorkspaceReadCommand(ui cli.Ui, reader workspaceReader) *WorkspaceReadCommand {
	return &WorkspaceReadCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: reader,
	}
}

func TestWorkspaceReadRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceReadCommand(ui, &mockWorkspaceReader{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org, got %d", code)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error, got %q", ui.ErrorWriter.String())
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing name, got %d", code)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected name error, got %q", ui.ErrorWriter.String())
	}
}

func TestWorkspaceReadHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	reader := &mockWorkspaceReader{err: errors.New("boom")}
	cmd := newWorkspaceReadCommand(ui, reader)

	code := cmd.Run([]string{"-organization=my-org", "-name=prod"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if reader.lastOrg != "my-org" || reader.lastName != "prod" {
		t.Fatalf("unexpected parameters: %#v", reader)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got %q", ui.ErrorWriter.String())
	}
}

func TestWorkspaceReadOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	reader := &mockWorkspaceReader{workspace: &tfe.Workspace{
		ID:               "ws-123",
		Name:             "prod",
		TerraformVersion: "1.7.0",
		AutoApply:        true,
		Project: &tfe.Project{
			ID:   "prj-123",
			Name: "platform-prod",
		},
		CreatedAt: time.Unix(0, 0),
		UpdatedAt: time.Unix(0, 0),
	}}
	cmd := newWorkspaceReadCommand(ui, reader)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=prod", "-output=json"})
	})

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if data["Name"] != "prod" || data["TerraformVersion"] != "1.7.0" {
		t.Fatalf("unexpected data: %#v", data)
	}
	if data["ProjectID"] != "prj-123" || data["ProjectName"] != "platform-prod" {
		t.Fatalf("expected project metadata in output: %#v", data)
	}
	if reader.lastReadOptions == nil || len(reader.lastReadOptions.Include) == 0 {
		t.Fatalf("expected read options include to be set")
	}
}
