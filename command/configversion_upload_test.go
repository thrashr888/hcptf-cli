package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestConfigVersionUploadRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ConfigVersionUploadCommand{
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

func TestConfigVersionUploadRequiresPath(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ConfigVersionUploadCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=cv-123"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-path") {
		t.Fatalf("expected path error, got %q", out)
	}
}

func TestConfigVersionUploadRequiresEmptyID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ConfigVersionUploadCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=", "-path=/tmp"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestConfigVersionUploadRequiresEmptyPath(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ConfigVersionUploadCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=cv-123", "-path="})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-path") {
		t.Fatalf("expected path error, got %q", out)
	}
}

func TestConfigVersionUploadRequiresPathExists(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &ConfigVersionUploadCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=cv-123", "-path=/nonexistent/path/that/does/not/exist"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "does not exist") {
		t.Fatalf("expected path existence error, got %q", out)
	}
}

func TestConfigVersionUploadHelp(t *testing.T) {
	cmd := &ConfigVersionUploadCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf configversion upload") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !strings.Contains(help, "Configuration version ID") {
		t.Error("Help should mention configuration version ID")
	}
	if !strings.Contains(help, "-path") {
		t.Error("Help should mention -path flag")
	}
	if !strings.Contains(help, "configuration directory or tar.gz file") {
		t.Error("Help should mention path types")
	}
	if !strings.Contains(help, "automatically archived") {
		t.Error("Help should mention automatic archiving")
	}
}

func TestConfigVersionUploadSynopsis(t *testing.T) {
	cmd := &ConfigVersionUploadCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Upload configuration files" {
		t.Errorf("expected 'Upload configuration files', got %q", synopsis)
	}
}

func TestConfigVersionUploadFlagParsing(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		expectedID         string
		expectedPath       string
	}{
		{
			name:               "directory path",
			args:               []string{"-id=cv-abc123", "-path=./terraform"},
			expectedID:         "cv-abc123",
			expectedPath:       "./terraform",
		},
		{
			name:               "tar.gz file path",
			args:               []string{"-id=cv-xyz789", "-path=./config.tar.gz"},
			expectedID:         "cv-xyz789",
			expectedPath:       "./config.tar.gz",
		},
		{
			name:               "absolute path",
			args:               []string{"-id=cv-123456", "-path=/home/user/configs"},
			expectedID:         "cv-123456",
			expectedPath:       "/home/user/configs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ConfigVersionUploadCommand{}

			flags := cmd.Meta.FlagSet("configversion upload")
			flags.StringVar(&cmd.configVersionID, "id", "", "Configuration version ID (required)")
			flags.StringVar(&cmd.path, "path", "", "Path to configuration directory or tar.gz file (required)")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the configVersionID was set correctly
			if cmd.configVersionID != tt.expectedID {
				t.Errorf("expected configVersionID %q, got %q", tt.expectedID, cmd.configVersionID)
			}

			// Verify the path was set correctly
			if cmd.path != tt.expectedPath {
				t.Errorf("expected path %q, got %q", tt.expectedPath, cmd.path)
			}
		})
	}
}
