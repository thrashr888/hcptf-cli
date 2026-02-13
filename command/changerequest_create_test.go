package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestChangeRequestCreateRequiresFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ChangeRequestCreateCommand{
		Meta: newTestMeta(ui),
	}

	// Test missing organization
	if code := cmd.Run(nil); code != 1 {
		t.Fatalf("expected exit 1 missing org, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-organization") {
		t.Fatalf("expected org error, got %q", ui.ErrorWriter.String())
	}

	// Test missing workspace
	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org"}); code != 1 {
		t.Fatalf("expected exit 1 missing workspace, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-workspace") {
		t.Fatalf("expected workspace error, got %q", ui.ErrorWriter.String())
	}

	// Test missing subject
	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod"}); code != 1 {
		t.Fatalf("expected exit 1 missing subject, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-subject") {
		t.Fatalf("expected subject error, got %q", ui.ErrorWriter.String())
	}

	// Test missing message
	ui.ErrorWriter.Reset()
	if code := cmd.Run([]string{"-organization=my-org", "-workspace=prod", "-subject=test"}); code != 1 {
		t.Fatalf("expected exit 1 missing message, got %d", code)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "-message") {
		t.Fatalf("expected message error, got %q", ui.ErrorWriter.String())
	}
}

func TestChangeRequestCreateHelp(t *testing.T) {
	cmd := &ChangeRequestCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf changerequest create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
	}
	if !strings.Contains(help, "-workspace") {
		t.Error("Help should mention -workspace flag")
	}
	if !strings.Contains(help, "-subject") {
		t.Error("Help should mention -subject flag")
	}
	if !strings.Contains(help, "-message") {
		t.Error("Help should mention -message flag")
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

func TestChangeRequestCreateSynopsis(t *testing.T) {
	cmd := &ChangeRequestCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a new change request for a workspace" {
		t.Errorf("expected 'Create a new change request for a workspace', got %q", synopsis)
	}
}

func TestChangeRequestCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedOrg      string
		expectedWorkspace string
		expectedSubject  string
		expectedMessage  string
		expectedFormat   string
	}{
		{
			name:             "all required flags with default format",
			args:             []string{"-organization=my-org", "-workspace=prod", "-subject=test subject", "-message=test message"},
			expectedOrg:      "my-org",
			expectedWorkspace: "prod",
			expectedSubject:  "test subject",
			expectedMessage:  "test message",
			expectedFormat:   "table",
		},
		{
			name:             "org alias flag",
			args:             []string{"-org=test-org", "-workspace=staging", "-subject=fix bug", "-message=urgent fix"},
			expectedOrg:      "test-org",
			expectedWorkspace: "staging",
			expectedSubject:  "fix bug",
			expectedMessage:  "urgent fix",
			expectedFormat:   "table",
		},
		{
			name:             "all flags with json format",
			args:             []string{"-organization=my-org", "-workspace=dev", "-subject=update", "-message=details", "-output=json"},
			expectedOrg:      "my-org",
			expectedWorkspace: "dev",
			expectedSubject:  "update",
			expectedMessage:  "details",
			expectedFormat:   "json",
		},
		{
			name:             "org alias with json format",
			args:             []string{"-org=acme", "-workspace=prod", "-subject=security", "-message=patch required", "-output=json"},
			expectedOrg:      "acme",
			expectedWorkspace: "prod",
			expectedSubject:  "security",
			expectedMessage:  "patch required",
			expectedFormat:   "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ChangeRequestCreateCommand{}

			flags := cmd.Meta.FlagSet("changerequest create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.workspace, "workspace", "", "Workspace name (required)")
			flags.StringVar(&cmd.subject, "subject", "", "Change request subject (required)")
			flags.StringVar(&cmd.message, "message", "", "Change request message (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the workspace was set correctly
			if cmd.workspace != tt.expectedWorkspace {
				t.Errorf("expected workspace %q, got %q", tt.expectedWorkspace, cmd.workspace)
			}

			// Verify the subject was set correctly
			if cmd.subject != tt.expectedSubject {
				t.Errorf("expected subject %q, got %q", tt.expectedSubject, cmd.subject)
			}

			// Verify the message was set correctly
			if cmd.message != tt.expectedMessage {
				t.Errorf("expected message %q, got %q", tt.expectedMessage, cmd.message)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
