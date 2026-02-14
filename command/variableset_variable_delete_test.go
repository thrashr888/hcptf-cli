package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVariableSetVariableDeleteRequiresVariableSetID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetVariableDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-variable-id=var-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-variableset-id") {
		t.Fatalf("expected variableset-id error, got %q", out)
	}
}

func TestVariableSetVariableDeleteRequiresVariableID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetVariableDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-variableset-id=varset-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-variable-id") {
		t.Fatalf("expected variable-id error, got %q", out)
	}
}

func TestVariableSetVariableDeleteHelp(t *testing.T) {
	cmd := &VariableSetVariableDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	if !strings.Contains(help, "hcptf variableset variable delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-variableset-id") {
		t.Error("Help should mention -variableset-id flag")
	}
	if !strings.Contains(help, "-variable-id") {
		t.Error("Help should mention -variable-id flag")
	}
}

func TestVariableSetVariableDeleteSynopsis(t *testing.T) {
	cmd := &VariableSetVariableDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a variable from a variable set" {
		t.Errorf("expected 'Delete a variable from a variable set', got %q", synopsis)
	}
}

func TestVariableSetVariableDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedVarSetID string
		expectedVarID    string
	}{
		{"basic flags", []string{"-variableset-id=varset-123", "-variable-id=var-456"}, "varset-123", "var-456"},
		{"different ids", []string{"-variableset-id=varset-abc", "-variable-id=var-def"}, "varset-abc", "var-def"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VariableSetVariableDeleteCommand{}

			flags := cmd.Meta.FlagSet("variableset-variable delete")
			flags.StringVar(&cmd.variableSetID, "variableset-id", "", "Variable set ID")
			flags.StringVar(&cmd.variableID, "variable-id", "", "Variable ID")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			if cmd.variableSetID != tt.expectedVarSetID {
				t.Errorf("expected variableSetID %q, got %q", tt.expectedVarSetID, cmd.variableSetID)
			}

			if cmd.variableID != tt.expectedVarID {
				t.Errorf("expected variableID %q, got %q", tt.expectedVarID, cmd.variableID)
			}
		})
	}
}
