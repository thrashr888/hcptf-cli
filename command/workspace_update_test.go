package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newWorkspaceUpdateCommand(ui cli.Ui, svc workspaceUpdater) *WorkspaceUpdateCommand {
	return &WorkspaceUpdateCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: svc,
	}
}

func TestWorkspaceUpdateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceUpdateCommand(ui, &mockWorkspaceUpdateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing name")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected name error")
	}
}

func TestWorkspaceUpdateValidatesAutoApply(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceUpdateCommand(ui, &mockWorkspaceUpdateService{})

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod", "-auto-apply=maybe"}); code != 1 {
		t.Fatalf("expected exit 1 for invalid auto-apply")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "auto-apply") {
		t.Fatalf("expected validation error")
	}
}

func TestWorkspaceUpdateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{err: errors.New("boom")}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod", "-auto-apply=true"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastOrg != "my-org" || svc.lastName != "prod" {
		t.Fatalf("unexpected parameters: %#v", svc)
	}
	if svc.lastOptions.AutoApply == nil || !*svc.lastOptions.AutoApply {
		t.Fatalf("expected auto apply true")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestWorkspaceUpdateOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{response: &tfe.Workspace{ID: "ws-1", Name: "new-name", TerraformVersion: "1.6.1", AutoApply: true}}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=prod", "-new-name=new", "-terraform-version=1.6.1", "-description=hello", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if svc.lastOptions.Name == nil || *svc.lastOptions.Name != "new" {
		t.Fatalf("expected new name option")
	}
	if svc.lastOptions.Description == nil || *svc.lastOptions.Description != "hello" {
		t.Fatalf("expected description option")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["Name"] != "new-name" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
