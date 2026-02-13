package command

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/hcptf-cli/internal/client"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/mitchellh/cli"
)

func TestQueryRunListRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &QueryRunListCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestQueryRunListHelp(t *testing.T) {
	cmd := &QueryRunListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf queryrun list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
	}
	if !strings.Contains(help, "-status") {
		t.Error("Help should mention -status flag")
	}
	if !strings.Contains(help, "-operation") {
		t.Error("Help should mention -operation flag")
	}
	if !strings.Contains(help, "-source") {
		t.Error("Help should mention -source flag")
	}
	if !strings.Contains(help, "-workspace") {
		t.Error("Help should mention -workspace flag")
	}
	if !strings.Contains(help, "-agent-pool") {
		t.Error("Help should mention -agent-pool flag")
	}
	if !strings.Contains(help, "-status-group") {
		t.Error("Help should mention -status-group flag")
	}
	if !strings.Contains(help, "-search-user") {
		t.Error("Help should mention -search-user flag")
	}
	if !strings.Contains(help, "-search-commit") {
		t.Error("Help should mention -search-commit flag")
	}
	if !strings.Contains(help, "-search-basic") {
		t.Error("Help should mention -search-basic flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
}

func TestQueryRunListSynopsis(t *testing.T) {
	cmd := &QueryRunListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Search runs across organization" {
		t.Errorf("expected 'Search runs across organization', got %q", synopsis)
	}
}

