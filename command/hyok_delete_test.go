package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestHYOKDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKDeleteCommand{
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

func TestHYOKDeletePromptsWithoutForce(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("no\n")
	cmd := &HYOKDeleteCommand{
		Meta: newTestMeta(ui),
	}
	cmd.Meta.Ui = ui

	code := cmd.Run([]string{"-id=hyokc-123456"})
	if code != 0 {
		t.Fatalf("expected exit 0 when cancelled, got %d", code)
	}

	if !strings.Contains(ui.OutputWriter.String(), "Deletion cancelled") {
		t.Fatalf("expected cancellation message")
	}
}

func TestHYOKDeleteHelp(t *testing.T) {
	cmd := &HYOKDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf hyok delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "cannot be undone") {
		t.Error("Help should warn that deletion cannot be undone")
	}
	if !strings.Contains(help, "HYOK") {
		t.Error("Help should mention HYOK")
	}
}

func TestHYOKDeleteSynopsis(t *testing.T) {
	cmd := &HYOKDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a HYOK configuration" {
		t.Errorf("expected 'Delete a HYOK configuration', got %q", synopsis)
	}
}

func TestHYOKDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id only, no force",
			args:          []string{"-id=hyokc-123456"},
			expectedID:    "hyokc-123456",
			expectedForce: false,
		},
		{
			name:          "id with force",
			args:          []string{"-id=hyokc-abcdef", "-force"},
			expectedID:    "hyokc-abcdef",
			expectedForce: true,
		},
		{
			name:          "different id with force",
			args:          []string{"-id=hyokc-xyz789", "-force=true"},
			expectedID:    "hyokc-xyz789",
			expectedForce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &HYOKDeleteCommand{}

			flags := cmd.Meta.FlagSet("hyok delete")
			flags.StringVar(&cmd.id, "id", "", "HYOK configuration ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the force was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
