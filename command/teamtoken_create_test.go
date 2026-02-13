package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestTeamTokenCreateRequiresTeamID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamTokenCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-description=test-token"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-team-id") {
		t.Fatalf("expected team-id error, got %q", out)
	}
}

func TestTeamTokenCreateRequiresDescription(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamTokenCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-team-id=team-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-description") {
		t.Fatalf("expected description error, got %q", out)
	}
}

func TestTeamTokenCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &TeamTokenCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "required") {
		t.Fatalf("expected required error, got %q", out)
	}
}

func TestTeamTokenCreateHelp(t *testing.T) {
	cmd := &TeamTokenCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf teamtoken create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-team-id") {
		t.Error("Help should mention -team-id flag")
	}
	if !strings.Contains(help, "-description") {
		t.Error("Help should mention -description flag")
	}
	if !strings.Contains(help, "-expired-at") {
		t.Error("Help should mention -expired-at flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestTeamTokenCreateSynopsis(t *testing.T) {
	cmd := &TeamTokenCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a team API token" {
		t.Errorf("expected 'Create a team API token', got %q", synopsis)
	}
}

func TestTeamTokenCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedTeamID string
		expectedDesc   string
		expectedExpiry string
		expectedFmt    string
	}{
		{
			name:           "all required flags, default format",
			args:           []string{"-team-id=team-123abc", "-description=Production team token"},
			expectedTeamID: "team-123abc",
			expectedDesc:   "Production team token",
			expectedExpiry: "",
			expectedFmt:    "table",
		},
		{
			name:           "required flags with table format",
			args:           []string{"-team-id=team-456def", "-description=Dev team", "-output=table"},
			expectedTeamID: "team-456def",
			expectedDesc:   "Dev team",
			expectedExpiry: "",
			expectedFmt:    "table",
		},
		{
			name:           "required flags with json format",
			args:           []string{"-team-id=team-789ghi", "-description=CI team token", "-output=json"},
			expectedTeamID: "team-789ghi",
			expectedDesc:   "CI team token",
			expectedExpiry: "",
			expectedFmt:    "json",
		},
		{
			name:           "with expiration date",
			args:           []string{"-team-id=team-abc", "-description=Temporary token", "-expired-at=2024-12-31T23:59:59Z"},
			expectedTeamID: "team-abc",
			expectedDesc:   "Temporary token",
			expectedExpiry: "2024-12-31T23:59:59Z",
			expectedFmt:    "table",
		},
		{
			name:           "all flags with json format",
			args:           []string{"-team-id=team-xyz", "-description=Full test", "-expired-at=2025-06-30T12:00:00Z", "-output=json"},
			expectedTeamID: "team-xyz",
			expectedDesc:   "Full test",
			expectedExpiry: "2025-06-30T12:00:00Z",
			expectedFmt:    "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &TeamTokenCreateCommand{}

			flags := cmd.Meta.FlagSet("teamtoken create")
			flags.StringVar(&cmd.teamID, "team-id", "", "Team ID (required)")
			flags.StringVar(&cmd.description, "description", "", "Token description (required)")
			flags.StringVar(&cmd.expiredAt, "expired-at", "", "Expiration date in ISO 8601 format (e.g., 2024-12-31T23:59:59Z)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the team ID was set correctly
			if cmd.teamID != tt.expectedTeamID {
				t.Errorf("expected teamID %q, got %q", tt.expectedTeamID, cmd.teamID)
			}

			// Verify the description was set correctly
			if cmd.description != tt.expectedDesc {
				t.Errorf("expected description %q, got %q", tt.expectedDesc, cmd.description)
			}

			// Verify the expired-at was set correctly
			if cmd.expiredAt != tt.expectedExpiry {
				t.Errorf("expected expiredAt %q, got %q", tt.expectedExpiry, cmd.expiredAt)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
