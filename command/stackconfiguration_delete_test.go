package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackConfigurationDeleteReturnsError(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackConfigurationDeleteCommand{
		Meta: newTestMeta(ui),
	}

	// Test with no arguments
	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "cannot be deleted") {
		t.Fatalf("expected delete error message, got %q", out)
	}
}

func TestStackConfigurationDeleteReturnsErrorWithID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackConfigurationDeleteCommand{
		Meta: newTestMeta(ui),
	}

	// Test with ID provided - should still error because deletion is not supported
	code := cmd.Run([]string{"-id=stc-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	out := ui.ErrorWriter.String()
	if !strings.Contains(out, "cannot be deleted") {
		t.Fatalf("expected delete error message, got %q", out)
	}
	if !strings.Contains(out, "managed by HCP Terraform") {
		t.Fatalf("expected managed message, got %q", out)
	}
}

func TestStackConfigurationDeleteReturnsErrorWithForce(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackConfigurationDeleteCommand{
		Meta: newTestMeta(ui),
	}

	// Test with force flag - should still error
	code := cmd.Run([]string{"-id=stc-123", "-force"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	out := ui.ErrorWriter.String()
	if !strings.Contains(out, "cannot be deleted") {
		t.Fatalf("expected delete error message, got %q", out)
	}
}

func TestStackConfigurationDeleteHelp(t *testing.T) {
	cmd := &StackConfigurationDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stackconfiguration delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "cannot be deleted") {
		t.Error("Help should explain configurations cannot be deleted")
	}
	if !strings.Contains(help, "stack delete") {
		t.Error("Help should mention alternative command")
	}
}

func TestStackConfigurationDeleteSynopsis(t *testing.T) {
	cmd := &StackConfigurationDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete stack configuration (not supported - managed by HCP Terraform)" {
		t.Errorf("expected 'Delete stack configuration (not supported - managed by HCP Terraform)', got %q", synopsis)
	}
}

func TestStackConfigurationDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id, default force",
			args:          []string{"-id=stc-123"},
			expectedID:    "stc-123",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=stc-456", "-force"},
			expectedID:    "stc-456",
			expectedForce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackConfigurationDeleteCommand{}

			flags := cmd.Meta.FlagSet("stackconfiguration delete")
			flags.StringVar(&cmd.configID, "id", "", "Stack configuration ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the config ID was set correctly
			if cmd.configID != tt.expectedID {
				t.Errorf("expected config ID %q, got %q", tt.expectedID, cmd.configID)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
