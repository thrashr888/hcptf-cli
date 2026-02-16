package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newVariableSetDeleteCommand(ui cli.Ui, svc variableSetDeleter) *VariableSetDeleteCommand {
	return &VariableSetDeleteCommand{
		Meta:           newTestMeta(ui),
		variableSetSvc: svc,
	}
}

func TestVariableSetDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newVariableSetDeleteCommand(ui, &mockVariableSetDeleteService{})

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestVariableSetDeleteRequiresEmptyID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newVariableSetDeleteCommand(ui, &mockVariableSetDeleteService{})

	code := cmd.Run([]string{"-id="})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestVariableSetDeleteHelp(t *testing.T) {
	cmd := &VariableSetDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf variableset delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "-y") {
		t.Error("Help should mention -y flag")
	}
}

func TestVariableSetDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockVariableSetDeleteService{err: errors.New("boom")}
	cmd := newVariableSetDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=varset-123", "-force"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastID != "varset-123" {
		t.Fatalf("expected delete id varset-123, got %q", svc.lastID)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestVariableSetDeleteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockVariableSetDeleteService{}
	cmd := newVariableSetDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-id=varset-123", "-force"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastID != "varset-123" {
		t.Fatalf("expected delete id varset-123, got %q", svc.lastID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success output")
	}
}

func TestVariableSetDeletePromptsWithoutForce(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("no\n")
	cmd := newVariableSetDeleteCommand(ui, &mockVariableSetDeleteService{})

	if code := cmd.Run([]string{"-id=varset-123"}); code != 0 {
		t.Fatalf("expected exit 0 on cancel, got %d", code)
	}
	if !strings.Contains(ui.OutputWriter.String(), "Deletion cancelled") {
		t.Fatalf("expected cancellation message")
	}
}

func TestVariableSetDeleteSynopsis(t *testing.T) {
	cmd := &VariableSetDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a variable set" {
		t.Errorf("expected 'Delete a variable set', got %q", synopsis)
	}
}

func TestVariableSetDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		expectedID string
	}{
		{
			name:       "delete with id",
			args:       []string{"-id=varset-12345"},
			expectedID: "varset-12345",
		},
		{
			name:       "delete with different id",
			args:       []string{"-id=varset-abc123"},
			expectedID: "varset-abc123",
		},
		{
			name:       "delete with long id",
			args:       []string{"-id=varset-xyz789def456"},
			expectedID: "varset-xyz789def456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VariableSetDeleteCommand{}

			flags := cmd.Meta.FlagSet("variableset delete")
			flags.StringVar(&cmd.id, "id", "", "Variable set ID (required)")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}
		})
	}
}
