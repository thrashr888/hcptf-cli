package command

import (
	"github.com/mitchellh/cli"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFeatureSetListCommand_RequiresNoArgs(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &FeatureSetListCommand{Meta: newTestMeta(ui)}

	if cmd == nil {
		t.Fatal("command should initialize")
	}
}

func TestFeatureSetListCommand_FetchesExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/ping":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/feature-sets":
			if r.Method != http.MethodGet {
				t.Fatalf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":[{"id":"feat-1","type":"feature-sets","attributes":{"name":"plan-locking"}}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &FeatureSetListCommand{Meta: Meta{Ui: ui, client: apiClient}}
	code := cmd.Run([]string{"-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	out := ui.OutputWriter.String()
	if !strings.Contains(out, "\"id\": \"feat-1\"") {
		t.Fatalf("expected JSON output to include id feat-1, got %q", out)
	}
}

func TestFeatureSetListCommand_APIError(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &FeatureSetListCommand{Meta: newTestMeta(ui)}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/ping" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.Path == "/api/v2/feature-sets" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"errors":[{"status":"404"}]}`))
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.Path)
	}))
	defer server.Close()

	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd = &FeatureSetListCommand{Meta: Meta{Ui: ui, client: apiClient}}
	if code := cmd.Run([]string{}); code != 1 {
		t.Fatalf("expected exit 1 on API error")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "API request failed with status 404") {
		t.Fatalf("expected API error output, got: %q", ui.ErrorWriter.String())
	}
}
