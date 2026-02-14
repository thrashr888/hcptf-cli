package command

import (
	"context"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

type mockTeamAccessUpdateService struct {
	response    *tfe.TeamAccess
	err         error
	lastID      string
	lastOptions tfe.TeamAccessUpdateOptions
}

func (m *mockTeamAccessUpdateService) Update(ctx context.Context, teamAccessID string, options tfe.TeamAccessUpdateOptions) (*tfe.TeamAccess, error) {
	m.lastID = teamAccessID
	m.lastOptions = options
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func newTeamAccessUpdateCommand(ui cli.Ui, svc teamAccessUpdater) *TeamAccessUpdateCommand {
	return &TeamAccessUpdateCommand{
		Meta:          newTestMeta(ui),
		teamAccessSvc: svc,
	}
}

func TestTeamAccessUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamAccessUpdateCommand(ui, &mockTeamAccessUpdateService{})

	code := cmd.Run([]string{"-access=read"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestTeamAccessUpdateRequiresAccess(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamAccessUpdateCommand(ui, &mockTeamAccessUpdateService{})

	code := cmd.Run([]string{"-id=tws-123abc"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-access") {
		t.Fatalf("expected access error, got %q", out)
	}
}

func TestTeamAccessUpdateRequiresBothFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamAccessUpdateCommand(ui, &mockTeamAccessUpdateService{})

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "required") {
		t.Fatalf("expected required error, got %q", out)
	}
}

func TestTeamAccessUpdateValidatesAccessLevel(t *testing.T) {
	tests := []struct {
		name        string
		accessLevel string
		expectError bool
	}{
		{
			name:        "valid read access",
			accessLevel: "read",
			expectError: false,
		},
		{
			name:        "valid plan access",
			accessLevel: "plan",
			expectError: false,
		},
		{
			name:        "valid write access",
			accessLevel: "write",
			expectError: false,
		},
		{
			name:        "valid admin access",
			accessLevel: "admin",
			expectError: false,
		},
		{
			name:        "valid custom access",
			accessLevel: "custom",
			expectError: false,
		},
		{
			name:        "invalid access level",
			accessLevel: "invalid",
			expectError: true,
		},
		{
			name:        "empty access level",
			accessLevel: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := cli.NewMockUi()
			svc := &mockTeamAccessUpdateService{
				response: &tfe.TeamAccess{
					ID: "tws-123abc",
					Team: &tfe.Team{
						ID: "team-1",
					},
					Workspace: &tfe.Workspace{
						ID: "ws-1",
					},
					Access: tfe.AccessType(tt.accessLevel),
				},
			}
			cmd := newTeamAccessUpdateCommand(ui, svc)

			args := []string{"-id=tws-123abc"}
			if tt.accessLevel != "" {
				args = append(args, "-access="+tt.accessLevel)
			}

			code := cmd.Run(args)

			if tt.expectError {
				if code != 1 {
					t.Fatalf("expected exit 1 for invalid access level, got %d", code)
				}
				errOut := ui.ErrorWriter.String()
				if !strings.Contains(errOut, "access") && !strings.Contains(errOut, "required") {
					t.Fatalf("expected access validation error, got %q", errOut)
				}
			} else {
				if code != 0 {
					t.Fatalf("expected exit 0 for valid access level, got %d. Error: %s", code, ui.ErrorWriter.String())
				}
			}
		})
	}
}

func TestTeamAccessUpdateHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamAccessUpdateService{err: errors.New("boom")}
	cmd := newTeamAccessUpdateCommand(ui, svc)

	code := cmd.Run([]string{"-id=tws-123abc", "-access=write"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if svc.lastID != "tws-123abc" {
		t.Fatalf("expected team access ID 'tws-123abc', got %q", svc.lastID)
	}

	if svc.lastOptions.Access == nil || *svc.lastOptions.Access != tfe.AccessWrite {
		t.Fatalf("expected access level write, got %v", svc.lastOptions.Access)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got %q", ui.ErrorWriter.String())
	}
}

func TestTeamAccessUpdateSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamAccessUpdateService{
		response: &tfe.TeamAccess{
			ID: "tws-123abc",
			Team: &tfe.Team{
				ID: "team-1",
			},
			Workspace: &tfe.Workspace{
				ID: "ws-1",
			},
			Access: tfe.AccessWrite,
		},
	}
	cmd := newTeamAccessUpdateCommand(ui, svc)

	code := cmd.Run([]string{"-id=tws-123abc", "-access=write"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d. Error: %s", code, ui.ErrorWriter.String())
	}

	if !strings.Contains(ui.OutputWriter.String(), "updated successfully") {
		t.Fatalf("expected success message, got %q", ui.OutputWriter.String())
	}
}

func TestTeamAccessUpdateOutputsJSON(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamAccessUpdateService{
		response: &tfe.TeamAccess{
			ID: "tws-123abc",
			Team: &tfe.Team{
				ID: "team-1",
			},
			Workspace: &tfe.Workspace{
				ID: "ws-1",
			},
			Access: tfe.AccessAdmin,
		},
	}
	cmd := newTeamAccessUpdateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=tws-123abc", "-access=admin", "-output=json"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d. Error: %s", code, ui.ErrorWriter.String())
	}

	if svc.lastOptions.Access == nil || *svc.lastOptions.Access != tfe.AccessAdmin {
		t.Fatalf("expected access level admin, got %v", svc.lastOptions.Access)
	}

	// The output includes both the success message and JSON output
	// Check that we have the success message from ui.Output
	uiOutput := ui.OutputWriter.String()
	if !strings.Contains(uiOutput, "updated successfully") {
		t.Fatalf("expected success message in UI output, got: %s", uiOutput)
	}

	// Check JSON output was printed to stdout (captured by captureStdout)
	// The output should contain JSON data for the team access
	if !strings.Contains(output, "tws-123abc") {
		t.Fatalf("expected team access ID in JSON output, got: %s", output)
	}
}

func TestTeamAccessUpdateCustomAccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockTeamAccessUpdateService{
		response: &tfe.TeamAccess{
			ID: "tws-123abc",
			Team: &tfe.Team{
				ID: "team-1",
			},
			Workspace: &tfe.Workspace{
				ID: "ws-1",
			},
			Access:           "custom",
			Runs:             tfe.RunsPermissionRead,
			Variables:        tfe.VariablesPermissionWrite,
			StateVersions:    tfe.StateVersionsPermissionReadOutputs,
			SentinelMocks:    tfe.SentinelMocksPermissionRead,
			WorkspaceLocking: false,
			RunTasks:         true,
		},
	}
	cmd := newTeamAccessUpdateCommand(ui, svc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=tws-123abc", "-access=custom"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d. Error: %s", code, ui.ErrorWriter.String())
	}

	uiOutput := ui.OutputWriter.String()
	if !strings.Contains(uiOutput, "updated successfully") {
		t.Fatalf("expected success message, got %q", uiOutput)
	}

	// Check that custom permissions are displayed (printed to stdout)
	if !strings.Contains(output, "Runs") {
		t.Errorf("expected custom permissions output to contain 'Runs', got %q", output)
	}
	if !strings.Contains(output, "Variables") {
		t.Errorf("expected custom permissions output to contain 'Variables', got %q", output)
	}
}

