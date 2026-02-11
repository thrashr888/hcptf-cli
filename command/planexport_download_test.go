package command

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/cli"
)

type mockPlanExportDownloader struct {
	downloadFunc func(ctx context.Context, planExportID string) ([]byte, error)
}

func (m *mockPlanExportDownloader) Download(ctx context.Context, planExportID string) ([]byte, error) {
	return m.downloadFunc(ctx, planExportID)
}

func TestPlanExportDownloadCommand_Run(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		args      []string
		mock      func() *mockPlanExportDownloader
		wantCode  int
		wantError string
		setup     func() string
	}{
		{
			name:      "missing id flag",
			args:      []string{},
			wantCode:  1,
			wantError: "Error: -id flag is required",
		},
		{
			name: "successful download with default path",
			args: []string{"-id=pe-abc123"},
			mock: func() *mockPlanExportDownloader {
				return &mockPlanExportDownloader{
					downloadFunc: func(ctx context.Context, planExportID string) ([]byte, error) {
						return []byte("mock export data"), nil
					},
				}
			},
			setup: func() string {
				// Change to temp directory so default file is created there
				oldDir, _ := os.Getwd()
				os.Chdir(tmpDir)
				t.Cleanup(func() { os.Chdir(oldDir) })
				return ""
			},
			wantCode: 0,
		},
		{
			name: "successful download with custom path",
			args: func() []string {
				return []string{"-id=pe-abc123", "-path=" + filepath.Join(tmpDir, "custom-export.tar.gz")}
			}(),
			mock: func() *mockPlanExportDownloader {
				return &mockPlanExportDownloader{
					downloadFunc: func(ctx context.Context, planExportID string) ([]byte, error) {
						return []byte("mock export data"), nil
					},
				}
			},
			wantCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			ui := cli.NewMockUi()
			cmd := &PlanExportDownloadCommand{
				Meta: testMeta(t, ui),
			}

			if tt.mock != nil {
				cmd.planExportSvc = tt.mock()
			}

			code := cmd.Run(tt.args)
			if code != tt.wantCode {
				t.Errorf("Run() = %v, want %v", code, tt.wantCode)
				t.Logf("Output: %s", ui.OutputWriter.String())
				t.Logf("Error: %s", ui.ErrorWriter.String())
			}

			if tt.wantError != "" {
				if output := ui.ErrorWriter.String(); output == "" {
					t.Errorf("Expected error containing %q, got no error", tt.wantError)
				}
			}
		})
	}
}

func TestPlanExportDownloadCommand_Help(t *testing.T) {
	cmd := &PlanExportDownloadCommand{}
	help := cmd.Help()

	if help == "" {
		t.Error("Help() returned empty string")
	}

	expectedStrings := []string{
		"Usage:",
		"planexport download",
		"-id",
		"-path",
	}

	for _, expected := range expectedStrings {
		if !containsString(help, expected) {
			t.Errorf("Help() missing expected string: %q", expected)
		}
	}
}

func TestPlanExportDownloadCommand_Synopsis(t *testing.T) {
	cmd := &PlanExportDownloadCommand{}
	synopsis := cmd.Synopsis()

	if synopsis == "" {
		t.Error("Synopsis() returned empty string")
	}

	if len(synopsis) > 80 {
		t.Errorf("Synopsis() too long: %d characters (max 80)", len(synopsis))
	}
}
