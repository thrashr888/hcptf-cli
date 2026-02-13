package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestChangeRequestListRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ChangeRequestListCommand{
		Meta: newTestMeta(ui),
	}

	// Test missing organization
	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error, got %q", ui.ErrorWriter.String())
	}

	// Test missing workspace
	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing workspace, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-workspace") {
		t.Fatalf("expected workspace error, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestListFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedOrg      string
		expectedWorkspace string
		expectedFormat   string
	}{
		{
			name:             "organization and workspace with default format",
			args:             []string{"-organization=my-org", "-workspace=prod"},
			expectedOrg:      "my-org",
			expectedWorkspace: "prod",
			expectedFormat:   "table",
		},
		{
			name:             "org alias flag",
			args:             []string{"-org=test-org", "-workspace=staging"},
			expectedOrg:      "test-org",
			expectedWorkspace: "staging",
			expectedFormat:   "table",
		},
		{
			name:             "organization and workspace with table format",
			args:             []string{"-organization=my-org", "-workspace=dev", "-output=table"},
			expectedOrg:      "my-org",
			expectedWorkspace: "dev",
			expectedFormat:   "table",
		},
		{
			name:             "organization and workspace with json format",
			args:             []string{"-organization=acme", "-workspace=prod", "-output=json"},
			expectedOrg:      "acme",
			expectedWorkspace: "prod",
			expectedFormat:   "json",
		},
		{
			name:             "org alias with json format",
			args:             []string{"-org=test-org", "-workspace=qa", "-output=json"},
			expectedOrg:      "test-org",
			expectedWorkspace: "qa",
			expectedFormat:   "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ChangeRequestListCommand{}

			flags := cmd.Meta.FlagSet("changerequest list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.workspace, "workspace", "", "Workspace name (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the workspace was set correctly
			if cmd.workspace != tt.expectedWorkspace {
				t.Errorf("expected workspace %q, got %q", tt.expectedWorkspace, cmd.workspace)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}

func TestChangeRequestListRunNoResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			_, _ = w.Write([]byte(`{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-workspace"}}}`))
		case "/api/v2/workspaces/ws-123/change-requests":
			_, _ = w.Write([]byte(`{"data":[]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestListCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-workspace=my-workspace"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if !strings.Contains(ui.OutputWriter.String(), "No change requests found") {
		t.Fatalf("expected no results output, got %q", ui.OutputWriter.String())
	}
}

func TestChangeRequestListRunHasResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			_, _ = w.Write([]byte(`{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-workspace"}}}`))
		case "/api/v2/workspaces/ws-123/change-requests":
			_, _ = w.Write([]byte(`{"data":[{"id":"cr-1","type":"change-requests","attributes":{"subject":"Upgrade","message":"upgrade", "archived-at":null, "created-at":"2024-01-01T00:00:00Z", "updated-at":"2024-01-01T00:00:00Z"},"relationships":{"workspace":{"data":{"id":"ws-123","type":"workspaces"}}}}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestListCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-workspace=my-workspace", "-output=table"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	out := ui.OutputWriter.String()
	if !strings.Contains(out, "cr-1") || !strings.Contains(out, "Upgrade") {
		t.Fatalf("expected table output with change request, got %q", out)
	}
}

func TestChangeRequestListRunWorkspaceReadError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/ping" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.Path == "/api/v2/organizations/my-org/workspaces/my-workspace" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"errors":[{"status":"404"}]}`))
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.Path)
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestListCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-workspace=my-workspace"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "Error reading workspace") {
		t.Fatalf("expected workspace read error, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestListRunAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			_, _ = w.Write([]byte(`{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-workspace"}}}`))
		case "/api/v2/workspaces/ws-123/change-requests":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"backend error"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestListCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-workspace=my-workspace"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "API request failed with status 500") {
		t.Fatalf("expected API error output, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestListRunJSONOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			_, _ = w.Write([]byte(`{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-workspace"}}}`))
		case "/api/v2/workspaces/ws-123/change-requests":
			_, _ = w.Write([]byte(`{"data":[{"id":"cr-1","type":"change-requests","attributes":{"subject":"Upgrade","message":"upgrade", "archived-at":null, "created-at":"2024-01-01T00:00:00Z", "updated-at":"2024-01-01T00:00:00Z"},"relationships":{"workspace":{"data":{"id":"ws-123","type":"workspaces"}}}},{"id":"cr-2","type":"change-requests","attributes":{"subject":"Audit","message":"audit log", "archived-at":"2024-01-02T00:00:00Z", "created-at":"2024-01-01T00:00:00Z", "updated-at":"2024-01-01T00:00:00Z"},"relationships":{"workspace":{"data":{"id":"ws-123","type":"workspaces"}}}}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestListCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-workspace=my-workspace", "-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	out := strings.TrimSpace(ui.OutputWriter.String())
	var rows []map[string]string
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("failed to decode json output: %v, output: %q", err, out)
	}
	if len(rows) != 2 {
		t.Fatalf("expected two rows, got %d", len(rows))
	}
	if rows[0]["ID"] == "" || rows[1]["ID"] == "" {
		t.Fatalf("expected IDs in JSON output rows, got %v", rows)
	}
}
