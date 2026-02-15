package command

import (
	"github.com/mitchellh/cli"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNoCodeUpdateCommand_RequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NoCodeUpdateCommand{Meta: newTestMeta(ui)}

	if code := cmd.Run([]string{"-payload={\"enabled\":false}"}); code == 0 {
		t.Fatal("expected non-zero exit code when -organization missing")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "organization") {
		t.Fatalf("expected organization error, got: %q", ui.ErrorWriter.String())
	}
}

func TestNoCodeUpdateCommand_RequiresPayload(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &NoCodeUpdateCommand{Meta: newTestMeta(ui)}

	if code := cmd.Run([]string{"-organization=my-org"}); code == 0 {
		t.Fatal("expected non-zero exit code when -payload missing")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "payload") {
		t.Fatalf("expected payload error, got: %q", ui.ErrorWriter.String())
	}
}

func TestNoCodeUpdateCommand_FetchesExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/ping" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.Method != http.MethodPatch {
			t.Fatalf("expected method %s, got %s", http.MethodPatch, r.Method)
		}
		if r.URL.Path != "/api/v2/organizations/my-org/no-code-provisioning" {
			t.Fatalf("expected path /api/v2/organizations/my-org/no-code-provisioning, got %s", r.URL.Path)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		if !strings.Contains(string(body), `"enabled":false`) {
			t.Fatalf("expected payload in request body, got %q", string(body))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"id":"nc-1","type":"no-code-provisioning","attributes":{"enabled":false}}}`))
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &NoCodeUpdateCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-organization=my-org", "-payload={\"enabled\":false}", "-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	out := ui.OutputWriter.String()
	if !strings.Contains(out, "\"id\": \"nc-1\"") {
		t.Fatalf("expected JSON output to include id nc-1, got %q", out)
	}
}
