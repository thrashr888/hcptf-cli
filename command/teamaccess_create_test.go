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

type mockTeamAccessCreateService struct {
	response    *tfe.TeamAccess
	err         error
	lastOptions tfe.TeamAccessAddOptions
}

func (m *mockTeamAccessCreateService) Add(_ context.Context, options tfe.TeamAccessAddOptions) (*tfe.TeamAccess, error) {
	m.lastOptions = options
	return m.response, m.err
}

func newTeamAccessCreateCommand(ui cli.Ui, svc teamAccessCreator) *TeamAccessCreateCommand {
	return &TeamAccessCreateCommand{
		Meta:          newTestMeta(ui),
		teamAccessSvc: svc,
	}
}

func TestTeamAccessCreateRequiresWorkspaceID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamAccessCreateCommand(ui, &mockTeamAccessCreateService{})

	if code := cmd.Run([]string{"-team-id=team-123", "-access=read"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace-id") {
		t.Fatalf("expected workspace-id error, got %q", out)
	}
}

func TestTeamAccessCreateRequiresTeamID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamAccessCreateCommand(ui, &mockTeamAccessCreateService{})

	if code := cmd.Run([]string{"-workspace-id=ws-123", "-access=read"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-team-id") {
		t.Fatalf("expected team-id error, got %q", out)
	}
}

func TestTeamAccessCreateRequiresAccess(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamAccessCreateCommand(ui, &mockTeamAccessCreateService{})

	if code := cmd.Run([]string{"-workspace-id=ws-123", "-team-id=team-123"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-access") {
		t.Fatalf("expected access error, got %q", out)
	}
}

func TestTeamAccessCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamAccessCreateCommand(ui, &mockTeamAccessCreateService{})

	if code := cmd.Run([]string{}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "required") {
		t.Fatalf("expected required error, got %q", out)
	}
}

func TestTeamAccessCreateValidatesAccessLevel(t *testing.T) {
	tests := []struct {
		name        string
		accessLevel string
		shouldFail  bool
	}{
		{
			name:        "valid read access",
			accessLevel: "read",
			shouldFail:  false,
		},
		{
			name:        "valid plan access",
			accessLevel: "plan",
			shouldFail:  false,
		},
		{
			name:        "valid write access",
			accessLevel: "write",
			shouldFail:  false,
		},
		{
			name:        "valid admin access",
			accessLevel: "admin",
			shouldFail:  false,
		},
		{
			name:        "valid custom access",
			accessLevel: "custom",
			shouldFail:  false,
		},
		{
			name:        "invalid access level",
			accessLevel: "invalid",
			shouldFail:  true,
		},
		{
			name:        "case insensitive read",
			accessLevel: "READ",
			shouldFail:  false,
		},
		{
			name:        "case insensitive write",
			accessLevel: "Write",
			shouldFail:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := cli.NewMockUi()
			svc := &mockTeamAccessCreateService{
				response: &tfe.TeamAccess{
					ID: "ta-123",
					Team: &tfe.Team{
						ID: "team-123",
					},
					Workspace: &tfe.Workspace{
						ID: "ws-123",
					},
					Access: tfe.AccessRead,
				},
			}
			cmd := newTeamAccessCreateCommand(ui, svc)

			code := cmd.Run([]string{
				"-workspace-id=ws-123",
				"-team-id=team-123",
				"-access=" + tt.accessLevel,
			})

			if tt.shouldFail {
				if code != 1 {
					t.Fatalf("expected exit 1, got %d", code)
				}
				if out := ui.ErrorWriter.String(); !strings.Contains(out, "invalid access level") {
					t.Fatalf("expected invalid access level error, got %q", out)
				}
			} else {
				if code != 0 {
					t.Fatalf("expected exit 0, got %d", code)
				}
			}
		})
	}
}

func TestTeamAccessCreateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamAccessCreateService{err: errors.New("api error")}
	cmd := newTeamAccessCreateCommand(ui, svc)

	code := cmd.Run([]string{
		"-workspace-id=ws-123",
		"-team-id=team-456",
		"-access=write",
	})

	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if svc.lastOptions.Workspace == nil || svc.lastOptions.Workspace.ID != "ws-123" {
		t.Fatalf("expected workspace ID ws-123, got %#v", svc.lastOptions.Workspace)
	}

	if svc.lastOptions.Team == nil || svc.lastOptions.Team.ID != "team-456" {
		t.Fatalf("expected team ID team-456, got %#v", svc.lastOptions.Team)
	}

	if svc.lastOptions.Access == nil || *svc.lastOptions.Access != tfe.AccessWrite {
		t.Fatalf("expected access write, got %#v", svc.lastOptions.Access)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "api error") {
		t.Fatalf("expected error output, got %q", out)
	}
}

func TestTeamAccessCreateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamAccessCreateService{
		response: &tfe.TeamAccess{
			ID: "ta-new123",
			Team: &tfe.Team{
				ID: "team-456",
			},
			Workspace: &tfe.Workspace{
				ID: "ws-123",
			},
			Access: tfe.AccessWrite,
		},
	}
	cmd := newTeamAccessCreateCommand(ui, svc)

	code := cmd.Run([]string{
		"-workspace-id=ws-123",
		"-team-id=team-456",
		"-access=write",
	})

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if out := ui.OutputWriter.String(); !strings.Contains(out, "successfully") {
		t.Fatalf("expected success message, got %q", out)
	}
}

func TestTeamAccessCreateOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamAccessCreateService{
		response: &tfe.TeamAccess{
			ID: "ta-123",
			Team: &tfe.Team{
				ID: "team-456",
			},
			Workspace: &tfe.Workspace{
				ID: "ws-789",
			},
			Access: tfe.AccessAdmin,
		},
	}
	cmd := newTeamAccessCreateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{
			"-workspace-id=ws-789",
			"-team-id=team-456",
			"-access=admin",
			"-output=json",
		})
	})

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if svc.lastOptions.Access == nil || *svc.lastOptions.Access != tfe.AccessAdmin {
		t.Fatalf("expected admin access")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if data["ID"] != "ta-123" {
		t.Fatalf("unexpected ID in data: %#v", data)
	}

	if data["TeamID"] != "team-456" {
		t.Fatalf("unexpected TeamID in data: %#v", data)
	}

	if data["WorkspaceID"] != "ws-789" {
		t.Fatalf("unexpected WorkspaceID in data: %#v", data)
	}

	if data["Access"] != string(tfe.AccessAdmin) {
		t.Fatalf("unexpected Access in data: %#v", data)
	}
}

func TestTeamAccessCreateAccessLevelMapping(t *testing.T) {
	tests := []struct {
		name           string
		inputAccess    string
		expectedAccess tfe.AccessType
	}{
		{
			name:           "read access",
			inputAccess:    "read",
			expectedAccess: tfe.AccessRead,
		},
		{
			name:           "plan access",
			inputAccess:    "plan",
			expectedAccess: tfe.AccessPlan,
		},
		{
			name:           "write access",
			inputAccess:    "write",
			expectedAccess: tfe.AccessWrite,
		},
		{
			name:           "admin access",
			inputAccess:    "admin",
			expectedAccess: tfe.AccessAdmin,
		},
		{
			name:           "custom access",
			inputAccess:    "custom",
			expectedAccess: tfe.AccessCustom,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := cli.NewMockUi()
			svc := &mockTeamAccessCreateService{
				response: &tfe.TeamAccess{
					ID: "ta-123",
					Team: &tfe.Team{
						ID: "team-123",
					},
					Workspace: &tfe.Workspace{
						ID: "ws-123",
					},
					Access: tt.expectedAccess,
				},
			}
			cmd := newTeamAccessCreateCommand(ui, svc)

			code := cmd.Run([]string{
				"-workspace-id=ws-123",
				"-team-id=team-123",
				"-access=" + tt.inputAccess,
			})

			if code != 0 {
				t.Fatalf("expected exit 0, got %d", code)
			}

			if svc.lastOptions.Access == nil || *svc.lastOptions.Access != tt.expectedAccess {
				t.Fatalf("expected access %v, got %v", tt.expectedAccess, svc.lastOptions.Access)
			}
		})
	}
}

