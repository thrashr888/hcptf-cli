package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newOrganizationUpdateCommand(ui cli.Ui, svc organizationUpdater) *OrganizationUpdateCommand {
	return &OrganizationUpdateCommand{
		Meta:   newTestMeta(ui),
		orgSvc: svc,
	}
}

func TestOrganizationUpdateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newOrganizationUpdateCommand(ui, &mockOrganizationUpdateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing name, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected name error, got %q", ui.ErrorWriter.String())
	}
}

func TestOrganizationUpdateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockOrganizationUpdateService{err: errors.New("boom")}
	cmd := newOrganizationUpdateCommand(ui, svc)

	if code := cmd.Run([]string{"-name=my-org", "-email=test@test.com"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastName != "my-org" {
		t.Fatalf("expected lastName my-org, got %q", svc.lastName)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got %q", ui.ErrorWriter.String())
	}
}

func TestOrganizationUpdateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockOrganizationUpdateService{
		response: &tfe.Organization{
			Name:  "my-org",
			Email: "test@test.com",
		},
	}
	cmd := newOrganizationUpdateCommand(ui, svc)

	if code := cmd.Run([]string{"-name=my-org", "-email=test@test.com"}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if svc.lastName != "my-org" {
		t.Fatalf("expected lastName my-org, got %q", svc.lastName)
	}
	if !strings.Contains(ui.OutputWriter.String(), "updated successfully") {
		t.Fatalf("expected success message, got %q", ui.OutputWriter.String())
	}
}
