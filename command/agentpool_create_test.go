package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAgentPoolCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AgentPoolCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=test-pool"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestAgentPoolCreateRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AgentPoolCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestAgentPoolCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AgentPoolCreateCommand{
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

func TestAgentPoolCreateWithOrgAlias(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AgentPoolCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-org=test-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestAgentPoolCreateHelp(t *testing.T) {
	cmd := &AgentPoolCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf agentpool create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestAgentPoolCreateSynopsis(t *testing.T) {
	cmd := &AgentPoolCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new agent pool" {
		t.Errorf("expected 'Create a new agent pool', got %q", synopsis)
	}
}

func TestAgentPoolCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOrg    string
		expectedName   string
		expectedScoped bool
		expectedFmt    string
	}{
		{
			name:           "all required flags, default format",
			args:           []string{"-organization=my-org", "-name=test-pool"},
			expectedOrg:    "my-org",
			expectedName:   "test-pool",
			expectedScoped: false,
			expectedFmt:    "table",
		},
		{
			name:           "org alias with name",
			args:           []string{"-org=my-org", "-name=production-pool"},
			expectedOrg:    "my-org",
			expectedName:   "production-pool",
			expectedScoped: false,
			expectedFmt:    "table",
		},
		{
			name:           "organization scoped pool",
			args:           []string{"-org=prod-org", "-name=shared-pool", "-organization-scoped=true"},
			expectedOrg:    "prod-org",
			expectedName:   "shared-pool",
			expectedScoped: true,
			expectedFmt:    "table",
		},
		{
			name:           "all flags with json format",
			args:           []string{"-org=test-org", "-name=ci-pool", "-organization-scoped=true", "-output=json"},
			expectedOrg:    "test-org",
			expectedName:   "ci-pool",
			expectedScoped: true,
			expectedFmt:    "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AgentPoolCreateCommand{}

			flags := cmd.Meta.FlagSet("agentpool create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.name, "name", "", "Agent pool name (required)")
			flags.BoolVar(&cmd.organizationScoped, "organization-scoped", false, "Make agent pool organization scoped")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the organization scoped flag was set correctly
			if cmd.organizationScoped != tt.expectedScoped {
				t.Errorf("expected organizationScoped %v, got %v", tt.expectedScoped, cmd.organizationScoped)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}

func TestAgentPoolCreateOrganizationScopedFalse(t *testing.T) {
	cmd := &AgentPoolCreateCommand{}

	flags := cmd.Meta.FlagSet("agentpool create")
	flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&cmd.name, "name", "", "Agent pool name (required)")
	flags.BoolVar(&cmd.organizationScoped, "organization-scoped", false, "Make agent pool organization scoped")
	flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

	if err := flags.Parse([]string{"-organization=test-org", "-name=pool1"}); err != nil {
		t.Fatalf("flag parsing failed: %v", err)
	}

	if cmd.organizationScoped {
		t.Error("expected organization-scoped to default to false")
	}
}

func TestAgentPoolCreateOrganizationScopedTrue(t *testing.T) {
	cmd := &AgentPoolCreateCommand{}

	flags := cmd.Meta.FlagSet("agentpool create")
	flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&cmd.name, "name", "", "Agent pool name (required)")
	flags.BoolVar(&cmd.organizationScoped, "organization-scoped", false, "Make agent pool organization scoped")
	flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

	if err := flags.Parse([]string{"-organization=test-org", "-name=pool1", "-organization-scoped"}); err != nil {
		t.Fatalf("flag parsing failed: %v", err)
	}

	if !cmd.organizationScoped {
		t.Error("expected organization-scoped to be true")
	}
}

func TestAgentPoolCreateOrganizationScopedExplicitTrue(t *testing.T) {
	cmd := &AgentPoolCreateCommand{}

	flags := cmd.Meta.FlagSet("agentpool create")
	flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&cmd.name, "name", "", "Agent pool name (required)")
	flags.BoolVar(&cmd.organizationScoped, "organization-scoped", false, "Make agent pool organization scoped")
	flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

	if err := flags.Parse([]string{"-organization=test-org", "-name=pool1", "-organization-scoped=true"}); err != nil {
		t.Fatalf("flag parsing failed: %v", err)
	}

	if !cmd.organizationScoped {
		t.Error("expected organization-scoped to be true")
	}
}

func TestAgentPoolCreateOrganizationScopedExplicitFalse(t *testing.T) {
	cmd := &AgentPoolCreateCommand{}

	flags := cmd.Meta.FlagSet("agentpool create")
	flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&cmd.name, "name", "", "Agent pool name (required)")
	flags.BoolVar(&cmd.organizationScoped, "organization-scoped", false, "Make agent pool organization scoped")
	flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

	if err := flags.Parse([]string{"-organization=test-org", "-name=pool1", "-organization-scoped=false"}); err != nil {
		t.Fatalf("flag parsing failed: %v", err)
	}

	if cmd.organizationScoped {
		t.Error("expected organization-scoped to be false")
	}
}

func TestAgentPoolCreateOrganizationScopedWithJSONOutput(t *testing.T) {
	cmd := &AgentPoolCreateCommand{}

	flags := cmd.Meta.FlagSet("agentpool create")
	flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&cmd.name, "name", "", "Agent pool name (required)")
	flags.BoolVar(&cmd.organizationScoped, "organization-scoped", false, "Make agent pool organization scoped")
	flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

	if err := flags.Parse([]string{"-organization=test-org", "-name=pool1", "-organization-scoped", "-output=json"}); err != nil {
		t.Fatalf("flag parsing failed: %v", err)
	}

	if !cmd.organizationScoped {
		t.Error("expected organization-scoped to be true")
	}
	if cmd.format != "json" {
		t.Errorf("expected format 'json', got %q", cmd.format)
	}
}

