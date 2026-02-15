package command

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAuditTrailReadHelp(t *testing.T) {
	cmd := &AuditTrailReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf audit trail read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "Audit trail event ID") {
		t.Error("Help should mention audit trail event ID")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "table (default) or json") {
		t.Error("Help should mention output formats")
	}
	if !strings.Contains(help, "organization token or audit trail token") {
		t.Error("Help should mention required token types")
	}
}

func TestAuditTrailReadSynopsis(t *testing.T) {
	cmd := &AuditTrailReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read audit trail event details" {
		t.Errorf("expected 'Read audit trail event details', got %q", synopsis)
	}
}

func TestAuditTrailReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "default output format",
			args:           []string{"-id=ae66e491-db59-457c-8445-9c908ee726ae"},
			expectedID:     "ae66e491-db59-457c-8445-9c908ee726ae",
			expectedFormat: "table",
		},
		{
			name:           "json output format",
			args:           []string{"-id=ae66e491-db59-457c-8445-9c908ee726ae", "-output=json"},
			expectedID:     "ae66e491-db59-457c-8445-9c908ee726ae",
			expectedFormat: "json",
		},
		{
			name:           "table output format explicitly set",
			args:           []string{"-id=test-id-123", "-output=table"},
			expectedID:     "test-id-123",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AuditTrailReadCommand{}

			flags := cmd.Meta.FlagSet("audit trail read")
			flags.StringVar(&cmd.id, "id", "", "Audit trail event ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}

func TestAuditTrailReadRunSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			w.WriteHeader(http.StatusNoContent)
			return
		case "/api/v2/organization/audit-trail":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"data": [
					{
						"id": "at-001",
						"version": "1",
						"type": "organization",
						"timestamp": "2026-01-01T00:00:00Z",
						"auth": {
							"accessor_id": "user-123",
							"description": "user@example.com",
							"type": "api-token",
							"organization_id": "org-001"
						},
						"request": {
							"id": "req-001"
						},
						"resource": {
							"id": "ws-001",
							"type": "workspaces",
							"action": "read",
							"meta": {
								"trace_id": "trace-001"
							}
						}
					},
					{
						"id": "at-002",
						"version": "1",
						"type": "organization",
						"timestamp": "2026-01-01T00:01:00Z",
						"auth": {
							"accessor_id": "user-456",
							"description": "other@example.com",
							"type": "api-token",
							"organization_id": "org-001"
						},
						"request": {
							"id": "req-002"
						},
						"resource": {
							"id": "wk-002",
							"type": "workspaces",
							"action": "write",
							"meta": {
								"trace_id": "trace-002"
							}
						}
					}
				],
				"pagination": {
					"current_page": 1,
					"prev_page": 0,
					"next_page": 0,
					"total_pages": 1,
					"total_count": 2
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
	cmd := &AuditTrailReadCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-id=at-001",
			"-output=json",
		})
	})

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(output, "at-001") {
		t.Fatalf("expected output to include audit trail id, got %q", output)
	}
	if !strings.Contains(output, "workspace") {
		t.Fatalf("expected output to include resource type, got %q", output)
	}
}

func TestAuditTrailReadRunNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			w.WriteHeader(http.StatusNoContent)
			return
		case "/api/v2/organization/audit-trail":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"data": [
					{
						"id": "at-other-1",
						"version": "1",
						"type": "organization",
						"timestamp": "2026-01-01T00:00:00Z",
						"auth": {
							"accessor_id": "user-123",
							"description": "user@example.com",
							"type": "api-token",
							"organization_id": "org-001"
						},
						"request": {"id": "req-001"},
						"resource": {
							"id": "ws-001",
							"type": "workspaces",
							"action": "read",
							"meta": {"trace_id": "trace-001"}
						}
					}
				],
				"pagination": {
					"current_page": 1,
					"prev_page": 0,
					"next_page": 0,
					"total_pages": 1,
					"total_count": 1
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
	cmd := &AuditTrailReadCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	code := cmd.Run([]string{"-id=at-missing"})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if out := ui.ErrorWriter.String(); !strings.Contains(out, "not found") {
		t.Fatalf("expected not found error output, got %q", out)
	}
}
