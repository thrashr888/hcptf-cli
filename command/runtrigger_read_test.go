package command

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

type mockRunTriggerReadService struct {
	response *tfe.RunTrigger
	err      error
	lastID   string
}

func (m *mockRunTriggerReadService) Read(_ context.Context, runTriggerID string) (*tfe.RunTrigger, error) {
	m.lastID = runTriggerID
	return m.response, m.err
}

func newRunTriggerReadCommand(ui cli.Ui, rtSvc runTriggerReader) *RunTriggerReadCommand {
	return &RunTriggerReadCommand{
		Meta:          newTestMeta(ui),
		runTriggerSvc: rtSvc,
	}
}

func TestRunTriggerReadCommandRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunTriggerReadCommand(ui, nil)

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestRunTriggerReadCommandHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	rtSvc := &mockRunTriggerReadService{err: errors.New("not found")}
	cmd := newRunTriggerReadCommand(ui, rtSvc)

	code := cmd.Run([]string{"-id=rt-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if rtSvc.lastID != "rt-123" {
		t.Fatalf("expected ID rt-123, got %s", rtSvc.lastID)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "not found") {
		t.Fatalf("expected error output, got %q", out)
	}
}

func TestRunTriggerReadCommandSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	rtSvc := &mockRunTriggerReadService{
		response: &tfe.RunTrigger{
			ID:             "rt-abc123",
			WorkspaceName:  "target-workspace",
			SourceableName: "source-workspace",
			CreatedAt:      time.Now(),
			Workspace: &tfe.Workspace{
				ID: "ws-target",
			},
			Sourceable: &tfe.Workspace{
				ID: "ws-source",
			},
		},
	}
	cmd := newRunTriggerReadCommand(ui, rtSvc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=rt-abc123"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if rtSvc.lastID != "rt-abc123" {
		t.Fatalf("expected ID rt-abc123, got %s", rtSvc.lastID)
	}

	if !strings.Contains(output, "rt-abc123") {
		t.Fatalf("expected output to contain ID, got %q", output)
	}
	if !strings.Contains(output, "target-workspace") {
		t.Fatalf("expected output to contain workspace name, got %q", output)
	}
	if !strings.Contains(output, "source-workspace") {
		t.Fatalf("expected output to contain sourceable name, got %q", output)
	}
}

func TestRunTriggerReadCommandWithoutRelationships(t *testing.T) {
	ui := cli.NewMockUi()
	rtSvc := &mockRunTriggerReadService{
		response: &tfe.RunTrigger{
			ID:             "rt-abc123",
			WorkspaceName:  "target-workspace",
			SourceableName: "source-workspace",
			CreatedAt:      time.Now(),
		},
	}
	cmd := newRunTriggerReadCommand(ui, rtSvc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=rt-abc123"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if !strings.Contains(output, "rt-abc123") {
		t.Fatalf("expected output to contain ID, got %q", output)
	}
}
