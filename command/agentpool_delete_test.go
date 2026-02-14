package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newAgentPoolDeleteCommand(ui *cli.MockUi, svc agentPoolDeleter) *AgentPoolDeleteCommand {
	cmd := &AgentPoolDeleteCommand{
		Meta:         newTestMeta(ui),
		agentPoolSvc: svc,
	}
	cmd.Meta.Ui = ui
	return cmd
}

func TestAgentPoolDeleteRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newAgentPoolDeleteCommand(ui, &mockAgentPoolDeleteService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing id, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error, got: %s", ui.ErrorWriter.String())
	}
}

func TestAgentPoolDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockAgentPoolDeleteService{err: errors.New("boom")}
	cmd := newAgentPoolDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=apool-123", "-force"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastID != "apool-123" {
		t.Fatalf("unexpected delete id: %s", svc.lastID)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got: %s", ui.ErrorWriter.String())
	}
}

func TestAgentPoolDeleteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockAgentPoolDeleteService{}
	cmd := newAgentPoolDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=apool-123", "-force"}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if svc.lastID != "apool-123" {
		t.Fatalf("unexpected delete id: %s", svc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message, got: %s", ui.OutputWriter.String())
	}
}
