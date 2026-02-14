package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAgentPoolUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AgentPoolUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=new-name"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestAgentPoolUpdateRequiresAtLeastOneUpdateFlag(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AgentPoolUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=apool-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "At least one update flag") {
		t.Fatalf("expected update flag error, got %q", out)
	}
}

func TestAgentPoolUpdateInvalidOrganizationScopedValue(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AgentPoolUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=apool-123", "-organization-scoped=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "must be 'true' or 'false'") {
		t.Fatalf("expected organization-scoped validation error, got %q", out)
	}
}

func TestAgentPoolUpdateHelp(t *testing.T) {
	cmd := &AgentPoolUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf agentpool update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "-organization-scoped") {
		t.Error("Help should mention -organization-scoped flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestAgentPoolUpdateSynopsis(t *testing.T) {
	cmd := &AgentPoolUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update an agent pool" {
		t.Errorf("expected 'Update an agent pool', got %q", synopsis)
	}
}

func TestAgentPoolUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedID   string
		expectedName string
		expectedFmt  string
	}{
		{
			name:         "id with name update",
			args:         []string{"-id=apool-123abc", "-name=new-pool-name"},
			expectedID:   "apool-123abc",
			expectedName: "new-pool-name",
			expectedFmt:  "table",
		},
		{
			name:         "id with name and table format",
			args:         []string{"-id=apool-456def", "-name=updated-pool", "-output=table"},
			expectedID:   "apool-456def",
			expectedName: "updated-pool",
			expectedFmt:  "table",
		},
		{
			name:         "id with name and json format",
			args:         []string{"-id=apool-789ghi", "-name=prod-pool", "-output=json"},
			expectedID:   "apool-789ghi",
			expectedName: "prod-pool",
			expectedFmt:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AgentPoolUpdateCommand{}

			flags := cmd.Meta.FlagSet("agentpool update")
			flags.StringVar(&cmd.id, "id", "", "Agent pool ID (required)")
			flags.StringVar(&cmd.name, "name", "", "Agent pool name")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}

func TestAgentPoolUpdateOrganizationScopedTrue(t *testing.T) {
	cmd := &AgentPoolUpdateCommand{}

	flags := cmd.Meta.FlagSet("agentpool update")
	flags.StringVar(&cmd.id, "id", "", "Agent pool ID (required)")
	flags.StringVar(&cmd.name, "name", "", "Agent pool name")
	flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

	var orgScopedFlag string
	flags.StringVar(&orgScopedFlag, "organization-scoped", "", "Make agent pool organization scoped (true/false)")

	if err := flags.Parse([]string{"-id=apool-123", "-organization-scoped=true"}); err != nil {
		t.Fatalf("flag parsing failed: %v", err)
	}

	if orgScopedFlag != "true" {
		t.Errorf("expected organization-scoped 'true', got %q", orgScopedFlag)
	}
}

func TestAgentPoolUpdateOrganizationScopedFalse(t *testing.T) {
	cmd := &AgentPoolUpdateCommand{}

	flags := cmd.Meta.FlagSet("agentpool update")
	flags.StringVar(&cmd.id, "id", "", "Agent pool ID (required)")
	flags.StringVar(&cmd.name, "name", "", "Agent pool name")
	flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

	var orgScopedFlag string
	flags.StringVar(&orgScopedFlag, "organization-scoped", "", "Make agent pool organization scoped (true/false)")

	if err := flags.Parse([]string{"-id=apool-123", "-organization-scoped=false"}); err != nil {
		t.Fatalf("flag parsing failed: %v", err)
	}

	if orgScopedFlag != "false" {
		t.Errorf("expected organization-scoped 'false', got %q", orgScopedFlag)
	}
}

func TestAgentPoolUpdateWithNameOnly(t *testing.T) {
	cmd := &AgentPoolUpdateCommand{}

	flags := cmd.Meta.FlagSet("agentpool update")
	flags.StringVar(&cmd.id, "id", "", "Agent pool ID (required)")
	flags.StringVar(&cmd.name, "name", "", "Agent pool name")
	flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

	if err := flags.Parse([]string{"-id=apool-123", "-name=new-name"}); err != nil {
		t.Fatalf("flag parsing failed: %v", err)
	}

	if cmd.name != "new-name" {
		t.Errorf("expected name 'new-name', got %q", cmd.name)
	}

	if cmd.organizationScoped != nil {
		t.Error("expected organizationScoped to be nil when not provided")
	}
}
