package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newPolicySetListCommand(ui cli.Ui, svc policySetLister) *PolicySetListCommand {
	return &PolicySetListCommand{
		Meta:         newTestMeta(ui),
		policySetSvc: svc,
	}
}

func TestPolicySetListComprehensiveRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newPolicySetListCommand(ui, &mockPolicySetListService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected organization error")
	}
}

func TestPolicySetListComprehensiveSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicySetListService{
		response: &tfe.PolicySetList{
			Items: []*tfe.PolicySet{
				{
					ID:             "polset-1",
					Name:           "test-policyset",
					Description:    "A test policy set",
					Global:         true,
					PolicyCount:    3,
					WorkspaceCount: 5,
				},
			},
		},
	}
	cmd := newPolicySetListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-org=my-org"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if svc.lastOrg != "my-org" {
		t.Fatalf("expected organization my-org, got %s", svc.lastOrg)
	}
	if !strings.Contains(output, "test-policyset") {
		t.Fatalf("expected policy set name in output, got: %s", output)
	}
	if !strings.Contains(output, "polset-1") {
		t.Fatalf("expected policy set ID in output, got: %s", output)
	}
}

func TestPolicySetListComprehensiveOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicySetListService{
		response: &tfe.PolicySetList{
			Items: []*tfe.PolicySet{
				{
					ID:             "polset-1",
					Name:           "test-policyset",
					Description:    "A test policy set",
					Global:         true,
					PolicyCount:    3,
					WorkspaceCount: 5,
				},
			},
		},
	}
	cmd := newPolicySetListCommand(ui, svc)

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
	if rows[0]["Name"] != "test-policyset" {
		t.Fatalf("unexpected row: %#v", rows)
	}
}

func TestPolicySetListComprehensiveHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicySetListService{err: errors.New("API error")}
	cmd := newPolicySetListCommand(ui, svc)

	if code := cmd.Run([]string{"-org=my-org"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "API error") {
		t.Fatalf("expected error output, got: %s", ui.ErrorWriter.String())
	}
}
