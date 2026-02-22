package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newVariableSetListCommand(ui cli.Ui, svc variableSetLister) *VariableSetListCommand {
	return &VariableSetListCommand{
		Meta:      newTestMeta(ui),
		varSetSvc: svc,
	}
}

func TestVariableSetListRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newVariableSetListCommand(ui, &mockVariableSetListService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected organization error")
	}
}

func TestVariableSetListHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockVariableSetListService{err: errors.New("boom")}
	cmd := newVariableSetListCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=org"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastOrg != "org" {
		t.Fatalf("unexpected org recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestVariableSetListOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockVariableSetListService{
		response: &tfe.VariableSetList{
			Items: []*tfe.VariableSet{
				{
					ID:          "varset-1",
					Name:        "varset1",
					Description: "First variable set",
					Global:      true,
					Variables:   []*tfe.VariableSetVariable{{ID: "var-1"}},
				},
				{
					ID:          "varset-2",
					Name:        "varset2",
					Description: "Second variable set",
					Global:      false,
					Variables:   []*tfe.VariableSetVariable{},
				},
			},
		},
	}
	cmd := newVariableSetListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=org", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var rows []map[string]string
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0]["ID"] != "varset-1" {
		t.Fatalf("unexpected row: %#v", rows[0])
	}
}

func TestVariableSetListPassesQueryAndInclude(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockVariableSetListService{
		response: &tfe.VariableSetList{Items: []*tfe.VariableSet{}},
	}
	cmd := newVariableSetListCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=org", "-query=prod", "-include=workspaces,projects"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastOpts == nil {
		t.Fatalf("expected options to be passed")
	}
	if svc.lastOpts.Query != "prod" {
		t.Fatalf("expected query option, got %#v", svc.lastOpts)
	}
	if svc.lastOpts.Include != "workspaces,projects" {
		t.Fatalf("expected include option, got %#v", svc.lastOpts)
	}
}
