package command

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

// --- Dry-run tests ---

func TestWorkspaceDeleteDryRun(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceDeleteService{}
	cmd := newWorkspaceDeleteCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=prod", "-dry-run"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; errors: %s", code, ui.ErrorWriter.String())
	}

	// API should NOT have been called
	if svc.lastOrg != "" || svc.lastName != "" {
		t.Fatalf("expected no API call during dry-run, but delete was called with org=%q name=%q", svc.lastOrg, svc.lastName)
	}

	// Output should be valid JSON with expected fields
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("expected valid JSON, got %q: %v", output, err)
	}
	if result["action"] != "delete" {
		t.Errorf("expected action=delete, got %v", result["action"])
	}
	if result["resource"] != "workspace" {
		t.Errorf("expected resource=workspace, got %v", result["resource"])
	}
	if result["organization"] != "my-org" {
		t.Errorf("expected organization=my-org, got %v", result["organization"])
	}
	if result["name"] != "prod" {
		t.Errorf("expected name=prod, got %v", result["name"])
	}
}

func TestVariableCreateDryRun(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableCreateService{}
	cmd := newVariableCreateCommand(ui, ws, vars)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org", "-workspace=prod",
			"-key=region", "-value=us-east-1",
			"-dry-run",
		})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; errors: %s", code, ui.ErrorWriter.String())
	}

	// API should NOT have been called
	if vars.lastWorkspace != "" {
		t.Fatalf("expected no API call during dry-run, but create was called")
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("expected valid JSON, got %q: %v", output, err)
	}
	if result["action"] != "create" {
		t.Errorf("expected action=create, got %v", result["action"])
	}
	if result["resource"] != "variable" {
		t.Errorf("expected resource=variable, got %v", result["resource"])
	}
	if result["workspace_id"] != "ws-1" {
		t.Errorf("expected workspace_id=ws-1, got %v", result["workspace_id"])
	}
}

func TestRunApplyDryRun(t *testing.T) {
	ui := cli.NewMockUi()
	runs := &mockRunApplyService{}
	cmd := newRunApplyCommand(ui, runs)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=run-abc123", "-dry-run"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; errors: %s", code, ui.ErrorWriter.String())
	}

	if runs.lastRun != "" {
		t.Fatalf("expected no API call during dry-run, but apply was called")
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("expected valid JSON, got %q: %v", output, err)
	}
	if result["action"] != "apply" {
		t.Errorf("expected action=apply, got %v", result["action"])
	}
	if result["id"] != "run-abc123" {
		t.Errorf("expected id=run-abc123, got %v", result["id"])
	}
}

// --- Input validation tests ---

func TestWorkspaceDeleteRejectsPathTraversal(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceDeleteService{}
	cmd := newWorkspaceDeleteCommand(ui, svc)

	code := cmd.Run([]string{"-organization=../../../etc", "-name=test", "-force"})
	if code != 1 {
		t.Fatalf("expected exit 1 for path traversal, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "path traversal") {
		t.Fatalf("expected path traversal error, got %q", ui.ErrorWriter.String())
	}
	if svc.lastOrg != "" {
		t.Fatal("expected no API call for invalid input")
	}
}

func TestWorkspaceDeleteRejectsQueryChars(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceDeleteService{}
	cmd := newWorkspaceDeleteCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=ws?admin=true", "-force"})
	if code != 1 {
		t.Fatalf("expected exit 1 for query chars, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "query character") {
		t.Fatalf("expected query character error, got %q", ui.ErrorWriter.String())
	}
}

func TestWorkspaceDeleteRejectsPathSeparator(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceDeleteService{}
	cmd := newWorkspaceDeleteCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=ws/admin", "-force"})
	if code != 1 {
		t.Fatalf("expected exit 1 for path separator, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "path separator") {
		t.Fatalf("expected path separator error, got %q", ui.ErrorWriter.String())
	}
}

func TestVariableCreateRejectsURLEncodedKey(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableCreateService{}
	cmd := newVariableCreateCommand(ui, ws, vars)

	code := cmd.Run([]string{
		"-organization=my-org", "-workspace=prod",
		"-key=%2fadmin", "-value=test",
	})
	if code != 1 {
		t.Fatalf("expected exit 1 for URL-encoded key, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "URL-encoded") {
		t.Fatalf("expected URL-encoded error, got %q", ui.ErrorWriter.String())
	}
}

func TestVariableCreateRejectsControlCharsInDescription(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableCreateService{}
	cmd := newVariableCreateCommand(ui, ws, vars)

	code := cmd.Run([]string{
		"-organization=my-org", "-workspace=prod",
		"-key=mykey", "-value=myval",
		"-description=" + string([]byte{'h', 'i', 0x00, 'x'}),
	})
	if code != 1 {
		t.Fatalf("expected exit 1 for control chars, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "control characters") {
		t.Fatalf("expected control char error, got %q", ui.ErrorWriter.String())
	}
}

// --- Schema tests ---

func TestSchemaReturnsValidJSON(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &SchemaCommand{Meta: newTestMeta(ui)}

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"workspace", "create"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; errors: %s", code, ui.ErrorWriter.String())
	}

	var schema schemaOutput
	if err := json.Unmarshal([]byte(output), &schema); err != nil {
		t.Fatalf("expected valid JSON schema, got %q: %v", output, err)
	}
	if schema.Command != "workspace create" {
		t.Errorf("expected command='workspace create', got %q", schema.Command)
	}
	if len(schema.Flags) == 0 {
		t.Fatal("expected at least one flag in schema")
	}

	// Check known flags are present
	flagNames := make(map[string]schemaFlag)
	for _, f := range schema.Flags {
		flagNames[f.Name] = f
	}
	if _, ok := flagNames["organization"]; !ok {
		t.Error("expected 'organization' flag in schema")
	}
	if f, ok := flagNames["organization"]; ok && !f.Required {
		t.Error("expected 'organization' to be required")
	}
	if _, ok := flagNames["name"]; !ok {
		t.Error("expected 'name' flag in schema")
	}
	if f, ok := flagNames["org"]; ok {
		if len(f.Aliases) == 0 || f.Aliases[0] != "organization" {
			t.Errorf("expected 'org' to alias 'organization', got aliases=%v", f.Aliases)
		}
	}
}

func TestSchemaVariableCreate(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &SchemaCommand{Meta: newTestMeta(ui)}

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"variable", "create"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	var schema schemaOutput
	if err := json.Unmarshal([]byte(output), &schema); err != nil {
		t.Fatalf("expected valid JSON: %v", err)
	}

	flagNames := make(map[string]bool)
	for _, f := range schema.Flags {
		flagNames[f.Name] = true
	}
	for _, expected := range []string{"organization", "workspace", "key", "value", "category"} {
		if !flagNames[expected] {
			t.Errorf("expected flag %q in schema", expected)
		}
	}
}

// --- Fields filtering tests ---

func TestVariableCreateFieldsFiltering(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableCreateService{
		response: &tfe.Variable{
			ID:          "var-1",
			Key:         "region",
			Value:       "us-east-1",
			Category:    tfe.CategoryTerraform,
			Sensitive:   false,
			HCL:         false,
			Description: "AWS region",
		},
	}
	cmd := newVariableCreateCommand(ui, ws, vars)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org", "-workspace=prod",
			"-key=region", "-value=us-east-1",
			"-output=json", "-fields=ID,Key",
		})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; errors: %s", code, ui.ErrorWriter.String())
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("expected valid JSON, got %q: %v", output, err)
	}
	if data["ID"] != "var-1" {
		t.Errorf("expected ID=var-1, got %v", data["ID"])
	}
	if data["Key"] != "region" {
		t.Errorf("expected Key=region, got %v", data["Key"])
	}
	// Fields not requested should be absent
	if _, ok := data["Value"]; ok {
		t.Error("expected Value to be filtered out")
	}
	if _, ok := data["Category"]; ok {
		t.Error("expected Category to be filtered out")
	}
	if _, ok := data["Description"]; ok {
		t.Error("expected Description to be filtered out")
	}
}

// --- JSON input tests ---

func TestVariableCreateJSONInput(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableCreateService{
		response: &tfe.Variable{
			ID:       "var-1",
			Key:      "region",
			Value:    "us-west-2",
			Category: tfe.CategoryTerraform,
		},
	}
	cmd := newVariableCreateCommand(ui, ws, vars)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org", "-workspace=prod",
			"-output=json",
			`-json-input={"key":"region","value":"us-west-2","category":"terraform"}`,
		})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; errors: %s", code, ui.ErrorWriter.String())
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("expected valid JSON, got %q: %v", output, err)
	}
	if data["ID"] != "var-1" {
		t.Errorf("expected ID=var-1, got %v", data["ID"])
	}
}

