package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicyDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicyDeleteCommand{
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

func TestPolicyDeleteHelp(t *testing.T) {
	cmd := &PolicyDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policy delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
}

func TestPolicyDeleteSynopsis(t *testing.T) {
	cmd := &PolicyDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a policy" {
		t.Errorf("expected 'Delete a policy', got %q", synopsis)
	}
}

func TestPolicyDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id only, default force",
			args:          []string{"-id=pol-abc123"},
			expectedID:    "pol-abc123",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=pol-xyz789", "-force"},
			expectedID:    "pol-xyz789",
			expectedForce: true,
		},
		{
			name:          "id with force=true",
			args:          []string{"-id=pol-def456", "-force=true"},
			expectedID:    "pol-def456",
			expectedForce: true,
		},
		{
			name:          "id with force=false",
			args:          []string{"-id=pol-ghi789", "-force=false"},
			expectedID:    "pol-ghi789",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicyDeleteCommand{}

			flags := cmd.Meta.FlagSet("policy delete")
			flags.StringVar(&cmd.policyID, "id", "", "Policy ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the policy ID was set correctly
			if cmd.policyID != tt.expectedID {
				t.Errorf("expected policyID %q, got %q", tt.expectedID, cmd.policyID)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
