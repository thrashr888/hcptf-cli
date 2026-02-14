package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationTagListRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationTagListCommand{
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

func TestOrganizationTagListHelp(t *testing.T) {
	cmd := &OrganizationTagListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organizationtag list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "tag") {
		t.Error("Help should describe tags")
	}
}

func TestOrganizationTagListSynopsis(t *testing.T) {
	cmd := &OrganizationTagListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List organization tags" {
		t.Errorf("expected 'List organization tags', got %q", synopsis)
	}
}

func TestOrganizationTagListFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOrg string
		expectedFmt string
	}{
		{
			name:        "organization flag",
			args:        []string{"-organization=my-org"},
			expectedOrg: "my-org",
			expectedFmt: "table",
		},
		{
			name:        "org alias flag",
			args:        []string{"-org=test-org"},
			expectedOrg: "test-org",
			expectedFmt: "table",
		},
		{
			name:        "with json output",
			args:        []string{"-organization=my-org", "-output=json"},
			expectedOrg: "my-org",
			expectedFmt: "json",
		},
		{
			name:        "with table output",
			args:        []string{"-org=test-org", "-output=table"},
			expectedOrg: "test-org",
			expectedFmt: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationTagListCommand{}

			flags := cmd.Meta.FlagSet("organizationtag list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
