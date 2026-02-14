package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestGCPOIDCCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-service-account-email=test@project.iam.gserviceaccount.com", "-workload-provider-name=projects/123/locations/global/workloadIdentityPools/pool/providers/provider", "-project-number=123456"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestGCPOIDCCreateRequiresServiceAccountEmail(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-workload-provider-name=projects/123/locations/global/workloadIdentityPools/pool/providers/provider", "-project-number=123456"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-service-account-email") {
		t.Fatalf("expected service-account-email error, got %q", out)
	}
}

func TestGCPOIDCCreateRequiresWorkloadProviderName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-service-account-email=test@project.iam.gserviceaccount.com", "-project-number=123456"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workload-provider-name") {
		t.Fatalf("expected workload-provider-name error, got %q", out)
	}
}

func TestGCPOIDCCreateRequiresProjectNumber(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-service-account-email=test@project.iam.gserviceaccount.com", "-workload-provider-name=projects/123/locations/global/workloadIdentityPools/pool/providers/provider"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-project-number") {
		t.Fatalf("expected project-number error, got %q", out)
	}
}

func TestGCPOIDCCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcCreateCommand{
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

func TestGCPOIDCCreateRequiresEmptyServiceAccount(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-service-account-email=", "-workload-provider-name=projects/123/locations/global/workloadIdentityPools/pool/providers/provider", "-project-number=123456"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-service-account-email") {
		t.Fatalf("expected service-account-email error, got %q", out)
	}
}

func TestGCPOIDCCreateRequiresEmptyWorkloadProvider(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-service-account-email=test@project.iam.gserviceaccount.com", "-workload-provider-name=", "-project-number=123456"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-workload-provider-name") {
		t.Fatalf("expected workload-provider-name error, got %q", out)
	}
}

func TestGCPOIDCCreateRequiresEmptyProjectNumber(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-service-account-email=test@project.iam.gserviceaccount.com", "-workload-provider-name=projects/123/locations/global/workloadIdentityPools/pool/providers/provider", "-project-number="})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-project-number") {
		t.Fatalf("expected project-number error, got %q", out)
	}
}

func TestGCPOIDCCreateHelp(t *testing.T) {
	cmd := &GCPoidcCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf gcpoidc create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag")
	}
	if !strings.Contains(help, "-service-account-email") {
		t.Error("Help should mention -service-account-email flag")
	}
	if !strings.Contains(help, "-workload-provider-name") {
		t.Error("Help should mention -workload-provider-name flag")
	}
	if !strings.Contains(help, "-project-number") {
		t.Error("Help should mention -project-number flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
}

func TestGCPOIDCCreateSynopsis(t *testing.T) {
	cmd := &GCPoidcCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a GCP OIDC configuration for dynamic credentials" {
		t.Errorf("expected 'Create a GCP OIDC configuration for dynamic credentials', got %q", synopsis)
	}
}

func TestGCPOIDCCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name                     string
		args                     []string
		expectedOrg              string
		expectedServiceAccount   string
		expectedWorkloadProvider string
		expectedProjectNumber    string
		expectedFormat           string
	}{
		{
			name:                     "all required flags, default format",
			args:                     []string{"-organization=my-org", "-service-account-email=terraform@my-project.iam.gserviceaccount.com", "-workload-provider-name=projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider", "-project-number=123456789"},
			expectedOrg:              "my-org",
			expectedServiceAccount:   "terraform@my-project.iam.gserviceaccount.com",
			expectedWorkloadProvider: "projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
			expectedProjectNumber:    "123456789",
			expectedFormat:           "table",
		},
		{
			name:                     "org alias with required flags",
			args:                     []string{"-org=test-org", "-service-account-email=sa@test-project.iam.gserviceaccount.com", "-workload-provider-name=projects/999/locations/global/workloadIdentityPools/test-pool/providers/test-provider", "-project-number=999"},
			expectedOrg:              "test-org",
			expectedServiceAccount:   "sa@test-project.iam.gserviceaccount.com",
			expectedWorkloadProvider: "projects/999/locations/global/workloadIdentityPools/test-pool/providers/test-provider",
			expectedProjectNumber:    "999",
			expectedFormat:           "table",
		},
		{
			name:                     "json output format",
			args:                     []string{"-org=prod-org", "-service-account-email=prod@prod-project.iam.gserviceaccount.com", "-workload-provider-name=projects/111/locations/global/workloadIdentityPools/prod-pool/providers/prod-provider", "-project-number=111", "-output=json"},
			expectedOrg:              "prod-org",
			expectedServiceAccount:   "prod@prod-project.iam.gserviceaccount.com",
			expectedWorkloadProvider: "projects/111/locations/global/workloadIdentityPools/prod-pool/providers/prod-provider",
			expectedProjectNumber:    "111",
			expectedFormat:           "json",
		},
		{
			name:                     "table output format",
			args:                     []string{"-organization=dev-org", "-service-account-email=dev@dev-project.iam.gserviceaccount.com", "-workload-provider-name=projects/222/locations/global/workloadIdentityPools/dev-pool/providers/dev-provider", "-project-number=222", "-output=table"},
			expectedOrg:              "dev-org",
			expectedServiceAccount:   "dev@dev-project.iam.gserviceaccount.com",
			expectedWorkloadProvider: "projects/222/locations/global/workloadIdentityPools/dev-pool/providers/dev-provider",
			expectedProjectNumber:    "222",
			expectedFormat:           "table",
		},
		{
			name:                     "long project number",
			args:                     []string{"-org=staging-org", "-service-account-email=staging@staging-project.iam.gserviceaccount.com", "-workload-provider-name=projects/987654321/locations/global/workloadIdentityPools/staging-pool/providers/staging-provider", "-project-number=987654321"},
			expectedOrg:              "staging-org",
			expectedServiceAccount:   "staging@staging-project.iam.gserviceaccount.com",
			expectedWorkloadProvider: "projects/987654321/locations/global/workloadIdentityPools/staging-pool/providers/staging-provider",
			expectedProjectNumber:    "987654321",
			expectedFormat:           "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &GCPoidcCreateCommand{}

			flags := cmd.Meta.FlagSet("gcpoidc create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.serviceAccountEmail, "service-account-email", "", "GCP service account email (required)")
			flags.StringVar(&cmd.workloadProviderName, "workload-provider-name", "", "Workload provider path (required)")
			flags.StringVar(&cmd.projectNumber, "project-number", "", "GCP project number (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the service account email was set correctly
			if cmd.serviceAccountEmail != tt.expectedServiceAccount {
				t.Errorf("expected serviceAccountEmail %q, got %q", tt.expectedServiceAccount, cmd.serviceAccountEmail)
			}

			// Verify the workload provider name was set correctly
			if cmd.workloadProviderName != tt.expectedWorkloadProvider {
				t.Errorf("expected workloadProviderName %q, got %q", tt.expectedWorkloadProvider, cmd.workloadProviderName)
			}

			// Verify the project number was set correctly
			if cmd.projectNumber != tt.expectedProjectNumber {
				t.Errorf("expected projectNumber %q, got %q", tt.expectedProjectNumber, cmd.projectNumber)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}

func TestGCPOIDCCreatePartialRequiredFlagsRun(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError string
	}{
		{
			name:        "missing workload-provider-name and project-number",
			args:        []string{"-organization=my-org", "-service-account-email=terraform@my-project.iam.gserviceaccount.com"},
			expectError: "-workload-provider-name",
		},
		{
			name: "missing project-number only",
			args: []string{
				"-organization=my-org",
				"-service-account-email=terraform@my-project.iam.gserviceaccount.com",
				"-workload-provider-name=projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
			},
			expectError: "-project-number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := cli.NewMockUi()
			cmd := &GCPoidcCreateCommand{
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
