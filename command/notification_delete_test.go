package command

import (
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newNotificationDeleteCommand(ui *cli.MockUi, svc notificationDeleter) *NotificationDeleteCommand {
	cmd := &NotificationDeleteCommand{
		Meta:     newTestMeta(ui),
		notifSvc: svc,
	}
	cmd.Meta.Ui = ui
	return cmd
}

func TestNotificationDeleteRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newNotificationDeleteCommand(ui, &mockNotificationDeleteService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing id, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error, got: %s", ui.ErrorWriter.String())
	}
}

func TestNotificationDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockNotificationDeleteService{err: errors.New("boom")}
	cmd := newNotificationDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=nc-123", "-force"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastID != "nc-123" {
		t.Fatalf("unexpected delete id: %s", svc.lastID)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got: %s", ui.ErrorWriter.String())
	}
}

func TestNotificationDeleteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockNotificationDeleteService{}
	cmd := newNotificationDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=nc-123", "-force"}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if svc.lastID != "nc-123" {
		t.Fatalf("unexpected delete id: %s", svc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success message, got: %s", ui.OutputWriter.String())
	}
}

func TestNotificationDeleteReadsNotificationForConfirmation(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("yes\n")
	svc := &mockNotificationDeleteService{
		response: &tfe.NotificationConfiguration{
			ID:   "nc-123",
			Name: "Deploy Alerts",
		},
	}
	cmd := newNotificationDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=nc-123"}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if svc.lastRead != "nc-123" {
		t.Fatalf("expected read id recorded: %s", svc.lastRead)
	}
	if svc.lastID != "nc-123" {
		t.Fatalf("unexpected delete id: %s", svc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "Deploy Alerts") {
		t.Fatalf("expected confirmation name in output, got: %s", ui.OutputWriter.String())
	}
}

func TestNotificationDeleteReadFails(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockNotificationDeleteService{
		response: &tfe.NotificationConfiguration{},
		readErr:  errors.New("read failed"),
	}
	cmd := newNotificationDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=nc-123"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastRead != "nc-123" {
		t.Fatalf("expected read id recorded: %s", svc.lastRead)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "Error reading notification configuration: read failed") {
		t.Fatalf("expected read error output, got: %s", ui.ErrorWriter.String())
	}
}
