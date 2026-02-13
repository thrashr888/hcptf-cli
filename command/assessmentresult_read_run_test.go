package command

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAssessmentResultReadRunShowsDriftDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/assessment-results/ar-1":
			response := `{
				"data": {
					"id": "ar-1",
					"type": "assessment-results",
					"attributes": {
						"drifted": true,
						"succeeded": true,
						"created-at": "2024-01-01T00:00:00Z",
						"error-msg": null
					},
					"links": {
						"json-output": "/plan",
						"log-output": "/log"
					}
				}
			}`
			_, _ = io.WriteString(w, response)
		case "/plan":
			_, _ = io.WriteString(w, `{"resource_drift":[{"address":"module.test.aws_instance.example","mode":"managed","type":"aws_instance","name":"example","provider_name":"registry.terraform.io/hashicorp/aws","change":{"actions":["update"],"before":{"name":"old","enabled":true},"after":{"name":"new","enabled":false}}}]}`)
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, "not found")
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	meta := Meta{Ui: ui}
	cmd := &AssessmentResultReadCommand{Meta: meta}

	t.Setenv("TFE_TOKEN", "test-token")
	t.Setenv("HCPTF_ADDRESS", server.URL)

	code := cmd.Run([]string{"-id=ar-1"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "DRIFT DETAILS") {
		t.Fatalf("expected drift details output, got: %q", output)
	}
	if !strings.Contains(output, "INTERPRETATION") {
		t.Fatalf("expected interpretation section, got: %q", output)
	}
}

func TestAssessmentResultReadRunNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"errors":[{"status":"404"}]}`)
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	meta := Meta{Ui: ui}
	cmd := &AssessmentResultReadCommand{Meta: meta}

	t.Setenv("TFE_TOKEN", "test-token")
	t.Setenv("HCPTF_ADDRESS", server.URL)

	code := cmd.Run([]string{"-id=missing"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if got := ui.ErrorWriter.String(); !strings.Contains(got, "API request failed with status 404") {
		t.Fatalf("expected 404 error output, got: %q", got)
	}
}

func TestAssessmentResultReadRunRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AssessmentResultReadCommand{Meta: Meta{Ui: ui}}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if got := ui.ErrorWriter.String(); !strings.Contains(got, "-id flag is required") {
		t.Fatalf("expected id required error, got: %q", got)
	}
}

func TestAssessmentResultReadRunForbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"errors":[{"status":"403"}]}`)
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	meta := Meta{Ui: ui}
	cmd := &AssessmentResultReadCommand{Meta: meta}

	t.Setenv("TFE_TOKEN", "test-token")
	t.Setenv("HCPTF_ADDRESS", server.URL)

	code := cmd.Run([]string{"-id=ar-forbidden"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	output := ui.ErrorWriter.String()
	if !strings.Contains(output, "Note: You may not have permission to view this assessment result.") {
		t.Fatalf("expected permission note, got %q", output)
	}
	if !strings.Contains(output, "You need at least read access") {
		t.Fatalf("expected permission message, got %q", output)
	}
}

func TestAssessmentResultReadRunNoDriftOutputsFetchHints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/assessment-results/ar-1":
			_, _ = io.WriteString(w, `{
				"data": {
					"id": "ar-1",
					"type": "assessment-results",
					"attributes": {
						"drifted": false,
						"succeeded": true,
						"created-at": "2024-01-01T00:00:00Z",
						"error-msg": null
					},
					"links": {
						"json-output": "/plan",
						"log-output": "/log"
					}
				}
			}`)
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, "not found")
		}
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	apiClient := newAssessmentResultTestClient(t, server.URL)
	cmd := &AssessmentResultReadCommand{Meta: Meta{Ui: ui, client: apiClient}}

	code := cmd.Run([]string{"-id=ar-1"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "To retrieve detailed assessment outputs, use curl with your token:") {
		t.Fatalf("expected curl guidance, got %q", output)
	}
}

func TestAssessmentResultReadRunSummaryOnlySkipsDriftDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/assessment-results/ar-1":
			_, _ = io.WriteString(w, `{
				"data": {
					"id": "ar-1",
					"type": "assessment-results",
					"attributes": {
						"drifted": true,
						"succeeded": true,
						"created-at": "2024-01-01T00:00:00Z",
						"error-msg": null
					},
					"links": {
						"json-output": "/plan",
						"log-output": "/log"
					}
				}
			}`)
			return
		case "/api/v2/ping":
			_, _ = io.WriteString(w, `{"ok":true}`)
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.Path)
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	cmd := &AssessmentResultReadCommand{Meta: Meta{Ui: ui}}
	t.Setenv("TFE_TOKEN", "test-token")
	t.Setenv("HCPTF_ADDRESS", server.URL)

	code := cmd.Run([]string{"-id=ar-1", "-summary-only=true"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	output := ui.OutputWriter.String()
	if strings.Contains(output, "DRIFT DETAILS") {
		t.Fatalf("summary-only should not include drift details, got %q", output)
	}
}

func TestAssessmentResultReadRunInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/assessment-results/ar-1" {
			_, _ = io.WriteString(w, `{invalid}`)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, "not found")
	}))
	defer server.Close()

	ui := cli.NewMockUi()
	cmd := &AssessmentResultReadCommand{Meta: Meta{Ui: ui}}
	t.Setenv("TFE_TOKEN", "test-token")
	t.Setenv("HCPTF_ADDRESS", server.URL)

	code := cmd.Run([]string{"-id=ar-1"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "Error parsing response") {
		t.Fatalf("expected parsing error, got %q", ui.ErrorWriter.String())
	}
}
