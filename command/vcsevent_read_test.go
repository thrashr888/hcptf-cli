package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVCSEventReadCommand_Run_NoID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VCSEventReadCommand{
		Meta: testMeta(t, ui),
	}

	code := cmd.Run([]string{})
	if code == 0 {
		t.Fatal("Expected non-zero exit code when ID is missing")
	}

	output := ui.ErrorWriter.String()
	if !contains(output, "id") {
		t.Error("Error should mention missing id flag")
	}
}

func TestVCSEventReadCommand_RunSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/api/v2/ping" || r.URL.RequestURI() == "/api/v2/ping?" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.Path == "/api/v2/vcs-events/ve-1" {
			_, _ = w.Write([]byte(`{"data":{"id":"ve-1","type":"vcs-events","attributes":{"created-at":"2024-01-01T00:00:00Z","level":"error","message":"token expired","organization-id":"org-123","suggested_action":"Rotate token"}}}`))
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.RequestURI())
	}))
	defer server.Close()

	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &VCSEventReadCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=ve-1"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, output=%q, err=%q", code, ui.OutputWriter.String(), ui.ErrorWriter.String())
	}
	if !contains(ui.OutputWriter.String(), "SuggestedAction") {
		t.Fatalf("expected suggested action in output, got %q", ui.OutputWriter.String())
	}
}

func TestVCSEventReadCommand_RunNotFound(t *testing.T) {
	ui := cli.NewMockUi()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/api/v2/ping" || r.URL.RequestURI() == "/api/v2/ping?" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.Path == "/api/v2/vcs-events/ve-1" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"errors":[{"status":"404"}]}`))
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.RequestURI())
	}))
	defer server.Close()

	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &VCSEventReadCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=ve-1"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !contains(ui.ErrorWriter.String(), "API request failed with status 404") {
		t.Fatalf("expected 404 output, got %q", ui.ErrorWriter.String())
	}
}

func TestVCSEventReadCommand_RunJSONOutput(t *testing.T) {
	ui := cli.NewMockUi()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/api/v2/ping" || r.URL.RequestURI() == "/api/v2/ping?" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.Path == "/api/v2/vcs-events/ve-1" {
			_, _ = w.Write([]byte(`{"data":{"id":"ve-1","type":"vcs-events","attributes":{"created-at":"2024-01-01T00:00:00Z","level":"error","message":"token expired","organization-id":"org-123","suggested_action":"Rotate token"}}}`))
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.RequestURI())
	}))
	defer server.Close()

	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &VCSEventReadCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=ve-1", "-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, output=%q, err=%q", code, ui.OutputWriter.String(), ui.ErrorWriter.String())
	}

	out := strings.TrimSpace(ui.OutputWriter.String())
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("failed to decode json output: %v, output: %q", err, out)
	}
	if data["ID"] != "ve-1" {
		t.Fatalf("expected ID in json output, got %v", data["ID"])
	}
}
