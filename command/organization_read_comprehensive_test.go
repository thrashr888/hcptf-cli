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

func newOrganizationShowCommand(ui cli.Ui, svc organizationReader) *OrganizationShowCommand {
	return &OrganizationShowCommand{
		Meta:   newTestMeta(ui),
		orgSvc: svc,
	}
}

func TestOrganizationShowRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newOrganizationShowCommand(ui, &mockOrganizationReadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected name error")
	}
}

func TestOrganizationShowHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockOrganizationReadService{err: errors.New("boom")}
	cmd := newOrganizationShowCommand(ui, svc)

	if code := cmd.Run([]string{"-name=test-org"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastName != "test-org" {
		t.Fatalf("unexpected org name recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestOrganizationShowOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockOrganizationReadService{
		response: &tfe.Organization{
			Name:      "test-org",
			Email:     "test@example.com",
			CreatedAt: time.Unix(0, 0),
		},
	}
	cmd := newOrganizationShowCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-name=test-org", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["Name"] != "test-org" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