func TestVariableCreateJSONInputFromFile(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableCreateService{
		response: &tfe.Variable{
			ID:       "var-1",
			Key:      "env",
			Value:    "production",
			Category: tfe.CategoryTerraform,
		},
	}
	cmd := newVariableCreateCommand(ui, ws, vars)

	// Write JSON to a temp file
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "input.json")
	if err := os.WriteFile(jsonFile, []byte(`{"key":"env","value":"production","category":"terraform"}`), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org", "-workspace=prod",
			"-output=json",
			"-json-input=@" + jsonFile,
		})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; errors: %s", code, ui.ErrorWriter.String())
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("expected valid JSON, got %q: %v", output, err)
	}
	if data["ID"] != "var-1" {
		t.Errorf("expected ID=var-1, got %v", data["ID"])
	}
}

func TestVariableCreateJSONInputInvalidJSON(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableCreateService{}
	cmd := newVariableCreateCommand(ui, ws, vars)

	code := cmd.Run([]string{
		"-organization=my-org", "-workspace=prod",
		"-json-input={not valid json}",
	})
	if code != 1 {
		t.Fatalf("expected exit 1 for invalid JSON, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "JSON input") {
		t.Fatalf("expected JSON parse error, got %q", ui.ErrorWriter.String())
	}
}

func TestVariableCreateJSONInputSkipsRequiredFlagValidation(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableCreateService{
		response: &tfe.Variable{ID: "var-1", Key: "k", Value: "v", Category: tfe.CategoryTerraform},
	}
	cmd := newVariableCreateCommand(ui, ws, vars)

	// With json-input, -key and -value flags are not required
	_, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org", "-workspace=prod",
			"-output=json",
			`-json-input={"key":"k","value":"v","category":"terraform"}`,
		})
	})
	if code != 0 {
		t.Fatalf("expected exit 0 (json-input should skip flag requirements), got %d; errors: %s", code, ui.ErrorWriter.String())
	}
}

// --- Dry-run combined with JSON input ---

func TestVariableCreateDryRunWithJSONInput(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	vars := &mockVariableCreateService{}
	cmd := newVariableCreateCommand(ui, ws, vars)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org", "-workspace=prod",
			"-dry-run",
			`-json-input={"key":"region","value":"eu-west-1","category":"terraform"}`,
		})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; errors: %s", code, ui.ErrorWriter.String())
	}
	if vars.lastWorkspace != "" {
		t.Fatal("expected no API call during dry-run with json-input")
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("expected valid JSON, got %q: %v", output, err)
	}
	if result["action"] != "create" {
		t.Errorf("expected action=create, got %v", result["action"])
	}
}
