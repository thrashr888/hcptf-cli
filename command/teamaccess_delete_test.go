package command

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

type mockTeamAccessDeleteService struct {
	err    error
	lastID string
}

func (m *mockTeamAccessDeleteService) Remove(_ context.Context, teamAccessID string) error {
	m.lastID = teamAccessID
	return m.err
}

func newTeamAccessDeleteCommand(ui cli.Ui, taSvc teamAccessDeleter) *TeamAccessDeleteCommand {
	return &TeamAccessDeleteCommand{
		Meta:          newTestMeta(ui),
		teamAccessSvc: taSvc,
	}
}

func TestTeamAccessDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newTeamAccessDeleteCommand(ui, nil)

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestTeamAccessDeleteHelp(t *testing.T) {
	cmd := &TeamAccessDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf teamaccess delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id flag is required")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "Example:") {
		t.Error("Help should include examples")
	}
}

func TestTeamAccessDeleteSynopsis(t *testing.T) {
	cmd := &TeamAccessDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Remove team access from a workspace" {
		t.Errorf("expected 'Remove team access from a workspace', got %q", synopsis)
	}
}

func TestTeamAccessDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id only",
			args:          []string{"-id=tws-123abc"},
			expectedID:    "tws-123abc",
			expectedForce: false,
		},
		{
			name:          "id with force",
			args:          []string{"-id=tws-xyz789", "-force"},
			expectedID:    "tws-xyz789",
			expectedForce: true,
		},
		{
			name:          "id with explicit force=true",
			args:          []string{"-id=tws-test123", "-force=true"},
			expectedID:    "tws-test123",
			expectedForce: true,
		},
		{
			name:          "id with explicit force=false",
			args:          []string{"-id=tws-prod456", "-force=false"},
			expectedID:    "tws-prod456",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &TeamAccessDeleteCommand{}

			flags := cmd.Meta.FlagSet("teamaccess delete")
			flags.StringVar(&cmd.id, "id", "", "Team access ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}

func TestTeamAccessDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	taSvc := &mockTeamAccessDeleteService{err: errors.New("not found")}
	cmd := newTeamAccessDeleteCommand(ui, taSvc)

	code := cmd.Run([]string{"-id=tws-123", "-force"})
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

func TestTeamAccessDeleteSuccessWithForce(t *testing.T) {
	ui := cli.NewMockUi()
	taSvc := &mockTeamAccessDeleteService{}
	cmd := newTeamAccessDeleteCommand(ui, taSvc)

	code := cmd.Run([]string{"-id=tws-abc123", "-force"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if taSvc.lastID != "tws-abc123" {
		t.Fatalf("expected ID tws-abc123, got %s", taSvc.lastID)
	}

	if out := ui.OutputWriter.String(); !strings.Contains(out, "successfully") {
		t.Fatalf("expected success message, got %q", out)
	}
}

func TestTeamAccessDeleteCancelsWithoutConfirmation(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("no\n")
	taSvc := &mockTeamAccessDeleteService{}
	cmd := newTeamAccessDeleteCommand(ui, taSvc)

	code := cmd.Run([]string{"-id=tws-abc123"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if taSvc.lastID != "" {
		t.Fatalf("expected no deletion, but got ID %s", taSvc.lastID)
	}

	if out := ui.OutputWriter.String(); !strings.Contains(out, "cancelled") {
		t.Fatalf("expected cancellation message, got %q", out)
	}
}

func TestTeamAccessDeleteSuccessWithConfirmation(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("yes\n")
	taSvc := &mockTeamAccessDeleteService{}
	cmd := newTeamAccessDeleteCommand(ui, taSvc)

	code := cmd.Run([]string{"-id=tws-abc123"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if taSvc.lastID != "tws-abc123" {
		t.Fatalf("expected ID tws-abc123, got %s", taSvc.lastID)
	}

	if out := ui.OutputWriter.String(); !strings.Contains(out, "successfully") {
		t.Fatalf("expected success message, got %q", out)
	}
}
