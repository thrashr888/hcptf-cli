package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicySetOutcomeListRequiresPolicyEvaluationID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetOutcomeListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-policy-evaluation-id") {
		t.Fatalf("expected policy-evaluation-id error, got %q", out)
	}
}

func TestPolicySetOutcomeListHelp(t *testing.T) {
	cmd := &PolicySetOutcomeListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policysetoutcome list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-policy-evaluation-id") {
		t.Error("Help should mention -policy-evaluation-id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "policy set outcome") {
		t.Error("Help should describe policy set outcomes")
	}
}

func TestPolicySetOutcomeListSynopsis(t *testing.T) {
	cmd := &PolicySetOutcomeListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List policy set outcomes for a policy evaluation" {
		t.Errorf("expected 'List policy set outcomes for a policy evaluation', got %q", synopsis)
	}
}

func TestPolicySetOutcomeListFlagParsing(t *testing.T) {
	tests := []struct {
		name                   string
		args                   []string
		expectedPolicyEvalID   string
		expectedFmt            string
	}{
		{
			name:                 "policy-evaluation-id flag",
			args:                 []string{"-policy-evaluation-id=poleval-abc123"},
			expectedPolicyEvalID: "poleval-abc123",
			expectedFmt:          "table",
		},
		{
			name:                 "with json output",
			args:                 []string{"-policy-evaluation-id=poleval-xyz789", "-output=json"},
			expectedPolicyEvalID: "poleval-xyz789",
			expectedFmt:          "json",
		},
		{
			name:                 "with table output",
			args:                 []string{"-policy-evaluation-id=poleval-test456", "-output=table"},
			expectedPolicyEvalID: "poleval-test456",
			expectedFmt:          "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicySetOutcomeListCommand{}

			flags := cmd.Meta.FlagSet("policysetoutcome list")
			flags.StringVar(&cmd.policyEvaluationID, "policy-evaluation-id", "", "Policy Evaluation ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the policy evaluation ID was set correctly
			if cmd.policyEvaluationID != tt.expectedPolicyEvalID {
				t.Errorf("expected policy-evaluation-id %q, got %q", tt.expectedPolicyEvalID, cmd.policyEvaluationID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
