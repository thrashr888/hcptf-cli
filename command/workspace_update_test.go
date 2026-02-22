package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func newWorkspaceUpdateCommand(ui cli.Ui, svc workspaceUpdater) *WorkspaceUpdateCommand {
	return &WorkspaceUpdateCommand{
		Meta:         newTestMeta(ui),
		workspaceSvc: svc,
	}
}

func TestWorkspaceUpdateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceUpdateCommand(ui, &mockWorkspaceUpdateService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error")
	}

	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing name")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-name") {
		t.Fatalf("expected name error")
	}
}

func TestWorkspaceUpdateValidatesAutoApply(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceUpdateCommand(ui, &mockWorkspaceUpdateService{})

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod", "-auto-apply=maybe"}); code != 1 {
		t.Fatalf("expected exit 1 for invalid auto-apply")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "auto-apply") {
		t.Fatalf("expected validation error")
	}
}

func TestWorkspaceUpdateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{err: errors.New("boom")}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	if code := cmd.Run([]string{"-organization=my-org", "-name=prod", "-auto-apply=true"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastOrg != "my-org" || svc.lastName != "prod" {
		t.Fatalf("unexpected parameters: %#v", svc)
	}
	if svc.lastOptions.AutoApply == nil || !*svc.lastOptions.AutoApply {
		t.Fatalf("expected auto apply true")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestWorkspaceUpdateOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{response: &tfe.Workspace{ID: "ws-1", Name: "new-name", TerraformVersion: "1.6.1", AutoApply: true}}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=prod", "-new-name=new", "-terraform-version=1.6.1", "-description=hello", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if svc.lastOptions.Name == nil || *svc.lastOptions.Name != "new" {
		t.Fatalf("expected new name option")
	}
	if svc.lastOptions.Description == nil || *svc.lastOptions.Description != "hello" {
		t.Fatalf("expected description option")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["Name"] != "new-name" {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestWorkspaceUpdateProjectID(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{response: &tfe.Workspace{
		ID:   "ws-1",
		Name: "my-workspace",
		Project: &tfe.Project{
			ID:   "prj-abc123",
			Name: "Gym",
		},
	}}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-organization=my-org", "-name=my-workspace", "-project-id=prj-abc123", "-output=json"})
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

func TestWorkspaceUpdateExecutionMode(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{response: &tfe.Workspace{
		ID:            "ws-1",
		Name:          "test",
		ExecutionMode: "agent",
	}}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org",
			"-name=test",
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
	cmd2 := newWorkspaceUpdateCommand(ui2, &mockWorkspaceUpdateService{})
	if code := cmd2.Run([]string{"-organization=my-org", "-name=test", "-execution-mode=invalid"}); code != 1 {
		t.Fatalf("expected exit 1 for invalid execution mode")
	}
	if !strings.Contains(ui2.ErrorWriter.String(), "execution-mode") {
		t.Fatalf("expected execution-mode validation error, got: %s", ui2.ErrorWriter.String())
	}
}

func TestWorkspaceUpdateWorkingDirectory(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{response: &tfe.Workspace{
		ID:               "ws-1",
		Name:             "test",
		WorkingDirectory: "infra/prod",
	}}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org",
			"-name=test",
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

func TestWorkspaceUpdateVCSRepo(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{response: &tfe.Workspace{
		ID:   "ws-1",
		Name: "test",
		VCSRepo: &tfe.VCSRepo{
			Identifier: "my-org/my-repo",
		},
	}}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org",
			"-name=test",
			"-vcs-identifier=my-org/my-repo",
			"-vcs-branch=main",
			"-vcs-oauth-token-id=ot-abc",
			"-vcs-ingress-submodules=true",
			"-vcs-tags-regex=v\\d+",
			"-vcs-gha-installation-id=gha-123",
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
	if vcs.IngressSubmodules == nil || !*vcs.IngressSubmodules {
		t.Fatalf("expected VCS ingress submodules true, got %v", vcs.IngressSubmodules)
	}
	if vcs.TagsRegex == nil || *vcs.TagsRegex != "v\\d+" {
		t.Fatalf("expected VCS tags regex, got %v", vcs.TagsRegex)
	}
	if vcs.GHAInstallationID == nil || *vcs.GHAInstallationID != "gha-123" {
		t.Fatalf("expected GHA installation ID 'gha-123', got %v", vcs.GHAInstallationID)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if data["VCSRepo"] != "my-org/my-repo" {
		t.Fatalf("expected VCSRepo 'my-org/my-repo' in output, got: %v", data["VCSRepo"])
	}
}

func TestWorkspaceUpdateRemoveVCS(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{response: &tfe.Workspace{
		ID:   "ws-1",
		Name: "test",
	}}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	code := cmd.Run([]string{
		"-organization=my-org",
		"-name=test",
		"-remove-vcs",
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	if svc.lastOptions.VCSRepo == nil {
		t.Fatal("expected VCSRepo to be set (empty struct for deletion)")
	}
	// Empty struct means all fields should be nil
	if svc.lastOptions.VCSRepo.Identifier != nil {
		t.Fatalf("expected VCSRepo.Identifier to be nil, got %v", svc.lastOptions.VCSRepo.Identifier)
	}
}

func TestWorkspaceUpdateRemoveVCSConflict(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newWorkspaceUpdateCommand(ui, &mockWorkspaceUpdateService{})

	code := cmd.Run([]string{
		"-organization=my-org",
		"-name=test",
		"-remove-vcs",
		"-vcs-identifier=org/repo",
	})
	if code != 1 {
		t.Fatalf("expected exit 1 for remove-vcs conflict")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-remove-vcs cannot be used together with -vcs-* flags") {
		t.Fatalf("expected conflict error, got: %s", ui.ErrorWriter.String())
	}
}

func TestWorkspaceUpdateTriggerPrefixes(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{response: &tfe.Workspace{
		ID:              "ws-1",
		Name:            "test",
		TriggerPrefixes: []string{"modules/", "configs/"},
	}}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-organization=my-org",
			"-name=test",
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

func TestWorkspaceUpdateBoolFlags(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{response: &tfe.Workspace{
		ID:                  "ws-1",
		Name:                "test",
		AllowDestroyPlan:    true,
		FileTriggersEnabled: false,
		QueueAllRuns:        true,
		SpeculativeEnabled:  false,
	}}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	code := cmd.Run([]string{
		"-organization=my-org",
		"-name=test",
		"-allow-destroy-plan=true",
		"-file-triggers-enabled=false",
		"-queue-all-runs=true",
		"-speculative-enabled=false",
		"-global-remote-state=true",
		"-structured-run-output-enabled=true",
		"-inherits-project-auto-destroy=false",
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	opts := svc.lastOptions
	if opts.AllowDestroyPlan == nil || !*opts.AllowDestroyPlan {
		t.Fatalf("expected AllowDestroyPlan true")
	}
	if opts.FileTriggersEnabled == nil || *opts.FileTriggersEnabled {
		t.Fatalf("expected FileTriggersEnabled false")
	}
	if opts.QueueAllRuns == nil || !*opts.QueueAllRuns {
		t.Fatalf("expected QueueAllRuns true")
	}
	if opts.SpeculativeEnabled == nil || *opts.SpeculativeEnabled {
		t.Fatalf("expected SpeculativeEnabled false")
	}
	if opts.GlobalRemoteState == nil || !*opts.GlobalRemoteState {
		t.Fatalf("expected GlobalRemoteState true")
	}
	if opts.StructuredRunOutputEnabled == nil || !*opts.StructuredRunOutputEnabled {
		t.Fatalf("expected StructuredRunOutputEnabled true")
	}
	if opts.InheritsProjectAutoDestroy == nil || *opts.InheritsProjectAutoDestroy {
		t.Fatalf("expected InheritsProjectAutoDestroy false")
	}
}

func TestWorkspaceUpdateAutoDestroyAt(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{response: &tfe.Workspace{
		ID:   "ws-1",
		Name: "test",
	}}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	destroyTime := "2025-12-31T23:59:59Z"
	code := cmd.Run([]string{
		"-organization=my-org",
		"-name=test",
		"-auto-destroy-at=" + destroyTime,
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	// Verify the option was set with the parsed time
	got, err := svc.lastOptions.AutoDestroyAt.Get()
	if err != nil {
		t.Fatalf("expected AutoDestroyAt to be set, got error: %v", err)
	}
	expected, _ := time.Parse(time.RFC3339, destroyTime)
	if !got.Equal(expected) {
		t.Fatalf("expected AutoDestroyAt %v, got %v", expected, got)
	}

	// Test invalid time
	ui2 := cli.NewMockUi()
	cmd2 := newWorkspaceUpdateCommand(ui2, &mockWorkspaceUpdateService{})
	if code := cmd2.Run([]string{"-organization=my-org", "-name=test", "-auto-destroy-at=not-a-time"}); code != 1 {
		t.Fatalf("expected exit 1 for invalid time")
	}
	if !strings.Contains(ui2.ErrorWriter.String(), "RFC3339") {
		t.Fatalf("expected RFC3339 error, got: %s", ui2.ErrorWriter.String())
	}
}

func TestWorkspaceUpdateAutoDestroyAtNone(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockWorkspaceUpdateService{response: &tfe.Workspace{
		ID:   "ws-1",
		Name: "test",
	}}
	cmd := newWorkspaceUpdateCommand(ui, svc)

	code := cmd.Run([]string{
		"-organization=my-org",
		"-name=test",
		"-auto-destroy-at=none",
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, ui.ErrorWriter.String())
	}

	// Verify the option was set as null (to clear)
	if !svc.lastOptions.AutoDestroyAt.IsNull() {
		t.Fatal("expected AutoDestroyAt to be null (cleared)")
	}
}
