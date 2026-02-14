package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newAuditTrailTokenDeleteCommand(ui cli.Ui, svc auditTrailTokenDeleter) *AuditTrailTokenDeleteCommand {
	return &AuditTrailTokenDeleteCommand{
		Meta:        newTestMeta(ui),
		orgTokenSvc: svc,
	}
}

func TestAuditTrailTokenDeleteRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newAuditTrailTokenDeleteCommand(ui, &mockAuditTrailTokenDeleteService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing organization, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected organization error, got %q", ui.ErrorWriter.String())
	}
}

func TestAuditTrailTokenDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockAuditTrailTokenDeleteService{err: errors.New("boom")}
	cmd := newAuditTrailTokenDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-org=my-org", "-force"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastOrg != "my-org" {
		t.Fatalf("expected lastOrg my-org, got %q", svc.lastOrg)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got %q", ui.ErrorWriter.String())
	}
}

func TestAuditTrailTokenDeleteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockAuditTrailTokenDeleteService{}
	cmd := newAuditTrailTokenDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-org=my-org", "-force"}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if svc.lastOrg != "my-org" {
		t.Fatalf("expected lastOrg my-org, got %q", svc.lastOrg)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message, got %q", ui.OutputWriter.String())
	}
}
