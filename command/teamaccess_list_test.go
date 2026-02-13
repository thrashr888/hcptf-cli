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

type mockTeamAccessListService struct {
	response    *tfe.TeamAccessList
	err         error
	lastOptions *tfe.TeamAccessListOptions
}

func (m *mockTeamAccessListService) List(_ context.Context, options *tfe.TeamAccessListOptions) (*tfe.TeamAccessList, error) {
	m.lastOptions = options
	return m.response, m.err
}

func newTeamAccessListCommand(ui cli.Ui, svc teamAccessLister) *TeamAccessListCommand {
	return &TeamAccessListCommand{
		Meta:          newTestMeta(ui),
		teamAccessSvc: svc,
	}
}

func TestTeamAccessListCommandRequiresWorkspaceID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamAccessListCommand(ui, &mockTeamAccessListService{})

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workspace-id") {
		t.Fatalf("expected workspace-id error, got %q", out)
	}
}

func TestTeamAccessListCommandHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamAccessListService{err: errors.New("api error")}
	cmd := newTeamAccessListCommand(ui, svc)

	code := cmd.Run([]string{"-workspace-id=ws-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if svc.lastOptions == nil || svc.lastOptions.WorkspaceID != "ws-123" {
		t.Fatalf("expected workspace ID ws-123, got %#v", svc.lastOptions)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "api error") {
		t.Fatalf("expected error output, got %q", out)
	}
}

func TestTeamAccessListCommandSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamAccessListService{
		response: &tfe.TeamAccessList{
			Items: []*tfe.TeamAccess{
				{
					ID:     "tws-abc123",
					Access: tfe.AccessRead,
					Team:   &tfe.Team{ID: "team-123"},
				},
			},
		},
	}
	cmd := newTeamAccessListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-workspace-id=ws-123"})
	})

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if svc.lastOptions.WorkspaceID != "ws-123" {
		t.Fatalf("expected workspace ID ws-123, got %s", svc.lastOptions.WorkspaceID)
	}

	if !strings.Contains(output, "tws-abc123") {
		t.Fatalf("expected team access ID in output, got %q", output)
	}
}

func TestTeamAccessListCommandOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamAccessListService{
		response: &tfe.TeamAccessList{
			Items: []*tfe.TeamAccess{
				{
					ID:     "tws-abc123",
					Access: tfe.AccessRead,
					Team:   &tfe.Team{ID: "team-123"},
				},
				{
					ID:     "tws-def456",
					Access: tfe.AccessWrite,
					Team:   &tfe.Team{ID: "team-456"},
				},
			},
		},
	}
	cmd := newTeamAccessListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-workspace-id=ws-123", "-output=json"})
	})

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	var rows []map[string]string
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}

	if rows[0]["ID"] != "tws-abc123" {
		t.Fatalf("unexpected first row: %#v", rows[0])
	}

	if rows[1]["Access Level"] != "write" {
		t.Fatalf("unexpected access level in second row: %#v", rows[1])
	}
}

func TestTeamAccessListCommandCustomAccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamAccessListService{
		response: &tfe.TeamAccessList{
			Items: []*tfe.TeamAccess{
				{
					ID:     "tws-custom",
					Access: tfe.AccessCustom,
					Team:   &tfe.Team{ID: "team-789"},
				},
			},
		},
	}
	cmd := newTeamAccessListCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-workspace-id=ws-123", "-output=json"})
	})

	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	var rows []map[string]string
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if rows[0]["Access Level"] != "custom" {
		t.Fatalf("expected custom access level, got %s", rows[0]["Access Level"])
	}
}

