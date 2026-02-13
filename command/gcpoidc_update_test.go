package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestGCPOIDCUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-service-account-email=test@project.iam.gserviceaccount.com"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestGCPOIDCUpdateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcUpdateCommand{
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

func TestGCPOIDCUpdateValidatesEmptyID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=", "-service-account-email=test@project.iam.gserviceaccount.com"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestGCPOIDCUpdateHelp(t *testing.T) {
	cmd := &GCPoidcUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf gcpoidc update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
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
		t.Error("Help should indicate -id flag is required")
	}
}

func TestGCPOIDCUpdateSynopsis(t *testing.T) {
	cmd := &GCPoidcUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update GCP OIDC configuration settings" {
		t.Errorf("expected 'Update GCP OIDC configuration settings', got %q", synopsis)
	}
}

func TestGCPOIDCUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name                     string
		args                     []string
		expectedID               string
		expectedServiceAccount   string
		expectedWorkloadProvider string
		expectedProjectNumber    string
		expectedFormat           string
	}{
		{
			name:                     "id only with default format",
			args:                     []string{"-id=gcpoidc-ABC123"},
			expectedID:               "gcpoidc-ABC123",
			expectedServiceAccount:   "",
			expectedWorkloadProvider: "",
			expectedProjectNumber:    "",
			expectedFormat:           "table",
		},
		{
			name:                     "id with service account email",
			args:                     []string{"-id=gcpoidc-ABC123", "-service-account-email=terraform@my-project.iam.gserviceaccount.com"},
			expectedID:               "gcpoidc-ABC123",
			expectedServiceAccount:   "terraform@my-project.iam.gserviceaccount.com",
			expectedWorkloadProvider: "",
			expectedProjectNumber:    "",
			expectedFormat:           "table",
		},
		{
			name:                     "id with workload provider name",
			args:                     []string{"-id=gcpoidc-XYZ789", "-workload-provider-name=projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider"},
			expectedID:               "gcpoidc-XYZ789",
			expectedServiceAccount:   "",
			expectedWorkloadProvider: "projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
			expectedProjectNumber:    "",
			expectedFormat:           "table",
		},
		{
			name:                     "id with project number",
			args:                     []string{"-id=gcpoidc-DEF456", "-project-number=987654321"},
			expectedID:               "gcpoidc-DEF456",
			expectedServiceAccount:   "",
			expectedWorkloadProvider: "",
			expectedProjectNumber:    "987654321",
			expectedFormat:           "table",
		},
		{
			name:                     "all flags with json format",
			args:                     []string{"-id=gcpoidc-GHI789", "-service-account-email=sa@test.iam.gserviceaccount.com", "-workload-provider-name=projects/111/locations/global/workloadIdentityPools/pool/providers/provider", "-project-number=111", "-output=json"},
			expectedID:               "gcpoidc-GHI789",
			expectedServiceAccount:   "sa@test.iam.gserviceaccount.com",
			expectedWorkloadProvider: "projects/111/locations/global/workloadIdentityPools/pool/providers/provider",
			expectedProjectNumber:    "111",
			expectedFormat:           "json",
		},
		{
			name:                     "all flags with table format",
			args:                     []string{"-id=gcpoidc-JKL012", "-service-account-email=prod@prod.iam.gserviceaccount.com", "-workload-provider-name=projects/222/locations/global/workloadIdentityPools/prod-pool/providers/prod-provider", "-project-number=222", "-output=table"},
			expectedID:               "gcpoidc-JKL012",
			expectedServiceAccount:   "prod@prod.iam.gserviceaccount.com",
			expectedWorkloadProvider: "projects/222/locations/global/workloadIdentityPools/prod-pool/providers/prod-provider",
			expectedProjectNumber:    "222",
			expectedFormat:           "table",
		},
		{
			name:                     "partial update with service account and project number",
			args:                     []string{"-id=gcpoidc-MNO345", "-service-account-email=staging@staging.iam.gserviceaccount.com", "-project-number=333"},
			expectedID:               "gcpoidc-MNO345",
			expectedServiceAccount:   "staging@staging.iam.gserviceaccount.com",
			expectedWorkloadProvider: "",
			expectedProjectNumber:    "333",
			expectedFormat:           "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &GCPoidcUpdateCommand{}

			flags := cmd.Meta.FlagSet("gcpoidc update")
			flags.StringVar(&cmd.id, "id", "", "GCP OIDC configuration ID (required)")
			flags.StringVar(&cmd.serviceAccountEmail, "service-account-email", "", "GCP service account email")
			flags.StringVar(&cmd.workloadProviderName, "workload-provider-name", "", "Workload provider path")
			flags.StringVar(&cmd.projectNumber, "project-number", "", "GCP project number")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
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
