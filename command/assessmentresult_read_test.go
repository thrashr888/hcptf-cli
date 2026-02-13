package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAssessmentResultReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AssessmentResultReadCommand{
		Meta: newTestMeta(ui),
	}

	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error")
	}
}

func TestAssessmentResultReadHelp(t *testing.T) {
	cmd := &AssessmentResultReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf assessmentresult read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "asmtres-") {
		t.Error("Help should mention assessment result ID format")
	}
	if !strings.Contains(help, "Health assessments") {
		t.Error("Help should describe health assessments")
	}
	if !strings.Contains(help, "DriftStatus") {
		t.Error("Help should describe drift status")
	}
	if !strings.Contains(help, "JSONOutput") {
		t.Error("Help should mention JSON output details")
	}
}

func TestAssessmentResultReadSynopsis(t *testing.T) {
	cmd := &AssessmentResultReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show details of a specific health assessment result" {
		t.Errorf("expected 'Show details of a specific health assessment result', got %q", synopsis)
	}
}

func TestAssessmentResultReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id, default format",
			args:        []string{"-id=asmtres-123"},
			expectedID:  "asmtres-123",
			expectedFmt: "table",
		},
		{
			name:        "id, table format",
			args:        []string{"-id=asmtres-abc", "-output=table"},
			expectedID:  "asmtres-abc",
			expectedFmt: "table",
		},
		{
			name:        "id, json format",
			args:        []string{"-id=asmtres-xyz", "-output=json"},
			expectedID:  "asmtres-xyz",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AssessmentResultReadCommand{}

			flags := cmd.Meta.FlagSet("assessmentresult read")
			flags.StringVar(&cmd.id, "id", "", "Assessment result ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the ID was set correctly
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