func TestAgentPoolCreateNonScopedWithJSONOutput(t *testing.T) {
	cmd := &AgentPoolCreateCommand{}

	flags := cmd.Meta.FlagSet("agentpool create")
	flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&cmd.name, "name", "", "Agent pool name (required)")
	flags.BoolVar(&cmd.organizationScoped, "organization-scoped", false, "Make agent pool organization scoped")
	flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

	if err := flags.Parse([]string{"-organization=test-org", "-name=pool1", "-output=json"}); err != nil {
		t.Fatalf("flag parsing failed: %v", err)
	}

	if cmd.organizationScoped {
		t.Error("expected organization-scoped to be false")
	}
	if cmd.format != "json" {
		t.Errorf("expected format 'json', got %q", cmd.format)
	}
}

func TestAgentPoolCreateAllFlagCombinations(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectScoped bool
		expectFormat string
	}{
		{
			name:         "default values",
			args:         []string{"-organization=test-org", "-name=pool1"},
			expectScoped: false,
			expectFormat: "table",
		},
		{
			name:         "organization-scoped true",
			args:         []string{"-organization=test-org", "-name=pool2", "-organization-scoped"},
			expectScoped: true,
			expectFormat: "table",
		},
		{
			name:         "organization-scoped false explicit",
			args:         []string{"-organization=test-org", "-name=pool3", "-organization-scoped=false"},
			expectScoped: false,
			expectFormat: "table",
		},
		{
			name:         "json output only",
			args:         []string{"-organization=test-org", "-name=pool4", "-output=json"},
			expectScoped: false,
			expectFormat: "json",
		},
		{
			name:         "organization-scoped and json",
			args:         []string{"-organization=test-org", "-name=pool5", "-organization-scoped", "-output=json"},
			expectScoped: true,
			expectFormat: "json",
		},
		{
			name:         "organization-scoped false and json",
			args:         []string{"-organization=test-org", "-name=pool6", "-organization-scoped=false", "-output=json"},
			expectScoped: false,
			expectFormat: "json",
		},
		{
			name:         "table output explicit",
			args:         []string{"-organization=test-org", "-name=pool7", "-output=table"},
			expectScoped: false,
			expectFormat: "table",
		},
		{
			name:         "organization-scoped and table",
			args:         []string{"-organization=test-org", "-name=pool8", "-organization-scoped", "-output=table"},
			expectScoped: true,
			expectFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AgentPoolCreateCommand{}

			flags := cmd.Meta.FlagSet("agentpool create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.name, "name", "", "Agent pool name (required)")
			flags.BoolVar(&cmd.organizationScoped, "organization-scoped", false, "Make agent pool organization scoped")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			if cmd.organizationScoped != tt.expectScoped {
				t.Errorf("expected organizationScoped %v, got %v", tt.expectScoped, cmd.organizationScoped)
			}
			if cmd.format != tt.expectFormat {
				t.Errorf("expected format %q, got %q", tt.expectFormat, cmd.format)
			}
		})
	}
}

func TestAgentPoolCreateOrgAliasWithOrganizationScoped(t *testing.T) {
	cmd := &AgentPoolCreateCommand{}

	flags := cmd.Meta.FlagSet("agentpool create")
	flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&cmd.name, "name", "", "Agent pool name (required)")
	flags.BoolVar(&cmd.organizationScoped, "organization-scoped", false, "Make agent pool organization scoped")
	flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

	if err := flags.Parse([]string{"-org=test-org", "-name=pool1", "-organization-scoped"}); err != nil {
		t.Fatalf("flag parsing failed: %v", err)
	}

	if cmd.organization != "test-org" {
		t.Errorf("expected organization 'test-org', got %q", cmd.organization)
	}
	if !cmd.organizationScoped {
		t.Error("expected organization-scoped to be true")
	}
}

func TestAgentPoolCreateMultiplePoolNames(t *testing.T) {
	tests := []struct {
		name     string
		poolName string
	}{
		{"simple name", "mypool"},
		{"hyphenated name", "my-pool"},
		{"underscored name", "my_pool"},
		{"numbered name", "pool123"},
		{"mixed case name", "MyPool"},
		{"long name", "my-very-long-agent-pool-name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AgentPoolCreateCommand{}

			flags := cmd.Meta.FlagSet("agentpool create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.name, "name", "", "Agent pool name (required)")
			flags.BoolVar(&cmd.organizationScoped, "organization-scoped", false, "Make agent pool organization scoped")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse([]string{"-organization=test-org", "-name=" + tt.poolName}); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			if cmd.name != tt.poolName {
				t.Errorf("expected name %q, got %q", tt.poolName, cmd.name)
			}
		})
	}
}

func TestAgentPoolCreateTableOutputDefault(t *testing.T) {
	cmd := &AgentPoolCreateCommand{}

	flags := cmd.Meta.FlagSet("agentpool create")
	flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&cmd.name, "name", "", "Agent pool name (required)")
	flags.BoolVar(&cmd.organizationScoped, "organization-scoped", false, "Make agent pool organization scoped")
	flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

	if err := flags.Parse([]string{"-organization=test-org", "-name=test-pool"}); err != nil {
		t.Fatalf("flag parsing failed: %v", err)
	}

	if cmd.format != "table" {
		t.Errorf("expected default format 'table', got %q", cmd.format)
	}
}
