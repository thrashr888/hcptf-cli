package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newProjectListCommand(ui cli.Ui, svc projectLister) *ProjectListCommand {
	return &ProjectListCommand{
		Meta:       newTestMeta(ui),
		projectSvc: svc,
	}
}

func TestProjectListComprehensiveRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newProjectListCommand(ui, &mockProjectListService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected organization error")
	}
}

func TestProjectListComprehensiveSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockProjectListService{
		response: &tfe.ProjectList{
			Items: []*tfe.Project{
				{
					ID:          "prj-1",
					Name:        "my-project",
					Description: "A test project",
				},
			},
		},
	}
	cmd := newProjectListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-org=my-org"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if svc.lastOrg != "my-org" {
		t.Fatalf("expected organization my-org, got %s", svc.lastOrg)
	}
	if !strings.Contains(output, "my-project") {
		t.Fatalf("expected project name in output, got: %s", output)
	}
	if !strings.Contains(output, "prj-1") {
		t.Fatalf("expected project ID in output, got: %s", output)
	}
}

func TestProjectListComprehensiveOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockProjectListService{
		response: &tfe.ProjectList{
			Items: []*tfe.Project{
				{
					ID:          "prj-1",
					Name:        "my-project",
					Description: "A test project",
				},
			},
		},
	}
	cmd := newProjectListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-org=my-org", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	var rows []map[string]string
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if rows[0]["Name"] != "my-project" {
		t.Fatalf("unexpected row: %#v", rows)
	}
}

func TestProjectListComprehensiveHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockProjectListService{err: errors.New("API error")}
	cmd := newProjectListCommand(ui, svc)

	if code := cmd.Run([]string{"-org=my-org"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "API error") {
		t.Fatalf("expected error output, got: %s", ui.ErrorWriter.String())
	}
}
