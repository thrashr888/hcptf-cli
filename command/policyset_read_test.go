package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestPolicySetReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &PolicySetReadCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestPolicySetReadHelp(t *testing.T) {
	cmd := &PolicySetReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf policyset read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestPolicySetReadSynopsis(t *testing.T) {
	cmd := &PolicySetReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read policy set details" {
		t.Errorf("expected 'Read policy set details', got %q", synopsis)
	}
}

func TestPolicySetReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id only, default format",
			args:        []string{"-id=polset-12345"},
			expectedID:  "polset-12345",
			expectedFmt: "table",
		},
		{
			name:        "id with table format",
			args:        []string{"-id=polset-67890", "-output=table"},
			expectedID:  "polset-67890",
			expectedFmt: "table",
		},
		{
			name:        "id with json format",
			args:        []string{"-id=polset-abcde", "-output=json"},
			expectedID:  "polset-abcde",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &PolicySetReadCommand{}

			flags := cmd.Meta.FlagSet("policyset read")
			flags.StringVar(&cmd.id, "id", "", "Policy set ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
