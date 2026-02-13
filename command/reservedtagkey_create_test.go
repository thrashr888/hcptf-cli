package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestReservedTagKeyCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ReservedTagKeyCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestReservedTagKeyCreateRequiresKey(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ReservedTagKeyCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-key") {
		t.Fatalf("expected key error, got %q", out)
	}
}

func TestReservedTagKeyCreateHelp(t *testing.T) {
	cmd := &ReservedTagKeyCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf reservedtagkey create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
	}
	if !strings.Contains(help, "-key") {
		t.Error("Help should mention -key flag")
	}
	if !strings.Contains(help, "-disable-overrides") {
		t.Error("Help should mention -disable-overrides flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "Reserved tag keys") {
		t.Error("Help should describe reserved tag keys")
	}
}

func TestReservedTagKeyCreateSynopsis(t *testing.T) {
	cmd := &ReservedTagKeyCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a reserved tag key" {
		t.Errorf("expected 'Create a reserved tag key', got %q", synopsis)
	}
}

func TestReservedTagKeyCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name                     string
		args                     []string
		expectedOrg              string
		expectedKey              string
		expectedDisableOverrides bool
		expectedFmt              string
	}{
		{
			name:                     "organization and key flags",
			args:                     []string{"-organization=my-org", "-key=environment"},
			expectedOrg:              "my-org",
			expectedKey:              "environment",
			expectedDisableOverrides: false,
			expectedFmt:              "table",
		},
		{
			name:                     "org alias flag",
			args:                     []string{"-org=test-org", "-key=cost-center"},
			expectedOrg:              "test-org",
			expectedKey:              "cost-center",
			expectedDisableOverrides: false,
			expectedFmt:              "table",
		},
		{
			name:                     "with disable-overrides flag",
			args:                     []string{"-organization=my-org", "-key=team", "-disable-overrides"},
			expectedOrg:              "my-org",
			expectedKey:              "team",
			expectedDisableOverrides: true,
			expectedFmt:              "table",
		},
		{
			name:                     "with json output",
			args:                     []string{"-organization=my-org", "-key=department", "-output=json"},
			expectedOrg:              "my-org",
			expectedKey:              "department",
			expectedDisableOverrides: false,
			expectedFmt:              "json",
		},
		{
			name:                     "all options",
			args:                     []string{"-org=test-org", "-key=region", "-disable-overrides", "-output=json"},
			expectedOrg:              "test-org",
			expectedKey:              "region",
			expectedDisableOverrides: true,
			expectedFmt:              "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ReservedTagKeyCreateCommand{}

			flags := cmd.Meta.FlagSet("reservedtagkey create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.key, "key", "", "Tag key to reserve (required)")
			flags.BoolVar(&cmd.disableOverrides, "disable-overrides", false, "Disable overriding inherited tags at workspace level")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the key was set correctly
			if cmd.key != tt.expectedKey {
				t.Errorf("expected key %q, got %q", tt.expectedKey, cmd.key)
			}

			// Verify the disable-overrides was set correctly
			if cmd.disableOverrides != tt.expectedDisableOverrides {
				t.Errorf("expected disableOverrides %v, got %v", tt.expectedDisableOverrides, cmd.disableOverrides)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
