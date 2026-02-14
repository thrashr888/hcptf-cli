package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/config"
	"github.com/mitchellh/cli"
)

func TestAssessmentResultFindChangedAttributes(t *testing.T) {
	cmd := AssessmentResultReadCommand{}
	before := map[string]interface{}{
		"name":  "old",
		"count": float64(1),
	}
	after := map[string]interface{}{
		"name": "new",
		"size": "large",
	}

	changed := cmd.findChangedAttributes(before, after)
	if len(changed) != 3 {
		t.Fatalf("expected 3 changed attributes, got %d (%v)", len(changed), changed)
	}
}

func TestAssessmentResultValuesEqual(t *testing.T) {
	cmd := AssessmentResultReadCommand{}

	if !cmd.valuesEqual("a", "a") {
		t.Fatal("expected string values to be equal")
	}
	if cmd.valuesEqual("a", "b") {
		t.Fatal("expected string values to differ")
	}
}

func TestAssessmentResultFormatValue(t *testing.T) {
	cmd := AssessmentResultReadCommand{}

	if got := cmd.formatValue(nil); got != "<nil>" {
		t.Fatalf("expected <nil>, got %q", got)
	}
	if got := cmd.formatValue("value"); got != "\"value\"" {
		t.Fatalf("unexpected formatted string %q", got)
	}
	if got := cmd.formatValue(strings.Repeat("a", 150)); !strings.HasSuffix(got, "...") {
		t.Fatalf("expected truncated string, got %q", got)
	}
}

func TestAssessmentResultShowDriftDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/plan":
			plan := TerraformPlan{
				ResourceDrift: []struct {
					Address  string `json:"address"`
					Mode     string `json:"mode"`
					Type     string `json:"type"`
					Name     string `json:"name"`
					Provider string `json:"provider_name"`
					Change   struct {
						Actions []string               `json:"actions"`
						Before  map[string]interface{} `json:"before"`
						After   map[string]interface{} `json:"after"`
					} `json:"change"`
				}{
					{
						Address:  "module.test.aws_instance.example",
						Type:     "aws_instance",
						Name:     "example",
						Provider: "registry.terraform.io/hashicorp/aws",
						Change: struct {
							Actions []string               `json:"actions"`
							Before  map[string]interface{} `json:"before"`
							After   map[string]interface{} `json:"after"`
						}{
							Actions: []string{"update"},
							Before:  map[string]interface{}{"name": "old", "enabled": true},
							After:   map[string]interface{}{"name": "new", "enabled": true},
						},
					},
				},
			}
			data, err := json.Marshal(plan)
			if err != nil {
				t.Fatalf("failed to marshal plan: %v", err)
			}
			_, _ = w.Write(data)
		case "/plan-empty":
			_, _ = w.Write([]byte(`{"resource_drift":[]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("not found"))
		}
	}))
	defer server.Close()

	apiClient := newAssessmentResultTestClient(t, server.URL)

	ui := cli.NewMockUi()
	cmd := &AssessmentResultReadCommand{Meta: newTestMeta(ui)}

	if err := cmd.showDriftDetails(apiClient, "/plan"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "Changed attributes:") {
		t.Fatalf("expected changed attributes output, got %q", output)
	}

	if err := cmd.showDriftDetails(apiClient, "/plan-empty"); err != nil {
		t.Fatalf("expected no error for empty plan, got %v", err)
	}

	output = ui.OutputWriter.String()
	if !strings.Contains(output, "No resource drift details available") {
		t.Fatalf("expected empty-drift output, got %q", output)
	}

	if err := cmd.showDriftDetails(apiClient, "/missing"); err == nil {
		t.Fatal("expected error for missing plan")
	}
}

func TestAssessmentResultShowDriftDetailsNoActualDrift(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/plan":
			plan := TerraformPlan{
				ResourceDrift: []struct {
					Address  string `json:"address"`
					Mode     string `json:"mode"`
					Type     string `json:"type"`
					Name     string `json:"name"`
					Provider string `json:"provider_name"`
					Change   struct {
						Actions []string               `json:"actions"`
						Before  map[string]interface{} `json:"before"`
						After   map[string]interface{} `json:"after"`
					} `json:"change"`
				}{
					{
						Address:  "mod.noop.example",
						Mode:     "managed",
						Type:     "null_resource",
						Name:     "noop",
						Provider: "registry.terraform.io/hashicorp/null",
						Change: struct {
							Actions []string               `json:"actions"`
							Before  map[string]interface{} `json:"before"`
							After   map[string]interface{} `json:"after"`
						}{
							Actions: []string{"no-op"},
							Before:  map[string]interface{}{"name": "a"},
							After:   map[string]interface{}{"name": "a"},
						},
					},
				},
			}
			data, err := json.Marshal(plan)
			if err != nil {
				t.Fatalf("failed to marshal plan: %v", err)
			}
			_, _ = w.Write(data)
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("not found"))
		}
	}))
	defer server.Close()

	apiClient := newAssessmentResultTestClient(t, server.URL)

	ui := cli.NewMockUi()
	cmd := &AssessmentResultReadCommand{Meta: newTestMeta(ui)}
	if err := cmd.showDriftDetails(apiClient, "/plan"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "No actual drift detected in resources (all resources are in sync).") {
		t.Fatalf("expected no actual drift output, got %q", output)
	}
}

func TestAssessmentResultShowDriftDetailsTooManyChangedAttributes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/plan" {
			plan := TerraformPlan{
				ResourceDrift: []struct {
					Address  string `json:"address"`
					Mode     string `json:"mode"`
					Type     string `json:"type"`
					Name     string `json:"name"`
					Provider string `json:"provider_name"`
					Change   struct {
						Actions []string               `json:"actions"`
						Before  map[string]interface{} `json:"before"`
						After   map[string]interface{} `json:"after"`
					} `json:"change"`
				}{
					{
						Address:  "module.test.aws_instance.example",
						Mode:     "managed",
						Type:     "aws_instance",
						Name:     "example",
						Provider: "registry.terraform.io/hashicorp/aws",
						Change: struct {
							Actions []string               `json:"actions"`
							Before  map[string]interface{} `json:"before"`
							After   map[string]interface{} `json:"after"`
						}{
							Actions: []string{"update"},
							Before:  map[string]interface{}{"a": "old", "b": "old", "c": "old", "d": "old", "e": "old", "f": "old", "g": "old", "h": "old", "i": "old", "j": "old", "k": "old"},
							After:   map[string]interface{}{"a": "new", "b": "new", "c": "new", "d": "new", "e": "new", "f": "new", "g": "new", "h": "new", "i": "new", "j": "new", "k": "new"},
						},
					},
				},
			}
			data, err := json.Marshal(plan)
			if err != nil {
				t.Fatalf("failed to marshal plan: %v", err)
			}
			_, _ = w.Write(data)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	}))
	defer server.Close()

	apiClient := newAssessmentResultTestClient(t, server.URL)

	ui := cli.NewMockUi()
	cmd := &AssessmentResultReadCommand{Meta: newTestMeta(ui)}
	if err := cmd.showDriftDetails(apiClient, "/plan"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "Changed attributes: 11 attributes changed (too many to display)") {
		t.Fatalf("expected too many attributes output, got %q", output)
	}
}

func TestAssessmentResultFormatValueAdditionalTypes(t *testing.T) {
	cmd := &AssessmentResultReadCommand{}

	if got := cmd.formatValue(true); got != "true" {
		t.Fatalf("expected true, got %q", got)
	}
	if got := cmd.formatValue(3.14); got != "3.14" {
		t.Fatalf("expected 3.14, got %q", got)
	}
	if got := cmd.formatValue(map[string]interface{}{}); got != "{}" {
		t.Fatalf("expected empty map output, got %q", got)
	}
	if got := cmd.formatValue([]interface{}{}); got != "[]" {
		t.Fatalf("expected empty array output, got %q", got)
	}
}

func newAssessmentResultTestClient(t *testing.T, serverURL string) *client.Client {
	t.Helper()

	t.Setenv("HCPTF_ADDRESS", serverURL)
	parsed, err := url.Parse(serverURL)
	if err != nil {
		t.Fatalf("failed to parse server URL: %v", err)
	}

	cfg := &config.Config{
		Credentials: map[string]*config.Credential{
			parsed.Hostname(): {
				Hostname: parsed.Hostname(),
				Token:    "test-token",
			},
		},
	}

	apiClient, err := client.New(cfg)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	return apiClient
}
