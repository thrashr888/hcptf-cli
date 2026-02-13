package command

import (
	"strings"
	"testing"
)

func TestCommentListHelp(t *testing.T) {
	cmd := &CommentListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf comment list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-run-id") {
		t.Error("Help should mention -run-id flag")
	}
	if !strings.Contains(help, "Run ID") {
		t.Error("Help should mention run ID")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "table (default) or json") {
		t.Error("Help should mention output formats")
	}
	if !strings.Contains(help, "team members") {
		t.Error("Help should mention team members")
	}
}

func TestCommentListSynopsis(t *testing.T) {
	cmd := &CommentListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List comments for a run" {
		t.Errorf("expected 'List comments for a run', got %q", synopsis)
	}
}

func TestCommentListFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedRunID  string
		expectedFormat string
	}{
		{
			name:           "default output format",
			args:           []string{"-run-id=run-ABC123"},
			expectedRunID:  "run-ABC123",
			expectedFormat: "table",
		},
		{
			name:           "json output format",
			args:           []string{"-run-id=run-XYZ789", "-output=json"},
			expectedRunID:  "run-XYZ789",
			expectedFormat: "json",
		},
		{
			name:           "table output format explicitly set",
			args:           []string{"-run-id=run-TEST456", "-output=table"},
			expectedRunID:  "run-TEST456",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &CommentListCommand{}

			flags := cmd.Meta.FlagSet("comment list")
			flags.StringVar(&cmd.runID, "run-id", "", "Run ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the runID was set correctly
			if cmd.runID != tt.expectedRunID {
				t.Errorf("expected runID %q, got %q", tt.expectedRunID, cmd.runID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
