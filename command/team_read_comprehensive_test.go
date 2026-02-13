package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newTeamShowCommand(ui cli.Ui, svc teamReader) *TeamShowCommand {
	return &TeamShowCommand{
		Meta:    newTestMeta(ui),
		teamSvc: svc,
	}
}

func TestTeamShowRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamShowCommand(ui, &mockTeamReadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected organization error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=org"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected name error")
	}
}

func TestTeamShowHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamReadService{err: errors.New("boom")}
	cmd := newTeamShowCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=org", "-name=team1"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastName != "team1" {
		t.Fatalf("unexpected team name recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestTeamShowOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamReadService{
		response: &tfe.Team{
			ID:   "team-123",
			Name: "team1",
		},
	}
	cmd := newTeamShowCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=org", "-name=team1", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "team-123" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
