package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVariableSetDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetDeleteCommand{
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

func TestVariableSetDeleteRequiresEmptyID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetDeleteCommand{
		Meta: newTestMeta(ui),
	}

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
