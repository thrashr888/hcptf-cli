package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newVariableSetCreateCommand(ui cli.Ui, svc variableSetCreator) *VariableSetCreateCommand {
	return &VariableSetCreateCommand{
		Meta:           newTestMeta(ui),
		variableSetSvc: svc,
	}
}

func TestVariableSetCreateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newVariableSetCreateCommand(ui, &mockVariableSetCreateService{})

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

func TestVariableSetCreateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockVariableSetCreateService{err: errors.New("boom")}
	cmd := newVariableSetCreateCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=my-varset"})
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

func TestVariableSetCreateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockVariableSetCreateService{response: &tfe.VariableSet{ID: "varset-1", Name: "my-varset", Global: false}}
	cmd := newVariableSetCreateCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=my-varset", "-description=Test set", "-global"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	if svc.lastOrg != "my-org" {
		t.Fatalf("expected organization my-org, got %s", svc.lastOrg)
	}
	if !strings.Contains(ui.OutputWriter.String(), "my-varset") {
		t.Fatalf("expected success output with varset name")
	}
}

func TestVariableSetCreateHelp(t *testing.T) {
	cmd := &VariableSetCreateCommand{}
	help := cmd.Help()
	if !strings.Contains(help, "variableset create") {
		t.Fatalf("expected help text, got: %s", help)
	}
}

func TestVariableSetCreateSynopsis(t *testing.T) {
	cmd := &VariableSetCreateCommand{}
	syn := cmd.Synopsis()
	if syn == "" {
		t.Fatal("expected non-empty synopsis")
	}
}
