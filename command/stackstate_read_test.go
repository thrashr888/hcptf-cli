package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestStackStateReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &StackStateReadCommand{
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

func TestStackStateReadHelp(t *testing.T) {
	cmd := &StackStateReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf stackstate read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
}

func TestStackStateReadSynopsis(t *testing.T) {
	cmd := &StackStateReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read stack state details" {
		t.Errorf("expected 'Read stack state details', got %q", synopsis)
	}
}

func TestStackStateReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id flag only",
			args:        []string{"-id=sts-abc123"},
			expectedID:  "sts-abc123",
			expectedFmt: "table",
		},
		{
			name:        "with json output",
			args:        []string{"-id=sts-xyz789", "-output=json"},
			expectedID:  "sts-xyz789",
			expectedFmt: "json",
		},
		{
			name:        "with table output",
			args:        []string{"-id=sts-def456", "-output=table"},
			expectedID:  "sts-def456",
			expectedFmt: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &StackStateReadCommand{}

			flags := cmd.Meta.FlagSet("stackstate read")
			flags.StringVar(&cmd.stateID, "id", "", "Stack state ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the state ID was set correctly
			if cmd.stateID != tt.expectedID {
				t.Errorf("expected state ID %q, got %q", tt.expectedID, cmd.stateID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
