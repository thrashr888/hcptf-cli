package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationMembershipReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationMembershipReadCommand{
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

func TestOrganizationMembershipReadRequiresFlagMessage(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationMembershipReadCommand{
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

func TestOrganizationMembershipReadHelp(t *testing.T) {
	cmd := &OrganizationMembershipReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organizationmembership read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id flag is required")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "Show organization membership details") {
		t.Error("Help should contain command description")
	}
	if !strings.Contains(help, "Example:") {
		t.Error("Help should contain examples")
	}
}

func TestOrganizationMembershipReadSynopsis(t *testing.T) {
	cmd := &OrganizationMembershipReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show organization membership details" {
		t.Errorf("expected 'Show organization membership details', got %q", synopsis)
	}
}

func TestOrganizationMembershipReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "id only with default format",
			args:           []string{"-id=ou-abc123xyz"},
			expectedID:     "ou-abc123xyz",
			expectedFormat: "table",
		},
		{
			name:           "id with different format",
			args:           []string{"-id=ou-test123"},
			expectedID:     "ou-test123",
			expectedFormat: "table",
		},
		{
			name:           "id with table format",
			args:           []string{"-id=ou-abc123xyz", "-output=table"},
			expectedID:     "ou-abc123xyz",
			expectedFormat: "table",
		},
		{
			name:           "id with json format",
			args:           []string{"-id=ou-xyz789", "-output=json"},
			expectedID:     "ou-xyz789",
			expectedFormat: "json",
		},
		{
			name:           "output before id",
			args:           []string{"-output=json", "-id=ou-test456"},
			expectedID:     "ou-test456",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationMembershipReadCommand{}

			flags := cmd.Meta.FlagSet("organizationmembership read")
			flags.StringVar(&cmd.id, "id", "", "Organization membership ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
