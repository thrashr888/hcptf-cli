package command

import (
	"strings"
	"testing"
)

func TestAuditTrailReadHelp(t *testing.T) {
	cmd := &AuditTrailReadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf audittrail read") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "Audit trail event ID") {
		t.Error("Help should mention audit trail event ID")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "table (default) or json") {
		t.Error("Help should mention output formats")
	}
	if !strings.Contains(help, "organization token or audit trail token") {
		t.Error("Help should mention required token types")
	}
}

func TestAuditTrailReadSynopsis(t *testing.T) {
	cmd := &AuditTrailReadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Read audit trail event details" {
		t.Errorf("expected 'Read audit trail event details', got %q", synopsis)
	}
}

func TestAuditTrailReadFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedFormat string
	}{
		{
			name:           "default output format",
			args:           []string{"-id=ae66e491-db59-457c-8445-9c908ee726ae"},
			expectedID:     "ae66e491-db59-457c-8445-9c908ee726ae",
			expectedFormat: "table",
		},
		{
			name:           "json output format",
			args:           []string{"-id=ae66e491-db59-457c-8445-9c908ee726ae", "-output=json"},
			expectedID:     "ae66e491-db59-457c-8445-9c908ee726ae",
			expectedFormat: "json",
		},
		{
			name:           "table output format explicitly set",
			args:           []string{"-id=test-id-123", "-output=table"},
			expectedID:     "test-id-123",
			expectedFormat: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AuditTrailReadCommand{}

			flags := cmd.Meta.FlagSet("audittrail read")
			flags.StringVar(&cmd.id, "id", "", "Audit trail event ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
