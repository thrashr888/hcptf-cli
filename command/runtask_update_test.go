package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newRunTaskUpdateCommand(ui cli.Ui, svc runTaskUpdater) *RunTaskUpdateCommand {
	return &RunTaskUpdateCommand{
		Meta:       newTestMeta(ui),
		runTaskSvc: svc,
	}
}

func TestRunTaskUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTaskUpdateCommand(ui, &mockRunTaskUpdateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestRunTaskUpdateInvalidCategory(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTaskUpdateCommand(ui, &mockRunTaskUpdateService{})

	code := cmd.Run([]string{"-id=task-1", "-category=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "category") {
		t.Fatalf("expected category error")
	}
}

func TestRunTaskUpdateInvalidEnabled(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTaskUpdateCommand(ui, &mockRunTaskUpdateService{})

	code := cmd.Run([]string{"-id=task-1", "-enabled=maybe"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "enabled") {
		t.Fatalf("expected enabled error")
	}
}

func TestRunTaskUpdateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunTaskUpdateService{err: errors.New("boom")}
	cmd := newRunTaskUpdateCommand(ui, svc)

	code := cmd.Run([]string{"-id=task-1", "-name=updated"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastID != "task-1" {
		t.Fatalf("expected id task-1, got %s", svc.lastID)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRunTaskUpdateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunTaskUpdateService{response: &tfe.RunTask{
		ID:       "task-1",
		Name:     "updated",
		URL:      "https://example.com",
		Category: "task",
		Enabled:  true,
	}}
	cmd := newRunTaskUpdateCommand(ui, svc)

	code := cmd.Run([]string{"-id=task-1", "-name=updated", "-enabled=true"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
	if svc.lastID != "task-1" {
		t.Fatalf("expected id task-1, got %s", svc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "updated") {
		t.Fatalf("expected success output with task name")
	}
}

func TestRunTaskUpdateHelp(t *testing.T) {
	cmd := &RunTaskUpdateCommand{}
	help := cmd.Help()
	if !strings.Contains(help, "runtask update") {
		t.Fatalf("expected help text, got: %s", help)
	}
}

func TestRunTaskUpdateSynopsis(t *testing.T) {
	cmd := &RunTaskUpdateCommand{}
	syn := cmd.Synopsis()
	if syn == "" {
		t.Fatal("expected non-empty synopsis")
	}
}
