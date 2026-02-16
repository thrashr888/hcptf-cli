package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func newPolicySetParameterDeleteCommand(ui cli.Ui, svc policySetParameterDeleter) *PolicySetParameterDeleteCommand {
	return &PolicySetParameterDeleteCommand{
		Meta:                  newTestMeta(ui),
		policySetParameterSvc: svc,
	}
}

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

func TestPolicySetParameterDeleteHandlesAPIError(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicySetParameterDeleteService{err: errors.New("boom")}
	cmd := newPolicySetParameterDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-policy-set-id=polset-123", "-id=var-123", "-force"}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if svc.lastPolicySetID != "polset-123" || svc.lastParameterID != "var-123" {
		t.Fatalf("unexpected IDs: %q/%q", svc.lastPolicySetID, svc.lastParameterID)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "boom") {
		t.Fatalf("expected error output, got %q", ui.ErrorWriter.String())
	}
}

func TestPolicySetParameterDeleteSuccessWithForce(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicySetParameterDeleteService{}
	cmd := newPolicySetParameterDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-policy-set-id=polset-123", "-id=var-123", "-force"}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if svc.lastPolicySetID != "polset-123" || svc.lastParameterID != "var-123" {
		t.Fatalf("unexpected IDs: %q/%q", svc.lastPolicySetID, svc.lastParameterID)
	}
	if !strings.Contains(ui.OutputWriter.String(), "deleted successfully") {
		t.Fatalf("expected success output, got %q", ui.OutputWriter.String())
	}
}

func TestPolicySetParameterDeleteSuccessWithYesFlag(t *testing.T) {
	ui := cli.NewMockUi()
	svc := &mockPolicySetParameterDeleteService{}
	cmd := newPolicySetParameterDeleteCommand(ui, svc)

	if code := cmd.Run([]string{"-policy-set-id=polset-123", "-id=var-123", "-y"}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if svc.lastPolicySetID != "polset-123" || svc.lastParameterID != "var-123" {
		t.Fatalf("unexpected IDs: %q/%q", svc.lastPolicySetID, svc.lastParameterID)
	}
}

func TestPolicySetParameterDeleteCancellation(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("no\n")
	svc := &mockPolicySetParameterDeleteService{}
	cmd := newPolicySetParameterDeleteCommand(ui, svc)

	code := cmd.Run([]string{"-policy-set-id=polset-123", "-id=var-123"})
	if code != 0 {
		t.Fatalf("expected exit 0 on cancel, got %d", code)
	}
	if svc.lastPolicySetID != "" || svc.lastParameterID != "" {
		t.Fatalf("expected no delete call, got %q/%q", svc.lastPolicySetID, svc.lastParameterID)
	}

	if !strings.Contains(ui.OutputWriter.String(), "Deletion cancelled") {
		t.Fatalf("expected cancellation message")
	}
}

func TestPolicySetParameterDeleteSuccessWithConfirmation(t *testing.T) {
	ui := cli.NewMockUi()
	ui.InputReader = strings.NewReader("yes\n")
	svc := &mockPolicySetParameterDeleteService{}
	cmd := newPolicySetParameterDeleteCommand(ui, svc)

	code := cmd.Run([]string{"-policy-set-id=polset-123", "-id=var-123"})
	if code != 0 {
		t.Fatalf("expected exit 0 on confirm, got %d", code)
	}
	if svc.lastPolicySetID != "polset-123" || svc.lastParameterID != "var-123" {
		t.Fatalf("expected IDs polset-123/var-123, got %q/%q", svc.lastPolicySetID, svc.lastParameterID)
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
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "-y") {
		t.Error("Help should mention -y flag")
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
		name              string
		args              []string
		expectedPolicySet string
		expectedID        string
		expectedForce     bool
		expectedYes       bool
	}{
		{
			name:              "required flags only, default values",
			args:              []string{"-policy-set-id=polset-abc123", "-id=var-xyz789"},
			expectedPolicySet: "polset-abc123",
			expectedID:        "var-xyz789",
		},
		{
			name:              "with force flag",
			args:              []string{"-policy-set-id=polset-123", "-id=var-456", "-force"},
			expectedPolicySet: "polset-123",
			expectedID:        "var-456",
			expectedForce:     true,
		},
		{
			name:              "with shorthand force flag",
			args:              []string{"-policy-set-id=polset-789", "-id=var-000", "-f"},
			expectedPolicySet: "polset-789",
			expectedID:        "var-000",
			expectedForce:     true,
		},
		{
			name:              "with yes flag",
			args:              []string{"-policy-set-id=polset-full", "-id=var-full", "-y"},
			expectedPolicySet: "polset-full",
			expectedID:        "var-full",
			expectedYes:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicySetParameterDeleteCommand{}

			flags := cmd.Meta.FlagSet("policysetparameter delete")
			flags.StringVar(&cmd.policySetID, "policy-set-id", "", "Policy Set ID (required)")
			flags.StringVar(&cmd.parameterID, "id", "", "Parameter ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")
			flags.BoolVar(&cmd.force, "f", false, "Shorthand for -force")
			flags.BoolVar(&cmd.yes, "y", false, "Confirm delete without prompt")

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

			// Verify confirmation flags were set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
			if cmd.yes != tt.expectedYes {
				t.Errorf("expected yes %v, got %v", tt.expectedYes, cmd.yes)
			}
		})
	}
}
