package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVariableSetApplyRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetApplyCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-workspaces=ws-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestVariableSetApplyRequiresTargets(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetApplyCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=varset-12345"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	out := ui.ErrorWriter.String()
	if !strings.Contains(out, "-workspaces") && !strings.Contains(out, "-projects") && !strings.Contains(out, "-stacks") {
		t.Fatalf("expected target error, got %q", out)
	}
}

func TestVariableSetApplyHelp(t *testing.T) {
	cmd := &VariableSetApplyCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf variableset apply") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-workspaces") {
		t.Error("Help should mention -workspaces flag")
	}
	if !strings.Contains(help, "-projects") {
		t.Error("Help should mention -projects flag")
	}
	if !strings.Contains(help, "-stacks") {
		t.Error("Help should mention -stacks flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestVariableSetApplySynopsis(t *testing.T) {
	cmd := &VariableSetApplyCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Apply a variable set to workspaces, projects, or stacks" {
		t.Errorf("expected updated synopsis, got %q", synopsis)
	}
}

func TestVariableSetApplyFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedID       string
		expectedWS       string
		expectedProjects string
		expectedStacks   string
	}{
		{
			name:             "apply to single workspace",
			args:             []string{"-id=varset-12345", "-workspaces=ws-abc123"},
			expectedID:       "varset-12345",
			expectedWS:       "ws-abc123",
			expectedProjects: "",
			expectedStacks:   "",
		},
		{
			name:             "apply to multiple workspaces",
			args:             []string{"-id=varset-xyz789", "-workspaces=ws-abc123,ws-def456,ws-ghi789"},
			expectedID:       "varset-xyz789",
			expectedWS:       "ws-abc123,ws-def456,ws-ghi789",
			expectedProjects: "",
			expectedStacks:   "",
		},
		{
			name:             "apply to single project",
			args:             []string{"-id=varset-abc123", "-projects=prj-xyz789"},
			expectedID:       "varset-abc123",
			expectedWS:       "",
			expectedProjects: "prj-xyz789",
			expectedStacks:   "",
		},
		{
			name:             "apply to workspaces projects and stacks",
			args:             []string{"-id=varset-def456", "-workspaces=ws-123,ws-456", "-projects=prj-abc,prj-def", "-stacks=stack-123"},
			expectedID:       "varset-def456",
			expectedWS:       "ws-123,ws-456",
			expectedProjects: "prj-abc,prj-def",
			expectedStacks:   "stack-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VariableSetApplyCommand{}

			flags := cmd.Meta.FlagSet("variableset apply")
			flags.StringVar(&cmd.id, "id", "", "Variable set ID (required)")
			flags.StringVar(&cmd.workspaces, "workspaces", "", "Comma-separated list of workspace IDs to apply to")
			flags.StringVar(&cmd.projects, "projects", "", "Comma-separated list of project IDs to apply to")
			flags.StringVar(&cmd.stacks, "stacks", "", "Comma-separated list of stack IDs to apply to")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the workspaces was set correctly
			if cmd.workspaces != tt.expectedWS {
				t.Errorf("expected workspaces %q, got %q", tt.expectedWS, cmd.workspaces)
			}

			// Verify the projects was set correctly
			if cmd.projects != tt.expectedProjects {
				t.Errorf("expected projects %q, got %q", tt.expectedProjects, cmd.projects)
			}
			if cmd.stacks != tt.expectedStacks {
				t.Errorf("expected stacks %q, got %q", tt.expectedStacks, cmd.stacks)
			}
		})
	}
}
