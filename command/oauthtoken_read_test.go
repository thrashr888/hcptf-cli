package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOAuthTokenReadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OAuthTokenReadCommand{
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

func TestOAuthTokenReadHelp(t *testing.T) {
	cmd := &OAuthTokenReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf oauthtoken read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
}

func TestOAuthTokenReadSynopsis(t *testing.T) {
	cmd := &OAuthTokenReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read OAuth token details" {
		t.Errorf("expected 'Read OAuth token details', got %q", synopsis)
	}
}

func TestOAuthTokenReadFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedFmt string
	}{
		{
			name:        "id only, default format",
			args:        []string{"-id=ot-hmAyP66qk2AMVdbJ"},
			expectedID:  "ot-hmAyP66qk2AMVdbJ",
			expectedFmt: "table",
		},
		{
			name:        "id with table format",
			args:        []string{"-id=ot-ABC123XYZ456", "-output=table"},
			expectedID:  "ot-ABC123XYZ456",
			expectedFmt: "table",
		},
		{
			name:        "id with json format",
			args:        []string{"-id=ot-test12345678", "-output=json"},
			expectedID:  "ot-test12345678",
			expectedFmt: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OAuthTokenReadCommand{}

			flags := cmd.Meta.FlagSet("oauthtoken read")
			flags.StringVar(&cmd.id, "id", "", "OAuth token ID (required)")
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
