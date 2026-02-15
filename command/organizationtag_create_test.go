package command

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationTagCreateCommand_RequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationTagCreateCommand{Meta: newTestMeta(ui)}

	if code := cmd.Run([]string{"-name=platform"}); code == 0 {
		t.Fatal("expected non-zero exit code when -organization missing")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "organization") {
		t.Fatalf("expected organization error, got: %q", ui.ErrorWriter.String())
	}
}

func TestOrganizationTagCreateCommand_RequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationTagCreateCommand{Meta: newTestMeta(ui)}

	if code := cmd.Run([]string{"-org=my-org"}); code == 0 {
		t.Fatal("expected non-zero exit code when -name missing")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "name") {
		t.Fatalf("expected name error, got: %q", ui.ErrorWriter.String())
	}
}

func TestOrganizationTagCreateCommand_FetchesExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/ping" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.Method != http.MethodPost {
			t.Fatalf("expected method %s, got %s", http.MethodPost, r.Method)
		}
		if r.URL.Path != "/api/v2/organizations/my-org/tags" {
			t.Fatalf("expected path /api/v2/organizations/my-org/tags, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"data":{"id":"tag-1","type":"tags","attributes":{"name":"platform"}}}`))
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &OrganizationTagCreateCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-org=my-org", "-name=platform", "-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	out := ui.OutputWriter.String()
	if !strings.Contains(out, "\"id\": \"tag-1\"") {
		t.Fatalf("expected JSON output to include id tag-1, got %q", out)
	}
}
