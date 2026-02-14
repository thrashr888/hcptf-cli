package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicyCheckOverrideRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicyCheckOverrideCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestPolicyCheckOverrideHelp(t *testing.T) {
	cmd := &PolicyCheckOverrideCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policycheck override") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "-auto-approve") {
		t.Error("Help should mention -auto-approve flag")
	}
	if !strings.Contains(help, "soft-mandatory") {
		t.Error("Help should mention soft-mandatory policy checks")
	}
}

func TestPolicyCheckOverrideSynopsis(t *testing.T) {
	cmd := &PolicyCheckOverrideCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Override a soft-mandatory policy check" {
		t.Errorf("expected 'Override a soft-mandatory policy check', got %q", synopsis)
	}
}

func TestPolicyCheckOverrideFlagParsing(t *testing.T) {
	tests := []struct {
		name                string
		args                []string
		expectedID          string
		expectedFormat      string
		expectedAutoApprove bool
	}{
		{
			name:                "policy check id, default format, no auto-approve",
			args:                []string{"-id=polchk-abc123"},
			expectedID:          "polchk-abc123",
			expectedFormat:      "table",
			expectedAutoApprove: false,
		},
		{
			name:                "policy check id, table format, no auto-approve",
			args:                []string{"-id=polchk-xyz789", "-output=table"},
			expectedID:          "polchk-xyz789",
			expectedFormat:      "table",
			expectedAutoApprove: false,
		},
		{
			name:                "policy check id, json format, no auto-approve",
			args:                []string{"-id=polchk-def456", "-output=json"},
			expectedID:          "polchk-def456",
			expectedFormat:      "json",
			expectedAutoApprove: false,
		},
		{
			name:                "policy check id with auto-approve",
			args:                []string{"-id=polchk-prod999", "-auto-approve"},
			expectedID:          "polchk-prod999",
			expectedFormat:      "table",
			expectedAutoApprove: true,
		},
		{
			name:                "policy check id with auto-approve and json format",
			args:                []string{"-id=polchk-staging", "-auto-approve", "-output=json"},
			expectedID:          "polchk-staging",
			expectedFormat:      "json",
			expectedAutoApprove: true,
		},
		{
			name:                "explicit auto-approve true",
			args:                []string{"-id=polchk-test", "-auto-approve=true"},
			expectedID:          "polchk-test",
			expectedFormat:      "table",
			expectedAutoApprove: true,
		},
		{
			name:                "explicit auto-approve false",
			args:                []string{"-id=polchk-dev", "-auto-approve=false"},
			expectedID:          "polchk-dev",
			expectedFormat:      "table",
			expectedAutoApprove: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicyCheckOverrideCommand{}

			flags := cmd.Meta.FlagSet("policycheck override")
			flags.StringVar(&cmd.policyCheckID, "id", "", "Policy Check ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")
			flags.BoolVar(&cmd.autoApprove, "auto-approve", false, "Skip confirmation prompt")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the policy check ID was set correctly
			if cmd.policyCheckID != tt.expectedID {
				t.Errorf("expected policyCheckID %q, got %q", tt.expectedID, cmd.policyCheckID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}

			// Verify the auto-approve flag was set correctly
			if cmd.autoApprove != tt.expectedAutoApprove {
				t.Errorf("expected autoApprove %v, got %v", tt.expectedAutoApprove, cmd.autoApprove)
			}
		})
	}
}

func TestPolicyCheckOverrideRunAutoApproveSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		case "/api/v2/policy-checks/pc-1":
			body := map[string]interface{}{
				"data": map[string]interface{}{
					"id":   "pc-1",
					"type": "policy-checks",
					"attributes": map[string]interface{}{
						"actions": map[string]interface{}{
							"is-overridable": true,
						},
						"permissions": map[string]interface{}{
							"can-override": true,
						},
						"scope":  "organization",
						"status": "passed",
						"result": map[string]interface{}{
							"soft-failed": 0,
						},
					},
				},
			}
			if encoded, err := json.Marshal(body); err == nil {
				_, _ = w.Write(encoded)
			}
			return
		case "/api/v2/policy-checks/pc-1/actions/override":
			body := map[string]interface{}{
				"data": map[string]interface{}{
					"id":   "pc-1",
					"type": "policy-checks",
					"attributes": map[string]interface{}{
						"actions": map[string]interface{}{
							"is-overridable": true,
						},
						"permissions": map[string]interface{}{
							"can-override": true,
						},
						"scope":  "organization",
						"status": "overridden",
						"result": map[string]interface{}{
							"soft-failed": 0,
						},
					},
				},
			}
			if encoded, err := json.Marshal(body); err == nil {
				_, _ = w.Write(encoded)
			}
			return
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &PolicyCheckOverrideCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=pc-1", "-auto-approve", "-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, output=%q, err=%q", code, ui.OutputWriter.String(), ui.ErrorWriter.String())
	}

	out := strings.TrimSpace(ui.OutputWriter.String())
	var data map[string]interface{}
	start := strings.Index(out, "{")
	end := strings.LastIndex(out, "}")
	if start == -1 || end == -1 || end <= start {
		t.Fatalf("expected JSON output in response, got %q", out)
	}
	jsonOutput := out[start : end+1]
	if err := json.Unmarshal([]byte(jsonOutput), &data); err != nil {
		t.Fatalf("failed to decode json output: %v, output: %q", err, jsonOutput)
	}
	if data["Status"] != "overridden" {
		t.Fatalf("expected status overridden, got %v", data["Status"])
	}
}

func TestPolicyCheckOverrideRunCannotOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		case "/api/v2/policy-checks/pc-1":
			body := map[string]interface{}{
				"data": map[string]interface{}{
					"id":   "pc-1",
					"type": "policy-checks",
					"attributes": map[string]interface{}{
						"actions": map[string]interface{}{
							"is-overridable": false,
						},
						"permissions": map[string]interface{}{
							"can-override": true,
						},
						"scope":  "organization",
						"status": "soft_failed",
						"result": map[string]interface{}{
							"soft-failed": 1,
						},
					},
				},
			}
			if encoded, err := json.Marshal(body); err == nil {
				_, _ = w.Write(encoded)
			}
			return
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &PolicyCheckOverrideCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=pc-1", "-auto-approve"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "cannot be overridden") {
		t.Fatalf("expected cannot be overridden error, got %q", ui.ErrorWriter.String())
	}
}
