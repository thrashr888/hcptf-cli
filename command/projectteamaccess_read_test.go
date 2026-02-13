package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestProjectTeamAccessReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectTeamAccessReadCommand{
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

func TestProjectTeamAccessReadHelp(t *testing.T) {
	cmd := &ProjectTeamAccessReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf projectteamaccess read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "table") {
		t.Error("Help should mention table format")
	}
	if !strings.Contains(help, "json") {
		t.Error("Help should mention json format")
	}
	if !strings.Contains(help, "tprj-") {
		t.Error("Help should contain example ID")
	}
}

func TestProjectTeamAccessReadSynopsis(t *testing.T) {
	cmd := &ProjectTeamAccessReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show project team access details" {
		t.Errorf("expected 'Show project team access details', got %q", synopsis)
	}
}

func TestProjectTeamAccessReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedOutput string
	}{
		{
			name:           "required flags only",
			args:           []string{"-id=tprj-123"},
			expectedID:     "tprj-123",
			expectedOutput: "table",
		},
		{
			name:           "with json output",
			args:           []string{"-id=tprj-456", "-output=json"},
			expectedID:     "tprj-456",
			expectedOutput: "json",
		},
		{
			name:           "with table output",
			args:           []string{"-id=tprj-abc", "-output=table"},
			expectedID:     "tprj-abc",
			expectedOutput: "table",
		},
		{
			name:           "different id format",
			args:           []string{"-id=tprj-xyz123abc"},
			expectedID:     "tprj-xyz123abc",
			expectedOutput: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ProjectTeamAccessReadCommand{}

			flags := cmd.Meta.FlagSet("projectteamaccess read")
			flags.StringVar(&cmd.id, "id", "", "Project team access ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify flags were set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}
			if cmd.format != tt.expectedOutput {
				t.Errorf("expected output %q, got %q", tt.expectedOutput, cmd.format)
			}
		})
	}
}
