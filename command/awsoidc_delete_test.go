package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAWSOIDCDeleteRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AWSoidcDeleteCommand{
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

func TestAWSOIDCDeleteHelp(t *testing.T) {
	cmd := &AWSoidcDeleteCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf awsoidc delete") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-force") {
		t.Error("Help should mention -force flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "AWS OIDC configuration") {
		t.Error("Help should describe AWS OIDC configuration")
	}
	if !strings.Contains(help, "WARNING") {
		t.Error("Help should contain a warning")
	}
}

func TestAWSOIDCDeleteSynopsis(t *testing.T) {
	cmd := &AWSoidcDeleteCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Delete an AWS OIDC configuration" {
		t.Errorf("expected 'Delete an AWS OIDC configuration', got %q", synopsis)
	}
}

func TestAWSOIDCDeleteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedID    string
		expectedForce bool
	}{
		{
			name:          "id flag only",
			args:          []string{"-id=awsoidc-ABC123"},
			expectedID:    "awsoidc-ABC123",
			expectedForce: false,
		},
		{
			name:          "id with force flag",
			args:          []string{"-id=awsoidc-XYZ789", "-force"},
			expectedID:    "awsoidc-XYZ789",
			expectedForce: true,
		},
		{
			name:          "force flag with id",
			args:          []string{"-force", "-id=awsoidc-DEF456"},
			expectedID:    "awsoidc-DEF456",
			expectedForce: true,
		},
		{
			name:          "id flag without force",
			args:          []string{"-id=awsoidc-GHI789"},
			expectedID:    "awsoidc-GHI789",
			expectedForce: false,
		},
		{
			name:          "different id format with force",
			args:          []string{"-id=awsoidc-12345abcde", "-force"},
			expectedID:    "awsoidc-12345abcde",
			expectedForce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AWSoidcDeleteCommand{}

			flags := cmd.Meta.FlagSet("awsoidc delete")
			flags.StringVar(&cmd.id, "id", "", "AWS OIDC configuration ID (required)")
			flags.BoolVar(&cmd.force, "force", false, "Force delete without confirmation")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the ID was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected ID %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the force flag was set correctly
			if cmd.force != tt.expectedForce {
				t.Errorf("expected force %v, got %v", tt.expectedForce, cmd.force)
			}
		})
	}
}
