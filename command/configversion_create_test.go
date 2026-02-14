package command

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

type mockConfigVersionCreateService struct {
	response      *tfe.ConfigurationVersion
	err           error
	lastWorkspace string
	lastOptions   tfe.ConfigurationVersionCreateOptions
}

func (m *mockConfigVersionCreateService) Create(_ context.Context, workspaceID string, options tfe.ConfigurationVersionCreateOptions) (*tfe.ConfigurationVersion, error) {
	m.lastWorkspace = workspaceID
	m.lastOptions = options
	return m.response, m.err
}

func newConfigVersionCreateCommand(ui cli.Ui, ws workspaceReader, cv configVersionCreator) *ConfigVersionCreateCommand {
	return &ConfigVersionCreateCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: ws,
		configVerSvc: cv,
	}
}

func TestConfigVersionCreateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newConfigVersionCreateCommand(ui, &mockWorkspaceReader{}, &mockConfigVersionCreateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing workspace")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-workspace") {
		t.Fatalf("expected workspace error")
	}
}

func TestConfigVersionCreateRequiresEmptyOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newConfigVersionCreateCommand(ui, &mockWorkspaceReader{}, &mockConfigVersionCreateService{})

	code := cmd.Run([]string{"-organization=", "-workspace=prod"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestConfigVersionCreateRequiresEmptyWorkspace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newConfigVersionCreateCommand(ui, &mockWorkspaceReader{}, &mockConfigVersionCreateService{})

	code := cmd.Run([]string{"-organization=my-org", "-workspace="})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace") {
		t.Fatalf("expected workspace error, got %q", out)
	}
}

func TestConfigVersionCreateHandlesWorkspaceError(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{err: errors.New("workspace not found")}
	cv := &mockConfigVersionCreateService{}
	cmd := newConfigVersionCreateCommand(ui, ws, cv)

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if ws.lastOrg != "my-org" || ws.lastName != "prod" {
		t.Fatalf("expected workspace read called with correct params")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "workspace not found") {
		t.Fatalf("expected error output")
	}
}

func TestConfigVersionCreateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	cv := &mockConfigVersionCreateService{err: errors.New("boom")}
	cmd := newConfigVersionCreateCommand(ui, ws, cv)

	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if cv.lastWorkspace != "ws-1" {
		t.Fatalf("expected workspace ID passed")
	}
	if cv.lastOptions.AutoQueueRuns == nil || !*cv.lastOptions.AutoQueueRuns {
		t.Fatalf("expected auto-queue-runs true by default")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestConfigVersionCreateWithFlags(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	cv := &mockConfigVersionCreateService{response: &tfe.ConfigurationVersion{
		ID:            "cv-1",
		Status:        tfe.ConfigurationPending,
		Source:        tfe.ConfigurationSourceAPI,
		Speculative:   false,
		Provisional:   false,
		AutoQueueRuns: true,
		UploadURL:     "https://upload.example.com",
	}}
	cmd := newConfigVersionCreateCommand(ui, ws, cv)

	code := cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-speculative", "-provisional", "-auto-queue-runs=false"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if cv.lastOptions.Speculative == nil || !*cv.lastOptions.Speculative {
		t.Fatalf("expected speculative true")
	}
	if cv.lastOptions.Provisional == nil || !*cv.lastOptions.Provisional {
		t.Fatalf("expected provisional true")
	}
	if cv.lastOptions.AutoQueueRuns == nil || *cv.lastOptions.AutoQueueRuns {
		t.Fatalf("expected auto-queue-runs false")
	}
}

func TestConfigVersionCreateOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	ws := &mockWorkspaceReader{workspace: &tfe.Workspace{ID: "ws-1"}}
	cv := &mockConfigVersionCreateService{response: &tfe.ConfigurationVersion{
		ID:            "cv-123",
		Status:        tfe.ConfigurationPending,
		Source:        tfe.ConfigurationSourceAPI,
		Speculative:   true,
		Provisional:   false,
		AutoQueueRuns: false,
		UploadURL:     "https://upload.example.com/path",
	}}
	cmd := newConfigVersionCreateCommand(ui, ws, cv)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-speculative", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0")
	}

	if cv.lastOptions.Speculative == nil || !*cv.lastOptions.Speculative {
		t.Fatalf("expected speculative true")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ID"] != "cv-123" {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestConfigVersionCreateHelp(t *testing.T) {
	cmd := &ConfigVersionCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf configversion create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-workspace") {
		t.Error("Help should mention -workspace flag")
	}
	if !strings.Contains(help, "-auto-queue-runs") {
		t.Error("Help should mention -auto-queue-runs flag")
	}
	if !strings.Contains(help, "-speculative") {
		t.Error("Help should mention -speculative flag")
	}
	if !strings.Contains(help, "-provisional") {
		t.Error("Help should mention -provisional flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestConfigVersionCreateSynopsis(t *testing.T) {
	cmd := &ConfigVersionCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new configuration version" {
		t.Errorf("expected 'Create a new configuration version', got %q", synopsis)
	}
}

func TestConfigVersionCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name                  string
		args                  []string
		expectedOrg           string
		expectedWorkspace     string
		expectedAutoQueueRuns bool
		expectedSpeculative   bool
		expectedProvisional   bool
		expectedFormat        string
	}{
		{
			name:                  "org and workspace only, defaults",
			args:                  []string{"-organization=test-org", "-workspace=test-ws"},
			expectedOrg:           "test-org",
			expectedWorkspace:     "test-ws",
			expectedAutoQueueRuns: true,
			expectedSpeculative:   false,
			expectedProvisional:   false,
			expectedFormat:        "table",
		},
		{
			name:                  "org alias and workspace",
			args:                  []string{"-org=my-org", "-workspace=prod"},
			expectedOrg:           "my-org",
			expectedWorkspace:     "prod",
			expectedAutoQueueRuns: true,
			expectedSpeculative:   false,
			expectedProvisional:   false,
			expectedFormat:        "table",
		},
		{
			name:                  "with speculative flag",
			args:                  []string{"-organization=test-org", "-workspace=test-ws", "-speculative"},
			expectedOrg:           "test-org",
			expectedWorkspace:     "test-ws",
			expectedAutoQueueRuns: true,
			expectedSpeculative:   true,
			expectedProvisional:   false,
			expectedFormat:        "table",
		},
		{
			name:                  "with provisional flag",
			args:                  []string{"-organization=test-org", "-workspace=test-ws", "-provisional"},
			expectedOrg:           "test-org",
			expectedWorkspace:     "test-ws",
			expectedAutoQueueRuns: true,
			expectedSpeculative:   false,
			expectedProvisional:   true,
			expectedFormat:        "table",
		},
		{
			name:                  "with auto-queue-runs false",
			args:                  []string{"-organization=test-org", "-workspace=test-ws", "-auto-queue-runs=false"},
			expectedOrg:           "test-org",
			expectedWorkspace:     "test-ws",
			expectedAutoQueueRuns: false,
			expectedSpeculative:   false,
			expectedProvisional:   false,
			expectedFormat:        "table",
		},
		{
			name:                  "all flags with json output",
			args:                  []string{"-organization=prod-org", "-workspace=prod-ws", "-speculative", "-provisional", "-auto-queue-runs=false", "-output=json"},
			expectedOrg:           "prod-org",
			expectedWorkspace:     "prod-ws",
			expectedAutoQueueRuns: false,
			expectedSpeculative:   true,
			expectedProvisional:   true,
			expectedFormat:        "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ConfigVersionCreateCommand{}

			flags := cmd.Meta.FlagSet("configversion create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.workspace, "workspace", "", "Workspace name (required)")
			flags.BoolVar(&cmd.autoQueueRuns, "auto-queue-runs", true, "Automatically queue runs when uploaded")
			flags.BoolVar(&cmd.speculative, "speculative", false, "Create a speculative configuration version")
			flags.BoolVar(&cmd.provisional, "provisional", false, "Create a provisional configuration version")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the workspace was set correctly
			if cmd.workspace != tt.expectedWorkspace {
				t.Errorf("expected workspace %q, got %q", tt.expectedWorkspace, cmd.workspace)
			}

			// Verify the auto-queue-runs was set correctly
			if cmd.autoQueueRuns != tt.expectedAutoQueueRuns {
				t.Errorf("expected autoQueueRuns %v, got %v", tt.expectedAutoQueueRuns, cmd.autoQueueRuns)
			}

			// Verify the speculative was set correctly
			if cmd.speculative != tt.expectedSpeculative {
				t.Errorf("expected speculative %v, got %v", tt.expectedSpeculative, cmd.speculative)
			}

			// Verify the provisional was set correctly
			if cmd.provisional != tt.expectedProvisional {
				t.Errorf("expected provisional %v, got %v", tt.expectedProvisional, cmd.provisional)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
