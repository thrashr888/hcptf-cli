package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newRunTaskCreateCommand(ui cli.Ui, svc runTaskCreator) *RunTaskCreateCommand {
	return &RunTaskCreateCommand{
		Meta:       newTestMeta(ui),
		runTaskSvc: svc,
	}
}

func TestRunTaskCreateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTaskCreateCommand(ui, &mockRunTaskCreateService{})

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

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org", "-name=scan"}); code != 1 {
		t.Fatalf("expected exit 1 missing url, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-url") {
		t.Fatalf("expected url error")
	}
}

func TestRunTaskCreateInvalidCategory(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTaskCreateCommand(ui, &mockRunTaskCreateService{})

	code := cmd.Run([]string{"-organization=my-org", "-name=scan", "-url=https://example.com", "-category=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "category") {
		t.Fatalf("expected category error")
	}
}

func TestRunTaskCreateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunTaskCreateService{err: errors.New("boom")}
	cmd := newRunTaskCreateCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=scan", "-url=https://example.com"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastOrg != "my-org" {
		t.Fatalf("expected organization my-org, got %s", svc.lastOrg)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRunTaskCreateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunTaskCreateService{response: &tfe.RunTask{
		ID:       "task-1",
		Name:     "scan",
		URL:      "https://example.com",
		Category: "task",
		Enabled:  true,
	}}
	cmd := newRunTaskCreateCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=scan", "-url=https://example.com"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
	if svc.lastOrg != "my-org" {
		t.Fatalf("expected organization my-org, got %s", svc.lastOrg)
	}
	if !strings.Contains(ui.OutputWriter.String(), "scan") {
		t.Fatalf("expected success output with task name")
	}
}

func TestRunTaskCreateHelp(t *testing.T) {
	cmd := &RunTaskCreateCommand{}
	help := cmd.Help()
	if !strings.Contains(help, "runtask create") {
		t.Fatalf("expected help text, got: %s", help)
	}
}

func TestRunTaskCreateSynopsis(t *testing.T) {
	cmd := &RunTaskCreateCommand{}
	syn := cmd.Synopsis()
	if syn == "" {
		t.Fatal("expected non-empty synopsis")
	}
}
