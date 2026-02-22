package command

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newPolicyUploadCommand(ui cli.Ui, svc policyUploader) *PolicyUploadCommand {
	return &PolicyUploadCommand{
		Meta:      newTestMeta(ui),
		policySvc: svc,
	}
}

func TestPolicyUploadRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newPolicyUploadCommand(ui, &mockPolicyUploadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestPolicyUploadHandlesServiceError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicyUploadService{err: errors.New("boom")}
	cmd := newPolicyUploadCommand(ui, svc)

	tempDir := t.TempDir()
	policyFile := filepath.Join(tempDir, "policy.sentinel")
	if err := os.WriteFile(policyFile, []byte("main = rule { true }"), 0o600); err != nil {
		t.Fatalf("failed to write policy file: %v", err)
	}

	if code := cmd.Run([]string{"-id=pol-1", "-policy-file=" + policyFile}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID != "pol-1" {
		t.Fatalf("expected policy id")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected service error")
	}
}

func TestPolicyUploadSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicyUploadService{}
	cmd := newPolicyUploadCommand(ui, svc)

	tempDir := t.TempDir()
	policyFile := filepath.Join(tempDir, "policy.sentinel")
	content := []byte("main = rule { true }")
	if err := os.WriteFile(policyFile, content, 0o600); err != nil {
		t.Fatalf("failed to write policy file: %v", err)
	}

	if code := cmd.Run([]string{"-id=pol-1", "-policy-file=" + policyFile}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastID != "pol-1" || string(svc.lastData) != string(content) {
		t.Fatalf("expected upload request to be captured")
	}
}
