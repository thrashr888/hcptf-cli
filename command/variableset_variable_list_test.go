package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVariableSetVariableListRequiresVariableSetID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetVariableListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-variableset-id") {
		t.Fatalf("expected variableset-id error, got %q", out)
	}
}

func TestVariableSetVariableListHelp(t *testing.T) {
	cmd := &VariableSetVariableListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	if !strings.Contains(help, "hcptf variableset variable list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-variableset-id") {
		t.Error("Help should mention -variableset-id flag")
	}
}

func TestVariableSetVariableListSynopsis(t *testing.T) {
	cmd := &VariableSetVariableListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List variables in a variable set" {
		t.Errorf("expected 'List variables in a variable set', got %q", synopsis)
	}
}

func TestVariableSetVariableListFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedVarSetID string
		expectedFormat   string
	}{
		{"default format", []string{"-variableset-id=varset-123"}, "varset-123", "table"},
		{"json format", []string{"-variableset-id=varset-456", "-output=json"}, "varset-456", "json"},
		{"table format", []string{"-variableset-id=varset-789", "-output=table"}, "varset-789", "table"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VariableSetVariableListCommand{}

			flags := cmd.Meta.FlagSet("variableset-variable list")
			flags.StringVar(&cmd.variableSetID, "variableset-id", "", "Variable set ID")
			flags.StringVar(&cmd.format, "output", "table", "Output format")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			if cmd.variableSetID != tt.expectedVarSetID {
				t.Errorf("expected variableSetID %q, got %q", tt.expectedVarSetID, cmd.variableSetID)
			}

			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
