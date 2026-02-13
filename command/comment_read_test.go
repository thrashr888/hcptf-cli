package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestCommentReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &CommentReadCommand{
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

func TestCommentReadHelp(t *testing.T) {
	cmd := &CommentReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf comment read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
}

func TestCommentReadSynopsis(t *testing.T) {
	cmd := &CommentReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show comment details" {
		t.Errorf("expected 'Show comment details', got %q", synopsis)
	}
}

func TestCommentReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "id only, default format",
			args:           []string{"-id=wsc-ABC123"},
			expectedID:     "wsc-ABC123",
			expectedFormat: "table",
		},
		{
			name:           "id with table format",
			args:           []string{"-id=wsc-XYZ789", "-output=table"},
			expectedID:     "wsc-XYZ789",
			expectedFormat: "table",
		},
		{
			name:           "id with json format",
			args:           []string{"-id=wsc-DEF456", "-output=json"},
			expectedID:     "wsc-DEF456",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &CommentReadCommand{}

			flags := cmd.Meta.FlagSet("comment read")
			flags.StringVar(&cmd.id, "id", "", "Comment ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
