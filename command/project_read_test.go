package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestProjectReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectReadCommand{
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

func TestProjectReadHelp(t *testing.T) {
	cmd := &ProjectReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf project read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestProjectReadSynopsis(t *testing.T) {
	cmd := &ProjectReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read project details" {
		t.Errorf("expected 'Read project details', got %q", synopsis)
	}
}

func TestProjectReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id, default format",
			args:        []string{"-id=prj-abc123"},
			expectedID:  "prj-abc123",
			expectedFmt: "table",
		},
		{
			name:        "id, table format",
			args:        []string{"-id=prj-xyz789", "-output=table"},
			expectedID:  "prj-xyz789",
			expectedFmt: "table",
		},
		{
			name:        "id, json format",
			args:        []string{"-id=prj-test456", "-output=json"},
			expectedID:  "prj-test456",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ProjectReadCommand{}

			flags := cmd.Meta.FlagSet("project read")
			flags.StringVar(&cmd.projectID, "id", "", "Project ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.projectID != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.projectID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
