package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newVariableCreateCommand(ui cli.Ui, ws workspaceReader, vars variableCreator) *VariableCreateCommand {
	return &VariableCreateCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: ws,
		variableSvc:  vars,
	}
}

func TestVariableCreateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newVariableCreateCommand(ui, &mockWorkspaceReader{}, &mockVariableCreateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing workspace")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-workspace") {
		t.Fatalf("expected workspace error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod"}); code != 1 {
		t.Fatalf("expected exit 1 missing key")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-key") {
		t.Fatalf("expected key error")
	}
}

func TestVariableCreateValidatesCategory(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newVariableCreateCommand(ui, &mockWorkspaceReader{}, &mockVariableCreateService{})

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-key=k", "-value=v", "-category=invalid"}); code != 1 {
		t.Fatalf("expected exit 1 invalid category")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "category") {
		t.Fatalf("expected category error")
	}
}

func TestVariableCreateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableCreateService{err: errors.New("boom")}
	cmd := newVariableCreateCommand(ui, ws, vars)

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-key=k", "-value=v"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if vars.lastWorkspace != "ws-1" {
		t.Fatalf("expected workspace ID passed")
	}
	if vars.lastOptions.Key == nil || *vars.lastOptions.Key != "k" {
		t.Fatalf("expected key option")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestVariableCreateOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableCreateService{response: &tfe.Variable{ID: "var-1", Key: "k", Category: tfe.CategoryTerraform}}
	cmd := newVariableCreateCommand(ui, ws, vars)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-key=k", "-value=v", "-category=env", "-description=desc", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	if vars.lastOptions.Category == nil || *vars.lastOptions.Category != tfe.CategoryEnv {
		t.Fatalf("expected env category")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "var-1" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
