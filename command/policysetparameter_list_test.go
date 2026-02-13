package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicySetParameterListRequiresPolicySetID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-policy-set-id") {
		t.Fatalf("expected policy-set-id error, got %q", out)
	}
}

func TestPolicySetParameterListRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "required") {
		t.Fatalf("expected required error, got %q", out)
	}
}

func TestPolicySetParameterListHelp(t *testing.T) {
	cmd := &PolicySetParameterListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policysetparameter list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-policy-set-id") {
		t.Error("Help should mention -policy-set-id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestPolicySetParameterListSynopsis(t *testing.T) {
	cmd := &PolicySetParameterListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List parameters for a policy set" {
		t.Errorf("expected 'List parameters for a policy set', got %q", synopsis)
	}
}

func TestPolicySetParameterListFlagParsing(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedPolicySet string
		expectedFormat    string
	}{
		{
			name:              "required flags only, default values",
			args:              []string{"-policy-set-id=polset-abc123"},
			expectedPolicySet: "polset-abc123",
			expectedFormat:    "table",
		},
		{
			name:              "with json output",
			args:              []string{"-policy-set-id=polset-xyz789", "-output=json"},
			expectedPolicySet: "polset-xyz789",
			expectedFormat:    "json",
		},
		{
			name:              "with table output explicitly",
			args:              []string{"-policy-set-id=polset-123", "-output=table"},
			expectedPolicySet: "polset-123",
			expectedFormat:    "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicySetParameterListCommand{}

			flags := cmd.Meta.FlagSet("policysetparameter list")
			flags.StringVar(&cmd.policySetID, "policy-set-id", "", "Policy Set ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the policy-set-id was set correctly
			if cmd.policySetID != tt.expectedPolicySet {
				t.Errorf("expected policySetID %q, got %q", tt.expectedPolicySet, cmd.policySetID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
