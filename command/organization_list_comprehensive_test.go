package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newOrganizationListCommand(ui cli.Ui, svc organizationLister) *OrganizationListCommand {
	return &OrganizationListCommand{
		Meta:   newTestMeta(ui),
		orgSvc: svc,
	}
}

func TestOrganizationListSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockOrganizationListService{
		response: &tfe.OrganizationList{
			Items: []*tfe.Organization{
				{
					Name:                   "org1",
					Email:                  "admin@org1.com",
					CollaboratorAuthPolicy: tfe.AuthPolicyPassword,
					CreatedAt:              time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	cmd := newOrganizationListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if !strings.Contains(output, "org1") {
		t.Fatalf("expected org1 in output, got: %s", output)
	}
	if !strings.Contains(output, "admin@org1.com") {
		t.Fatalf("expected email in output, got: %s", output)
	}
}

func TestOrganizationListOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockOrganizationListService{
		response: &tfe.OrganizationList{
			Items: []*tfe.Organization{
				{
					Name:                   "org1",
					Email:                  "admin@org1.com",
					CollaboratorAuthPolicy: tfe.AuthPolicyPassword,
					CreatedAt:              time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	cmd := newOrganizationListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	var rows []map[string]string
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if rows[0]["Name"] != "org1" {
		t.Fatalf("unexpected row: %#v", rows)
	}
}

func TestOrganizationListHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockOrganizationListService{err: errors.New("API error")}
	cmd := newOrganizationListCommand(ui, svc)

	if code := cmd.Run([]string{}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "API error") {
		t.Fatalf("expected error output, got: %s", ui.ErrorWriter.String())
	}
}
