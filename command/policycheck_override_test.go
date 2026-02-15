package command

import (
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
	if !strings.Contains(help, "hcptf policy check override") {
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

			flags := cmd.Meta.FlagSet("policy check override")
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

func TestPolicyCheckOverrideRunSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			w.WriteHeader(http.StatusNoContent)
			return
		case "/api/v2/policy-checks/pchk-override":
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"data": {
					"id": "pchk-override",
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
		case "/api/v2/policy-checks/pchk-override/actions/override":
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"data": {
					"id": "pchk-override",
					"type": "policy-checks",
					"attributes": {
						"status": "overridden",
						"scope": "workspace",
						"actions": {
							"is-overridable": true
						},
						"permissions": {
							"can-override": true
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
	cmd := &PolicyCheckOverrideCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-id=pchk-override",
			"-auto-approve",
			"-output=json",
		})
	})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(output, "pchk-override") {
		t.Fatalf("expected output to include policy check id, got %q", output)
	}
	if !strings.Contains(output, "overridden") {
		t.Fatalf("expected output to include overridden status, got %q", output)
	}
}

func TestPolicyCheckOverrideRunCannotOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			w.WriteHeader(http.StatusNoContent)
			return
		case "/api/v2/policy-checks/pchk-blocked":
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"data": {
					"id": "pchk-blocked",
					"type": "policy-checks",
					"attributes": {
						"status": "hard_failed",
						"scope": "workspace",
						"actions": {
							"is-overridable": false
						},
						"permissions": {
							"can-override": true
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
	cmd := &PolicyCheckOverrideCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	code := cmd.Run([]string{
		"-id=pchk-blocked",
		"-auto-approve",
	})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if out := ui.ErrorWriter.String(); !strings.Contains(out, "cannot be overridden") {
		t.Fatalf("expected cannot override error, got %q", out)
	}
}

func TestPolicyCheckOverrideRunNoPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			w.WriteHeader(http.StatusNoContent)
			return
		case "/api/v2/policy-checks/pchk-noperm":
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"data": {
					"id": "pchk-noperm",
					"type": "policy-checks",
					"attributes": {
						"status": "soft_failed",
						"scope": "workspace",
						"actions": {
							"is-overridable": true
						},
						"permissions": {
							"can-override": false
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
	cmd := &PolicyCheckOverrideCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	code := cmd.Run([]string{
		"-id=pchk-noperm",
		"-auto-approve",
	})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if out := ui.ErrorWriter.String(); !strings.Contains(out, "You do not have permission") {
		t.Fatalf("expected permission error, got %q", out)
	}
}
