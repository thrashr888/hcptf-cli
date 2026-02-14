package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVariableSetUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetUpdateCommand{
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

func TestVariableSetUpdateRequiresEmptyID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=", "-name=new-name"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestVariableSetUpdateHelp(t *testing.T) {
	cmd := &VariableSetUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf variableset update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestVariableSetUpdateSynopsis(t *testing.T) {
	cmd := &VariableSetUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update a variable set's settings" {
		t.Errorf("expected 'Update a variable set's settings', got %q", synopsis)
	}
}

func TestVariableSetUpdateValidatesGlobal(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VariableSetUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=varset-123", "-global=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "global") && !strings.Contains(out, "'true' or 'false'") {
		t.Fatalf("expected global validation error, got %q", out)
	}
}

func TestVariableSetUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedName   string
		expectedDesc   string
		expectedGlobal string
		expectedFormat string
	}{
		{
			name:           "id only, default values",
			args:           []string{"-id=varset-12345"},
			expectedID:     "varset-12345",
			expectedName:   "",
			expectedDesc:   "",
			expectedGlobal: "",
			expectedFormat: "table",
		},
		{
			name:           "update name",
			args:           []string{"-id=varset-abc123", "-name=new-name"},
			expectedID:     "varset-abc123",
			expectedName:   "new-name",
			expectedDesc:   "",
			expectedGlobal: "",
			expectedFormat: "table",
		},
		{
			name:           "update global to true with json output",
			args:           []string{"-id=varset-xyz789", "-global=true", "-output=json"},
			expectedID:     "varset-xyz789",
			expectedName:   "",
			expectedDesc:   "",
			expectedGlobal: "true",
			expectedFormat: "json",
		},
		{
			name:           "update all fields",
			args:           []string{"-id=varset-def456", "-name=updated-varset", "-description=Updated description", "-global=false", "-output=json"},
			expectedID:     "varset-def456",
			expectedName:   "updated-varset",
			expectedDesc:   "Updated description",
			expectedGlobal: "false",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VariableSetUpdateCommand{}

			flags := cmd.Meta.FlagSet("variableset update")
			flags.StringVar(&cmd.id, "id", "", "Variable set ID (required)")
			flags.StringVar(&cmd.name, "name", "", "Variable set name")
			flags.StringVar(&cmd.description, "description", "", "Variable set description")
			flags.StringVar(&cmd.global, "global", "", "Apply to all workspaces (true or false)")
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

			// Verify the description was set correctly
			if cmd.description != tt.expectedDesc {
				t.Errorf("expected description %q, got %q", tt.expectedDesc, cmd.description)
			}

			// Verify the global was set correctly
			if cmd.global != tt.expectedGlobal {
				t.Errorf("expected global %q, got %q", tt.expectedGlobal, cmd.global)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
