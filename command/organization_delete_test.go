package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationDeleteRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestOrganizationDeleteHelp(t *testing.T) {
	cmd := &OrganizationDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organization delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -name is required")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "cannot be undone") {
		t.Error("Help should warn that action cannot be undone")
	}
}

func TestOrganizationDeleteSynopsis(t *testing.T) {
	cmd := &OrganizationDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete an organization" {
		t.Errorf("expected 'Delete an organization', got %q", synopsis)
	}
}

func TestOrganizationDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedName  string
		expectedForce bool
	}{
		{
			name:          "name without force",
			args:          []string{"-name=test-org"},
			expectedName:  "test-org",
			expectedForce: false,
		},
		{
			name:          "name with force",
			args:          []string{"-name=old-org", "-force"},
			expectedName:  "old-org",
			expectedForce: true,
		},
		{
			name:          "name with force=true",
			args:          []string{"-name=deprecated-org", "-force=true"},
			expectedName:  "deprecated-org",
			expectedForce: true,
		},
		{
			name:          "name with force=false",
			args:          []string{"-name=keep-org", "-force=false"},
			expectedName:  "keep-org",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationDeleteCommand{}

			flags := cmd.Meta.FlagSet("organization delete")
			flags.StringVar(&cmd.name, "name", "", "Organization name (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
