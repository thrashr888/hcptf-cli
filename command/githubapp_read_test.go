package command

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestGitHubAppReadCommand_RequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GitHubAppReadCommand{Meta: newTestMeta(ui)}

	if code := cmd.Run(nil); code == 0 {
		t.Fatal("expected non-zero exit code when -id missing")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "id") {
		t.Fatalf("expected id error, got: %q", ui.ErrorWriter.String())
	}
}

func TestGitHubAppReadCommand_FetchesExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/ping" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.Method != http.MethodGet {
			t.Fatalf("expected method %s, got %s", http.MethodGet, r.Method)
		}
		if r.URL.Path != "/api/v2/github-app-installations/gha-1" {
			t.Fatalf("expected path /api/v2/github-app-installations/gha-1, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"id":"gha-1","type":"github-app-installations","attributes":{"name":"demo"}}}`))
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &GitHubAppReadCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=gha-1", "-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	out := ui.OutputWriter.String()
	if !strings.Contains(out, "\"id\": \"gha-1\"") {
		t.Fatalf("expected JSON output to include id gha-1, got %q", out)
	}
}
