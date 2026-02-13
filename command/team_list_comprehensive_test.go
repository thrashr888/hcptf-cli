package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newTeamListCommand(ui cli.Ui, svc teamLister) *TeamListCommand {
	return &TeamListCommand{
		Meta:    newTestMeta(ui),
		teamSvc: svc,
	}
}

func TestTeamListComprehensiveRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamListCommand(ui, &mockTeamListService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected organization error")
	}
}

func TestTeamListComprehensiveSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamListService{
		response: &tfe.TeamList{
			Items: []*tfe.Team{
				{
					ID:         "team-1",
					Name:       "developers",
					Visibility: "organization",
				},
			},
		},
	}
	cmd := newTeamListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-org=my-org"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if svc.lastOrg != "my-org" {
		t.Fatalf("expected organization my-org, got %s", svc.lastOrg)
	}
	if !strings.Contains(output, "developers") {
		t.Fatalf("expected team name in output, got: %s", output)
	}
	if !strings.Contains(output, "team-1") {
		t.Fatalf("expected team ID in output, got: %s", output)
	}
}

func TestTeamListComprehensiveOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamListService{
		response: &tfe.TeamList{
			Items: []*tfe.Team{
				{
					ID:         "team-1",
					Name:       "developers",
					Visibility: "organization",
				},
			},
		},
	}
	cmd := newTeamListCommand(ui, svc)

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
	if rows[0]["Name"] != "developers" {
		t.Fatalf("unexpected row: %#v", rows)
	}
}

func TestTeamListComprehensiveHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamListService{err: errors.New("API error")}
	cmd := newTeamListCommand(ui, svc)

	if code := cmd.Run([]string{"-org=my-org"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "API error") {
		t.Fatalf("expected error output, got: %s", ui.ErrorWriter.String())
	}
}
