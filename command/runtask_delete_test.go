package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestRunTaskDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &RunTaskDeleteCommand{
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

func TestRunTaskDeleteHelp(t *testing.T) {
	cmd := &RunTaskDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf runtask delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestRunTaskDeleteSynopsis(t *testing.T) {
	cmd := &RunTaskDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a run task" {
		t.Errorf("expected 'Delete a run task', got %q", synopsis)
	}
}

func TestRunTaskDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id only, no force",
			args:          []string{"-id=task-ABC123"},
			expectedID:    "task-ABC123",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=task-XYZ789", "-force"},
			expectedID:    "task-XYZ789",
			expectedForce: true,
		},
		{
			name:          "id with force=true",
			args:          []string{"-id=task-DEF456", "-force=true"},
			expectedID:    "task-DEF456",
			expectedForce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &RunTaskDeleteCommand{}

			flags := cmd.Meta.FlagSet("runtask delete")
			flags.StringVar(&cmd.id, "id", "", "Run task ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the force was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
