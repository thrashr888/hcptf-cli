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

func newPolicyListCommand(ui cli.Ui, svc policyLister) *PolicyListCommand {
	return &PolicyListCommand{
		Meta:      newTestMeta(ui),
		policySvc: svc,
	}
}

func TestPolicyListComprehensiveRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newPolicyListCommand(ui, &mockPolicyListService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected organization error")
	}
}

func TestPolicyListComprehensiveSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicyListService{
		response: &tfe.PolicyList{
			Items: []*tfe.Policy{
				{
					ID:               "pol-1",
					Name:             "test-policy",
					EnforcementLevel: tfe.EnforcementAdvisory,
					PolicySetCount:   2,
					UpdatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	cmd := newPolicyListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-org=my-org"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if svc.lastOrg != "my-org" {
		t.Fatalf("expected organization my-org, got %s", svc.lastOrg)
	}
	if !strings.Contains(output, "test-policy") {
		t.Fatalf("expected policy name in output, got: %s", output)
	}
	if !strings.Contains(output, "pol-1") {
		t.Fatalf("expected policy ID in output, got: %s", output)
	}
}

func TestPolicyListComprehensiveOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicyListService{
		response: &tfe.PolicyList{
			Items: []*tfe.Policy{
				{
					ID:               "pol-1",
					Name:             "test-policy",
					EnforcementLevel: tfe.EnforcementAdvisory,
					PolicySetCount:   2,
					UpdatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	cmd := newPolicyListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-org=my-org", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	var rows []map[string]string
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if rows[0]["Name"] != "test-policy" {
		t.Fatalf("unexpected row: %#v", rows)
	}
}

func TestPolicyListComprehensiveHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicyListService{err: errors.New("API error")}
	cmd := newPolicyListCommand(ui, svc)

	if code := cmd.Run([]string{"-org=my-org"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "API error") {
		t.Fatalf("expected error output, got: %s", ui.ErrorWriter.String())
	}
}
