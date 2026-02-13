package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationMemberReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationMemberReadCommand{
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

func TestOrganizationMemberReadHelp(t *testing.T) {
	cmd := &OrganizationMemberReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organizationmember read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "organization member") {
		t.Error("Help should describe organization member information")
	}
}

func TestOrganizationMemberReadSynopsis(t *testing.T) {
	cmd := &OrganizationMemberReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show detailed organization member information" {
		t.Errorf("expected 'Show detailed organization member information', got %q", synopsis)
	}
}

func TestOrganizationMemberReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id flag",
			args:        []string{"-id=ou-abc123xyz"},
			expectedID:  "ou-abc123xyz",
			expectedFmt: "table",
		},
		{
			name:        "with json output",
			args:        []string{"-id=ou-test123", "-output=json"},
			expectedID:  "ou-test123",
			expectedFmt: "json",
		},
		{
			name:        "with table output",
			args:        []string{"-id=ou-xyz789", "-output=table"},
			expectedID:  "ou-xyz789",
			expectedFmt: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationMemberReadCommand{}

			flags := cmd.Meta.FlagSet("organizationmember read")
			flags.StringVar(&cmd.id, "id", "", "Organization membership ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the ID was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
