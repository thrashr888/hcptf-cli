package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newWorkspaceCreateCommand(ui cli.Ui, svc workspaceCreator) *WorkspaceCreateCommand {
	return &WorkspaceCreateCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: svc,
	}
}

func TestWorkspaceCreateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceCreateCommand(ui, &mockWorkspaceCreateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing name, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected name error")
	}
}

func TestWorkspaceCreateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceCreateService{err: errors.New("boom")}
	cmd := newWorkspaceCreateCommand(ui, svc)

	code := cmd.Run([]string{"-organization=my-org", "-name=prod"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if svc.lastOrg != "my-org" {
		t.Fatalf("expected organization my-org, got %s", svc.lastOrg)
	}
	if svc.lastOptions.Name == nil || *svc.lastOptions.Name != "prod" {
		t.Fatalf("expected workspace name prod")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestWorkspaceCreateOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceCreateService{response: &tfe.Workspace{ID: "ws-1", Name: "prod", TerraformVersion: "1.6.0", AutoApply: true}}
	cmd := newWorkspaceCreateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=prod", "-auto-apply", "-terraform-version=1.6.0", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if svc.lastOptions.Name == nil || *svc.lastOptions.Name != "prod" {
		t.Fatalf("expected name in options")
	}
	if svc.lastOptions.AutoApply == nil || !*svc.lastOptions.AutoApply {
		t.Fatalf("expected auto apply true")
	}
	if svc.lastOptions.TerraformVersion == nil || *svc.lastOptions.TerraformVersion != "1.6.0" {
		t.Fatalf("expected terraform version set")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["Name"] != "prod" {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestWorkspaceCreateProjectID(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceCreateService{response: &tfe.Workspace{
		ID:   "ws-1",
		Name: "prod",
		Project: &tfe.Project{
			ID:   "prj-abc123",
			Name: "MyProject",
		},
	}}
	cmd := newWorkspaceCreateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=prod", "-project-id=prj-abc123", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	if svc.lastOptions.Project == nil || svc.lastOptions.Project.ID != "prj-abc123" {
		t.Fatalf("expected project option to be set with ID prj-abc123, got %+v", svc.lastOptions.Project)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ProjectID"] != "prj-abc123" {
		t.Fatalf("expected ProjectID in output, got: %#v", data)
	}
}

func TestWorkspaceCreateExecutionMode(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceCreateService{response: &tfe.Workspace{
		ID:            "ws-1",
		Name:          "prod",
		ExecutionMode: "agent",
	}}
	cmd := newWorkspaceCreateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org",
			"-name=prod",
			"-execution-mode=agent",
			"-agent-pool-id=apool-123",
			"-output=json",
		})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	if svc.lastOptions.ExecutionMode == nil || *svc.lastOptions.ExecutionMode != "agent" {
		t.Fatalf("expected execution mode 'agent', got %v", svc.lastOptions.ExecutionMode)
	}
	if svc.lastOptions.AgentPoolID == nil || *svc.lastOptions.AgentPoolID != "apool-123" {
		t.Fatalf("expected agent pool ID 'apool-123', got %v", svc.lastOptions.AgentPoolID)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["ExecutionMode"] != "agent" {
		t.Fatalf("expected ExecutionMode 'agent' in output, got: %v", data["ExecutionMode"])
	}

	// Test invalid execution mode
	ui2 := cli.NewMockUi()
	cmd2 := newWorkspaceCreateCommand(ui2, &mockWorkspaceCreateService{})
	if code := cmd2.Run([]string{"-organization=my-org", "-name=prod", "-execution-mode=invalid"}); code != 1 {
		t.Fatalf("expected exit 1 for invalid execution mode")
	}
	if !strings.Contains(ui2.ErrorWriter.String(), "execution-mode") {
		t.Fatalf("expected execution-mode validation error, got: %s", ui2.ErrorWriter.String())
	}
}

func TestWorkspaceCreateVCSRepo(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceCreateService{response: &tfe.Workspace{
		ID:   "ws-1",
		Name: "prod",
		VCSRepo: &tfe.VCSRepo{
			Identifier: "my-org/my-repo",
		},
	}}
	cmd := newWorkspaceCreateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org",
			"-name=prod",
			"-vcs-identifier=my-org/my-repo",
			"-vcs-branch=main",
			"-vcs-oauth-token-id=ot-abc",
			"-output=json",
		})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	vcs := svc.lastOptions.VCSRepo
	if vcs == nil {
		t.Fatal("expected VCSRepo options to be set")
	}
	if vcs.Identifier == nil || *vcs.Identifier != "my-org/my-repo" {
		t.Fatalf("expected VCS identifier 'my-org/my-repo', got %v", vcs.Identifier)
	}
	if vcs.Branch == nil || *vcs.Branch != "main" {
		t.Fatalf("expected VCS branch 'main', got %v", vcs.Branch)
	}
	if vcs.OAuthTokenID == nil || *vcs.OAuthTokenID != "ot-abc" {
		t.Fatalf("expected VCS OAuth token ID 'ot-abc', got %v", vcs.OAuthTokenID)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["VCSRepo"] != "my-org/my-repo" {
		t.Fatalf("expected VCSRepo 'my-org/my-repo' in output, got: %v", data["VCSRepo"])
	}
}

func TestWorkspaceCreateTags(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceCreateService{response: &tfe.Workspace{
		ID:       "ws-1",
		Name:     "prod",
		TagNames: []string{"env:prod", "team:infra"},
	}}
	cmd := newWorkspaceCreateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org",
			"-name=prod",
			"-tags=env:prod,team:infra",
			"-output=json",
		})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	if len(svc.lastOptions.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(svc.lastOptions.Tags))
	}
	if svc.lastOptions.Tags[0].Name != "env:prod" {
		t.Fatalf("expected tag 'env:prod', got %s", svc.lastOptions.Tags[0].Name)
	}
	if svc.lastOptions.Tags[1].Name != "team:infra" {
		t.Fatalf("expected tag 'team:infra', got %s", svc.lastOptions.Tags[1].Name)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["TagNames"] != "env:prod,team:infra" {
		t.Fatalf("expected TagNames in output, got: %v", data["TagNames"])
	}
}

