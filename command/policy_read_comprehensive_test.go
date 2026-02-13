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

type mockPolicyReadDownloadService struct {
	mockPolicyReadService
	mockPolicyDownloadService
}

func newPolicyReadCommand(ui cli.Ui, svc *mockPolicyReadDownloadService) *PolicyReadCommand {
	return &PolicyReadCommand{
		Meta:        newTestMeta(ui),
		policySvc:   &svc.mockPolicyReadService,
		downloadSvc: &svc.mockPolicyDownloadService,
	}
}

func TestPolicyReadRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newPolicyReadCommand(ui, &mockPolicyReadDownloadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestPolicyReadHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicyReadDownloadService{}
	svc.mockPolicyReadService.err = errors.New("boom")
	cmd := newPolicyReadCommand(ui, svc)

	if code := cmd.Run([]string{"-id=pol-123"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.mockPolicyReadService.lastID != "pol-123" {
		t.Fatalf("unexpected policy id recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestPolicyReadOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicyReadDownloadService{}
	svc.mockPolicyReadService.response = &tfe.Policy{
		ID:               "pol-123",
		Name:             "test-policy",
		Description:      "A test policy",
		EnforcementLevel: tfe.EnforcementAdvisory,
		UpdatedAt:        time.Unix(0, 0),
	}
	svc.mockPolicyDownloadService.content = []byte("policy content")
	cmd := newPolicyReadCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=pol-123", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "pol-123" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
