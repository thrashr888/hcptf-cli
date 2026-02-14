package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestProjectTeamAccessUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectTeamAccessUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-access=write"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestProjectTeamAccessUpdateRequiresAccess(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectTeamAccessUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=tprj-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-access") {
		t.Fatalf("expected access error, got %q", out)
	}
}

func TestProjectTeamAccessUpdateInvalidAccessLevel(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectTeamAccessUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=tprj-123", "-access=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "invalid access level") {
		t.Fatalf("expected invalid access level error, got %q", out)
	}
}

func TestProjectTeamAccessUpdateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectTeamAccessUpdateCommand{
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

func TestProjectTeamAccessUpdateHelp(t *testing.T) {
	cmd := &ProjectTeamAccessUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf projectteamaccess update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
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

func TestProjectTeamAccessUpdateSynopsis(t *testing.T) {
	cmd := &ProjectTeamAccessUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update team project permissions" {
		t.Errorf("expected 'Update team project permissions', got %q", synopsis)
	}
}

func TestProjectTeamAccessUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedAccess string
		expectedOutput string
	}{
		{
			name:           "read access",
			args:           []string{"-id=tprj-123", "-access=read"},
			expectedID:     "tprj-123",
			expectedAccess: "read",
			expectedOutput: "table",
		},
		{
			name:           "write access with json",
			args:           []string{"-id=tprj-456", "-access=write", "-output=json"},
			expectedID:     "tprj-456",
			expectedAccess: "write",
			expectedOutput: "json",
		},
		{
			name:           "maintain access",
			args:           []string{"-id=tprj-abc", "-access=maintain"},
			expectedID:     "tprj-abc",
			expectedAccess: "maintain",
			expectedOutput: "table",
		},
		{
			name:           "admin access",
			args:           []string{"-id=tprj-xyz", "-access=admin"},
			expectedID:     "tprj-xyz",
			expectedAccess: "admin",
			expectedOutput: "table",
		},
		{
			name:           "custom access",
			args:           []string{"-id=tprj-custom", "-access=custom", "-output=json"},
			expectedID:     "tprj-custom",
			expectedAccess: "custom",
			expectedOutput: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ProjectTeamAccessUpdateCommand{}

			flags := cmd.Meta.FlagSet("projectteamaccess update")
			flags.StringVar(&cmd.id, "id", "", "Project team access ID (required)")
			flags.StringVar(&cmd.access, "access", "", "Access level")
			flags.StringVar(&cmd.format, "output", "table", "Output format")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify flags were set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
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
