package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestProjectDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectDeleteCommand{
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

func TestProjectDeleteHelp(t *testing.T) {
	cmd := &ProjectDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf project delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
}

func TestProjectDeleteSynopsis(t *testing.T) {
	cmd := &ProjectDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a project" {
		t.Errorf("expected 'Delete a project', got %q", synopsis)
	}
}

func TestProjectDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id without force",
			args:          []string{"-id=prj-abc123"},
			expectedID:    "prj-abc123",
			expectedForce: false,
		},
		{
			name:          "id with force",
			args:          []string{"-id=prj-old123", "-force"},
			expectedID:    "prj-old123",
			expectedForce: true,
		},
		{
			name:          "id with force=true",
			args:          []string{"-id=prj-deprecated", "-force=true"},
			expectedID:    "prj-deprecated",
			expectedForce: true,
		},
		{
			name:          "id with force=false",
			args:          []string{"-id=prj-keep123", "-force=false"},
			expectedID:    "prj-keep123",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ProjectDeleteCommand{}

			flags := cmd.Meta.FlagSet("project delete")
			flags.StringVar(&cmd.projectID, "id", "", "Project ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.projectID != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.projectID)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
