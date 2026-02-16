package command

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicyCheckReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicyCheckReadCommand{
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

func TestPolicyCheckReadHelp(t *testing.T) {
	cmd := &PolicyCheckReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policy check read") {
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
}

func TestPolicyCheckReadSynopsis(t *testing.T) {
	cmd := &PolicyCheckReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read policy check details" {
		t.Errorf("expected 'Read policy check details', got %q", synopsis)
	}
}

func TestPolicyCheckReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "policy check id, default format",
			args:           []string{"-id=polchk-abc123"},
			expectedID:     "polchk-abc123",
			expectedFormat: "table",
		},
		{
			name:           "policy check id, table format",
			args:           []string{"-id=polchk-xyz789", "-output=table"},
			expectedID:     "polchk-xyz789",
			expectedFormat: "table",
		},
		{
			name:           "policy check id, json format",
			args:           []string{"-id=polchk-def456", "-output=json"},
			expectedID:     "polchk-def456",
			expectedFormat: "json",
		},
		{
			name:           "different policy check id format",
			args:           []string{"-id=polchk-prod999", "-output=table"},
			expectedID:     "polchk-prod999",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicyCheckReadCommand{}

			flags := cmd.Meta.FlagSet("policy check read")
			flags.StringVar(&cmd.policyCheckID, "id", "", "Policy Check ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

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
		})
	}
}

func TestPolicyCheckReadRunSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			w.WriteHeader(http.StatusNoContent)
			return
		case "/api/v2/policy-checks/pchk-001":
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"data": {
					"id": "pchk-001",
					"type": "policy-checks",
					"attributes": {
						"status": "soft_failed",
						"scope": "workspace",
						"actions": {
							"is-overridable": true
						},
						"permissions": {
							"can-override": true
						},
						"result": {
							"passed": 0,
							"total-failed": 1,
							"hard-failed": 0,
							"soft-failed": 1,
							"advisory-failed": 0,
							"duration": 100,
							"result": false
						}
					}
				}
			}`))
			return
		default:
			t.Fatalf("unexpected request: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	t.Setenv("HCPTF_ADDRESS", server.URL)
	t.Setenv("TFE_TOKEN", "test-token")

	ui := cli.NewMockUi()
	cmd := &PolicyCheckReadCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	code := cmd.Run([]string{
		"-id=pchk-001",
		"-output=json",
	})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	output := ui.OutputWriter.String()
	if !strings.Contains(output, "pchk-001") {
		t.Fatalf("expected output to include policy check id, got %q", output)
	}
}

func TestPolicyCheckReadRunNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			w.WriteHeader(http.StatusNoContent)
			return
		case "/api/v2/policy-checks/pchk-missing":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"errors":[{"status":"404","detail":"not found"}]}`))
			return
		default:
			t.Fatalf("unexpected request: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	t.Setenv("HCPTF_ADDRESS", server.URL)
	t.Setenv("TFE_TOKEN", "test-token")

	ui := cli.NewMockUi()
	cmd := &PolicyCheckReadCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	code := cmd.Run([]string{"-id=pchk-missing", "-output=json"})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if out := ui.ErrorWriter.String(); !strings.Contains(out, "not found") && !strings.Contains(out, "404") {
		t.Fatalf("expected not found error output, got %q", out)
	}
}
