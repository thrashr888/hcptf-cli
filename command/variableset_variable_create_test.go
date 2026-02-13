package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVariableSetVariableCreateRequiresVariableSetID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetVariableCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-key=test", "-value=val"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-variableset-id") {
		t.Fatalf("expected variableset-id error, got %q", out)
	}
}

func TestVariableSetVariableCreateRequiresKey(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetVariableCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-variableset-id=varset-123", "-value=val"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-key") {
		t.Fatalf("expected key error, got %q", out)
	}
}

func TestVariableSetVariableCreateRequiresValue(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetVariableCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-variableset-id=varset-123", "-key=test"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-value") {
		t.Fatalf("expected value error, got %q", out)
	}
}

func TestVariableSetVariableCreateHelp(t *testing.T) {
	cmd := &VariableSetVariableCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	if !strings.Contains(help, "hcptf variableset variable create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-variableset-id") {
		t.Error("Help should mention -variableset-id flag")
	}
	if !strings.Contains(help, "-key") {
		t.Error("Help should mention -key flag")
	}
	if !strings.Contains(help, "-value") {
		t.Error("Help should mention -value flag")
	}
}

func TestVariableSetVariableCreateSynopsis(t *testing.T) {
	cmd := &VariableSetVariableCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Add a variable to a variable set" {
		t.Errorf("expected 'Add a variable to a variable set', got %q", synopsis)
	}
}

func TestVariableSetVariableCreateValidatesCategory(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetVariableCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-variableset-id=varset-123", "-key=test", "-value=val", "-category=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "category") && !strings.Contains(out, "'terraform' or 'env'") {
		t.Fatalf("expected category validation error, got %q", out)
	}
}

func TestVariableSetVariableCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedVarSetID  string
		expectedKey       string
		expectedValue     string
		expectedCategory  string
		expectedSensitive bool
	}{
		{
			name:              "all required flags",
			args:              []string{"-variableset-id=varset-123", "-key=FOO", "-value=bar"},
			expectedVarSetID:  "varset-123",
			expectedKey:       "FOO",
			expectedValue:     "bar",
			expectedCategory:  "terraform",
			expectedSensitive: false,
		},
		{
			name:              "env category",
			args:              []string{"-variableset-id=varset-456", "-key=PATH", "-value=/usr/bin", "-category=env"},
			expectedVarSetID:  "varset-456",
			expectedKey:       "PATH",
			expectedValue:     "/usr/bin",
			expectedCategory:  "env",
			expectedSensitive: false,
		},
		{
			name:              "sensitive variable",
			args:              []string{"-variableset-id=varset-789", "-key=SECRET", "-value=xxx", "-sensitive=true"},
			expectedVarSetID:  "varset-789",
			expectedKey:       "SECRET",
			expectedValue:     "xxx",
			expectedCategory:  "terraform",
			expectedSensitive: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VariableSetVariableCreateCommand{}

			flags := cmd.Meta.FlagSet("variableset-variable create")
			flags.StringVar(&cmd.variableSetID, "variableset-id", "", "Variable set ID")
			flags.StringVar(&cmd.key, "key", "", "Variable key")
			flags.StringVar(&cmd.value, "value", "", "Variable value")
			flags.StringVar(&cmd.category, "category", "terraform", "Category")
			flags.BoolVar(&cmd.sensitive, "sensitive", false, "Sensitive")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			if cmd.variableSetID != tt.expectedVarSetID {
				t.Errorf("expected variableSetID %q, got %q", tt.expectedVarSetID, cmd.variableSetID)
			}

			if cmd.key != tt.expectedKey {
				t.Errorf("expected key %q, got %q", tt.expectedKey, cmd.key)
			}

			if cmd.value != tt.expectedValue {
				t.Errorf("expected value %q, got %q", tt.expectedValue, cmd.value)
			}

			if cmd.category != tt.expectedCategory {
				t.Errorf("expected category %q, got %q", tt.expectedCategory, cmd.category)
			}

			if cmd.sensitive != tt.expectedSensitive {
				t.Errorf("expected sensitive %v, got %v", tt.expectedSensitive, cmd.sensitive)
			}
		})
	}
}