func TestTeamAccessUpdateHelp(t *testing.T) {
	cmd := &TeamAccessUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf teamaccess update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-access") {
		t.Error("Help should mention -access flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
	if !strings.Contains(help, "read") {
		t.Error("Help should mention read access level")
	}
	if !strings.Contains(help, "plan") {
		t.Error("Help should mention plan access level")
	}
	if !strings.Contains(help, "write") {
		t.Error("Help should mention write access level")
	}
	if !strings.Contains(help, "admin") {
		t.Error("Help should mention admin access level")
	}
	if !strings.Contains(help, "custom") {
		t.Error("Help should mention custom access level")
	}
}

func TestTeamAccessUpdateSynopsis(t *testing.T) {
	cmd := &TeamAccessUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update team workspace permissions" {
		t.Errorf("expected 'Update team workspace permissions', got %q", synopsis)
	}
}

func TestTeamAccessUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedAccess string
		expectedFormat string
	}{
		{
			name:           "id and access with default format",
			args:           []string{"-id=tws-123abc", "-access=read"},
			expectedID:     "tws-123abc",
			expectedAccess: "read",
			expectedFormat: "table",
		},
		{
			name:           "id and access with json format",
			args:           []string{"-id=tws-456def", "-access=write", "-output=json"},
			expectedID:     "tws-456def",
			expectedAccess: "write",
			expectedFormat: "json",
		},
		{
			name:           "id and plan access",
			args:           []string{"-id=tws-789ghi", "-access=plan"},
			expectedID:     "tws-789ghi",
			expectedAccess: "plan",
			expectedFormat: "table",
		},
		{
			name:           "id and admin access with json",
			args:           []string{"-id=tws-admin01", "-access=admin", "-output=json"},
			expectedID:     "tws-admin01",
			expectedAccess: "admin",
			expectedFormat: "json",
		},
		{
			name:           "id and custom access",
			args:           []string{"-id=tws-custom01", "-access=custom"},
			expectedID:     "tws-custom01",
			expectedAccess: "custom",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &TeamAccessUpdateCommand{}

			flags := cmd.Meta.FlagSet("teamaccess update")
			flags.StringVar(&cmd.id, "id", "", "Team access ID (required)")
			flags.StringVar(&cmd.access, "access", "", "Access level: read, plan, write, admin, or custom (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the access was set correctly
			if cmd.access != tt.expectedAccess {
				t.Errorf("expected access %q, got %q", tt.expectedAccess, cmd.access)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
