package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newVariableUpdateCommand(ui cli.Ui, ws workspaceReader, vars variableUpdater) *VariableUpdateCommand {
	return &VariableUpdateCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: ws,
		variableSvc:  vars,
	}
}

func TestVariableUpdateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newVariableUpdateCommand(ui, &mockWorkspaceReader{}, &mockVariableUpdateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 workspace")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod"}); code != 1 {
		t.Fatalf("expected exit 1 id")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestVariableUpdateValidatesBools(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newVariableUpdateCommand(ui, &mockWorkspaceReader{}, &mockVariableUpdateService{})

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-id=var", "-sensitive=maybe"}); code != 1 {
		t.Fatalf("expected exit 1 sensitive")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "sensitive") {
		t.Fatalf("expected sensitive error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-id=var", "-hcl=maybe"}); code != 1 {
		t.Fatalf("expected exit 1 hcl")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "hcl") {
		t.Fatalf("expected hcl error")
	}
}

func TestVariableUpdateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableUpdateService{err: errors.New("boom")}
	cmd := newVariableUpdateCommand(ui, ws, vars)

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-id=var", "-value=hi", "-sensitive=true"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if vars.lastWorkspace != "ws-1" || vars.lastID != "var" {
		t.Fatalf("unexpected request")
	}
	if vars.lastOptions.Sensitive == nil || !*vars.lastOptions.Sensitive {
		t.Fatalf("expected sensitive true")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestVariableUpdateOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableUpdateService{response: &tfe.Variable{ID: "var-1", Key: "k", Category: tfe.CategoryTerraform}}
	cmd := newVariableUpdateCommand(ui, ws, vars)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-id=var-1", "-key=new", "-description=desc", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	if vars.lastOptions.Key == nil || *vars.lastOptions.Key != "new" {
		t.Fatalf("expected key option")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "var-1" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
