package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicyReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicyReadCommand{
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

func TestPolicyReadHelp(t *testing.T) {
	cmd := &PolicyReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policy read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestPolicyReadSynopsis(t *testing.T) {
	cmd := &PolicyReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read policy details" {
		t.Errorf("expected 'Read policy details', got %q", synopsis)
	}
}

func TestPolicyReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "policy id, default format",
			args:           []string{"-id=pol-abc123"},
			expectedID:     "pol-abc123",
			expectedFormat: "table",
		},
		{
			name:           "policy id, table format",
			args:           []string{"-id=pol-xyz789", "-output=table"},
			expectedID:     "pol-xyz789",
			expectedFormat: "table",
		},
		{
			name:           "policy id, json format",
			args:           []string{"-id=pol-def456", "-output=json"},
			expectedID:     "pol-def456",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicyReadCommand{}

			flags := cmd.Meta.FlagSet("policy read")
			flags.StringVar(&cmd.policyID, "id", "", "Policy ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the policy ID was set correctly
			if cmd.policyID != tt.expectedID {
				t.Errorf("expected policyID %q, got %q", tt.expectedID, cmd.policyID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
