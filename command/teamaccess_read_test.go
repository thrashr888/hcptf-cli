package command

import (
	"context"
	"errors"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

type mockTeamAccessReadService struct {
	response *tfe.TeamAccess
	err      error
	lastID   string
}

func (m *mockTeamAccessReadService) Read(_ context.Context, teamAccessID string) (*tfe.TeamAccess, error) {
	m.lastID = teamAccessID
	return m.response, m.err
}

func newTeamAccessReadCommand(ui cli.Ui, taSvc teamAccessReader) *TeamAccessReadCommand {
	return &TeamAccessReadCommand{
		Meta:          newTestMeta(ui),
		teamAccessSvc: taSvc,
	}
}

func TestTeamAccessReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamAccessReadCommand(ui, nil)

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestTeamAccessReadHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	taSvc := &mockTeamAccessReadService{err: errors.New("not found")}
	cmd := newTeamAccessReadCommand(ui, taSvc)

	code := cmd.Run([]string{"-id=tws-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if taSvc.lastID != "tws-123" {
		t.Fatalf("expected ID tws-123, got %s", taSvc.lastID)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "not found") {
		t.Fatalf("expected error output, got %q", out)
	}
}

func TestTeamAccessReadSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	taSvc := &mockTeamAccessReadService{
		response: &tfe.TeamAccess{
			ID:     "tws-abc123",
			Access: "read",
			Team: &tfe.Team{
				ID: "team-123",
			},
			Workspace: &tfe.Workspace{
				ID: "ws-456",
			},
		},
	}
	cmd := newTeamAccessReadCommand(ui, taSvc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=tws-abc123"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if taSvc.lastID != "tws-abc123" {
		t.Fatalf("expected ID tws-abc123, got %s", taSvc.lastID)
	}

	if !strings.Contains(output, "tws-abc123") {
		t.Fatalf("expected output to contain ID, got %q", output)
	}
	if !strings.Contains(output, "team-123") {
		t.Fatalf("expected output to contain team ID, got %q", output)
	}
	if !strings.Contains(output, "ws-456") {
		t.Fatalf("expected output to contain workspace ID, got %q", output)
	}
	if !strings.Contains(output, "read") {
		t.Fatalf("expected output to contain access level, got %q", output)
	}
}

func TestTeamAccessReadSuccessCustomAccess(t *testing.T) {
	ui := cli.NewMockUi()
	taSvc := &mockTeamAccessReadService{
		response: &tfe.TeamAccess{
			ID:               "tws-custom123",
			Access:           "custom",
			Runs:             "apply",
			Variables:        "write",
			StateVersions:    "read",
			SentinelMocks:    "read",
			WorkspaceLocking: true,
			RunTasks:         false,
			Team: &tfe.Team{
				ID: "team-789",
			},
			Workspace: &tfe.Workspace{
				ID: "ws-012",
			},
		},
	}
	cmd := newTeamAccessReadCommand(ui, taSvc)

	output, code := captureStdout(t, func() int {
		return cmd.Run([]string{"-id=tws-custom123"})
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if taSvc.lastID != "tws-custom123" {
		t.Fatalf("expected ID tws-custom123, got %s", taSvc.lastID)
	}

	// Check that all custom permissions are included
	if !strings.Contains(output, "custom") {
		t.Fatalf("expected output to contain custom access, got %q", output)
	}
	if !strings.Contains(output, "apply") {
		t.Fatalf("expected output to contain runs permission, got %q", output)
	}
	if !strings.Contains(output, "write") {
		t.Fatalf("expected output to contain variables permission, got %q", output)
	}
	if !strings.Contains(output, "read") {
		t.Fatalf("expected output to contain read permissions, got %q", output)
	}
}

func TestTeamAccessReadHelp(t *testing.T) {
	cmd := &TeamAccessReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf teamaccess read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
}

func TestTeamAccessReadSynopsis(t *testing.T) {
	cmd := &TeamAccessReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show team access details" {
		t.Errorf("expected 'Show team access details', got %q", synopsis)
	}
}

func TestTeamAccessReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id only, default format",
			args:        []string{"-id=tws-abc123"},
			expectedID:  "tws-abc123",
			expectedFmt: "table",
		},
		{
			name:        "id with table format",
			args:        []string{"-id=tws-xyz789", "-output=table"},
			expectedID:  "tws-xyz789",
			expectedFmt: "table",
		},
		{
			name:        "id with json format",
			args:        []string{"-id=tws-test456", "-output=json"},
			expectedID:  "tws-test456",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &TeamAccessReadCommand{}

			flags := cmd.Meta.FlagSet("teamaccess read")
			flags.StringVar(&cmd.id, "id", "", "Team access ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
