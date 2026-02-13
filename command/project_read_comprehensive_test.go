package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newProjectReadCommand(ui cli.Ui, svc projectReader) *ProjectReadCommand {
	return &ProjectReadCommand{
		Meta:       newTestMeta(ui),
		projectSvc: svc,
	}
}

func TestProjectReadRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newProjectReadCommand(ui, &mockProjectReadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestProjectReadHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockProjectReadService{err: errors.New("boom")}
	cmd := newProjectReadCommand(ui, svc)

	if code := cmd.Run([]string{"-id=prj-123"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID != "prj-123" {
		t.Fatalf("unexpected project id recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestProjectReadOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockProjectReadService{
		response: &tfe.Project{
			ID:          "prj-123",
			Name:        "test-project",
			Description: "A test project",
		},
	}
	cmd := newProjectReadCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=prj-123", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "prj-123" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
