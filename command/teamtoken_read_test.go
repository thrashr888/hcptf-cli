package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestTeamTokenReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamTokenReadCommand{
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

func TestTeamTokenReadHelp(t *testing.T) {
	cmd := &TeamTokenReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf teamtoken read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestTeamTokenReadSynopsis(t *testing.T) {
	cmd := &TeamTokenReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show team token details" {
		t.Errorf("expected 'Show team token details', got %q", synopsis)
	}
}

func TestTeamTokenReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id, default format",
			args:        []string{"-id=at-123abc"},
			expectedID:  "at-123abc",
			expectedFmt: "table",
		},
		{
			name:        "id, table format",
			args:        []string{"-id=at-456def", "-output=table"},
			expectedID:  "at-456def",
			expectedFmt: "table",
		},
		{
			name:        "id, json format",
			args:        []string{"-id=at-789ghi", "-output=json"},
			expectedID:  "at-789ghi",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &TeamTokenReadCommand{}

			flags := cmd.Meta.FlagSet("teamtoken read")
			flags.StringVar(&cmd.id, "id", "", "Team token ID (required)")
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
