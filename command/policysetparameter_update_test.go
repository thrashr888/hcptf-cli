package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicySetParameterUpdateRequiresPolicySetID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=var-123", "-key=new-key"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-policy-set-id") {
		t.Fatalf("expected policy-set-id error, got %q", out)
	}
}

func TestPolicySetParameterUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-policy-set-id=polset-123", "-key=new-key"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestPolicySetParameterUpdateRequiresAtLeastOneUpdateField(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-policy-set-id=polset-123", "-id=var-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "At least one of") {
		t.Fatalf("expected 'at least one of' error, got %q", out)
	}
}

func TestPolicySetParameterUpdateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetParameterUpdateCommand{
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

func TestPolicySetParameterUpdateHelp(t *testing.T) {
	cmd := &PolicySetParameterUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policysetparameter update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-policy-set-id") {
		t.Error("Help should mention -policy-set-id flag")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
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

func TestPolicySetParameterUpdateSynopsis(t *testing.T) {
	cmd := &PolicySetParameterUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update a policy set parameter" {
		t.Errorf("expected 'Update a policy set parameter', got %q", synopsis)
	}
}

func TestPolicySetParameterUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedPolicySet string
		expectedID        string
		expectedKey       string
		expectedValue     string
		expectedFormat    string
		hasSensitiveFlag  bool
	}{
		{
			name:              "update key only",
			args:              []string{"-policy-set-id=polset-abc123", "-id=var-xyz789", "-key=new_key"},
			expectedPolicySet: "polset-abc123",
			expectedID:        "var-xyz789",
			expectedKey:       "new_key",
			expectedValue:     "",
			expectedFormat:    "table",
			hasSensitiveFlag:  false,
		},
		{
			name:              "update value only",
			args:              []string{"-policy-set-id=polset-abc123", "-id=var-xyz789", "-value=new_value"},
			expectedPolicySet: "polset-abc123",
			expectedID:        "var-xyz789",
			expectedKey:       "",
			expectedValue:     "new_value",
			expectedFormat:    "table",
			hasSensitiveFlag:  false,
		},
		{
			name:              "update key and value",
			args:              []string{"-policy-set-id=polset-123", "-id=var-456", "-key=updated_key", "-value=updated_value"},
			expectedPolicySet: "polset-123",
			expectedID:        "var-456",
			expectedKey:       "updated_key",
			expectedValue:     "updated_value",
			expectedFormat:    "table",
			hasSensitiveFlag:  false,
		},
		{
			name:              "update with sensitive flag",
			args:              []string{"-policy-set-id=polset-123", "-id=var-456", "-value=secret", "-sensitive"},
			expectedPolicySet: "polset-123",
			expectedID:        "var-456",
			expectedKey:       "",
			expectedValue:     "secret",
			expectedFormat:    "table",
			hasSensitiveFlag:  true,
		},
		{
			name:              "all flags with json output",
			args:              []string{"-policy-set-id=polset-full", "-id=var-full", "-key=full_key", "-value=full_value", "-sensitive", "-output=json"},
			expectedPolicySet: "polset-full",
			expectedID:        "var-full",
			expectedKey:       "full_key",
			expectedValue:     "full_value",
			expectedFormat:    "json",
			hasSensitiveFlag:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicySetParameterUpdateCommand{}

			flags := cmd.Meta.FlagSet("policysetparameter update")
			flags.StringVar(&cmd.policySetID, "policy-set-id", "", "Policy Set ID (required)")
			flags.StringVar(&cmd.parameterID, "id", "", "Parameter ID (required)")
			flags.StringVar(&cmd.key, "key", "", "Parameter key")
			flags.StringVar(&cmd.value, "value", "", "Parameter value")
			sensitiveFlag := flags.Bool("sensitive", false, "Mark parameter as sensitive")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

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

			// Verify the key was set correctly
			if cmd.key != tt.expectedKey {
				t.Errorf("expected key %q, got %q", tt.expectedKey, cmd.key)
			}

			// Verify the value was set correctly
			if cmd.value != tt.expectedValue {
				t.Errorf("expected value %q, got %q", tt.expectedValue, cmd.value)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}

			// Note: We can't directly check cmd.sensitive as it's only set via flags.Visit
			// but we can verify the sensitiveFlag pointer was parsed correctly
			if tt.hasSensitiveFlag && sensitiveFlag != nil && !*sensitiveFlag {
				t.Errorf("expected sensitive flag to be parsed as true")
			}
		})
	}
}