func TestTeamAccessListCommandEmptyList(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamAccessListService{
		response: &tfe.TeamAccessList{Items: []*tfe.TeamAccess{}},
	}
	cmd := newTeamAccessListCommand(ui, svc)

	code := cmd.Run([]string{"-workspace-id=ws-123"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if out := ui.OutputWriter.String(); !strings.Contains(out, "No team access found") {
		t.Fatalf("expected empty message, got %q", out)
	}
}

func TestTeamAccessListHelp(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamAccessListCommand(ui, &mockTeamAccessListService{})

	help := cmd.Help()

	// Check for usage line
	if !strings.Contains(help, "Usage: hcptf teamaccess list") {
		t.Errorf("Help should contain usage line")
	}

	// Check for description
	if !strings.Contains(help, "List team access for a workspace") {
		t.Errorf("Help should contain description")
	}

	// Check for required flag
	if !strings.Contains(help, "-workspace-id") {
		t.Errorf("Help should document -workspace-id flag")
	}

	// Check for output flag
	if !strings.Contains(help, "-output") {
		t.Errorf("Help should document -output flag")
	}

	// Check for examples
	if !strings.Contains(help, "Example:") {
		t.Errorf("Help should contain examples")
	}

	// Check for table and json formats
	if !strings.Contains(help, "table") || !strings.Contains(help, "json") {
		t.Errorf("Help should mention table and json output formats")
	}
}

func TestTeamAccessListSynopsis(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamAccessListCommand(ui, &mockTeamAccessListService{})

	synopsis := cmd.Synopsis()

	expected := "List team access for a workspace"
	if synopsis != expected {
		t.Errorf("Synopsis() = %q, want %q", synopsis, expected)
	}
}

func TestTeamAccessListFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorText   string
		checkFunc   func(*testing.T, *mockTeamAccessListService)
	}{
		{
			name:        "missing workspace-id flag",
			args:        []string{},
			expectError: true,
			errorText:   "-workspace-id",
		},
		{
			name:        "valid workspace-id",
			args:        []string{"-workspace-id=ws-123"},
			expectError: false,
			checkFunc: func(t *testing.T, svc *mockTeamAccessListService) {
				if svc.lastOptions == nil || svc.lastOptions.WorkspaceID != "ws-123" {
					t.Errorf("expected workspace ID ws-123, got %#v", svc.lastOptions)
				}
			},
		},
		{
			name:        "table output format",
			args:        []string{"-workspace-id=ws-123", "-output=table"},
			expectError: false,
			checkFunc: func(t *testing.T, svc *mockTeamAccessListService) {
				if svc.lastOptions.WorkspaceID != "ws-123" {
					t.Errorf("expected workspace ID ws-123")
				}
			},
		},
		{
			name:        "json output format",
			args:        []string{"-workspace-id=ws-123", "-output=json"},
			expectError: false,
			checkFunc: func(t *testing.T, svc *mockTeamAccessListService) {
				if svc.lastOptions.WorkspaceID != "ws-123" {
					t.Errorf("expected workspace ID ws-123")
				}
			},
		},
		{
			name:        "workspace-id with special characters",
			args:        []string{"-workspace-id=ws-abc-123-def"},
			expectError: false,
			checkFunc: func(t *testing.T, svc *mockTeamAccessListService) {
				if svc.lastOptions == nil || svc.lastOptions.WorkspaceID != "ws-abc-123-def" {
					t.Errorf("expected workspace ID ws-abc-123-def, got %#v", svc.lastOptions)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := cli.NewMockUi()
			svc := &mockTeamAccessListService{
				response: &tfe.TeamAccessList{
					Items: []*tfe.TeamAccess{
						{
							ID:     "tws-test",
							Access: tfe.AccessRead,
							Team:   &tfe.Team{ID: "team-test"},
						},
					},
				},
			}
			cmd := newTeamAccessListCommand(ui, svc)

			code := cmd.Run(tt.args)

			if tt.expectError {
				if code == 0 {
					t.Errorf("expected non-zero exit code, got 0")
				}
				if tt.errorText != "" {
					errOutput := ui.ErrorWriter.String()
					if !strings.Contains(errOutput, tt.errorText) {
						t.Errorf("expected error containing %q, got %q", tt.errorText, errOutput)
					}
				}
			} else {
				if code != 0 {
					t.Errorf("expected exit code 0, got %d. Error: %s", code, ui.ErrorWriter.String())
				}
				if tt.checkFunc != nil {
					tt.checkFunc(t, svc)
				}
			}
		})
	}
}
