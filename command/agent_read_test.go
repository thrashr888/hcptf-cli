package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAgentReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AgentReadCommand{
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

func TestAgentReadHelp(t *testing.T) {
	cmd := &AgentReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf agent read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestAgentReadSynopsis(t *testing.T) {
	cmd := &AgentReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read agent details and status" {
		t.Errorf("expected 'Read agent details and status', got %q", synopsis)
	}
}

func TestAgentReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id with default format",
			args:        []string{"-id=agent-123abc"},
			expectedID:  "agent-123abc",
			expectedFmt: "table",
		},
		{
			name:        "id with table format",
			args:        []string{"-id=agent-456def", "-output=table"},
			expectedID:  "agent-456def",
			expectedFmt: "table",
		},
		{
			name:        "id with json format",
			args:        []string{"-id=agent-789ghi", "-output=json"},
			expectedID:  "agent-789ghi",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AgentReadCommand{}

			flags := cmd.Meta.FlagSet("agent read")
			flags.StringVar(&cmd.id, "id", "", "Agent ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
