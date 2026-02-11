package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newVariableDeleteCommand(ui *cli.MockUi, ws workspaceReader, vars variableDeleter) *VariableDeleteCommand {
	cmd := &VariableDeleteCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: ws,
		variableSvc:  vars,
	}
	cmd.Meta.Ui = ui
	return cmd
}

func TestVariableDeleteRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newVariableDeleteCommand(ui, &mockWorkspaceReader{}, &mockVariableDeleteService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected workspace error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod"}); code != 1 {
		t.Fatalf("expected id error")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error message")
	}
}

func TestVariableDeleteCancellation(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("no\n")
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableDeleteService{}
	cmd := newVariableDeleteCommand(ui, ws, vars)

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-id=var-1"}); code != 0 {
		t.Fatalf("expected exit 0 on cancel")
	}
	if vars.lastID != "" {
		t.Fatalf("expected delete not called")
	}
	if !strings.Contains(ui.OutputWriter.String(), "Deletion cancelled") {
		t.Fatalf("expected cancellation message")
	}
}

func TestVariableDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("yes\n")
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableDeleteService{err: errors.New("boom")}
	cmd := newVariableDeleteCommand(ui, ws, vars)

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-id=var-1"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if vars.lastWorkspace != "ws-1" || vars.lastID != "var-1" {
		t.Fatalf("unexpected delete args")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestVariableDeleteForce(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableDeleteService{}
	cmd := newVariableDeleteCommand(ui, ws, vars)

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-id=var-1", "-force"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if vars.lastWorkspace != "ws-1" || vars.lastID != "var-1" {
		t.Fatalf("unexpected delete args")
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message")
	}
}
