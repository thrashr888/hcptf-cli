package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestUserTokenReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &UserTokenReadCommand{
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

func TestUserTokenReadHelp(t *testing.T) {
	cmd := &UserTokenReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf usertoken read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestUserTokenReadSynopsis(t *testing.T) {
	cmd := &UserTokenReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Show user token details" {
		t.Errorf("expected 'Show user token details', got %q", synopsis)
	}
}

func TestUserTokenReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "required flag, default format",
			args:        []string{"-id=at-abc123xyz"},
			expectedID:  "at-abc123xyz",
			expectedFmt: "table",
		},
		{
			name:        "required flag with table format",
			args:        []string{"-id=at-def456uvw", "-output=table"},
			expectedID:  "at-def456uvw",
			expectedFmt: "table",
		},
		{
			name:        "required flag with json format",
			args:        []string{"-id=at-ghi789rst", "-output=json"},
			expectedID:  "at-ghi789rst",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &UserTokenReadCommand{
				Meta: newTestMeta(cli.NewMockUi()),
			}

			flags := cmd.Meta.FlagSet("usertoken read")
			flags.StringVar(&cmd.id, "id", "", "User token ID (required)")
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
