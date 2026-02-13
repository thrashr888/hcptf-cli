package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicySetDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-force"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestPolicySetDeleteHelp(t *testing.T) {
	cmd := &PolicySetDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policyset delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestPolicySetDeleteSynopsis(t *testing.T) {
	cmd := &PolicySetDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a policy set" {
		t.Errorf("expected 'Delete a policy set', got %q", synopsis)
	}
}

func TestPolicySetDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id only, no force",
			args:          []string{"-id=polset-12345"},
			expectedID:    "polset-12345",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=polset-67890", "-force"},
			expectedID:    "polset-67890",
			expectedForce: true,
		},
		{
			name:          "id with force=true",
			args:          []string{"-id=polset-abcde", "-force=true"},
			expectedID:    "polset-abcde",
			expectedForce: true,
		},
		{
			name:          "id with force=false",
			args:          []string{"-id=polset-xyz", "-force=false"},
			expectedID:    "polset-xyz",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicySetDeleteCommand{}

			flags := cmd.Meta.FlagSet("policyset delete")
			flags.StringVar(&cmd.id, "id", "", "Policy set ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
