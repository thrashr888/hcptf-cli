package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestHYOKUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKUpdateCommand{
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

func TestHYOKUpdateHelp(t *testing.T) {
	cmd := &HYOKUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf hyok update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "-name") {
		t.Error("Help should mention -name flag")
	}
	if !strings.Contains(help, "-kek-id") {
		t.Error("Help should mention -kek-id flag")
	}
	if !strings.Contains(help, "-primary") {
		t.Error("Help should mention -primary flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate -id is required")
	}
	if !strings.Contains(help, "HYOK") {
		t.Error("Help should mention HYOK")
	}
}

func TestHYOKUpdateSynopsis(t *testing.T) {
	cmd := &HYOKUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update a HYOK configuration" {
		t.Errorf("expected 'Update a HYOK configuration', got %q", synopsis)
	}
}

func TestHYOKUpdateValidatesPrimary(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &HYOKUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=hyok-123", "-primary=invalid"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "primary") {
		t.Fatalf("expected primary validation error, got %q", out)
	}
}

func TestHYOKUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedID     string
		expectedName   string
		expectedKEKID  string
		expectedPrim   string
		expectedRegion string
		expectedLoc    string
		expectedRing   string
		expectedFormat string
	}{
		{
			name:           "id only, default format",
			args:           []string{"-id=hyokc-123456"},
			expectedID:     "hyokc-123456",
			expectedName:   "",
			expectedKEKID:  "",
			expectedPrim:   "",
			expectedRegion: "",
			expectedLoc:    "",
			expectedRing:   "",
			expectedFormat: "table",
		},
		{
			name:           "id with name",
			args:           []string{"-id=hyokc-abcdef", "-name=updated-hyok"},
			expectedID:     "hyokc-abcdef",
			expectedName:   "updated-hyok",
			expectedKEKID:  "",
			expectedPrim:   "",
			expectedRegion: "",
			expectedLoc:    "",
			expectedRing:   "",
			expectedFormat: "table",
		},
		{
			name:           "id with kek-id",
			args:           []string{"-id=hyokc-xyz789", "-kek-id=new-key-123"},
			expectedID:     "hyokc-xyz789",
			expectedName:   "",
			expectedKEKID:  "new-key-123",
			expectedPrim:   "",
			expectedRegion: "",
			expectedLoc:    "",
			expectedRing:   "",
			expectedFormat: "table",
		},
		{
			name:           "id with primary true",
			args:           []string{"-id=hyokc-111222", "-primary=true"},
			expectedID:     "hyokc-111222",
			expectedName:   "",
			expectedKEKID:  "",
			expectedPrim:   "true",
			expectedRegion: "",
			expectedLoc:    "",
			expectedRing:   "",
			expectedFormat: "table",
		},
		{
			name:           "id with primary false",
			args:           []string{"-id=hyokc-333444", "-primary=false"},
			expectedID:     "hyokc-333444",
			expectedName:   "",
			expectedKEKID:  "",
			expectedPrim:   "false",
			expectedRegion: "",
			expectedLoc:    "",
			expectedRing:   "",
			expectedFormat: "table",
		},
		{
			name:           "id with aws key region",
			args:           []string{"-id=hyokc-555666", "-key-region=us-east-1"},
			expectedID:     "hyokc-555666",
			expectedName:   "",
			expectedKEKID:  "",
			expectedPrim:   "",
			expectedRegion: "us-east-1",
			expectedLoc:    "",
			expectedRing:   "",
			expectedFormat: "table",
		},
		{
			name:           "id with gcp key location and ring",
			args:           []string{"-id=hyokc-777888", "-key-location=us-west1", "-key-ring-id=my-ring"},
			expectedID:     "hyokc-777888",
			expectedName:   "",
			expectedKEKID:  "",
			expectedPrim:   "",
			expectedRegion: "",
			expectedLoc:    "us-west1",
			expectedRing:   "my-ring",
			expectedFormat: "table",
		},
		{
			name:           "id with all optional fields and json output",
			args:           []string{"-id=hyokc-999000", "-name=full-update", "-kek-id=full-key", "-primary=true", "-key-region=eu-west-1", "-output=json"},
			expectedID:     "hyokc-999000",
			expectedName:   "full-update",
			expectedKEKID:  "full-key",
			expectedPrim:   "true",
			expectedRegion: "eu-west-1",
			expectedLoc:    "",
			expectedRing:   "",
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &HYOKUpdateCommand{}

			flags := cmd.Meta.FlagSet("hyok update")
			flags.StringVar(&cmd.id, "id", "", "HYOK configuration ID (required)")
			flags.StringVar(&cmd.name, "name", "", "HYOK configuration name")
			flags.StringVar(&cmd.kekID, "kek-id", "", "Key Encryption Key ID from your KMS")
			flags.StringVar(&cmd.primary, "primary", "", "Set as primary HYOK configuration (true/false)")
			flags.StringVar(&cmd.keyRegion, "key-region", "", "AWS KMS key region (for AWS KMS only)")
			flags.StringVar(&cmd.keyLocation, "key-location", "", "GCP key location (for GCP Cloud KMS only)")
			flags.StringVar(&cmd.keyRingID, "key-ring-id", "", "GCP key ring ID (for GCP Cloud KMS only)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the id was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected id %q, got %q", tt.expectedID, cmd.id)
			}

			// Verify the name was set correctly
			if cmd.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, cmd.name)
			}

			// Verify the kek-id was set correctly
			if cmd.kekID != tt.expectedKEKID {
				t.Errorf("expected kekID %q, got %q", tt.expectedKEKID, cmd.kekID)
			}

			// Verify the primary was set correctly
			if cmd.primary != tt.expectedPrim {
				t.Errorf("expected primary %q, got %q", tt.expectedPrim, cmd.primary)
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
