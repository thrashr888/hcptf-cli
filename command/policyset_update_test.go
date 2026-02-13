package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicySetUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetUpdateCommand{
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

func TestPolicySetUpdateHelp(t *testing.T) {
	cmd := &PolicySetUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policyset update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestPolicySetUpdateSynopsis(t *testing.T) {
	cmd := &PolicySetUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update policy set settings" {
		t.Errorf("expected 'Update policy set settings', got %q", synopsis)
	}
}

func TestPolicySetUpdateFlagParsing(t *testing.T) {
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
			name:           "id with name update",
			args:           []string{"-id=polset-12345", "-name=new-name"},
			expectedID:     "polset-12345",
			expectedName:   "new-name",
			expectedDesc:   "",
			expectedGlobal: "",
			expectedFormat: "table",
		},
		{
			name:           "id with description update",
			args:           []string{"-id=polset-67890", "-description=Updated description"},
			expectedID:     "polset-67890",
			expectedName:   "",
			expectedDesc:   "Updated description",
			expectedGlobal: "",
			expectedFormat: "table",
		},
		{
			name:           "id with global true",
			args:           []string{"-id=polset-abcde", "-global=true", "-output=json"},
			expectedID:     "polset-abcde",
			expectedName:   "",
			expectedDesc:   "",
			expectedGlobal: "true",
			expectedFormat: "json",
		},
		{
			name:           "all flags",
			args:           []string{"-id=polset-xyz", "-name=updated-name", "-description=Updated desc", "-global=false"},
			expectedID:     "polset-xyz",
			expectedName:   "updated-name",
			expectedDesc:   "Updated desc",
			expectedGlobal: "false",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicySetUpdateCommand{}

			flags := cmd.Meta.FlagSet("policyset update")
			flags.StringVar(&cmd.id, "id", "", "Policy set ID (required)")
			flags.StringVar(&cmd.name, "name", "", "Policy set name")
			flags.StringVar(&cmd.description, "description", "", "Policy set description")
			flags.StringVar(&cmd.global, "global", "", "Apply to all workspaces (true/false)")
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

			// Verify the global flag was set correctly
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
