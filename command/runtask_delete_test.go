package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newRunTaskDeleteCommand(ui cli.Ui, svc runTaskDeleterReader) *RunTaskDeleteCommand {
	return &RunTaskDeleteCommand{
		Meta:       newTestMeta(ui),
		runTaskSvc: svc,
	}
}

func TestRunTaskDeleteRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunTaskDeleteReaderService{
		readResponse: &tfe.RunTask{Name: "test-task"},
	}
	cmd := newRunTaskDeleteCommand(ui, svc)

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing id, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error, got %q", ui.ErrorWriter.String())
	}
}

func TestRunTaskDeleteReadError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunTaskDeleteReaderService{
		readErr: errors.New("not found"),
	}
	cmd := newRunTaskDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=task-123", "-force"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastReadID != "task-123" {
		t.Fatalf("expected lastReadID task-123, got %q", svc.lastReadID)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "not found") {
		t.Fatalf("expected read error output, got %q", ui.ErrorWriter.String())
	}
}

func TestRunTaskDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunTaskDeleteReaderService{
		readResponse: &tfe.RunTask{Name: "my-task"},
		deleteErr:    errors.New("boom"),
	}
	cmd := newRunTaskDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=task-123", "-force"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastDeleteID != "task-123" {
		t.Fatalf("expected lastDeleteID task-123, got %q", svc.lastDeleteID)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got %q", ui.ErrorWriter.String())
	}
}

func TestRunTaskDeleteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunTaskDeleteReaderService{
		readResponse: &tfe.RunTask{Name: "my-task"},
	}
	cmd := newRunTaskDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=task-123", "-force"}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if svc.lastReadID != "task-123" {
		t.Fatalf("expected lastReadID task-123, got %q", svc.lastReadID)
	}
	if svc.lastDeleteID != "task-123" {
		t.Fatalf("expected lastDeleteID task-123, got %q", svc.lastDeleteID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message, got %q", ui.OutputWriter.String())
	}
}
