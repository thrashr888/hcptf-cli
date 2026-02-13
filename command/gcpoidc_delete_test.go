package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestGCPOIDCDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcDeleteCommand{
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

func TestGCPOIDCDeleteRequiresAllFlags(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &GCPoidcDeleteCommand{
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

func TestGCPOIDCDeleteHelp(t *testing.T) {
	cmd := &GCPoidcDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf gcpoidc delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id flag is required")
	}
	if !strings.Contains(help, "WARNING") || !strings.Contains(help, "Delete") {
		t.Error("Help should contain warning about deletion")
	}
}

func TestGCPOIDCDeleteSynopsis(t *testing.T) {
	cmd := &GCPoidcDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete a GCP OIDC configuration" {
		t.Errorf("expected 'Delete a GCP OIDC configuration', got %q", synopsis)
	}
}

func TestGCPOIDCDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id only without force",
			args:          []string{"-id=gcpoidc-ABC123"},
			expectedID:    "gcpoidc-ABC123",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=gcpoidc-XYZ789", "-force"},
			expectedID:    "gcpoidc-XYZ789",
			expectedForce: true,
		},
		{
			name:          "long id without force",
			args:          []string{"-id=gcpoidc-1234567890ABCDEF"},
			expectedID:    "gcpoidc-1234567890ABCDEF",
			expectedForce: false,
		},
		{
			name:          "long id with force",
			args:          []string{"-id=gcpoidc-9876543210FEDCBA", "-force"},
			expectedID:    "gcpoidc-9876543210FEDCBA",
			expectedForce: true,
		},
		{
			name:          "id with underscores and force",
			args:          []string{"-id=gcpoidc-test_config_456", "-force"},
			expectedID:    "gcpoidc-test_config_456",
			expectedForce: true,
		},
		{
			name:          "force true explicit",
			args:          []string{"-id=gcpoidc-DEF456", "-force=true"},
			expectedID:    "gcpoidc-DEF456",
			expectedForce: true,
		},
		{
			name:          "force false explicit",
			args:          []string{"-id=gcpoidc-GHI789", "-force=false"},
			expectedID:    "gcpoidc-GHI789",
			expectedForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &GCPoidcDeleteCommand{}

			flags := cmd.Meta.FlagSet("gcpoidc delete")
			flags.StringVar(&cmd.id, "id", "", "GCP OIDC configuration ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
