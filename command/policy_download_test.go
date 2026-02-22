package command

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newPolicyDownloadCommand(ui cli.Ui, svc policyDownloader) *PolicyDownloadCommand {
	return &PolicyDownloadCommand{
		Meta:      newTestMeta(ui),
		policySvc: svc,
	}
}

func TestPolicyDownloadRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newPolicyDownloadCommand(ui, &mockPolicyDownloadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestPolicyDownloadHandlesServiceError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicyDownloadService{err: errors.New("boom")}
	cmd := newPolicyDownloadCommand(ui, svc)

	tempDir := t.TempDir()
	output := filepath.Join(tempDir, "policy.sentinel")

	if code := cmd.Run([]string{"-id=pol-1", "-output=" + output}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID != "pol-1" {
		t.Fatalf("expected policy id")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected service error")
	}
}

func TestPolicyDownloadSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	content := []byte("main = rule { true }")
	svc := &mockPolicyDownloadService{content: content}
	cmd := newPolicyDownloadCommand(ui, svc)

	tempDir := t.TempDir()
	output := filepath.Join(tempDir, "policy.sentinel")

	if code := cmd.Run([]string{"-id=pol-1", "-output=" + output}); code != 0 {
		t.Fatalf("expected exit 0")
	}

	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if string(got) != string(content) {
		t.Fatalf("unexpected output content: %s", string(got))
	}
}
