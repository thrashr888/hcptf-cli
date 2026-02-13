package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestChangeRequestReadRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ChangeRequestReadCommand{
		Meta: newTestMeta(ui),
	}

	// Test missing id
	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing id, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "id with default format",
			args:           []string{"-id=wscr-abc123"},
			expectedID:     "wscr-abc123",
			expectedFormat: "table",
		},
		{
			name:           "id with table format",
			args:           []string{"-id=wscr-xyz789", "-output=table"},
			expectedID:     "wscr-xyz789",
			expectedFormat: "table",
		},
		{
			name:           "id with json format",
			args:           []string{"-id=wscr-test456", "-output=json"},
			expectedID:     "wscr-test456",
			expectedFormat: "json",
		},
		{
			name:           "different id format",
			args:           []string{"-id=wscr-prod999"},
			expectedID:     "wscr-prod999",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ChangeRequestReadCommand{}

			flags := cmd.Meta.FlagSet("changerequest read")
			flags.StringVar(&cmd.id, "id", "", "Change request ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}

func TestChangeRequestReadRunSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/api/v2/ping", "/api/v2/ping?":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/change-requests/cr-1":
			_, _ = w.Write([]byte(`{"data":{"id":"cr-1","type":"change-requests","attributes":{"subject":"Fix","message":"Please update","archived-by":null,"archived-at":null,"created-at":"2024-01-01T00:00:00Z","updated-at":"2024-01-02T00:00:00Z"},"relationships":{"workspace":{"data":{"id":"ws-123","type":"workspaces"}}}}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.RequestURI())
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestReadCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=cr-1"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, output=%q, err=%q", code, ui.OutputWriter.String(), ui.ErrorWriter.String())
	}

	out := strings.TrimSpace(ui.OutputWriter.String())
	if !strings.Contains(out, "ID") || !strings.Contains(out, "cr-1") {
		t.Fatalf("expected ID output, got %q", out)
	}
}

func TestChangeRequestReadRunNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/api/v2/ping", "/api/v2/ping?":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/change-requests/cr-1":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"errors":[{"status":"404"}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.RequestURI())
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestReadCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=cr-1"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "API request failed with status 404") {
		t.Fatalf("expected 404 output, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestReadRunJSONOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/api/v2/ping", "/api/v2/ping?":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/change-requests/cr-1":
			_, _ = w.Write([]byte(`{"data":{"id":"cr-1","type":"change-requests","attributes":{"subject":"Fix","message":"Please update","archived-by":null,"archived-at":null,"created-at":"2024-01-01T00:00:00Z","updated-at":"2024-01-02T00:00:00Z"},"relationships":{"workspace":{"data":{"id":"ws-123","type":"workspaces"}}}}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.RequestURI())
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestReadCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=cr-1", "-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, output=%q, err=%q", code, ui.OutputWriter.String(), ui.ErrorWriter.String())
	}

	out := strings.TrimSpace(ui.OutputWriter.String())
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("failed to decode json output: %v, output: %q", err, out)
	}
	if data["ID"] != "cr-1" {
		t.Fatalf("expected ID in json output, got %v", data["ID"])
	}
}

func TestChangeRequestReadRunArchived(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/api/v2/ping", "/api/v2/ping?":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/api/v2/change-requests/cr-archived":
			_, _ = w.Write([]byte(`{"data":{"id":"cr-archived","type":"change-requests","attributes":{"subject":"Fix","message":"Please update","archived-by":"user-1","archived-at":"2024-01-05T00:00:00Z","created-at":"2024-01-01T00:00:00Z","updated-at":"2024-01-02T00:00:00Z"},"relationships":{"workspace":{"data":{"id":"ws-123","type":"workspaces"}}}}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.RequestURI())
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestReadCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=cr-archived"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, output=%q, err=%q", code, ui.OutputWriter.String(), ui.ErrorWriter.String())
	}

	if !strings.Contains(ui.OutputWriter.String(), "ArchivedBy") {
		t.Fatalf("expected archived fields in output, got %q", ui.OutputWriter.String())
	}
}
