package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newWorkspaceDeleteCommand(ui *cli.MockUi, svc workspaceDeleter) *WorkspaceDeleteCommand {
	cmd := &WorkspaceDeleteCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: svc,
	}
	cmd.Meta.Ui = ui
	return cmd
}

func TestWorkspaceDeleteRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceDeleteCommand(ui, &mockWorkspaceDeleteService{})

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

func TestWorkspaceDeletePromptsWithoutForce(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("no\n")
	svc := &mockWorkspaceDeleteService{}
	cmd := newWorkspaceDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod"}); code != 0 {
		t.Fatalf("expected exit 0 when cancelled, got %d", code)
	}
	if svc.lastName != "" {
		t.Fatalf("expected delete not called")
	}
	if !strings.Contains(ui.OutputWriter.String(), "Deletion cancelled") {
		t.Fatalf("expected cancellation message")
	}
}

func TestWorkspaceDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("yes\n")
	svc := &mockWorkspaceDeleteService{err: errors.New("boom")}
	cmd := newWorkspaceDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastOrg != "my-org" || svc.lastName != "prod" {
		t.Fatalf("unexpected delete args: %#v", svc)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestWorkspaceDeleteForceBypassesPrompt(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceDeleteService{}
	cmd := newWorkspaceDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod", "-force"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastOrg != "my-org" || svc.lastName != "prod" {
		t.Fatalf("unexpected delete args")
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message")
	}
}