func TestQueryRunListFlagParsing(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedOrg     string
		expectedStatus  string
		expectedOp      string
		expectedSource  string
		expectedWS      string
		expectedAgent   string
		expectedSGroup  string
		expectedSUser   string
		expectedSCommit string
		expectedSBasic  string
		expectedFmt     string
	}{
		{
			name:        "organization flag only",
			args:        []string{"-organization=my-org"},
			expectedOrg: "my-org",
			expectedFmt: "table",
		},
		{
			name:        "org alias flag",
			args:        []string{"-org=test-org"},
			expectedOrg: "test-org",
			expectedFmt: "table",
		},
		{
			name:           "with status filter",
			args:           []string{"-org=my-org", "-status=applied,applying"},
			expectedOrg:    "my-org",
			expectedStatus: "applied,applying",
			expectedFmt:    "table",
		},
		{
			name:        "with operation filter",
			args:        []string{"-org=my-org", "-operation=plan_and_apply"},
			expectedOrg: "my-org",
			expectedOp:  "plan_and_apply",
			expectedFmt: "table",
		},
		{
			name:           "with source filter",
			args:           []string{"-org=my-org", "-source=tfe-api"},
			expectedOrg:    "my-org",
			expectedSource: "tfe-api",
			expectedFmt:    "table",
		},
		{
			name:        "with workspace filter",
			args:        []string{"-org=my-org", "-workspace=prod-app"},
			expectedOrg: "my-org",
			expectedWS:  "prod-app",
			expectedFmt: "table",
		},
		{
			name:          "with agent pool filter",
			args:          []string{"-org=my-org", "-agent-pool=prod-agents"},
			expectedOrg:   "my-org",
			expectedAgent: "prod-agents",
			expectedFmt:   "table",
		},
		{
			name:          "with status group filter",
			args:          []string{"-org=my-org", "-status-group=final"},
			expectedOrg:   "my-org",
			expectedSGroup: "final",
			expectedFmt:   "table",
		},
		{
			name:          "with search user",
			args:          []string{"-org=my-org", "-search-user=alice"},
			expectedOrg:   "my-org",
			expectedSUser: "alice",
			expectedFmt:   "table",
		},
		{
			name:            "with search commit",
			args:            []string{"-org=my-org", "-search-commit=abc123"},
			expectedOrg:     "my-org",
			expectedSCommit: "abc123",
			expectedFmt:     "table",
		},
		{
			name:           "with search basic",
			args:           []string{"-org=my-org", "-search-basic=test"},
			expectedOrg:    "my-org",
			expectedSBasic: "test",
			expectedFmt:    "table",
		},
		{
			name:        "with json output",
			args:        []string{"-org=my-org", "-output=json"},
			expectedOrg: "my-org",
			expectedFmt: "json",
		},
		{
			name:           "multiple filters",
			args:           []string{"-org=my-org", "-status=applied", "-workspace=prod", "-output=json"},
			expectedOrg:    "my-org",
			expectedStatus: "applied",
			expectedWS:     "prod",
			expectedFmt:    "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &QueryRunListCommand{}

			flags := cmd.Meta.FlagSet("queryrun list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.status, "status", "", "Filter by run status")
			flags.StringVar(&cmd.operation, "operation", "", "Filter by operation type")
			flags.StringVar(&cmd.source, "source", "", "Filter by run source")
			flags.StringVar(&cmd.workspace, "workspace", "", "Filter by workspace name")
			flags.StringVar(&cmd.agentPool, "agent-pool", "", "Filter by agent pool name")
			flags.StringVar(&cmd.statusGroup, "status-group", "", "Filter by status group")
			flags.StringVar(&cmd.searchUser, "search-user", "", "Search by VCS username")
			flags.StringVar(&cmd.searchCommit, "search-commit", "", "Search by commit SHA")
			flags.StringVar(&cmd.searchBasic, "search-basic", "", "Basic search")
			flags.StringVar(&cmd.format, "output", "table", "Output format")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the status was set correctly
			if cmd.status != tt.expectedStatus {
				t.Errorf("expected status %q, got %q", tt.expectedStatus, cmd.status)
			}

			// Verify the operation was set correctly
			if cmd.operation != tt.expectedOp {
				t.Errorf("expected operation %q, got %q", tt.expectedOp, cmd.operation)
			}

			// Verify the source was set correctly
			if cmd.source != tt.expectedSource {
				t.Errorf("expected source %q, got %q", tt.expectedSource, cmd.source)
			}

			// Verify the workspace was set correctly
			if cmd.workspace != tt.expectedWS {
				t.Errorf("expected workspace %q, got %q", tt.expectedWS, cmd.workspace)
			}

			// Verify the agent pool was set correctly
			if cmd.agentPool != tt.expectedAgent {
				t.Errorf("expected agent pool %q, got %q", tt.expectedAgent, cmd.agentPool)
			}

			// Verify the status group was set correctly
			if cmd.statusGroup != tt.expectedSGroup {
				t.Errorf("expected status group %q, got %q", tt.expectedSGroup, cmd.statusGroup)
			}

			// Verify the search user was set correctly
			if cmd.searchUser != tt.expectedSUser {
				t.Errorf("expected search user %q, got %q", tt.expectedSUser, cmd.searchUser)
			}

			// Verify the search commit was set correctly
			if cmd.searchCommit != tt.expectedSCommit {
				t.Errorf("expected search commit %q, got %q", tt.expectedSCommit, cmd.searchCommit)
			}

			// Verify the search basic was set correctly
			if cmd.searchBasic != tt.expectedSBasic {
				t.Errorf("expected search basic %q, got %q", tt.expectedSBasic, cmd.searchBasic)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}

func TestQueryRunListRunNoRuns(t *testing.T) {
	ui := cli.NewMockUi()
	meta := newTestMeta(ui)
	meta.client = &client.Client{
		Client: &tfe.Client{
			Runs: &mockRunsService{
				listForOrganizationResponse: &tfe.OrganizationRunList{
					Items: []*tfe.Run{},
				},
			},
		},
	}

	cmd := &QueryRunListCommand{Meta: meta}
	code := cmd.Run([]string{"-organization=my-org"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if !strings.Contains(ui.OutputWriter.String(), "No runs found") {
		t.Fatalf("expected no runs output, got %q", ui.OutputWriter.String())
	}
}

func TestQueryRunListRunListsRuns(t *testing.T) {
	ui := cli.NewMockUi()
	meta := newTestMeta(ui)
	meta.client = &client.Client{
		Client: &tfe.Client{
			Runs: &mockRunsService{
				listForOrganizationResponse: &tfe.OrganizationRunList{
					Items: []*tfe.Run{
						{
							ID:       "run-001",
							Message:  "initial apply",
							Status:   tfe.RunApplied,
							Source:   tfe.RunSourceUI,
							CreatedAt: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
							Workspace: &tfe.Workspace{
								Name: "demo-workspace",
							},
						},
					},
				},
			},
		},
	}

	cmd := &QueryRunListCommand{Meta: meta}
	code := cmd.Run([]string{"-organization=my-org"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	if !strings.Contains(ui.OutputWriter.String(), "run-001") {
		t.Fatalf("expected run output, got %q", ui.OutputWriter.String())
	}
}

func TestQueryRunListRunError(t *testing.T) {
	ui := cli.NewMockUi()
	meta := newTestMeta(ui)
	meta.client = &client.Client{
		Client: &tfe.Client{
			Runs: &mockRunsService{
				listForOrganizationErr: errors.New("backend failure"),
			},
		},
	}

	cmd := &QueryRunListCommand{Meta: meta}
	code := cmd.Run([]string{"-organization=my-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if !strings.Contains(ui.ErrorWriter.String(), "Error listing runs: backend failure") {
		t.Fatalf("expected list error output, got %q", ui.ErrorWriter.String())
	}
}

type mockRunsService struct {
	listForOrganizationResponse *tfe.OrganizationRunList
	listForOrganizationErr      error
}

func (m *mockRunsService) List(_ context.Context, workspaceID string, options *tfe.RunListOptions) (*tfe.RunList, error) {
	return nil, nil
}

func (m *mockRunsService) ListForOrganization(_ context.Context, organization string, _ *tfe.RunListForOrganizationOptions) (*tfe.OrganizationRunList, error) {
	return m.listForOrganizationResponse, m.listForOrganizationErr
}

func (m *mockRunsService) Create(_ context.Context, options tfe.RunCreateOptions) (*tfe.Run, error) {
	return nil, nil
}

func (m *mockRunsService) Read(_ context.Context, runID string) (*tfe.Run, error) {
	return nil, nil
}

func (m *mockRunsService) ReadWithOptions(_ context.Context, runID string, options *tfe.RunReadOptions) (*tfe.Run, error) {
	return nil, nil
}

func (m *mockRunsService) Apply(_ context.Context, runID string, options tfe.RunApplyOptions) error {
	return nil
}

func (m *mockRunsService) Cancel(_ context.Context, runID string, options tfe.RunCancelOptions) error {
	return nil
}

func (m *mockRunsService) ForceCancel(_ context.Context, runID string, options tfe.RunForceCancelOptions) error {
	return nil
}

func (m *mockRunsService) ForceExecute(_ context.Context, runID string) error {
	return nil
}

func (m *mockRunsService) Discard(_ context.Context, runID string, options tfe.RunDiscardOptions) error {
	return nil
}
