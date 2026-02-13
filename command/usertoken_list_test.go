package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestUserTokenListHelp(t *testing.T) {
	cmd := &UserTokenListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf usertoken list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
}

func TestUserTokenListSynopsis(t *testing.T) {
	cmd := &UserTokenListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "List user API tokens" {
		t.Errorf("expected 'List user API tokens', got %q", synopsis)
	}
}

func TestUserTokenListFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedFmt string
	}{
		{
			name:        "default format",
			args:        []string{},
			expectedFmt: "table",
		},
		{
			name:        "table format",
			args:        []string{"-output=table"},
			expectedFmt: "table",
		},
		{
			name:        "json format",
			args:        []string{"-output=json"},
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &UserTokenListCommand{
				Meta: newTestMeta(cli.NewMockUi()),
			}

			flags := cmd.Meta.FlagSet("usertoken list")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}

