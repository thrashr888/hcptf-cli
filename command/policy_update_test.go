package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicyUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicyUpdateCommand{
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

func TestPolicyUpdateHelp(t *testing.T) {
	cmd := &PolicyUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policy update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "-description") {
		t.Error("Help should mention -description flag")
	}
	if !strings.Contains(help, "-enforce") {
		t.Error("Help should mention -enforce flag")
	}
	if !strings.Contains(help, "-policy-file") {
		t.Error("Help should mention -policy-file flag")
	}
}

func TestPolicyUpdateSynopsis(t *testing.T) {
	cmd := &PolicyUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update policy settings" {
		t.Errorf("expected 'Update policy settings', got %q", synopsis)
	}
}

func TestPolicyUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedDesc   string
		expectedEnf    string
		expectedFile   string
		expectedFormat string
	}{
		{
			name:           "id only, default format",
			args:           []string{"-id=pol-abc123"},
			expectedID:     "pol-abc123",
			expectedDesc:   "",
			expectedEnf:    "",
			expectedFile:   "",
			expectedFormat: "table",
		},
		{
			name:           "id with description and enforcement",
			args:           []string{"-id=pol-xyz789", "-description=Updated policy", "-enforce=hard-mandatory"},
			expectedID:     "pol-xyz789",
			expectedDesc:   "Updated policy",
			expectedEnf:    "hard-mandatory",
			expectedFile:   "",
			expectedFormat: "table",
		},
		{
			name:           "id with policy file and json format",
			args:           []string{"-id=pol-def456", "-policy-file=updated.sentinel", "-output=json"},
			expectedID:     "pol-def456",
			expectedDesc:   "",
			expectedEnf:    "",
			expectedFile:   "updated.sentinel",
			expectedFormat: "json",
		},
		{
			name:           "id with all optional flags",
			args:           []string{"-id=pol-ghi789", "-description=Complete update", "-enforce=soft-mandatory", "-policy-file=new.sentinel", "-output=json"},
			expectedID:     "pol-ghi789",
			expectedDesc:   "Complete update",
			expectedEnf:    "soft-mandatory",
			expectedFile:   "new.sentinel",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicyUpdateCommand{}

			flags := cmd.Meta.FlagSet("policy update")
			flags.StringVar(&cmd.policyID, "id", "", "Policy ID (required)")
			flags.StringVar(&cmd.description, "description", "", "Policy description")
			flags.StringVar(&cmd.enforce, "enforce", "", "Enforcement level: advisory, soft-mandatory, or hard-mandatory")
			flags.StringVar(&cmd.policyFile, "policy-file", "", "Path to policy file")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the policy ID was set correctly
			if cmd.policyID != tt.expectedID {
				t.Errorf("expected policyID %q, got %q", tt.expectedID, cmd.policyID)
			}

			// Verify the description was set correctly
			if cmd.description != tt.expectedDesc {
				t.Errorf("expected description %q, got %q", tt.expectedDesc, cmd.description)
			}

			// Verify the enforce was set correctly
			if cmd.enforce != tt.expectedEnf {
				t.Errorf("expected enforce %q, got %q", tt.expectedEnf, cmd.enforce)
			}

			// Verify the policy-file was set correctly
			if cmd.policyFile != tt.expectedFile {
				t.Errorf("expected policyFile %q, got %q", tt.expectedFile, cmd.policyFile)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
