package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicySetParameterDeleteRequiresPolicySetID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=var-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-policy-set-id") {
		t.Fatalf("expected policy-set-id error, got %q", out)
	}
}

func TestPolicySetParameterDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterDeleteCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-policy-set-id=polset-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestPolicySetParameterDeleteRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterDeleteCommand{
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

func TestPolicySetParameterDeleteCancellation(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("no\n")
	cmd := &PolicySetParameterDeleteCommand{
		Meta: newTestMeta(ui),
	}
	cmd.Meta.Ui = ui

	code := cmd.Run([]string{"-policy-set-id=polset-123", "-id=var-123"})
	if code != 0 {
		t.Fatalf("expected exit 0 on cancel, got %d", code)
	}

	if !strings.Contains(ui.OutputWriter.String(), "Deletion cancelled") {
		t.Fatalf("expected cancellation message")
	}
}

func TestPolicySetParameterDeleteHelp(t *testing.T) {
	cmd := &PolicySetParameterDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policysetparameter delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-policy-set-id") {
		t.Error("Help should mention -policy-set-id flag")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-auto-approve") {
		t.Error("Help should mention -auto-approve flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestPolicySetParameterDeleteSynopsis(t *testing.T) {
	cmd := &PolicySetParameterDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a policy set parameter" {
		t.Errorf("expected 'Delete a policy set parameter', got %q", synopsis)
	}
}

func TestPolicySetParameterDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name                string
		args                []string
		expectedPolicySet   string
		expectedID          string
		expectedAutoApprove bool
	}{
		{
			name:                "required flags only, default values",
			args:                []string{"-policy-set-id=polset-abc123", "-id=var-xyz789"},
			expectedPolicySet:   "polset-abc123",
			expectedID:          "var-xyz789",
			expectedAutoApprove: false,
		},
		{
			name:                "with auto-approve flag",
			args:                []string{"-policy-set-id=polset-123", "-id=var-456", "-auto-approve"},
			expectedPolicySet:   "polset-123",
			expectedID:          "var-456",
			expectedAutoApprove: true,
		},
		{
			name:                "different policy set and parameter",
			args:                []string{"-policy-set-id=polset-prod", "-id=var-prod-123"},
			expectedPolicySet:   "polset-prod",
			expectedID:          "var-prod-123",
			expectedAutoApprove: false,
		},
		{
			name:                "all flags set",
			args:                []string{"-policy-set-id=polset-full", "-id=var-full", "-auto-approve"},
			expectedPolicySet:   "polset-full",
			expectedID:          "var-full",
			expectedAutoApprove: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicySetParameterDeleteCommand{}

			flags := cmd.Meta.FlagSet("policysetparameter delete")
			flags.StringVar(&cmd.policySetID, "policy-set-id", "", "Policy Set ID (required)")
			flags.StringVar(&cmd.parameterID, "id", "", "Parameter ID (required)")
			flags.BoolVar(&cmd.autoApprove, "auto-approve", false, "Skip confirmation prompt")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the policy-set-id was set correctly
			if cmd.policySetID != tt.expectedPolicySet {
				t.Errorf("expected policySetID %q, got %q", tt.expectedPolicySet, cmd.policySetID)
			}

			// Verify the id was set correctly
			if cmd.parameterID != tt.expectedID {
				t.Errorf("expected parameterID %q, got %q", tt.expectedID, cmd.parameterID)
			}

			// Verify the auto-approve flag was set correctly
			if cmd.autoApprove != tt.expectedAutoApprove {
				t.Errorf("expected autoApprove %v, got %v", tt.expectedAutoApprove, cmd.autoApprove)
			}
		})
	}
}
