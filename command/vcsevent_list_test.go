package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVCSEventListCommand_Run_NoOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VCSEventListCommand{
		Meta: testMeta(t, ui),
	}

	code := cmd.Run([]string{})
	if code == 0 {
		t.Fatal("Expected non-zero exit code when organization is missing")
	}

	output := ui.ErrorWriter.String()
	if !contains(output, "organization") {
		t.Error("Error should mention missing organization flag")
	}
}

func TestVCSEventListCommand_Run_InvalidLevel(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VCSEventListCommand{
		Meta: testMeta(t, ui),
	}

	code := cmd.Run([]string{"-org=test-org", "-level=invalid"})
	if code == 0 {
		t.Fatal("Expected non-zero exit code when level is invalid")
	}

	output := ui.ErrorWriter.String()
	if !contains(output, "level must be either 'info' or 'error'") {
		t.Error("Error should mention valid level values")
	}
}

func TestVCSEventListCommand_RunNoResults(t *testing.T) {
	ui := cli.NewMockUi()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/api/v2/ping" || r.URL.RequestURI() == "/api/v2/ping?" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.Path == "/api/v2/organizations/my-org/vcs-events" && r.URL.RawQuery == "" {
			_, _ = w.Write([]byte(`{"data":[]}`))
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.RequestURI())
	}))
	defer server.Close()

	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &VCSEventListCommand{Meta: Meta{Ui: ui, client: apiClient}}
	code := cmd.Run([]string{"-organization=my-org"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if !contains(ui.OutputWriter.String(), "No VCS events found") {
		t.Fatalf("expected no VCS events output, got %q", ui.OutputWriter.String())
	}
}

func TestVCSEventListCommand_RunHasResults(t *testing.T) {
	ui := cli.NewMockUi()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/api/v2/ping" || r.URL.RequestURI() == "/api/v2/ping?" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.Path == "/api/v2/organizations/my-org/vcs-events" && r.URL.RawQuery == "" {
			_, _ = w.Write([]byte(`{"data":[{"id":"ve-1","type":"vcs-events","attributes":{"created-at":"2024-01-01T00:00:00Z","level":"error","message":"token expired","organization-id":"org-123"}},{"id":"ve-2","type":"vcs-events","attributes":{"created-at":"2024-01-01T01:00:00Z","level":"info","message":"webhook received","organization-id":"org-123"}}]}`))
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.RequestURI())
	}))
	defer server.Close()

	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &VCSEventListCommand{Meta: Meta{Ui: ui, client: apiClient}}
	code := cmd.Run([]string{"-organization=my-org"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	out := strings.TrimSpace(ui.OutputWriter.String())
	if out == "" {
		t.Fatal("expected output from table render")
	}
}

func TestVCSEventListCommand_RunAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/api/v2/ping" || r.URL.RequestURI() == "/api/v2/ping?" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.Path == "/api/v2/organizations/my-org/vcs-events" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"errors":[{"status":"404"}]}`))
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.RequestURI())
	}))
	defer server.Close()

	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &VCSEventListCommand{Meta: Meta{Ui: ui, client: apiClient}}
	code := cmd.Run([]string{"-organization=my-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !contains(ui.ErrorWriter.String(), "API request failed with status 404") {
		t.Fatalf("expected API error output, got %q", ui.ErrorWriter.String())
	}
}

func TestVCSEventListCommand_RunQueryFilters(t *testing.T) {
	ui := cli.NewMockUi()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/api/v2/ping" || r.URL.RequestURI() == "/api/v2/ping?" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.Path == "/api/v2/organizations/my-org/vcs-events" {
			if !contains(r.URL.RawQuery, "filter[levels]=error") || !contains(r.URL.RawQuery, "filter[from]=2024-01-01T00:00:00Z") {
				t.Fatalf("expected level and from query params, got %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"data":[]}`))
			return
		}
			t.Fatalf("unexpected path: %s", r.URL.RequestURI())
	}))
	defer server.Close()

	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &VCSEventListCommand{Meta: Meta{Ui: ui, client: apiClient}}
	code := cmd.Run([]string{"-organization=my-org", "-level=error", "-from=2024-01-01T00:00:00Z"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
}

func TestVCSEventListCommand_RunJSONOutput(t *testing.T) {
	ui := cli.NewMockUi()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/api/v2/ping" || r.URL.RequestURI() == "/api/v2/ping?" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.Path == "/api/v2/organizations/my-org/vcs-events" && r.URL.RawQuery == "" {
			_, _ = w.Write([]byte(`{"data":[{"id":"ve-1","type":"vcs-events","attributes":{"created-at":"2024-01-01T00:00:00Z","level":"error","message":"token expired","organization-id":"org-123"}}]}`))
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.RequestURI())
	}))
	defer server.Close()

	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &VCSEventListCommand{Meta: Meta{Ui: ui, client: apiClient}}
	code := cmd.Run([]string{"-organization=my-org", "-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	out := strings.TrimSpace(ui.OutputWriter.String())
	var rows []map[string]string
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("failed to decode json output: %v, output: %q", err, out)
	}
	if len(rows) != 1 || rows[0]["ID"] != "ve-1" {
		t.Fatalf("expected single json row with id ve-1, got %v", rows)
	}
}
