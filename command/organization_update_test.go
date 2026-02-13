package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationUpdateRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-email=test@example.com"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestOrganizationUpdateHelp(t *testing.T) {
	cmd := &OrganizationUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf organization update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -name is required")
	}
}

func TestOrganizationUpdateSynopsis(t *testing.T) {
	cmd := &OrganizationUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update organization settings" {
		t.Errorf("expected 'Update organization settings', got %q", synopsis)
	}
}

func TestOrganizationUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name                  string
		args                  []string
		expectedName          string
		expectedEmail         string
		expectedSessionTime   int
		expectedSessionRemem  int
		expectedCostEst       string
		expectedFormat        string
	}{
		{
			name:         "name only",
			args:         []string{"-name=test-org"},
			expectedName: "test-org",
			expectedFormat: "table",
		},
		{
			name:          "name and email",
			args:          []string{"-name=my-org", "-email=new@example.com"},
			expectedName:  "my-org",
			expectedEmail: "new@example.com",
			expectedFormat: "table",
		},
		{
			name:                "name and session timeout",
			args:                []string{"-name=my-org", "-session-timeout=20160"},
			expectedName:        "my-org",
			expectedSessionTime: 20160,
			expectedFormat: "table",
		},
		{
			name:                 "name and session remember",
			args:                 []string{"-name=my-org", "-session-remember=20160"},
			expectedName:         "my-org",
			expectedSessionRemem: 20160,
			expectedFormat: "table",
		},
		{
			name:            "name and cost estimation",
			args:            []string{"-name=my-org", "-cost-estimation=true"},
			expectedName:    "my-org",
			expectedCostEst: "true",
			expectedFormat: "table",
		},
		{
			name:           "name with json output",
			args:           []string{"-name=my-org", "-output=json"},
			expectedName:   "my-org",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &OrganizationUpdateCommand{}

			flags := cmd.Meta.FlagSet("organization update")
			flags.StringVar(&cmd.name, "name", "", "Organization name (required)")
			flags.StringVar(&cmd.email, "email", "", "Admin email address")
			flags.IntVar(&cmd.sessionTimeout, "session-timeout", 0, "Session timeout in minutes")
			flags.IntVar(&cmd.sessionRemember, "session-remember", 0, "Session remember duration in minutes")
			flags.StringVar(&cmd.costEstimationEnabled, "cost-estimation", "", "Enable cost estimation (true/false)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the email was set correctly
			if cmd.email != tt.expectedEmail {
				t.Errorf("expected email %q, got %q", tt.expectedEmail, cmd.email)
			}

			// Verify session timeout
			if cmd.sessionTimeout != tt.expectedSessionTime {
				t.Errorf("expected session timeout %d, got %d", tt.expectedSessionTime, cmd.sessionTimeout)
			}

			// Verify session remember
			if cmd.sessionRemember != tt.expectedSessionRemem {
				t.Errorf("expected session remember %d, got %d", tt.expectedSessionRemem, cmd.sessionRemember)
			}

			// Verify cost estimation
			if cmd.costEstimationEnabled != tt.expectedCostEst {
				t.Errorf("expected cost estimation %q, got %q", tt.expectedCostEst, cmd.costEstimationEnabled)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
