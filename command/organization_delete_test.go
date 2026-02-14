package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newOrganizationDeleteCommand(ui cli.Ui, svc organizationDeleter) *OrganizationDeleteCommand {
	return &OrganizationDeleteCommand{
		Meta:   newTestMeta(ui),
		orgSvc: svc,
	}
}

func TestOrganizationDeleteRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newOrganizationDeleteCommand(ui, &mockOrganizationDeleteService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing name, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected name error, got %q", ui.ErrorWriter.String())
	}
}

func TestOrganizationDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockOrganizationDeleteService{err: errors.New("boom")}
	cmd := newOrganizationDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-name=my-org", "-force"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastName != "my-org" {
		t.Fatalf("expected lastName my-org, got %q", svc.lastName)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got %q", ui.ErrorWriter.String())
	}
}

func TestOrganizationDeleteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockOrganizationDeleteService{}
	cmd := newOrganizationDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-name=my-org", "-force"}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if svc.lastName != "my-org" {
		t.Fatalf("expected lastName my-org, got %q", svc.lastName)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message, got %q", ui.OutputWriter.String())
	}
}
