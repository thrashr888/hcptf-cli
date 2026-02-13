package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestHYOKCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-name=test-hyok", "-kek-id=test-key", "-agent-pool-id=apool-123", "-oidc-config-id=oidc-456", "-oidc-type=aws"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestHYOKCreateRequiresName(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-kek-id=test-key", "-agent-pool-id=apool-123", "-oidc-config-id=oidc-456", "-oidc-type=aws"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-name") {
		t.Fatalf("expected name error, got %q", out)
	}
}

func TestHYOKCreateRequiresKEKID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-name=test-hyok", "-agent-pool-id=apool-123", "-oidc-config-id=oidc-456", "-oidc-type=aws"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-kek-id") {
		t.Fatalf("expected kek-id error, got %q", out)
	}
}

func TestHYOKCreateRequiresAgentPoolID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-name=test-hyok", "-kek-id=test-key", "-oidc-config-id=oidc-456", "-oidc-type=aws"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-agent-pool-id") {
		t.Fatalf("expected agent-pool-id error, got %q", out)
	}
}

func TestHYOKCreateRequiresOIDCConfigID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-name=test-hyok", "-kek-id=test-key", "-agent-pool-id=apool-123", "-oidc-type=aws"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-oidc-config-id") {
		t.Fatalf("expected oidc-config-id error, got %q", out)
	}
}

func TestHYOKCreateRequiresOIDCType(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=test-org", "-name=test-hyok", "-kek-id=test-key", "-agent-pool-id=apool-123", "-oidc-config-id=oidc-456"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-oidc-type") {
		t.Fatalf("expected oidc-type error, got %q", out)
	}
}

func TestHYOKCreateRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKCreateCommand{
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

func TestHYOKCreateHelp(t *testing.T) {
	cmd := &HYOKCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf hyok create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "-kek-id") {
		t.Error("Help should mention -kek-id flag")
	}
	if !strings.Contains(help, "-agent-pool-id") {
		t.Error("Help should mention -agent-pool-id flag")
	}
	if !strings.Contains(help, "-oidc-config-id") {
		t.Error("Help should mention -oidc-config-id flag")
	}
	if !strings.Contains(help, "-oidc-type") {
		t.Error("Help should mention -oidc-type flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate flags are required")
	}
	if !strings.Contains(help, "HYOK") {
		t.Error("Help should mention HYOK")
	}
}

func TestHYOKCreateSynopsis(t *testing.T) {
	cmd := &HYOKCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create a HYOK configuration" {
		t.Errorf("expected 'Create a HYOK configuration', got %q", synopsis)
	}
}

func TestHYOKCreateValidatesOIDCType(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{
		"-organization=test-org",
		"-name=test-hyok",
		"-kek-id=kek-123",
		"-agent-pool-id=apool-123",
		"-oidc-config-id=oidc-123",
		"-oidc-type=invalid",
	})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "oidc-type") {
		t.Fatalf("expected oidc-type validation error, got %q", out)
	}
}

func TestHYOKCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedOrg     string
		expectedName    string
		expectedKEKID   string
		expectedAgentID string
		expectedOIDCID  string
		expectedOIDC    string
		expectedRegion  string
		expectedLoc     string
		expectedRing    string
		expectedFormat  string
	}{
		{
			name:            "required flags only, default values",
			args:            []string{"-organization=my-org", "-name=test-hyok", "-kek-id=kek-123", "-agent-pool-id=apool-456", "-oidc-config-id=oidc-789", "-oidc-type=aws"},
			expectedOrg:     "my-org",
			expectedName:    "test-hyok",
			expectedKEKID:   "kek-123",
			expectedAgentID: "apool-456",
			expectedOIDCID:  "oidc-789",
			expectedOIDC:    "aws",
			expectedRegion:  "",
			expectedLoc:     "",
			expectedRing:    "",
			expectedFormat:  "table",
		},
		{
			name:            "aws with key region",
			args:            []string{"-organization=prod-org", "-name=aws-hyok", "-kek-id=arn:aws:kms:us-west-2:123456789012:key/uuid", "-agent-pool-id=apool-111", "-oidc-config-id=oidc-222", "-oidc-type=aws", "-key-region=us-west-2"},
			expectedOrg:     "prod-org",
			expectedName:    "aws-hyok",
			expectedKEKID:   "arn:aws:kms:us-west-2:123456789012:key/uuid",
			expectedAgentID: "apool-111",
			expectedOIDCID:  "oidc-222",
			expectedOIDC:    "aws",
			expectedRegion:  "us-west-2",
			expectedLoc:     "",
			expectedRing:    "",
			expectedFormat:  "table",
		},
		{
			name:            "gcp with key location and ring",
			args:            []string{"-organization=test-org", "-name=gcp-hyok", "-kek-id=my-key", "-agent-pool-id=apool-333", "-oidc-config-id=oidc-444", "-oidc-type=gcp", "-key-location=us-central1", "-key-ring-id=my-keyring"},
			expectedOrg:     "test-org",
			expectedName:    "gcp-hyok",
			expectedKEKID:   "my-key",
			expectedAgentID: "apool-333",
			expectedOIDCID:  "oidc-444",
			expectedOIDC:    "gcp",
			expectedRegion:  "",
			expectedLoc:     "us-central1",
			expectedRing:    "my-keyring",
			expectedFormat:  "table",
		},
		{
			name:            "azure with json output",
			args:            []string{"-organization=dev-org", "-name=azure-hyok", "-kek-id=azure-key", "-agent-pool-id=apool-555", "-oidc-config-id=oidc-666", "-oidc-type=azure", "-output=json"},
			expectedOrg:     "dev-org",
			expectedName:    "azure-hyok",
			expectedKEKID:   "azure-key",
			expectedAgentID: "apool-555",
			expectedOIDCID:  "oidc-666",
			expectedOIDC:    "azure",
			expectedRegion:  "",
			expectedLoc:     "",
			expectedRing:    "",
			expectedFormat:  "json",
		},
		{
			name:            "vault type",
			args:            []string{"-organization=vault-org", "-name=vault-hyok", "-kek-id=vault-key", "-agent-pool-id=apool-777", "-oidc-config-id=oidc-888", "-oidc-type=vault"},
			expectedOrg:     "vault-org",
			expectedName:    "vault-hyok",
			expectedKEKID:   "vault-key",
			expectedAgentID: "apool-777",
			expectedOIDCID:  "oidc-888",
			expectedOIDC:    "vault",
			expectedRegion:  "",
			expectedLoc:     "",
			expectedRing:    "",
			expectedFormat:  "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &HYOKCreateCommand{}

			flags := cmd.Meta.FlagSet("hyok create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.name, "name", "", "HYOK configuration name (required)")
			flags.StringVar(&cmd.kekID, "kek-id", "", "Key Encryption Key ID from your KMS (required)")
			flags.StringVar(&cmd.agentPoolID, "agent-pool-id", "", "Agent pool ID (required)")
			flags.StringVar(&cmd.oidcConfigID, "oidc-config-id", "", "OIDC configuration ID (required)")
			flags.StringVar(&cmd.oidcType, "oidc-type", "", "OIDC type: aws, azure, gcp, or vault (required)")
			flags.StringVar(&cmd.keyRegion, "key-region", "", "AWS KMS key region (for AWS KMS only)")
			flags.StringVar(&cmd.keyLocation, "key-location", "", "GCP key location (for GCP Cloud KMS only)")
			flags.StringVar(&cmd.keyRingID, "key-ring-id", "", "GCP key ring ID (for GCP Cloud KMS only)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the kek-id was set correctly
			if cmd.kekID != tt.expectedKEKID {
				t.Errorf("expected kekID %q, got %q", tt.expectedKEKID, cmd.kekID)
			}

			// Verify the agent-pool-id was set correctly
			if cmd.agentPoolID != tt.expectedAgentID {
				t.Errorf("expected agentPoolID %q, got %q", tt.expectedAgentID, cmd.agentPoolID)
			}

			// Verify the oidc-config-id was set correctly
			if cmd.oidcConfigID != tt.expectedOIDCID {
				t.Errorf("expected oidcConfigID %q, got %q", tt.expectedOIDCID, cmd.oidcConfigID)
			}

			// Verify the oidc-type was set correctly
			if cmd.oidcType != tt.expectedOIDC {
				t.Errorf("expected oidcType %q, got %q", tt.expectedOIDC, cmd.oidcType)
			}

			// Verify the key-region was set correctly
			if cmd.keyRegion != tt.expectedRegion {
				t.Errorf("expected keyRegion %q, got %q", tt.expectedRegion, cmd.keyRegion)
			}

			// Verify the key-location was set correctly
			if cmd.keyLocation != tt.expectedLoc {
				t.Errorf("expected keyLocation %q, got %q", tt.expectedLoc, cmd.keyLocation)
			}

			// Verify the key-ring-id was set correctly
			if cmd.keyRingID != tt.expectedRing {
				t.Errorf("expected keyRingID %q, got %q", tt.expectedRing, cmd.keyRingID)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFormat {
				t.Errorf("expected format %q, got %q", tt.expectedFormat, cmd.format)
			}
		})
	}
}