func TestTeamAccessCreateHelp(t *testing.T) {
	cmd := &TeamAccessCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf teamaccess create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-workspace-id") {
		t.Error("Help should mention -workspace-id flag")
	}
	if !strings.Contains(help, "-team-id") {
		t.Error("Help should mention -team-id flag")
	}
	if !strings.Contains(help, "-access") {
		t.Error("Help should mention -access flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
	if !strings.Contains(help, "Access Levels:") {
		t.Error("Help should contain Access Levels section")
	}
	if !strings.Contains(help, "read") || !strings.Contains(help, "plan") ||
		!strings.Contains(help, "write") || !strings.Contains(help, "admin") ||
		!strings.Contains(help, "custom") {
		t.Error("Help should list all access levels")
	}
	if !strings.Contains(help, "Example:") {
		t.Error("Help should contain examples")
	}
}

func TestTeamAccessCreateSynopsis(t *testing.T) {
	cmd := &TeamAccessCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Grant team access to a workspace" {
		t.Errorf("expected 'Grant team access to a workspace', got %q", synopsis)
	}
}

func TestTeamAccessCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedWSID     string
		expectedTeamID   string
		expectedAccess   string
		expectedFmt      string
	}{
		{
			name:             "all required flags with defaults",
			args:             []string{"-workspace-id=ws-123", "-team-id=team-456", "-access=read"},
			expectedWSID:     "ws-123",
			expectedTeamID:   "team-456",
			expectedAccess:   "read",
			expectedFmt:      "table",
		},
		{
			name:             "all flags with json output",
			args:             []string{"-workspace-id=ws-abc", "-team-id=team-xyz", "-access=write", "-output=json"},
			expectedWSID:     "ws-abc",
			expectedTeamID:   "team-xyz",
			expectedAccess:   "write",
			expectedFmt:      "json",
		},
		{
			name:             "plan access level",
			args:             []string{"-workspace-id=ws-111", "-team-id=team-222", "-access=plan"},
			expectedWSID:     "ws-111",
			expectedTeamID:   "team-222",
			expectedAccess:   "plan",
			expectedFmt:      "table",
		},
		{
			name:             "admin access level with json",
			args:             []string{"-workspace-id=ws-aaa", "-team-id=team-bbb", "-access=admin", "-output=json"},
			expectedWSID:     "ws-aaa",
			expectedTeamID:   "team-bbb",
			expectedAccess:   "admin",
			expectedFmt:      "json",
		},
		{
			name:             "custom access level",
			args:             []string{"-workspace-id=ws-custom", "-team-id=team-custom", "-access=custom"},
			expectedWSID:     "ws-custom",
			expectedTeamID:   "team-custom",
			expectedAccess:   "custom",
			expectedFmt:      "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &TeamAccessCreateCommand{}

			flags := cmd.Meta.FlagSet("teamaccess create")
			flags.StringVar(&cmd.workspaceID, "workspace-id", "", "Workspace ID (required)")
			flags.StringVar(&cmd.teamID, "team-id", "", "Team ID (required)")
			flags.StringVar(&cmd.access, "access", "", "Access level: read, plan, write, admin, or custom (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the workspace-id was set correctly
			if cmd.workspaceID != tt.expectedWSID {
				t.Errorf("expected workspace-id %q, got %q", tt.expectedWSID, cmd.workspaceID)
			}

			// Verify the team-id was set correctly
			if cmd.teamID != tt.expectedTeamID {
				t.Errorf("expected team-id %q, got %q", tt.expectedTeamID, cmd.teamID)
			}

			// Verify the access was set correctly
			if cmd.access != tt.expectedAccess {
				t.Errorf("expected access %q, got %q", tt.expectedAccess, cmd.access)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
