package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicyCheckReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicyCheckReadCommand{
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

func TestPolicyCheckReadHelp(t *testing.T) {
	cmd := &PolicyCheckReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policycheck read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
}

func TestPolicyCheckReadSynopsis(t *testing.T) {
	cmd := &PolicyCheckReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read policy check details" {
		t.Errorf("expected 'Read policy check details', got %q", synopsis)
	}
}

func TestPolicyCheckReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "policy check id, default format",
			args:           []string{"-id=polchk-abc123"},
			expectedID:     "polchk-abc123",
			expectedFormat: "table",
		},
		{
			name:           "policy check id, table format",
			args:           []string{"-id=polchk-xyz789", "-output=table"},
			expectedID:     "polchk-xyz789",
			expectedFormat: "table",
		},
		{
			name:           "policy check id, json format",
			args:           []string{"-id=polchk-def456", "-output=json"},
			expectedID:     "polchk-def456",
			expectedFormat: "json",
		},
		{
			name:           "different policy check id format",
			args:           []string{"-id=polchk-prod999", "-output=table"},
			expectedID:     "polchk-prod999",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicyCheckReadCommand{}

			flags := cmd.Meta.FlagSet("policycheck read")
			flags.StringVar(&cmd.policyCheckID, "id", "", "Policy Check ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the policy check ID was set correctly
			if cmd.policyCheckID != tt.expectedID {
				t.Errorf("expected policyCheckID %q, got %q", tt.expectedID, cmd.policyCheckID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
