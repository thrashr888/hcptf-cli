package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicySetParameterCreateRequiresPolicySetID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-key=test-key", "-value=test-value"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-policy-set-id") {
		t.Fatalf("expected policy-set-id error, got %q", out)
	}
}

func TestPolicySetParameterCreateRequiresKey(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-policy-set-id=polset-123", "-value=test-value"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-key") {
		t.Fatalf("expected key error, got %q", out)
	}
}

func TestPolicySetParameterCreateRequiresValue(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-policy-set-id=polset-123", "-key=test-key"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-value") {
		t.Fatalf("expected value error, got %q", out)
	}
}

func TestPolicySetParameterCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterCreateCommand{
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

func TestPolicySetParameterCreateHelp(t *testing.T) {
	cmd := &PolicySetParameterCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policysetparameter create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-policy-set-id") {
		t.Error("Help should mention -policy-set-id flag")
	}
	if !strings.Contains(help, "-key") {
		t.Error("Help should mention -key flag")
	}
	if !strings.Contains(help, "-value") {
		t.Error("Help should mention -value flag")
	}
	if !strings.Contains(help, "-sensitive") {
		t.Error("Help should mention -sensitive flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestPolicySetParameterCreateSynopsis(t *testing.T) {
	cmd := &PolicySetParameterCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a policy set parameter" {
		t.Errorf("expected 'Create a policy set parameter', got %q", synopsis)
	}
}

func TestPolicySetParameterCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedPolicySet string
		expectedKey       string
		expectedValue     string
		expectedSensitive bool
		expectedFormat    string
	}{
		{
			name:              "required flags only, default values",
			args:              []string{"-policy-set-id=polset-abc123", "-key=max_cost", "-value=1000"},
			expectedPolicySet: "polset-abc123",
			expectedKey:       "max_cost",
			expectedValue:     "1000",
			expectedSensitive: false,
			expectedFormat:    "table",
		},
		{
			name:              "with sensitive flag",
			args:              []string{"-policy-set-id=polset-xyz789", "-key=api_key", "-value=secret123", "-sensitive"},
			expectedPolicySet: "polset-xyz789",
			expectedKey:       "api_key",
			expectedValue:     "secret123",
			expectedSensitive: true,
			expectedFormat:    "table",
		},
		{
			name:              "with json output",
			args:              []string{"-policy-set-id=polset-123", "-key=param1", "-value=value1", "-output=json"},
			expectedPolicySet: "polset-123",
			expectedKey:       "param1",
			expectedValue:     "value1",
			expectedSensitive: false,
			expectedFormat:    "json",
		},
		{
			name:              "all flags set",
			args:              []string{"-policy-set-id=polset-full", "-key=full_key", "-value=full_value", "-sensitive", "-output=json"},
			expectedPolicySet: "polset-full",
			expectedKey:       "full_key",
			expectedValue:     "full_value",
			expectedSensitive: true,
			expectedFormat:    "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicySetParameterCreateCommand{}

			flags := cmd.Meta.FlagSet("policysetparameter create")
			flags.StringVar(&cmd.policySetID, "policy-set-id", "", "Policy Set ID (required)")
			flags.StringVar(&cmd.key, "key", "", "Parameter key (required)")
			flags.StringVar(&cmd.value, "value", "", "Parameter value (required)")
			flags.BoolVar(&cmd.sensitive, "sensitive", false, "Mark parameter as sensitive")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the policy-set-id was set correctly
			if cmd.policySetID != tt.expectedPolicySet {
				t.Errorf("expected policySetID %q, got %q", tt.expectedPolicySet, cmd.policySetID)
			}

			// Verify the key was set correctly
			if cmd.key != tt.expectedKey {
				t.Errorf("expected key %q, got %q", tt.expectedKey, cmd.key)
			}

			// Verify the value was set correctly
			if cmd.value != tt.expectedValue {
				t.Errorf("expected value %q, got %q", tt.expectedValue, cmd.value)
			}

			// Verify the sensitive flag was set correctly
			if cmd.sensitive != tt.expectedSensitive {
				t.Errorf("expected sensitive %v, got %v", tt.expectedSensitive, cmd.sensitive)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
