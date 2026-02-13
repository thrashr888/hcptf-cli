package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAuditTrailTokenCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AuditTrailTokenCreateCommand{
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

func TestAuditTrailTokenCreateHelp(t *testing.T) {
	cmd := &AuditTrailTokenCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf audittrailtoken create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag")
	}
	if !strings.Contains(help, "-expired-at") {
		t.Error("Help should mention -expired-at flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -organization is required")
	}
	if !strings.Contains(help, "audit trail token") {
		t.Error("Help should describe audit trail tokens")
	}
}

func TestAuditTrailTokenCreateSynopsis(t *testing.T) {
	cmd := &AuditTrailTokenCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create an audit trail token for an organization" {
		t.Errorf("expected 'Create an audit trail token for an organization', got %q", synopsis)
	}
}

func TestAuditTrailTokenCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedOrg      string
		expectedExpiredAt string
		expectedFormat   string
	}{
		{
			name:             "organization and default format",
			args:             []string{"-organization=test-org"},
			expectedOrg:      "test-org",
			expectedExpiredAt: "",
			expectedFormat:   "table",
		},
		{
			name:             "org alias and default format",
			args:             []string{"-org=my-org"},
			expectedOrg:      "my-org",
			expectedExpiredAt: "",
			expectedFormat:   "table",
		},
		{
			name:             "organization with expiration",
			args:             []string{"-organization=test-org", "-expired-at=2025-12-31T23:59:59.000Z"},
			expectedOrg:      "test-org",
			expectedExpiredAt: "2025-12-31T23:59:59.000Z",
			expectedFormat:   "table",
		},
		{
			name:             "organization with json format",
			args:             []string{"-org=prod-org", "-output=json"},
			expectedOrg:      "prod-org",
			expectedExpiredAt: "",
			expectedFormat:   "json",
		},
		{
			name:             "all flags",
			args:             []string{"-org=full-org", "-expired-at=2026-01-01T00:00:00.000Z", "-output=json"},
			expectedOrg:      "full-org",
			expectedExpiredAt: "2026-01-01T00:00:00.000Z",
			expectedFormat:   "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AuditTrailTokenCreateCommand{}

			flags := cmd.Meta.FlagSet("audittrailtoken create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.expiredAt, "expired-at", "", "Token expiration date (ISO8601 format)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the expiredAt was set correctly
			if cmd.expiredAt != tt.expectedExpiredAt {
				t.Errorf("expected expiredAt %q, got %q", tt.expectedExpiredAt, cmd.expiredAt)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
