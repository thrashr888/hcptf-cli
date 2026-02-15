package command

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

type mockPlanExportDeleter struct {
	lastID     string
	deleteFunc func(ctx context.Context, planExportID string) error
}

func (m *mockPlanExportDeleter) Delete(ctx context.Context, planExportID string) error {
	m.lastID = planExportID
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, planExportID)
	}
	return nil
}

func TestPlanExportDeleteCommand_RequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PlanExportDeleteCommand{Meta: newTestMeta(ui)}

	if code := cmd.Run(nil); code == 0 {
		t.Fatal("expected non-zero exit code when -id missing")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error, got: %q", ui.ErrorWriter.String())
	}
}

func TestPlanExportDeleteCommand_SuccessWithForce(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPlanExportDeleter{}
	cmd := &PlanExportDeleteCommand{Meta: testMeta(t, ui), planExportSvc: svc}

	if code := cmd.Run([]string{"-id=pe-1", "-force"}); code != 0 {
		t.Fatalf("expected exit 0, got %d; err=%q", code, ui.ErrorWriter.String())
	}
	if svc.lastID != "pe-1" {
		t.Fatalf("expected delete called with pe-1, got %q", svc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success output, got: %q", ui.OutputWriter.String())
	}
}

func TestPlanExportDeleteCommand_DeleteError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPlanExportDeleter{
		deleteFunc: func(ctx context.Context, planExportID string) error {
			return errors.New("boom")
		},
	}
	cmd := &PlanExportDeleteCommand{Meta: testMeta(t, ui), planExportSvc: svc}

	if code := cmd.Run([]string{"-id=pe-1", "-y"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got: %q", ui.ErrorWriter.String())
	}
}
