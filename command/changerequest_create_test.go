package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestChangeRequestCreateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ChangeRequestCreateCommand{
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

	// Test missing subject
	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod"}); code != 1 {
		t.Fatalf("expected exit 1 missing subject, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-subject") {
		t.Fatalf("expected subject error, got %q", ui.ErrorWriter.String())
	}

	// Test missing message
	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-subject=test"}); code != 1 {
		t.Fatalf("expected exit 1 missing message, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-message") {
		t.Fatalf("expected message error, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedOrg      string
		expectedWorkspace string
		expectedSubject  string
		expectedMessage  string
		expectedFormat   string
	}{
		{
			name:             "all required flags with default format",
			args:             []string{"-organization=my-org", "-workspace=prod", "-subject=test subject", "-message=test message"},
			expectedOrg:      "my-org",
			expectedWorkspace: "prod",
			expectedSubject:  "test subject",
			expectedMessage:  "test message",
			expectedFormat:   "table",
		},
		{
			name:             "org alias flag",
			args:             []string{"-org=test-org", "-workspace=staging", "-subject=fix bug", "-message=urgent fix"},
			expectedOrg:      "test-org",
			expectedWorkspace: "staging",
			expectedSubject:  "fix bug",
			expectedMessage:  "urgent fix",
			expectedFormat:   "table",
		},
		{
			name:             "all flags with json format",
			args:             []string{"-organization=my-org", "-workspace=dev", "-subject=update", "-message=details", "-output=json"},
			expectedOrg:      "my-org",
			expectedWorkspace: "dev",
			expectedSubject:  "update",
			expectedMessage:  "details",
			expectedFormat:   "json",
		},
		{
			name:             "org alias with json format",
			args:             []string{"-org=acme", "-workspace=prod", "-subject=security", "-message=patch required", "-output=json"},
			expectedOrg:      "acme",
			expectedWorkspace: "prod",
			expectedSubject:  "security",
			expectedMessage:  "patch required",
			expectedFormat:   "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ChangeRequestCreateCommand{}

			flags := cmd.Meta.FlagSet("changerequest create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.workspace, "workspace", "", "Workspace name (required)")
			flags.StringVar(&cmd.subject, "subject", "", "Change request subject (required)")
			flags.StringVar(&cmd.message, "message", "", "Change request message (required)")
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

			// Verify the subject was set correctly
			if cmd.subject != tt.expectedSubject {
				t.Errorf("expected subject %q, got %q", tt.expectedSubject, cmd.subject)
			}

			// Verify the message was set correctly
			if cmd.message != tt.expectedMessage {
				t.Errorf("expected message %q, got %q", tt.expectedMessage, cmd.message)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}

func TestChangeRequestCreateRunSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/api/v2/ping", "/api/v2/ping?":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			_, _ = w.Write([]byte(`{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-workspace"}}}`))
		case "/api/v2/organizations/my-org/explorer/bulk-actions":
			_, _ = w.Write([]byte(`{"data":{"id":"ba-123","type":"bulk-actions","attributes":{"organization_id":"my-org","action_type":"change_requests","action_inputs":{"subject":"Fix","message":"Please update"}}}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.RequestURI())
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestCreateCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-workspace=my-workspace", "-subject=Fix", "-message=Please update"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, output=%q, err=%q", code, ui.OutputWriter.String(), ui.ErrorWriter.String())
	}

	out := ui.OutputWriter.String()
	if !strings.Contains(out, "Change request created successfully via bulk action 'ba-123'") {
		t.Fatalf("expected success output, got %q", out)
	}
}

func TestChangeRequestCreateRunWorkspaceReadError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/api/v2/ping", "/api/v2/ping?":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"errors":[{"status":"404"}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.RequestURI())
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestCreateCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-workspace=my-workspace", "-subject=Fix", "-message=Please update"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "Error reading workspace") {
		t.Fatalf("expected workspace read error output, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestCreateRunAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/api/v2/ping", "/api/v2/ping?":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			_, _ = w.Write([]byte(`{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-workspace"}}}`))
		case "/api/v2/organizations/my-org/explorer/bulk-actions":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"errors":[{"status":"404"}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.RequestURI())
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestCreateCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-workspace=my-workspace", "-subject=Fix", "-message=Please update"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "API request failed with status 404") {
		t.Fatalf("expected API error output, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestCreateRunJSONOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/api/v2/ping", "/api/v2/ping?":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/organizations/my-org/workspaces/my-workspace":
			_, _ = w.Write([]byte(`{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-workspace"}}}`))
		case "/api/v2/organizations/my-org/explorer/bulk-actions":
			_, _ = w.Write([]byte(`{"data":{"id":"ba-123","type":"bulk-actions","attributes":{"organization_id":"my-org","action_type":"change_requests","action_inputs":{"subject":"Fix","message":"Please update"}}}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.RequestURI())
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestCreateCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-workspace=my-workspace", "-subject=Fix", "-message=Please update", "-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, output=%q, err=%q", code, ui.OutputWriter.String(), ui.ErrorWriter.String())
	}

	output := strings.TrimSpace(ui.OutputWriter.String())
	start := strings.Index(output, "{")
	end := strings.LastIndex(output, "}")
	if start == -1 || end == -1 || end <= start {
		t.Fatalf("expected JSON output in response, got: %q", output)
	}
	out := output[start : end+1]

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("failed to decode json output: %v, output: %q", err, out)
	}
	if data["BulkActionID"] != "ba-123" {
		t.Fatalf("expected bulk action id in output, got %v", data["BulkActionID"])
	}
}
