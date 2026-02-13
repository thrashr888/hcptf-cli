package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAgentPoolTokenCreateRequiresAgentPoolID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AgentPoolTokenCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-description=test-token"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-agent-pool-id") {
		t.Fatalf("expected agent-pool-id error, got %q", out)
	}
}

func TestAgentPoolTokenCreateRequiresDescription(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AgentPoolTokenCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-agent-pool-id=apool-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-description") {
		t.Fatalf("expected description error, got %q", out)
	}
}

func TestAgentPoolTokenCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AgentPoolTokenCreateCommand{
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

func TestAgentPoolTokenCreateHelp(t *testing.T) {
	cmd := &AgentPoolTokenCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf agentpool token-create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-agent-pool-id") {
		t.Error("Help should mention -agent-pool-id flag")
	}
	if !strings.Contains(help, "-description") {
		t.Error("Help should mention -description flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestAgentPoolTokenCreateSynopsis(t *testing.T) {
	cmd := &AgentPoolTokenCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create an agent token for an agent pool" {
		t.Errorf("expected 'Create an agent token for an agent pool', got %q", synopsis)
	}
}

func TestAgentPoolTokenCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedDesc string
		expectedFmt string
	}{
		{
			name:        "all required flags, default format",
			args:        []string{"-agent-pool-id=apool-123abc", "-description=Production agent token"},
			expectedID:  "apool-123abc",
			expectedDesc: "Production agent token",
			expectedFmt: "table",
		},
		{
			name:        "required flags with table format",
			args:        []string{"-agent-pool-id=apool-456def", "-description=Dev agent", "-output=table"},
			expectedID:  "apool-456def",
			expectedDesc: "Dev agent",
			expectedFmt: "table",
		},
		{
			name:        "required flags with json format",
			args:        []string{"-agent-pool-id=apool-789ghi", "-description=CI agent token", "-output=json"},
			expectedID:  "apool-789ghi",
			expectedDesc: "CI agent token",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AgentPoolTokenCreateCommand{}

			flags := cmd.Meta.FlagSet("agentpool token-create")
			flags.StringVar(&cmd.agentPoolID, "agent-pool-id", "", "Agent pool ID (required)")
			flags.StringVar(&cmd.description, "description", "", "Agent token description (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the agent pool ID was set correctly
			if cmd.agentPoolID != tt.expectedID {
				t.Errorf("expected agentPoolID %q, got %q", tt.expectedID, cmd.agentPoolID)
			}

			// Verify the description was set correctly
			if cmd.description != tt.expectedDesc {
				t.Errorf("expected description %q, got %q", tt.expectedDesc, cmd.description)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
