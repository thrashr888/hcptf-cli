package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newWorkspaceCreateCommand(ui cli.Ui, svc workspaceCreator) *WorkspaceCreateCommand {
	return &WorkspaceCreateCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: svc,
	}
}

func TestWorkspaceCreateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceCreateCommand(ui, &mockWorkspaceCreateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing name, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected name error")
	}
}

func TestWorkspaceCreateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceCreateService{err: errors.New("boom")}
	cmd := newWorkspaceCreateCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=prod"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if svc.lastOrg != "my-org" {
		t.Fatalf("expected organization my-org, got %s", svc.lastOrg)
	}
	if svc.lastOptions.Name == nil || *svc.lastOptions.Name != "prod" {
		t.Fatalf("expected workspace name prod")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestWorkspaceCreateOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceCreateService{response: &tfe.Workspace{ID: "ws-1", Name: "prod", TerraformVersion: "1.6.0", AutoApply: true}}
	cmd := newWorkspaceCreateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=prod", "-auto-apply", "-terraform-version=1.6.0", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if svc.lastOptions.Name == nil || *svc.lastOptions.Name != "prod" {
		t.Fatalf("expected name in options")
	}
	if svc.lastOptions.AutoApply == nil || !*svc.lastOptions.AutoApply {
		t.Fatalf("expected auto apply true")
	}
	if svc.lastOptions.TerraformVersion == nil || *svc.lastOptions.TerraformVersion != "1.6.0" {
		t.Fatalf("expected terraform version set")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["Name"] != "prod" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
