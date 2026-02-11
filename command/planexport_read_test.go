package command

import (
	"context"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

type mockPlanExportReader struct {
	readFunc func(ctx context.Context, planExportID string) (*tfe.PlanExport, error)
}

func (m *mockPlanExportReader) Read(ctx context.Context, planExportID string) (*tfe.PlanExport, error) {
	return m.readFunc(ctx, planExportID)
}

func TestPlanExportReadCommand_Run(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		mock      func() *mockPlanExportReader
		wantCode  int
		wantError string
	}{
		{
			name:      "missing id flag",
			args:      []string{},
			wantCode:  1,
			wantError: "Error: -id flag is required",
		},
		{
			name: "successful read",
			args: []string{"-id=pe-abc123"},
			mock: func() *mockPlanExportReader {
				return &mockPlanExportReader{
					readFunc: func(ctx context.Context, planExportID string) (*tfe.PlanExport, error) {
						return &tfe.PlanExport{
							ID:       "pe-abc123",
							DataType: tfe.PlanExportSentinelMockBundleV0,
							Status:   "finished",
							StatusTimestamps: &tfe.PlanExportStatusTimestamps{
								QueuedAt:   time.Now().Add(-5 * time.Minute),
								FinishedAt: time.Now(),
							},
						}, nil
					},
				}
			},
			wantCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := cli.NewMockUi()
			cmd := &PlanExportReadCommand{
				Meta: testMeta(t, ui),
			}

			if tt.mock != nil {
				cmd.planExportSvc = tt.mock()
			}

			code := cmd.Run(tt.args)
			if code != tt.wantCode {
				t.Errorf("Run() = %v, want %v", code, tt.wantCode)
			}

			if tt.wantError != "" {
				if output := ui.ErrorWriter.String(); output == "" {
					t.Errorf("Expected error containing %q, got no error", tt.wantError)
				}
			}
		})
	}
}

func TestPlanExportReadCommand_Help(t *testing.T) {
	cmd := &PlanExportReadCommand{}
	help := cmd.Help()

	if help == "" {
		t.Error("Help() returned empty string")
	}

	expectedStrings := []string{
		"Usage:",
		"planexport read",
		"-id",
		"-output",
	}

	for _, expected := range expectedStrings {
		if !containsString(help, expected) {
			t.Errorf("Help() missing expected string: %q", expected)
		}
	}
}

func TestPlanExportReadCommand_Synopsis(t *testing.T) {
	cmd := &PlanExportReadCommand{}
	synopsis := cmd.Synopsis()

	if synopsis == "" {
		t.Error("Synopsis() returned empty string")
	}

	if len(synopsis) > 80 {
		t.Errorf("Synopsis() too long: %d characters (max 80)", len(synopsis))
	}
}
