package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationMembershipDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationMembershipDeleteCommand{
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

func TestOrganizationMembershipDeleteRequiresFlagMessage(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationMembershipDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run(nil)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "required") {
		t.Fatalf("expected required error, got %q", out)
	}
}

func TestOrganizationMembershipDeleteHelp(t *testing.T) {
	cmd := &OrganizationMembershipDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organizationmembership delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id flag is required")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "Remove a user from an organization") {
		t.Error("Help should contain command description")
	}
	if !strings.Contains(help, "Example:") {
		t.Error("Help should contain examples")
	}
}

func TestOrganizationMembershipDeleteSynopsis(t *testing.T) {
	cmd := &OrganizationMembershipDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Remove a user from an organization" {
		t.Errorf("expected 'Remove a user from an organization', got %q", synopsis)
	}
}

func TestOrganizationMembershipDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id only, no force",
			args:          []string{"-id=ou-abc123xyz"},
			expectedID:    "ou-abc123xyz",
			expectedForce: false,
		},
		{
			name:          "id with different format",
			args:          []string{"-id=ou-test123"},
			expectedID:    "ou-test123",
			expectedForce: false,
		},
		{
			name:          "id with force",
			args:          []string{"-id=ou-abc123xyz", "-force"},
			expectedID:    "ou-abc123xyz",
			expectedForce: true,
		},
		{
			name:          "id with explicit force=true",
			args:          []string{"-id=ou-xyz789", "-force=true"},
			expectedID:    "ou-xyz789",
			expectedForce: true,
		},
		{
			name:          "force before id",
			args:          []string{"-force", "-id=ou-test456"},
			expectedID:    "ou-test456",
			expectedForce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationMembershipDeleteCommand{}

			flags := cmd.Meta.FlagSet("organizationmembership delete")
			flags.StringVar(&cmd.id, "id", "", "Organization membership ID (required)")
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
