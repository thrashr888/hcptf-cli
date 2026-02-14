package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVariableSetVariableUpdateRequiresVariableSetID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetVariableUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-variable-id=var-123", "-key=test"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-variableset-id") {
		t.Fatalf("expected variableset-id error, got %q", out)
	}
}

func TestVariableSetVariableUpdateRequiresVariableID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetVariableUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-variableset-id=varset-123", "-key=test"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-variable-id") {
		t.Fatalf("expected variable-id error, got %q", out)
	}
}

func TestVariableSetVariableUpdateHelp(t *testing.T) {
	cmd := &VariableSetVariableUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	if !strings.Contains(help, "hcptf variableset variable update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-variableset-id") {
		t.Error("Help should mention -variableset-id flag")
	}
	if !strings.Contains(help, "-variable-id") {
		t.Error("Help should mention -variable-id flag")
	}
}

func TestVariableSetVariableUpdateSynopsis(t *testing.T) {
	cmd := &VariableSetVariableUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update a variable in a variable set" {
		t.Errorf("expected 'Update a variable in a variable set', got %q", synopsis)
	}
}

func TestVariableSetVariableUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedVarSetID string
		expectedVarID    string
		expectedKey      string
	}{
		{"basic flags", []string{"-variableset-id=varset-123", "-variable-id=var-456", "-key=FOO"}, "varset-123", "var-456", "FOO"},
		{"different ids", []string{"-variableset-id=varset-abc", "-variable-id=var-def", "-key=BAR"}, "varset-abc", "var-def", "BAR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VariableSetVariableUpdateCommand{}

			flags := cmd.Meta.FlagSet("variableset-variable update")
			flags.StringVar(&cmd.variableSetID, "variableset-id", "", "Variable set ID")
			flags.StringVar(&cmd.variableID, "variable-id", "", "Variable ID")
			flags.StringVar(&cmd.key, "key", "", "Variable key")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			if cmd.variableSetID != tt.expectedVarSetID {
				t.Errorf("expected variableSetID %q, got %q", tt.expectedVarSetID, cmd.variableSetID)
			}

			if cmd.variableID != tt.expectedVarID {
				t.Errorf("expected variableID %q, got %q", tt.expectedVarID, cmd.variableID)
			}

			if cmd.key != tt.expectedKey {
				t.Errorf("expected key %q, got %q", tt.expectedKey, cmd.key)
			}
		})
	}
}
