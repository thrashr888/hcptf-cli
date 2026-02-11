package command

import (
	"context"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

type mockPlanExportCreator struct {
	createFunc func(ctx context.Context, options tfe.PlanExportCreateOptions) (*tfe.PlanExport, error)
}

func (m *mockPlanExportCreator) Create(ctx context.Context, options tfe.PlanExportCreateOptions) (*tfe.PlanExport, error) {
	return m.createFunc(ctx, options)
}

func TestPlanExportCreateCommand_Run(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		mock      func() *mockPlanExportCreator
		wantCode  int
		wantError string
	}{
		{
			name:      "missing plan-id flag",
			args:      []string{},
			wantCode:  1,
			wantError: "Error: -plan-id flag is required",
		},
		{
			name: "successful create",
			args: []string{"-plan-id=plan-abc123"},
			mock: func() *mockPlanExportCreator {
				return &mockPlanExportCreator{
					createFunc: func(ctx context.Context, options tfe.PlanExportCreateOptions) (*tfe.PlanExport, error) {
						return &tfe.PlanExport{
							ID:               "pe-test123",
							DataType:         tfe.PlanExportSentinelMockBundleV0,
							Status:           "queued",
							StatusTimestamps: &tfe.PlanExportStatusTimestamps{},
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
			cmd := &PlanExportCreateCommand{
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

func TestPlanExportCreateCommand_Help(t *testing.T) {
	cmd := &PlanExportCreateCommand{}
	help := cmd.Help()

	if help == "" {
		t.Error("Help() returned empty string")
	}

	// Check for key elements in help text
	expectedStrings := []string{
		"Usage:",
		"planexport create",
		"-plan-id",
		"-data-type",
		"-output",
	}

	for _, expected := range expectedStrings {
		if !containsString(help, expected) {
			t.Errorf("Help() missing expected string: %q", expected)
		}
	}
}

func TestPlanExportCreateCommand_Synopsis(t *testing.T) {
	cmd := &PlanExportCreateCommand{}
	synopsis := cmd.Synopsis()

	if synopsis == "" {
		t.Error("Synopsis() returned empty string")
	}

	if len(synopsis) > 80 {
		t.Errorf("Synopsis() too long: %d characters (max 80)", len(synopsis))
	}
}

func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
