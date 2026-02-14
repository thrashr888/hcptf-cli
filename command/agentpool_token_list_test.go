package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAgentPoolTokenListRequiresAgentPoolID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AgentPoolTokenListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-agent-pool-id") {
		t.Fatalf("expected agent-pool-id error, got %q", out)
	}
}

func TestAgentPoolTokenListHelp(t *testing.T) {
	cmd := &AgentPoolTokenListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf agentpool token-list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-agent-pool-id") {
		t.Error("Help should mention -agent-pool-id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -agent-pool-id is required")
	}
}

func TestAgentPoolTokenListSynopsis(t *testing.T) {
	cmd := &AgentPoolTokenListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List agent tokens in an agent pool" {
		t.Errorf("expected 'List agent tokens in an agent pool', got %q", synopsis)
	}
}

func TestAgentPoolTokenListFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "agent pool id, default format",
			args:        []string{"-agent-pool-id=apool-123abc"},
			expectedID:  "apool-123abc",
			expectedFmt: "table",
		},
		{
			name:        "agent pool id, table format",
			args:        []string{"-agent-pool-id=apool-456def", "-output=table"},
			expectedID:  "apool-456def",
			expectedFmt: "table",
		},
		{
			name:        "agent pool id, json format",
			args:        []string{"-agent-pool-id=apool-789ghi", "-output=json"},
			expectedID:  "apool-789ghi",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AgentPoolTokenListCommand{}

			flags := cmd.Meta.FlagSet("agentpool token-list")
			flags.StringVar(&cmd.agentPoolID, "agent-pool-id", "", "Agent pool ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the agent pool ID was set correctly
			if cmd.agentPoolID != tt.expectedID {
				t.Errorf("expected agentPoolID %q, got %q", tt.expectedID, cmd.agentPoolID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
