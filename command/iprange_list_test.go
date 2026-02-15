package command

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestIPRangeListCommand_FetchesExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/ping" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.Method != http.MethodGet {
			t.Fatalf("expected method %s, got %s", http.MethodGet, r.Method)
		}
		if r.URL.Path != "/api/v2/meta/ip-ranges" {
			t.Fatalf("expected path /api/v2/meta/ip-ranges, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"id":"ip-1","type":"ip-ranges","attributes":{"name":"public"}}]}`))
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &IPRangeListCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	out := ui.OutputWriter.String()
	if !strings.Contains(out, "\"id\": \"ip-1\"") {
		t.Fatalf("expected JSON output to include id ip-1, got %q", out)
	}
}
