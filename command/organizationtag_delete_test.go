package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationTagDeleteRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationTagDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestOrganizationTagDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationTagDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestOrganizationTagDeleteRequiresEmptyOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationTagDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=", "-id=tag-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestOrganizationTagDeleteRequiresEmptyID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationTagDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-id="})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestOrganizationTagDeleteHelp(t *testing.T) {
	cmd := &OrganizationTagDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organizationtag delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
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
	if !strings.Contains(help, "Delete") {
		t.Error("Help should describe deletion")
	}
}

func TestOrganizationTagDeleteSynopsis(t *testing.T) {
	cmd := &OrganizationTagDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete an organization tag" {
		t.Errorf("expected 'Delete an organization tag', got %q", synopsis)
	}
}

func TestOrganizationTagDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedOrg   string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "organization and id flags",
			args:          []string{"-organization=my-org", "-id=tag-ABC123"},
			expectedOrg:   "my-org",
			expectedID:    "tag-ABC123",
			expectedForce: false,
		},
		{
			name:          "org alias flag",
			args:          []string{"-org=test-org", "-id=tag-XYZ789"},
			expectedOrg:   "test-org",
			expectedID:    "tag-XYZ789",
			expectedForce: false,
		},
		{
			name:          "with force flag",
			args:          []string{"-organization=my-org", "-id=tag-DEF456", "-force"},
			expectedOrg:   "my-org",
			expectedID:    "tag-DEF456",
			expectedForce: true,
		},
		{
			name:          "org alias with force",
			args:          []string{"-org=test-org", "-id=tag-GHI789", "-force=true"},
			expectedOrg:   "test-org",
			expectedID:    "tag-GHI789",
			expectedForce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationTagDeleteCommand{}

			flags := cmd.Meta.FlagSet("organizationtag delete")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.id, "id", "", "Tag ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the ID was set correctly
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
