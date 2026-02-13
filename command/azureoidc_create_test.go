package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAzureOIDCCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AzureoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestAzureOIDCCreateRequiresClientID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AzureoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-client-id") {
		t.Fatalf("expected client-id error, got %q", out)
	}
}

func TestAzureOIDCCreateRequiresSubscriptionID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AzureoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{
		"-organization=my-org",
		"-client-id=12345678-1234-1234-1234-123456789012",
	})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-subscription-id") {
		t.Fatalf("expected subscription-id error, got %q", out)
	}
}

func TestAzureOIDCCreateRequiresTenantID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AzureoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{
		"-organization=my-org",
		"-client-id=12345678-1234-1234-1234-123456789012",
		"-subscription-id=87654321-4321-4321-4321-210987654321",
	})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-tenant-id") {
		t.Fatalf("expected tenant-id error, got %q", out)
	}
}

func TestAzureOIDCCreateRequiresEmptyClientID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AzureoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-client-id=", "-subscription-id=87654321-4321-4321-4321-210987654321", "-tenant-id=abcd-abcd"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-client-id") {
		t.Fatalf("expected client-id error, got %q", out)
	}
}

func TestAzureOIDCCreateRequiresEmptyTenantID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AzureoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-client-id=12345678-1234-1234-1234-123456789012", "-subscription-id=87654321-4321-4321-4321-210987654321", "-tenant-id="})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-tenant-id") {
		t.Fatalf("expected tenant-id error, got %q", out)
	}
}

func TestAzureOIDCCreateHelp(t *testing.T) {
	cmd := &AzureoidcCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf azureoidc create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
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

func TestAzureOIDCCreateSynopsis(t *testing.T) {
	cmd := &AzureoidcCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create an Azure OIDC configuration for dynamic credentials" {
		t.Errorf("expected 'Create an Azure OIDC configuration for dynamic credentials', got %q", synopsis)
	}
}

func TestAzureOIDCCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOrg    string
		expectedClient string
		expectedSub    string
		expectedTenant string
		expectedFmt    string
	}{
		{
			name: "all flags",
			args: []string{
				"-organization=my-org",
				"-client-id=12345678-1234-1234-1234-123456789012",
				"-subscription-id=87654321-4321-4321-4321-210987654321",
				"-tenant-id=abcdefab-abcd-abcd-abcd-abcdefabcdef",
			},
			expectedOrg:    "my-org",
			expectedClient: "12345678-1234-1234-1234-123456789012",
			expectedSub:    "87654321-4321-4321-4321-210987654321",
			expectedTenant: "abcdefab-abcd-abcd-abcd-abcdefabcdef",
			expectedFmt:    "table",
		},
		{
			name: "org alias flag",
			args: []string{
				"-org=test-org",
				"-client-id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				"-subscription-id=11111111-2222-3333-4444-555555555555",
				"-tenant-id=66666666-7777-8888-9999-000000000000",
			},
			expectedOrg:    "test-org",
			expectedClient: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			expectedSub:    "11111111-2222-3333-4444-555555555555",
			expectedTenant: "66666666-7777-8888-9999-000000000000",
			expectedFmt:    "table",
		},
		{
			name: "with json output",
			args: []string{
				"-organization=my-org",
				"-client-id=12345678-1234-1234-1234-123456789012",
				"-subscription-id=87654321-4321-4321-4321-210987654321",
				"-tenant-id=abcdefab-abcd-abcd-abcd-abcdefabcdef",
				"-output=json",
			},
			expectedOrg:    "my-org",
			expectedClient: "12345678-1234-1234-1234-123456789012",
			expectedSub:    "87654321-4321-4321-4321-210987654321",
			expectedTenant: "abcdefab-abcd-abcd-abcd-abcdefabcdef",
			expectedFmt:    "json",
		},
		{
			name: "with table output",
			args: []string{
				"-org=test-org",
				"-client-id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				"-subscription-id=11111111-2222-3333-4444-555555555555",
				"-tenant-id=66666666-7777-8888-9999-000000000000",
				"-output=table",
			},
			expectedOrg:    "test-org",
			expectedClient: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			expectedSub:    "11111111-2222-3333-4444-555555555555",
			expectedTenant: "66666666-7777-8888-9999-000000000000",
			expectedFmt:    "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AzureoidcCreateCommand{}

			flags := cmd.Meta.FlagSet("azureoidc create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.clientID, "client-id", "", "Azure application (client) ID (required)")
			flags.StringVar(&cmd.subscriptionID, "subscription-id", "", "Azure subscription ID (required)")
			flags.StringVar(&cmd.tenantID, "tenant-id", "", "Azure tenant (directory) ID (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
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







func TestAzureOIDCCreatePartialRequiredFlagsRun(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError string
	}{
		{
			name:        "missing subscription-id and tenant-id",
			args:        []string{"-organization=my-org", "-client-id=12345678-1234-1234-1234-123456789012"},
			expectError: "-subscription-id",
		},
		{
			name: "missing tenant-id only",
			args: []string{
				"-organization=my-org",
				"-client-id=12345678-1234-1234-1234-123456789012",
				"-subscription-id=87654321-4321-4321-4321-210987654321",
			},
			expectError: "-tenant-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := cli.NewMockUi()
			cmd := &AzureoidcCreateCommand{
				Meta: newTestMeta(ui),
			}

			code := cmd.Run(tt.args)
			if code != 1 {
				t.Fatalf("expected exit 1, got %d", code)
			}

			if out := ui.ErrorWriter.String(); !strings.Contains(out, tt.expectError) {
				t.Fatalf("expected error containing %q, got %q", tt.expectError, out)
			}
		})
	}
}

