package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestChangeRequestUpdateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ChangeRequestUpdateCommand{
		Meta: newTestMeta(ui),
	}

	// Test missing id
	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing id, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error, got %q", ui.ErrorWriter.String())
	}

	// Test missing archive flag
	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-id=wscr-123"}); code != 1 {
		t.Fatalf("expected exit 1 missing archive, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-archive") {
		t.Fatalf("expected archive error, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedID      string
		expectedArchive bool
		expectedFormat  string
	}{
		{
			name:            "id and archive with default format",
			args:            []string{"-id=wscr-abc123", "-archive"},
			expectedID:      "wscr-abc123",
			expectedArchive: true,
			expectedFormat:  "table",
		},
		{
			name:            "id and archive with table format",
			args:            []string{"-id=wscr-xyz789", "-archive", "-output=table"},
			expectedID:      "wscr-xyz789",
			expectedArchive: true,
			expectedFormat:  "table",
		},
		{
			name:            "id and archive with json format",
			args:            []string{"-id=wscr-test456", "-archive", "-output=json"},
			expectedID:      "wscr-test456",
			expectedArchive: true,
			expectedFormat:  "json",
		},
		{
			name:            "archive with different id format",
			args:            []string{"-id=wscr-prod999", "-archive"},
			expectedID:      "wscr-prod999",
			expectedArchive: true,
			expectedFormat:  "table",
		},
		{
			name:            "explicit archive true with json",
			args:            []string{"-id=wscr-staging", "-archive=true", "-output=json"},
			expectedID:      "wscr-staging",
			expectedArchive: true,
			expectedFormat:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ChangeRequestUpdateCommand{}

			flags := cmd.Meta.FlagSet("changerequest update")
			flags.StringVar(&cmd.id, "id", "", "Change request ID (required)")
			flags.BoolVar(&cmd.archive, "archive", false, "Archive the change request")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the archive flag was set correctly
			if cmd.archive != tt.expectedArchive {
				t.Errorf("expected archive %v, got %v", tt.expectedArchive, cmd.archive)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}

func TestChangeRequestUpdateRunSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/api/v2/ping" || r.URL.RequestURI() == "/api/v2/ping?" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}

		if r.URL.Path != "/api/v2/change-requests/cr-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		_, _ = w.Write([]byte(`{"data":{"id":"cr-1","type":"change-requests","attributes":{"subject":"Fix","message":"Please update","archived-by":"user-99","archived-at":"2024-01-03T00:00:00Z","created-at":"2024-01-01T00:00:00Z","updated-at":"2024-01-03T00:00:00Z"},"relationships":{"workspace":{"data":{"id":"ws-123","type":"workspaces"}}}}}`))
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestUpdateCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=cr-1", "-archive"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, output=%q, err=%q", code, ui.OutputWriter.String(), ui.ErrorWriter.String())
	}

	if !strings.Contains(ui.OutputWriter.String(), "Change request 'cr-1' archived successfully") {
		t.Fatalf("expected success message, got %q", ui.OutputWriter.String())
	}
	if !strings.Contains(ui.OutputWriter.String(), "ArchivedBy") {
		t.Fatalf("expected archived metadata in output, got %q", ui.OutputWriter.String())
	}
}

func TestChangeRequestUpdateRunJSONOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/api/v2/ping" || r.URL.RequestURI() == "/api/v2/ping?" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}

		if r.URL.Path != "/api/v2/change-requests/cr-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		_, _ = w.Write([]byte(`{"data":{"id":"cr-1","type":"change-requests","attributes":{"subject":"Fix","message":"Please update","archived-by":"user-99","archived-at":"2024-01-03T00:00:00Z","created-at":"2024-01-01T00:00:00Z","updated-at":"2024-01-03T00:00:00Z"},"relationships":{"workspace":{"data":{"id":"ws-123","type":"workspaces"}}}}}`))
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestUpdateCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=cr-1", "-archive", "-output=json"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, output=%q, err=%q", code, ui.OutputWriter.String(), ui.ErrorWriter.String())
	}

	output := strings.TrimSpace(ui.OutputWriter.String())
	start := strings.Index(output, "{")
	end := strings.LastIndex(output, "}")
	if start == -1 || end == -1 || end <= start {
		t.Fatalf("expected JSON output in response, got: %q", output)
	}
	jsonOutput := output[start : end+1]

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOutput), &data); err != nil {
		t.Fatalf("failed to decode json output: %v, output: %q", err, jsonOutput)
	}

	if data["ID"] != "cr-1" {
		t.Fatalf("expected ID cr-1, got %v", data["ID"])
	}
	if data["Status"] != "Archived" {
		t.Fatalf("expected Status Archived, got %v", data["Status"])
	}
}

func TestChangeRequestUpdateRunNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/api/v2/ping" || r.URL.RequestURI() == "/api/v2/ping?" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}

		if r.URL.Path != "/api/v2/change-requests/cr-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors":[{"status":"404"}]}`))
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestUpdateCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=cr-1", "-archive"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "API request failed with status 404") {
		t.Fatalf("expected API failure output, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestUpdateRunInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() == "/api/v2/ping" || r.URL.RequestURI() == "/api/v2/ping?" {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.Path != "/api/v2/change-requests/cr-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		_, _ = w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &ChangeRequestUpdateCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=cr-1", "-archive"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "Error parsing response") {
		t.Fatalf("expected parse error output, got %q", ui.ErrorWriter.String())
	}
}
