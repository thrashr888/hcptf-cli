package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestProjectTeamAccessCreateRequiresProjectID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectTeamAccessCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-team-id=team-456", "-access=read"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-project-id") {
		t.Fatalf("expected project-id error, got %q", out)
	}
}

func TestProjectTeamAccessCreateRequiresTeamID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectTeamAccessCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-project-id=prj-123", "-access=read"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-team-id") {
		t.Fatalf("expected team-id error, got %q", out)
	}
}

func TestProjectTeamAccessCreateRequiresAccess(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectTeamAccessCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-project-id=prj-123", "-team-id=team-456"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-access") {
		t.Fatalf("expected access error, got %q", out)
	}
}

func TestProjectTeamAccessCreateInvalidAccessLevel(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectTeamAccessCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-project-id=prj-123", "-team-id=team-456", "-access=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "invalid access level") {
		t.Fatalf("expected invalid access level error, got %q", out)
	}
}

func TestProjectTeamAccessCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectTeamAccessCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "required") {
		t.Fatalf("expected required error, got %q", out)
	}
}

func TestProjectTeamAccessCreateHelp(t *testing.T) {
	cmd := &ProjectTeamAccessCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf projectteamaccess create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-project-id") {
		t.Error("Help should mention -project-id flag")
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
	if !strings.Contains(help, "read") {
		t.Error("Help should mention read access level")
	}
	if !strings.Contains(help, "write") {
		t.Error("Help should mention write access level")
	}
	if !strings.Contains(help, "maintain") {
		t.Error("Help should mention maintain access level")
	}
	if !strings.Contains(help, "admin") {
		t.Error("Help should mention admin access level")
	}
	if !strings.Contains(help, "custom") {
		t.Error("Help should mention custom access level")
	}
}

func TestProjectTeamAccessCreateSynopsis(t *testing.T) {
	cmd := &ProjectTeamAccessCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Grant team access to a project" {
		t.Errorf("expected 'Grant team access to a project', got %q", synopsis)
	}
}

func TestProjectTeamAccessCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedProjID string
		expectedTeamID string
		expectedAccess string
		expectedOutput string
	}{
		{
			name:           "all required flags",
			args:           []string{"-project-id=prj-123", "-team-id=team-456", "-access=read"},
			expectedProjID: "prj-123",
			expectedTeamID: "team-456",
			expectedAccess: "read",
			expectedOutput: "table",
		},
		{
			name:           "with output flag",
			args:           []string{"-project-id=prj-123", "-team-id=team-456", "-access=write", "-output=json"},
			expectedProjID: "prj-123",
			expectedTeamID: "team-456",
			expectedAccess: "write",
			expectedOutput: "json",
		},
		{
			name:           "maintain access level",
			args:           []string{"-project-id=prj-abc", "-team-id=team-def", "-access=maintain"},
			expectedProjID: "prj-abc",
			expectedTeamID: "team-def",
			expectedAccess: "maintain",
			expectedOutput: "table",
		},
		{
			name:           "admin access level",
			args:           []string{"-project-id=prj-xyz", "-team-id=team-123", "-access=admin"},
			expectedProjID: "prj-xyz",
			expectedTeamID: "team-123",
			expectedAccess: "admin",
			expectedOutput: "table",
		},
		{
			name:           "custom access level",
			args:           []string{"-project-id=prj-custom", "-team-id=team-custom", "-access=custom"},
			expectedProjID: "prj-custom",
			expectedTeamID: "team-custom",
			expectedAccess: "custom",
			expectedOutput: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ProjectTeamAccessCreateCommand{}

			flags := cmd.Meta.FlagSet("projectteamaccess create")
			flags.StringVar(&cmd.projectID, "project-id", "", "Project ID (required)")
			flags.StringVar(&cmd.teamID, "team-id", "", "Team ID (required)")
			flags.StringVar(&cmd.access, "access", "", "Access level")
			flags.StringVar(&cmd.format, "output", "table", "Output format")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify flags were set correctly
			if cmd.projectID != tt.expectedProjID {
				t.Errorf("expected project-id %q, got %q", tt.expectedProjID, cmd.projectID)
			}
			if cmd.teamID != tt.expectedTeamID {
				t.Errorf("expected team-id %q, got %q", tt.expectedTeamID, cmd.teamID)
			}
			if cmd.access != tt.expectedAccess {
				t.Errorf("expected access %q, got %q", tt.expectedAccess, cmd.access)
			}
			if cmd.format != tt.expectedOutput {
				t.Errorf("expected output %q, got %q", tt.expectedOutput, cmd.format)
			}
		})
	}
}
