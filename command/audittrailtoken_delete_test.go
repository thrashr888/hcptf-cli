package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAuditTrailTokenDeleteRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AuditTrailTokenDeleteCommand{
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

func TestAuditTrailTokenDeleteHelp(t *testing.T) {
	cmd := &AuditTrailTokenDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf audittrailtoken delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -organization is required")
	}
	if !strings.Contains(help, "audit trail token") {
		t.Error("Help should describe audit trail tokens")
	}
}

func TestAuditTrailTokenDeleteSynopsis(t *testing.T) {
	cmd := &AuditTrailTokenDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete an audit trail token" {
		t.Errorf("expected 'Delete an audit trail token', got %q", synopsis)
	}
}

func TestAuditTrailTokenDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedOrg   string
		expectedForce bool
	}{
		{
			name:          "organization without force",
			args:          []string{"-organization=test-org"},
			expectedOrg:   "test-org",
			expectedForce: false,
		},
		{
			name:          "org alias without force",
			args:          []string{"-org=my-org"},
			expectedOrg:   "my-org",
			expectedForce: false,
		},
		{
			name:          "organization with force",
			args:          []string{"-organization=old-org", "-force"},
			expectedOrg:   "old-org",
			expectedForce: true,
		},
		{
			name:          "organization with force=true",
			args:          []string{"-org=deprecated-org", "-force=true"},
			expectedOrg:   "deprecated-org",
			expectedForce: true,
		},
		{
			name:          "organization with force=false",
			args:          []string{"-org=keep-org", "-force=false"},
			expectedOrg:   "keep-org",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AuditTrailTokenDeleteCommand{}

			flags := cmd.Meta.FlagSet("audittrailtoken delete")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
