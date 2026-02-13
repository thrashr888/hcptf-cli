package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestProjectCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=test-project"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestProjectCreateRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectCreateCommand{
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

func TestProjectCreateRequiresBothFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectCreateCommand{
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

func TestProjectCreateHelp(t *testing.T) {
	cmd := &ProjectCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf project create") {
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

func TestProjectCreateSynopsis(t *testing.T) {
	cmd := &ProjectCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new project" {
		t.Errorf("expected 'Create a new project', got %q", synopsis)
	}
}

func TestProjectCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		expectedOrg        string
		expectedName       string
		expectedDesc       string
		expectedFmt        string
	}{
		{
			name:         "org and name, default format",
			args:         []string{"-organization=my-org", "-name=test-project"},
			expectedOrg:  "my-org",
			expectedName: "test-project",
			expectedDesc: "",
			expectedFmt:  "table",
		},
		{
			name:         "org alias and name with description",
			args:         []string{"-org=my-org", "-name=infra", "-description=Infrastructure project"},
			expectedOrg:  "my-org",
			expectedName: "infra",
			expectedDesc: "Infrastructure project",
			expectedFmt:  "table",
		},
		{
			name:         "org, name, json format",
			args:         []string{"-org=prod-org", "-name=platform", "-output=json"},
			expectedOrg:  "prod-org",
			expectedName: "platform",
			expectedDesc: "",
			expectedFmt:  "json",
		},
		{
			name:         "org, name, description, json format",
			args:         []string{"-org=test-org", "-name=services", "-description=Microservices", "-output=json"},
			expectedOrg:  "test-org",
			expectedName: "services",
			expectedDesc: "Microservices",
			expectedFmt:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ProjectCreateCommand{}

			flags := cmd.Meta.FlagSet("project create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.name, "name", "", "Project name (required)")
			flags.StringVar(&cmd.description, "description", "", "Project description")
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
