package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestChangeRequestUpdateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ChangeRequestUpdateCommand{
		Meta: newTestMeta(ui),
	}

	// Test missing id
	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing id, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-id") {
		t.Fatalf("expected id error, got %q", ui.ErrorWriter.String())
	}

	// Test missing archive flag
	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-id=wscr-123"}); code != 1 {
		t.Fatalf("expected exit 1 missing archive, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-archive") {
		t.Fatalf("expected archive error, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestUpdateHelp(t *testing.T) {
	cmd := &ChangeRequestUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf changerequest update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-archive") {
		t.Error("Help should mention -archive flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "HCP Terraform Plus or Enterprise") {
		t.Error("Help should mention plan requirements")
	}
}

func TestChangeRequestUpdateSynopsis(t *testing.T) {
	cmd := &ChangeRequestUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update a change request (archive)" {
		t.Errorf("expected 'Update a change request (archive)', got %q", synopsis)
	}
}

func TestChangeRequestUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedID      string
		expectedArchive bool
		expectedFormat  string
	}{
		{
			name:            "id and archive with default format",
			args:            []string{"-id=wscr-abc123", "-archive"},
			expectedID:      "wscr-abc123",
			expectedArchive: true,
			expectedFormat:  "table",
		},
		{
			name:            "id and archive with table format",
			args:            []string{"-id=wscr-xyz789", "-archive", "-output=table"},
			expectedID:      "wscr-xyz789",
			expectedArchive: true,
			expectedFormat:  "table",
		},
		{
			name:            "id and archive with json format",
			args:            []string{"-id=wscr-test456", "-archive", "-output=json"},
			expectedID:      "wscr-test456",
			expectedArchive: true,
			expectedFormat:  "json",
		},
		{
			name:            "archive with different id format",
			args:            []string{"-id=wscr-prod999", "-archive"},
			expectedID:      "wscr-prod999",
			expectedArchive: true,
			expectedFormat:  "table",
		},
		{
			name:            "explicit archive true with json",
			args:            []string{"-id=wscr-staging", "-archive=true", "-output=json"},
			expectedID:      "wscr-staging",
			expectedArchive: true,
			expectedFormat:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ChangeRequestUpdateCommand{}

			flags := cmd.Meta.FlagSet("changerequest update")
			flags.StringVar(&cmd.id, "id", "", "Change request ID (required)")
			flags.BoolVar(&cmd.archive, "archive", false, "Archive the change request")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the archive flag was set correctly
			if cmd.archive != tt.expectedArchive {
				t.Errorf("expected archive %v, got %v", tt.expectedArchive, cmd.archive)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
