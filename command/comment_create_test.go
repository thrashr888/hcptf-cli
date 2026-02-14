package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestCommentCreateRequiresRunID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &CommentCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-body=Test comment"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-run-id") {
		t.Fatalf("expected run-id error, got %q", out)
	}
}

func TestCommentCreateRequiresBody(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &CommentCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-run-id=run-ABC123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-body") {
		t.Fatalf("expected body error, got %q", out)
	}
}

func TestCommentCreateRequiresBothFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &CommentCreateCommand{
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

func TestCommentCreateHelp(t *testing.T) {
	cmd := &CommentCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf comment create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-run-id") {
		t.Error("Help should mention -run-id flag")
	}
	if !strings.Contains(help, "-body") {
		t.Error("Help should mention -body flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "Example:") {
		t.Error("Help should contain examples")
	}
}

func TestCommentCreateSynopsis(t *testing.T) {
	cmd := &CommentCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a comment on a run" {
		t.Errorf("expected 'Create a comment on a run', got %q", synopsis)
	}
}

func TestCommentCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedRunID  string
		expectedBody   string
		expectedFormat string
	}{
		{
			name:           "all required flags, default format",
			args:           []string{"-run-id=run-ABC123", "-body=Approved for production"},
			expectedRunID:  "run-ABC123",
			expectedBody:   "Approved for production",
			expectedFormat: "table",
		},
		{
			name:           "run-id and body with json format",
			args:           []string{"-run-id=run-XYZ789", "-body=Need to review security implications", "-output=json"},
			expectedRunID:  "run-XYZ789",
			expectedBody:   "Need to review security implications",
			expectedFormat: "json",
		},
		{
			name:           "run-id and body with table format explicit",
			args:           []string{"-run-id=run-TEST001", "-body=LGTM", "-output=table"},
			expectedRunID:  "run-TEST001",
			expectedBody:   "LGTM",
			expectedFormat: "table",
		},
		{
			name:           "run-id and multi-word body",
			args:           []string{"-run-id=run-MULTI123", "-body=This is a longer comment with multiple words"},
			expectedRunID:  "run-MULTI123",
			expectedBody:   "This is a longer comment with multiple words",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &CommentCreateCommand{}

			flags := cmd.Meta.FlagSet("comment create")
			flags.StringVar(&cmd.runID, "run-id", "", "Run ID (required)")
			flags.StringVar(&cmd.body, "body", "", "Comment body (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the run-id was set correctly
			if cmd.runID != tt.expectedRunID {
				t.Errorf("expected runID %q, got %q", tt.expectedRunID, cmd.runID)
			}

			// Verify the body was set correctly
			if cmd.body != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, cmd.body)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
