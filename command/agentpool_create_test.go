package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newAgentPoolCreateCommand(ui cli.Ui, svc agentPoolCreator) *AgentPoolCreateCommand {
	return &AgentPoolCreateCommand{
		Meta:         newTestMeta(ui),
		agentPoolSvc: svc,
	}
}

func TestAgentPoolCreateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newAgentPoolCreateCommand(ui, &mockAgentPoolCreateService{})

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

func TestAgentPoolCreateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockAgentPoolCreateService{err: errors.New("quota exceeded")}
	cmd := newAgentPoolCreateCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=my-pool"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastOrg != "my-org" {
		t.Fatalf("expected org my-org, got %s", svc.lastOrg)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "quota exceeded") {
		t.Fatalf("expected error output")
	}
}

func TestAgentPoolCreateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockAgentPoolCreateService{response: &tfe.AgentPool{
		ID:   "apool-1",
		Name: "my-pool",
	}}
	cmd := newAgentPoolCreateCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=my-pool"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}
	if !strings.Contains(ui.OutputWriter.String(), "my-pool") {
		t.Fatalf("expected success output with pool name")
	}
}

func TestAgentPoolCreateHelp(t *testing.T) {
	cmd := &AgentPoolCreateCommand{}
	if !strings.Contains(cmd.Help(), "agentpool create") {
		t.Fatal("expected help text")
	}
}

func TestAgentPoolCreateSynopsis(t *testing.T) {
	cmd := &AgentPoolCreateCommand{}
	if cmd.Synopsis() == "" {
		t.Fatal("expected non-empty synopsis")
	}
}
