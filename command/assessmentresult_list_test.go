package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAssessmentResultListRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AssessmentResultListCommand{
		Meta: newTestMeta(ui),
	}

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 org")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 workspace")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected name error")
	}
}

func TestAssessmentResultListHelp(t *testing.T) {
	cmd := &AssessmentResultListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	if !strings.Contains(help, "hcptf workspace run assessmentresult list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "-workspace") {
		t.Error("Help should mention -workspace flag alias")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
}

func TestAssessmentResultListSynopsis(t *testing.T) {
	cmd := &AssessmentResultListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List health assessment results for a workspace" {
		t.Errorf("unexpected synopsis: %q", synopsis)
	}
}

func TestAssessmentResultListFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOrg string
		expectedWS  string
		expectedFmt string
	}{
		{
			name:        "organization and workspace, default format",
			args:        []string{"-organization=my-org", "-name=my-ws"},
			expectedOrg: "my-org",
			expectedWS:  "my-ws",
			expectedFmt: "table",
		},
		{
			name:        "org alias",
			args:        []string{"-org=test-org", "-name=test-ws"},
			expectedOrg: "test-org",
			expectedWS:  "test-ws",
			expectedFmt: "table",
		},
		{
			name:        "workspace alias",
			args:        []string{"-org=my-org", "-workspace=prod"},
			expectedOrg: "my-org",
			expectedWS:  "prod",
			expectedFmt: "table",
		},
		{
			name:        "json format",
			args:        []string{"-org=my-org", "-name=staging", "-output=json"},
			expectedOrg: "my-org",
			expectedWS:  "staging",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AssessmentResultListCommand{}

			flags := cmd.Meta.FlagSet("assessmentresult list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.workspace, "name", "", "Workspace name (required)")
			flags.StringVar(&cmd.workspace, "workspace", "", "Workspace name (alias)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}
			if cmd.workspace != tt.expectedWS {
				t.Errorf("expected workspace %q, got %q", tt.expectedWS, cmd.workspace)
			}
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}

func TestAssessmentResultListRunNoResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			_, _ = w.Write([]byte(`{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-workspace"}}}`))
		case "/api/v2/workspaces/ws-123/assessment-results":
			_, _ = w.Write([]byte(`{"data":[]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &AssessmentResultListCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-name=my-workspace"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(ui.OutputWriter.String(), "No assessment results found") {
		t.Fatalf("expected no results output, got %q", ui.OutputWriter.String())
	}
}

func TestAssessmentResultListRunHasResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			_, _ = w.Write([]byte(`{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-workspace"}}}`))
		case "/api/v2/workspaces/ws-123/assessment-results":
			_, _ = w.Write([]byte(`{"data":[{"id":"ar-1","type":"assessment-results","attributes":{"drifted":true,"succeeded":false,"created-at":"2024-01-01T00:00:00Z","error-msg":"Drift detected"},{"id":"ar-2","type":"assessment-results","attributes":{"drifted":false,"succeeded":true,"created-at":"2024-02-01T00:00:00Z"}}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &AssessmentResultListCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-workspace=my-workspace", "-output=table"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	out := ui.OutputWriter.String()
	if !strings.Contains(out, "ar-1") || !strings.Contains(out, "ar-2") {
		t.Fatalf("expected table output with assessment results, got %q", out)
	}
}

func TestAssessmentResultListRunJSONOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			_, _ = w.Write([]byte(`{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-workspace"}}}`))
		case "/api/v2/workspaces/ws-123/assessment-results":
			_, _ = w.Write([]byte(`{"data":[{"id":"ar-1","type":"assessment-results","attributes":{"drifted":false,"succeeded":true,"created-at":"2024-01-01T00:00:00Z"}}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &AssessmentResultListCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-name=my-workspace", "-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	out := strings.TrimSpace(ui.OutputWriter.String())
	var rows []map[string]string
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("failed to decode json output: %v, output: %q", err, out)
	}
	if len(rows) != 1 || rows[0]["ID"] != "ar-1" || rows[0]["Status"] == "" {
		t.Fatalf("expected one json row with data, got %v", rows)
	}
}

func TestAssessmentResultListRunAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			_, _ = w.Write([]byte(`{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-workspace"}}}`))
		case "/api/v2/workspaces/ws-123/assessment-results":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"backend error"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &AssessmentResultListCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-org=my-org", "-name=my-workspace"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "API request failed with status 500") {
		t.Fatalf("expected API error output, got %q", ui.ErrorWriter.String())
	}
}
