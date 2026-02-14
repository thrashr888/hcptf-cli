package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicySetRemovePolicyRequiresPolicySetID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetRemovePolicyCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-policy-id=pol-12345"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-policyset-id") {
		t.Fatalf("expected policyset-id error, got %q", out)
	}
}

func TestPolicySetRemovePolicyRequiresPolicyID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetRemovePolicyCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-policyset-id=polset-12345"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-policy-id") {
		t.Fatalf("expected policy-id error, got %q", out)
	}
}

func TestPolicySetRemovePolicyRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetRemovePolicyCommand{
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

func TestPolicySetRemovePolicyHelp(t *testing.T) {
	cmd := &PolicySetRemovePolicyCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policyset remove-policy") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-policyset-id") {
		t.Error("Help should mention -policyset-id flag")
	}
	if !strings.Contains(help, "-policy-id") {
		t.Error("Help should mention -policy-id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestPolicySetRemovePolicySynopsis(t *testing.T) {
	cmd := &PolicySetRemovePolicyCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Remove a policy from a policy set" {
		t.Errorf("expected 'Remove a policy from a policy set', got %q", synopsis)
	}
}

func TestPolicySetRemovePolicyFlagParsing(t *testing.T) {
	tests := []struct {
		name                string
		args                []string
		expectedPolicySetID string
		expectedPolicyID    string
	}{
		{
			name:                "both required flags",
			args:                []string{"-policyset-id=polset-12345", "-policy-id=pol-67890"},
			expectedPolicySetID: "polset-12345",
			expectedPolicyID:    "pol-67890",
		},
		{
			name:                "different order",
			args:                []string{"-policy-id=pol-abcde", "-policyset-id=polset-xyz"},
			expectedPolicySetID: "polset-xyz",
			expectedPolicyID:    "pol-abcde",
		},
		{
			name:                "long policy set id",
			args:                []string{"-policyset-id=polset-long-id-123", "-policy-id=pol-456"},
			expectedPolicySetID: "polset-long-id-123",
			expectedPolicyID:    "pol-456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicySetRemovePolicyCommand{}

			flags := cmd.Meta.FlagSet("policyset remove-policy")
			flags.StringVar(&cmd.policySetID, "policyset-id", "", "Policy set ID (required)")
			flags.StringVar(&cmd.policyID, "policy-id", "", "Policy ID to remove (required)")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the policyset-id was set correctly
			if cmd.policySetID != tt.expectedPolicySetID {
				t.Errorf("expected policySetID %q, got %q", tt.expectedPolicySetID, cmd.policySetID)
			}

			// Verify the policy-id was set correctly
			if cmd.policyID != tt.expectedPolicyID {
				t.Errorf("expected policyID %q, got %q", tt.expectedPolicyID, cmd.policyID)
			}
		})
	}
}
