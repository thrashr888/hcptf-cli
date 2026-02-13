package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicySetOutcomeReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetOutcomeReadCommand{
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

func TestPolicySetOutcomeReadHelp(t *testing.T) {
	cmd := &PolicySetOutcomeReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policysetoutcome read") {
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
	if !strings.Contains(help, "policy set outcome") {
		t.Error("Help should describe policy set outcome details")
	}
}

func TestPolicySetOutcomeReadSynopsis(t *testing.T) {
	cmd := &PolicySetOutcomeReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read policy set outcome details" {
		t.Errorf("expected 'Read policy set outcome details', got %q", synopsis)
	}
}

func TestPolicySetOutcomeReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id flag",
			args:        []string{"-id=psout-abc123"},
			expectedID:  "psout-abc123",
			expectedFmt: "table",
		},
		{
			name:        "with json output",
			args:        []string{"-id=psout-xyz789", "-output=json"},
			expectedID:  "psout-xyz789",
			expectedFmt: "json",
		},
		{
			name:        "with table output",
			args:        []string{"-id=psout-test456", "-output=table"},
			expectedID:  "psout-test456",
			expectedFmt: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicySetOutcomeReadCommand{}

			flags := cmd.Meta.FlagSet("policysetoutcome read")
			flags.StringVar(&cmd.policySetOutcomeID, "id", "", "Policy Set Outcome ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the ID was set correctly
			if cmd.policySetOutcomeID != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.policySetOutcomeID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
