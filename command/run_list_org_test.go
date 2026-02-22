package command

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

type mockRunOrgListService struct {
	response *tfe.OrganizationRunList
	err      error
	lastOrg  string
	lastOpts *tfe.RunListForOrganizationOptions
}

func (m *mockRunOrgListService) ListForOrganization(_ context.Context, organization string, options *tfe.RunListForOrganizationOptions) (*tfe.OrganizationRunList, error) {
	m.lastOrg = organization
	if options != nil {
		copy := *options
		m.lastOpts = &copy
	}
	return m.response, m.err
}

func newRunListOrgCommand(ui cli.Ui, svc runOrgLister) *RunListOrgCommand {
	return &RunListOrgCommand{
		Meta:   newTestMeta(ui),
		runSvc: svc,
	}
}

func TestRunListOrgRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newRunListOrgCommand(ui, &mockRunOrgListService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected organization error")
	}
}

func TestRunListOrgHandlesError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunOrgListService{err: errors.New("boom")}
	cmd := newRunListOrgCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestRunListOrgPassesOptionsAndOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockRunOrgListService{
		response: &tfe.OrganizationRunList{
			Items: []*tfe.Run{{
				ID:        "run-1",
				Status:    tfe.RunApplied,
				Source:    tfe.RunSourceAPI,
				Message:   "ok",
				CreatedAt: time.Unix(0, 0),
				Workspace: &tfe.Workspace{Name: "prod"},
			}},
		},
	}
	cmd := newRunListOrgCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org",
			"-search=drift",
			"-status=planned,applied",
			"-include=workspace,plan",
			"-output=json",
		})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	if svc.lastOrg != "my-org" || svc.lastOpts == nil {
		t.Fatalf("expected request captured")
	}
	if svc.lastOpts.Basic != "drift" || svc.lastOpts.Status != "planned,applied" {
		t.Fatalf("expected options passed, got %#v", svc.lastOpts)
	}
	if len(svc.lastOpts.Include) != 2 {
		t.Fatalf("expected include options, got %#v", svc.lastOpts.Include)
	}

	var rows []map[string]string
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if len(rows) != 1 || rows[0]["ID"] != "run-1" {
		t.Fatalf("unexpected output rows: %#v", rows)
	}
}
