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

func newPolicySetReadCommand(ui cli.Ui, svc policySetReader) *PolicySetReadCommand {
	return &PolicySetReadCommand{
		Meta:         newTestMeta(ui),
		policySetSvc: svc,
	}
}

func TestPolicySetReadRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newPolicySetReadCommand(ui, &mockPolicySetReadService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestPolicySetReadHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicySetReadService{err: errors.New("boom")}
	cmd := newPolicySetReadCommand(ui, svc)

	if code := cmd.Run([]string{"-id=polset-123"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID != "polset-123" {
		t.Fatalf("unexpected policy set id recorded")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestPolicySetReadOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicySetReadService{
		response: &tfe.PolicySet{
			ID:          "polset-123",
			Name:        "test-policyset",
			Description: "A test policy set",
			Global:      true,
			PolicyCount: 5,
			CreatedAt:   time.Unix(0, 0),
			UpdatedAt:   time.Unix(0, 0),
		},
	}
	cmd := newPolicySetReadCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=polset-123", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "polset-123" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
