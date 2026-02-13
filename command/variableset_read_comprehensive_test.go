package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newVariableSetReadCommand(ui cli.Ui, svc variableSetReader) *VariableSetReadCommand {
	return &VariableSetReadCommand{
		Meta:      newTestMeta(ui),
		varSetSvc: svc,
	}
}

func TestVariableSetReadRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newVariableSetReadCommand(ui, &mockVariableSetReadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestVariableSetReadHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockVariableSetReadService{err: errors.New("boom")}
	cmd := newVariableSetReadCommand(ui, svc)

	if code := cmd.Run([]string{"-id=varset-123"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID != "varset-123" {
		t.Fatalf("unexpected variable set id recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestVariableSetReadOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockVariableSetReadService{
		response: &tfe.VariableSet{
			ID:          "varset-123",
			Name:        "test-varset",
			Description: "A test variable set",
			Global:      true,
			Variables:   []*tfe.VariableSetVariable{{ID: "var-1"}, {ID: "var-2"}},
		},
	}
	cmd := newVariableSetReadCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=varset-123", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "varset-123" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
