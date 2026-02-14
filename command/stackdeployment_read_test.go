package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackDeploymentReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackDeploymentReadCommand{
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

func TestStackDeploymentReadHelp(t *testing.T) {
	cmd := &StackDeploymentReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stackdeployment read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "Example:") {
		t.Error("Help should contain examples")
	}
}

func TestStackDeploymentReadSynopsis(t *testing.T) {
	cmd := &StackDeploymentReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read stack deployment details and status" {
		t.Errorf("expected 'Read stack deployment details and status', got %q", synopsis)
	}
}

func TestStackDeploymentReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id, default format",
			args:        []string{"-id=sdr-abc123"},
			expectedID:  "sdr-abc123",
			expectedFmt: "table",
		},
		{
			name:        "id, table format",
			args:        []string{"-id=sdr-xyz789", "-output=table"},
			expectedID:  "sdr-xyz789",
			expectedFmt: "table",
		},
		{
			name:        "id, json format",
			args:        []string{"-id=sdr-test456", "-output=json"},
			expectedID:  "sdr-test456",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackDeploymentReadCommand{}

			flags := cmd.Meta.FlagSet("stackdeployment read")
			flags.StringVar(&cmd.deploymentRunID, "id", "", "Stack deployment run ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.deploymentRunID != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.deploymentRunID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
