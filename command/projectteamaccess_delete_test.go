package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestProjectTeamAccessDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectTeamAccessDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-force"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestProjectTeamAccessDeleteHelp(t *testing.T) {
	cmd := &ProjectTeamAccessDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf projectteamaccess delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "confirmation") {
		t.Error("Help should mention confirmation")
	}
	if !strings.Contains(help, "tprj-") {
		t.Error("Help should contain example ID")
	}
}

func TestProjectTeamAccessDeleteSynopsis(t *testing.T) {
	cmd := &ProjectTeamAccessDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Remove team access from a project" {
		t.Errorf("expected 'Remove team access from a project', got %q", synopsis)
	}
}

func TestProjectTeamAccessDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "required flags only",
			args:          []string{"-id=tprj-123"},
			expectedID:    "tprj-123",
			expectedForce: false,
		},
		{
			name:          "with force flag",
			args:          []string{"-id=tprj-456", "-force"},
			expectedID:    "tprj-456",
			expectedForce: true,
		},
		{
			name:          "force flag set to true",
			args:          []string{"-id=tprj-abc", "-force=true"},
			expectedID:    "tprj-abc",
			expectedForce: true,
		},
		{
			name:          "force flag set to false",
			args:          []string{"-id=tprj-xyz", "-force=false"},
			expectedID:    "tprj-xyz",
			expectedForce: false,
		},
		{
			name:          "different id format",
			args:          []string{"-id=tprj-xyz123abc"},
			expectedID:    "tprj-xyz123abc",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ProjectTeamAccessDeleteCommand{}

			flags := cmd.Meta.FlagSet("projectteamaccess delete")
			flags.StringVar(&cmd.id, "id", "", "Project team access ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify flags were set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
