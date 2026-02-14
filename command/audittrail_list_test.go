package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAuditTrailListRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AuditTrailListCommand{
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

func TestAuditTrailListHelp(t *testing.T) {
	cmd := &AuditTrailListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf audittrail list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -organization is required")
	}
	if !strings.Contains(help, "-since") {
		t.Error("Help should mention -since flag")
	}
	if !strings.Contains(help, "-page-number") {
		t.Error("Help should mention -page-number flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
}

func TestAuditTrailListSynopsis(t *testing.T) {
	cmd := &AuditTrailListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List audit trail events for an organization" {
		t.Errorf("expected 'List audit trail events for an organization', got %q", synopsis)
	}
}

func TestAuditTrailListFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedOrg      string
		expectedSince    string
		expectedPageNum  int
		expectedPageSize int
		expectedFormat   string
	}{
		{
			name:             "organization with defaults",
			args:             []string{"-organization=my-org"},
			expectedOrg:      "my-org",
			expectedSince:    "",
			expectedPageNum:  1,
			expectedPageSize: 100,
			expectedFormat:   "table",
		},
		{
			name:             "org alias",
			args:             []string{"-org=my-org"},
			expectedOrg:      "my-org",
			expectedSince:    "",
			expectedPageNum:  1,
			expectedPageSize: 100,
			expectedFormat:   "table",
		},
		{
			name:             "with since date",
			args:             []string{"-organization=my-org", "-since=2024-01-01T00:00:00.000Z"},
			expectedOrg:      "my-org",
			expectedSince:    "2024-01-01T00:00:00.000Z",
			expectedPageNum:  1,
			expectedPageSize: 100,
			expectedFormat:   "table",
		},
		{
			name:             "with pagination",
			args:             []string{"-organization=my-org", "-page-number=2", "-page-size=50"},
			expectedOrg:      "my-org",
			expectedSince:    "",
			expectedPageNum:  2,
			expectedPageSize: 50,
			expectedFormat:   "table",
		},
		{
			name:             "json output",
			args:             []string{"-organization=my-org", "-output=json"},
			expectedOrg:      "my-org",
			expectedSince:    "",
			expectedPageNum:  1,
			expectedPageSize: 100,
			expectedFormat:   "json",
		},
		{
			name:             "all flags",
			args:             []string{"-organization=my-org", "-since=2024-01-01T00:00:00.000Z", "-page-number=3", "-page-size=25", "-output=json"},
			expectedOrg:      "my-org",
			expectedSince:    "2024-01-01T00:00:00.000Z",
			expectedPageNum:  3,
			expectedPageSize: 25,
			expectedFormat:   "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AuditTrailListCommand{}

			flags := cmd.Meta.FlagSet("audittrail list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.since, "since", "", "Return audit events since this date")
			flags.IntVar(&cmd.pageNumber, "page-number", 1, "Page number")
			flags.IntVar(&cmd.pageSize, "page-size", 100, "Number of items per page")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the since was set correctly
			if cmd.since != tt.expectedSince {
				t.Errorf("expected since %q, got %q", tt.expectedSince, cmd.since)
			}

			// Verify the page number was set correctly
			if cmd.pageNumber != tt.expectedPageNum {
				t.Errorf("expected pageNumber %d, got %d", tt.expectedPageNum, cmd.pageNumber)
			}

			// Verify the page size was set correctly
			if cmd.pageSize != tt.expectedPageSize {
				t.Errorf("expected pageSize %d, got %d", tt.expectedPageSize, cmd.pageSize)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
