package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackUpdateCommand{
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

func TestStackUpdateValidatesSpeculativeEnabled(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=st-abc123", "-speculative-enabled=maybe"})
	if code != 1 {
		t.Fatalf("expected exit 1 for invalid speculative-enabled, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "speculative-enabled") {
		t.Fatalf("expected speculative-enabled error, got %q", out)
	}
}

func TestStackUpdateHelp(t *testing.T) {
	cmd := &StackUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stack update") {
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
	if !strings.Contains(help, "-speculative-enabled") {
		t.Error("Help should mention -speculative-enabled flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
}

func TestStackUpdateSynopsis(t *testing.T) {
	cmd := &StackUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update stack settings" {
		t.Errorf("expected 'Update stack settings', got %q", synopsis)
	}
}

func TestStackUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name                      string
		args                      []string
		expectedID                string
		expectedName              string
		expectedDesc              string
		expectedSpeculativeEnable *bool
		expectedFmt               string
	}{
		{
			name:                      "id and name",
			args:                      []string{"-id=st-abc123", "-name=new-name"},
			expectedID:                "st-abc123",
			expectedName:              "new-name",
			expectedDesc:              "",
			expectedSpeculativeEnable: nil,
			expectedFmt:               "table",
		},
		{
			name:                      "id and description",
			args:                      []string{"-id=st-xyz789", "-description=Updated description"},
			expectedID:                "st-xyz789",
			expectedName:              "",
			expectedDesc:              "Updated description",
			expectedSpeculativeEnable: nil,
			expectedFmt:               "table",
		},
		{
			name:                      "id, name, description",
			args:                      []string{"-id=st-test456", "-name=updated-stack", "-description=New description"},
			expectedID:                "st-test456",
			expectedName:              "updated-stack",
			expectedDesc:              "New description",
			expectedSpeculativeEnable: nil,
			expectedFmt:               "table",
		},
		{
			name:                      "id, name, json format",
			args:                      []string{"-id=st-prod123", "-name=prod-stack", "-output=json"},
			expectedID:                "st-prod123",
			expectedName:              "prod-stack",
			expectedDesc:              "",
			expectedSpeculativeEnable: nil,
			expectedFmt:               "json",
		},
		{
			name:                      "id and speculative-enabled true",
			args:                      []string{"-id=st-spec123", "-speculative-enabled=true"},
			expectedID:                "st-spec123",
			expectedName:              "",
			expectedDesc:              "",
			expectedSpeculativeEnable: func() *bool { b := true; return &b }(),
			expectedFmt:               "table",
		},
		{
			name:                      "id and speculative-enabled false",
			args:                      []string{"-id=st-spec456", "-speculative-enabled=false"},
			expectedID:                "st-spec456",
			expectedName:              "",
			expectedDesc:              "",
			expectedSpeculativeEnable: func() *bool { b := false; return &b }(),
			expectedFmt:               "table",
		},
		{
			name:                      "all flags",
			args:                      []string{"-id=st-all123", "-name=complete-stack", "-description=All flags test", "-speculative-enabled=true", "-output=json"},
			expectedID:                "st-all123",
			expectedName:              "complete-stack",
			expectedDesc:              "All flags test",
			expectedSpeculativeEnable: func() *bool { b := true; return &b }(),
			expectedFmt:               "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackUpdateCommand{}

			flags := cmd.Meta.FlagSet("stack update")
			flags.StringVar(&cmd.stackID, "id", "", "Stack ID (required)")
			flags.StringVar(&cmd.name, "name", "", "New stack name")
			flags.StringVar(&cmd.description, "description", "", "New stack description")

			// Use a string flag for boolean to distinguish between set and unset
			var speculativeEnabledStr string
			flags.StringVar(&speculativeEnabledStr, "speculative-enabled", "", "Enable/disable speculative plans (true/false)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Parse speculative-enabled flag if provided
			if speculativeEnabledStr != "" {
				if speculativeEnabledStr == "true" {
					val := true
					cmd.speculativeEnabled = &val
				} else if speculativeEnabledStr == "false" {
					val := false
					cmd.speculativeEnabled = &val
				}
			}

			// Verify the id was set correctly
			if cmd.stackID != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.stackID)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the description was set correctly
			if cmd.description != tt.expectedDesc {
				t.Errorf("expected description %q, got %q", tt.expectedDesc, cmd.description)
			}

			// Verify the speculative-enabled was set correctly
			if tt.expectedSpeculativeEnable == nil {
				if cmd.speculativeEnabled != nil {
					t.Errorf("expected speculative-enabled to be nil, got %v", *cmd.speculativeEnabled)
				}
			} else {
				if cmd.speculativeEnabled == nil {
					t.Errorf("expected speculative-enabled to be %v, got nil", *tt.expectedSpeculativeEnable)
				} else if *cmd.speculativeEnabled != *tt.expectedSpeculativeEnable {
					t.Errorf("expected speculative-enabled %v, got %v", *tt.expectedSpeculativeEnable, *cmd.speculativeEnabled)
				}
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
