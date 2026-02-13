package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestProjectUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ProjectUpdateCommand{
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

func TestProjectUpdateHelp(t *testing.T) {
	cmd := &ProjectUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf project update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "-description") {
		t.Error("Help should mention -description flag")
	}
}

func TestProjectUpdateSynopsis(t *testing.T) {
	cmd := &ProjectUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update project settings" {
		t.Errorf("expected 'Update project settings', got %q", synopsis)
	}
}

func TestProjectUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedID   string
		expectedName string
		expectedDesc string
		expectedFmt  string
	}{
		{
			name:         "id and name",
			args:         []string{"-id=prj-abc123", "-name=new-name"},
			expectedID:   "prj-abc123",
			expectedName: "new-name",
			expectedDesc: "",
			expectedFmt:  "table",
		},
		{
			name:         "id and description",
			args:         []string{"-id=prj-xyz789", "-description=Updated description"},
			expectedID:   "prj-xyz789",
			expectedName: "",
			expectedDesc: "Updated description",
			expectedFmt:  "table",
		},
		{
			name:         "id, name, description",
			args:         []string{"-id=prj-test456", "-name=updated-project", "-description=New description"},
			expectedID:   "prj-test456",
			expectedName: "updated-project",
			expectedDesc: "New description",
			expectedFmt:  "table",
		},
		{
			name:         "id, name, json format",
			args:         []string{"-id=prj-prod123", "-name=prod-project", "-output=json"},
			expectedID:   "prj-prod123",
			expectedName: "prod-project",
			expectedDesc: "",
			expectedFmt:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ProjectUpdateCommand{}

			flags := cmd.Meta.FlagSet("project update")
			flags.StringVar(&cmd.projectID, "id", "", "Project ID (required)")
			flags.StringVar(&cmd.name, "name", "", "New project name")
			flags.StringVar(&cmd.description, "description", "", "New project description")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.projectID != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.projectID)
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
