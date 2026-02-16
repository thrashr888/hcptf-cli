package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newVariableSetVariableDeleteCommand(ui cli.Ui, svc variableSetVariableDeleter) *VariableSetVariableDeleteCommand {
	return &VariableSetVariableDeleteCommand{
		Meta:                   newTestMeta(ui),
		variableSetVariableSvc: svc,
	}
}

func TestVariableSetVariableDeleteRequiresVariableSetID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := newVariableSetVariableDeleteCommand(ui, &mockVariableSetVariableDeleteService{})

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
	cmd := newVariableSetVariableDeleteCommand(ui, &mockVariableSetVariableDeleteService{})

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
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "-y") {
		t.Error("Help should mention -y flag")
	}
}

func TestVariableSetVariableDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockVariableSetVariableDeleteService{err: errors.New("boom")}
	cmd := newVariableSetVariableDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-variableset-id=varset-123", "-variable-id=var-456", "-force"}); code != 1 {
		t.Fatalf("expected exit 1")
	}
	if svc.lastSetID != "varset-123" || svc.lastVariable != "var-456" {
		t.Fatalf("unexpected ids %q %q", svc.lastSetID, svc.lastVariable)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output")
	}
}

func TestVariableSetVariableDeleteSuccess(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockVariableSetVariableDeleteService{}
	cmd := newVariableSetVariableDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-variableset-id=varset-123", "-variable-id=var-456", "-y"}); code != 0 {
		t.Fatalf("expected exit 0")
	}
	if svc.lastSetID != "varset-123" || svc.lastVariable != "var-456" {
		t.Fatalf("unexpected ids %q %q", svc.lastSetID, svc.lastVariable)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success output")
	}
}

func TestVariableSetVariableDeletePromptsWithoutForce(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("no\n")
	cmd := newVariableSetVariableDeleteCommand(ui, &mockVariableSetVariableDeleteService{})

	if code := cmd.Run([]string{"-variableset-id=varset-123", "-variable-id=var-456"}); code != 0 {
		t.Fatalf("expected exit 0 on cancel, got %d", code)
	}
	if !strings.Contains(ui.OutputWriter.String(), "Deletion cancelled") {
		t.Fatalf("expected cancellation message")
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
