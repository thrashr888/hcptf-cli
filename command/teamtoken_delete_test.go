package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestTeamTokenDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamTokenDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestTeamTokenDeleteHelp(t *testing.T) {
	cmd := &TeamTokenDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf teamtoken delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestTeamTokenDeleteSynopsis(t *testing.T) {
	cmd := &TeamTokenDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a team API token" {
		t.Errorf("expected 'Delete a team API token', got %q", synopsis)
	}
}

func TestTeamTokenDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id without force",
			args:          []string{"-id=at-123abc"},
			expectedID:    "at-123abc",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=at-456def", "-force"},
			expectedID:    "at-456def",
			expectedForce: true,
		},
		{
			name:          "id with force=true",
			args:          []string{"-id=at-789ghi", "-force=true"},
			expectedID:    "at-789ghi",
			expectedForce: true,
		},
		{
			name:          "id with force=false",
			args:          []string{"-id=at-abc123", "-force=false"},
			expectedID:    "at-abc123",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &TeamTokenDeleteCommand{}

			flags := cmd.Meta.FlagSet("teamtoken delete")
			flags.StringVar(&cmd.id, "id", "", "Team token ID (required)")
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
