package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAzureOIDCUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AzureoidcUpdateCommand{
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

func TestAzureOIDCUpdateValidatesEmptyID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AzureoidcUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=", "-client-id=12345678-1234-1234-1234-123456789012"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestAzureOIDCUpdateHelp(t *testing.T) {
	cmd := &AzureoidcUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf azureoidc update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-client-id") {
		t.Error("Help should mention -client-id flag")
	}
	if !strings.Contains(help, "-subscription-id") {
		t.Error("Help should mention -subscription-id flag")
	}
	if !strings.Contains(help, "-tenant-id") {
		t.Error("Help should mention -tenant-id flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "Azure OIDC") {
		t.Error("Help should describe Azure OIDC configuration")
	}
}

func TestAzureOIDCUpdateSynopsis(t *testing.T) {
	cmd := &AzureoidcUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update Azure OIDC configuration settings" {
		t.Errorf("expected 'Update Azure OIDC configuration settings', got %q", synopsis)
	}
}

func TestAzureOIDCUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedClient string
		expectedSub    string
		expectedTenant string
		expectedFmt    string
	}{
		{
			name:           "id only",
			args:           []string{"-id=azoidc-ABC123"},
			expectedID:     "azoidc-ABC123",
			expectedClient: "",
			expectedSub:    "",
			expectedTenant: "",
			expectedFmt:    "table",
		},
		{
			name: "id with client-id",
			args: []string{
				"-id=azoidc-ABC123",
				"-client-id=12345678-1234-1234-1234-123456789012",
			},
			expectedID:     "azoidc-ABC123",
			expectedClient: "12345678-1234-1234-1234-123456789012",
			expectedSub:    "",
			expectedTenant: "",
			expectedFmt:    "table",
		},
		{
			name: "id with subscription-id",
			args: []string{
				"-id=azoidc-XYZ789",
				"-subscription-id=87654321-4321-4321-4321-210987654321",
			},
			expectedID:     "azoidc-XYZ789",
			expectedClient: "",
			expectedSub:    "87654321-4321-4321-4321-210987654321",
			expectedTenant: "",
			expectedFmt:    "table",
		},
		{
			name: "id with tenant-id",
			args: []string{
				"-id=azoidc-DEF456",
				"-tenant-id=abcdefab-abcd-abcd-abcd-abcdefabcdef",
			},
			expectedID:     "azoidc-DEF456",
			expectedClient: "",
			expectedSub:    "",
			expectedTenant: "abcdefab-abcd-abcd-abcd-abcdefabcdef",
			expectedFmt:    "table",
		},
		{
			name: "all flags with json output",
			args: []string{
				"-id=azoidc-GHI789",
				"-client-id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				"-subscription-id=11111111-2222-3333-4444-555555555555",
				"-tenant-id=66666666-7777-8888-9999-000000000000",
				"-output=json",
			},
			expectedID:     "azoidc-GHI789",
			expectedClient: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			expectedSub:    "11111111-2222-3333-4444-555555555555",
			expectedTenant: "66666666-7777-8888-9999-000000000000",
			expectedFmt:    "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AzureoidcUpdateCommand{}

			flags := cmd.Meta.FlagSet("azureoidc update")
			flags.StringVar(&cmd.id, "id", "", "Azure OIDC configuration ID (required)")
			flags.StringVar(&cmd.clientID, "client-id", "", "Azure application (client) ID")
			flags.StringVar(&cmd.subscriptionID, "subscription-id", "", "Azure subscription ID")
			flags.StringVar(&cmd.tenantID, "tenant-id", "", "Azure tenant (directory) ID")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the ID was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected ID %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the client ID was set correctly
			if cmd.clientID != tt.expectedClient {
				t.Errorf("expected client ID %q, got %q", tt.expectedClient, cmd.clientID)
			}

			// Verify the subscription ID was set correctly
			if cmd.subscriptionID != tt.expectedSub {
				t.Errorf("expected subscription ID %q, got %q", tt.expectedSub, cmd.subscriptionID)
			}

			// Verify the tenant ID was set correctly
			if cmd.tenantID != tt.expectedTenant {
				t.Errorf("expected tenant ID %q, got %q", tt.expectedTenant, cmd.tenantID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
