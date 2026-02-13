package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestReservedTagKeyDeleteRequiresId(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ReservedTagKeyDeleteCommand{
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

func TestReservedTagKeyDeleteHelp(t *testing.T) {
	cmd := &ReservedTagKeyDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf reservedtagkey delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "Delete a reserved tag key") {
		t.Error("Help should describe deleting reserved tag key")
	}
}

func TestReservedTagKeyDeleteSynopsis(t *testing.T) {
	cmd := &ReservedTagKeyDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a reserved tag key" {
		t.Errorf("expected 'Delete a reserved tag key', got %q", synopsis)
	}
}

func TestReservedTagKeyDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedId    string
		expectedForce bool
	}{
		{
			name:          "id flag only",
			args:          []string{"-id=rtk-ABC123"},
			expectedId:    "rtk-ABC123",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=rtk-XYZ789", "-force"},
			expectedId:    "rtk-XYZ789",
			expectedForce: true,
		},
		{
			name:          "different id",
			args:          []string{"-id=rtk-DEF456"},
			expectedId:    "rtk-DEF456",
			expectedForce: false,
		},
		{
			name:          "force without explicit value",
			args:          []string{"-id=rtk-GHI123", "-force"},
			expectedId:    "rtk-GHI123",
			expectedForce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ReservedTagKeyDeleteCommand{}

			flags := cmd.Meta.FlagSet("reservedtagkey delete")
			flags.StringVar(&cmd.id, "id", "", "Reserved tag key ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedId {
				t.Errorf("expected id %q, got %q", tt.expectedId, cmd.id)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
