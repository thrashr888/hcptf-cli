package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newConfigVersionListCommand(ui cli.Ui, ws workspaceReader, cvs configVersionLister) *ConfigVersionListCommand {
	return &ConfigVersionListCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: ws,
		configVerSvc: cvs,
	}
}

func TestConfigVersionListRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newConfigVersionListCommand(ui, &mockWorkspaceReader{}, &mockConfigVersionListService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 org")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 workspace")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-workspace") {
		t.Fatalf("expected workspace error")
	}
}

func TestConfigVersionListHandlesWorkspaceError(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{err: errors.New("boom")}
	cmd := newConfigVersionListCommand(ui, ws, &mockConfigVersionListService{})

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected workspace error output")
	}
}

func TestConfigVersionListHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	cvs := &mockConfigVersionListService{err: errors.New("fail")}
	cmd := newConfigVersionListCommand(ui, ws, cvs)

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if cvs.lastWorkspace != "ws-1" {
		t.Fatalf("expected workspace id recorded")
	}
	if cvs.lastOptions == nil || cvs.lastOptions.ListOptions.PageSize != 50 {
		t.Fatalf("expected list options set")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "fail") {
		t.Fatalf("expected error output")
	}
}

func TestConfigVersionListOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	cvs := &mockConfigVersionListService{response: &tfe.ConfigurationVersionList{Items: []*tfe.ConfigurationVersion{{
		ID:          "cv-1",
		Status:      tfe.ConfigurationUploaded,
		Source:      tfe.ConfigurationSourceAPI,
		Speculative: true,
		Provisional: false,
	}}}}
	cmd := newConfigVersionListCommand(ui, ws, cvs)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if len(rows) != 1 || rows[0]["ID"] != "cv-1" {
		t.Fatalf("unexpected rows: %#v", rows)
	}
}