func TestWorkspaceCreateBoolFlags(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceCreateService{response: &tfe.Workspace{
		ID:               "ws-1",
		Name:             "prod",
		AllowDestroyPlan: false,
		QueueAllRuns:     true,
	}}
	cmd := newWorkspaceCreateCommand(ui, svc)

	code := cmd.Run([]string{
		"-organization=my-org",
		"-name=prod",
		"-allow-destroy-plan=false",
		"-queue-all-runs=true",
		"-speculative-enabled=true",
		"-global-remote-state=false",
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	opts := svc.lastOptions
	if opts.AllowDestroyPlan == nil || *opts.AllowDestroyPlan {
		t.Fatalf("expected AllowDestroyPlan false")
	}
	if opts.QueueAllRuns == nil || !*opts.QueueAllRuns {
		t.Fatalf("expected QueueAllRuns true")
	}
	if opts.SpeculativeEnabled == nil || !*opts.SpeculativeEnabled {
		t.Fatalf("expected SpeculativeEnabled true")
	}
	if opts.GlobalRemoteState == nil || *opts.GlobalRemoteState {
		t.Fatalf("expected GlobalRemoteState false")
	}
}

func TestWorkspaceCreateWorkingDirectory(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceCreateService{response: &tfe.Workspace{
		ID:               "ws-1",
		Name:             "prod",
		WorkingDirectory: "infra/prod",
	}}
	cmd := newWorkspaceCreateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org",
			"-name=prod",
			"-working-directory=infra/prod",
			"-output=json",
		})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	if svc.lastOptions.WorkingDirectory == nil || *svc.lastOptions.WorkingDirectory != "infra/prod" {
		t.Fatalf("expected working directory 'infra/prod', got %v", svc.lastOptions.WorkingDirectory)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["WorkingDirectory"] != "infra/prod" {
		t.Fatalf("expected WorkingDirectory 'infra/prod' in output, got: %v", data["WorkingDirectory"])
	}
}

func TestWorkspaceCreateSourceFields(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceCreateService{response: &tfe.Workspace{
		ID:   "ws-1",
		Name: "prod",
	}}
	cmd := newWorkspaceCreateCommand(ui, svc)

	code := cmd.Run([]string{
		"-organization=my-org",
		"-name=prod",
		"-source-name=hcptf-cli",
		"-source-url=https://github.com/thrashr888/hcptf-cli",
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	if svc.lastOptions.SourceName == nil || *svc.lastOptions.SourceName != "hcptf-cli" {
		t.Fatalf("expected source name 'hcptf-cli', got %v", svc.lastOptions.SourceName)
	}
	if svc.lastOptions.SourceURL == nil || *svc.lastOptions.SourceURL != "https://github.com/thrashr888/hcptf-cli" {
		t.Fatalf("expected source URL, got %v", svc.lastOptions.SourceURL)
	}
}

func TestWorkspaceCreateTriggerPrefixes(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceCreateService{response: &tfe.Workspace{
		ID:              "ws-1",
		Name:            "prod",
		TriggerPrefixes: []string{"modules/", "configs/"},
	}}
	cmd := newWorkspaceCreateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org",
			"-name=prod",
			"-trigger-prefixes=modules/,configs/",
			"-output=json",
		})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	if len(svc.lastOptions.TriggerPrefixes) != 2 {
		t.Fatalf("expected 2 trigger prefixes, got %d: %v", len(svc.lastOptions.TriggerPrefixes), svc.lastOptions.TriggerPrefixes)
	}
	if svc.lastOptions.TriggerPrefixes[0] != "modules/" || svc.lastOptions.TriggerPrefixes[1] != "configs/" {
		t.Fatalf("unexpected trigger prefixes: %v", svc.lastOptions.TriggerPrefixes)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["TriggerPrefixes"] != "modules/,configs/" {
		t.Fatalf("expected TriggerPrefixes 'modules/,configs/' in output, got: %v", data["TriggerPrefixes"])
	}
}
