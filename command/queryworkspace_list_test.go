package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestQueryWorkspaceListRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &QueryWorkspaceListCommand{
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

func TestQueryWorkspaceListHelp(t *testing.T) {
	cmd := &QueryWorkspaceListCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf queryworkspace list") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
	}
	if !strings.Contains(help, "-search") {
		t.Error("Help should mention -search flag")
	}
	if !strings.Contains(help, "-tags") {
		t.Error("Help should mention -tags flag")
	}
	if !strings.Contains(help, "-exclude-tags") {
		t.Error("Help should mention -exclude-tags flag")
	}
	if !strings.Contains(help, "-wildcard") {
		t.Error("Help should mention -wildcard flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
}

func TestQueryWorkspaceListSynopsis(t *testing.T) {
	cmd := &QueryWorkspaceListCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Search workspaces across organization" {
		t.Errorf("expected 'Search workspaces across organization', got %q", synopsis)
	}
}

func TestQueryWorkspaceListFlagParsing(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedOrg       string
		expectedSearch    string
		expectedTags      string
		expectedExclTags  string
		expectedWildcard  string
		expectedFmt       string
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
			name:           "with search filter",
			args:           []string{"-org=my-org", "-search=production"},
			expectedOrg:    "my-org",
			expectedSearch: "production",
			expectedFmt:    "table",
		},
		{
			name:         "with tags filter",
			args:         []string{"-org=my-org", "-tags=env:prod,team:platform"},
			expectedOrg:  "my-org",
			expectedTags: "env:prod,team:platform",
			expectedFmt:  "table",
		},
		{
			name:             "with exclude tags filter",
			args:             []string{"-org=my-org", "-exclude-tags=archived,deprecated"},
			expectedOrg:      "my-org",
			expectedExclTags: "archived,deprecated",
			expectedFmt:      "table",
		},
		{
			name:             "with wildcard filter",
			args:             []string{"-org=my-org", "-wildcard=prod-*"},
			expectedOrg:      "my-org",
			expectedWildcard: "prod-*",
			expectedFmt:      "table",
		},
		{
			name:        "with json output",
			args:        []string{"-org=my-org", "-output=json"},
			expectedOrg: "my-org",
			expectedFmt: "json",
		},
		{
			name:           "multiple filters",
			args:           []string{"-org=my-org", "-search=app", "-tags=env:prod", "-output=json"},
			expectedOrg:    "my-org",
			expectedSearch: "app",
			expectedTags:   "env:prod",
			expectedFmt:    "json",
		},
		{
			name:             "all filters",
			args:             []string{"-org=my-org", "-search=app", "-tags=env:prod", "-exclude-tags=archived", "-wildcard=app-*"},
			expectedOrg:      "my-org",
			expectedSearch:   "app",
			expectedTags:     "env:prod",
			expectedExclTags: "archived",
			expectedWildcard: "app-*",
			expectedFmt:      "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &QueryWorkspaceListCommand{}

			flags := cmd.Meta.FlagSet("queryworkspace list")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.search, "search", "", "Search query for workspace name")
			flags.StringVar(&cmd.tags, "tags", "", "Filter by tags")
			flags.StringVar(&cmd.excludeTags, "exclude-tags", "", "Exclude workspaces with tags")
			flags.StringVar(&cmd.wildcard, "wildcard", "", "Wildcard filter for workspace name")
			flags.StringVar(&cmd.format, "output", "table", "Output format")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the search was set correctly
			if cmd.search != tt.expectedSearch {
				t.Errorf("expected search %q, got %q", tt.expectedSearch, cmd.search)
			}

			// Verify the tags was set correctly
			if cmd.tags != tt.expectedTags {
				t.Errorf("expected tags %q, got %q", tt.expectedTags, cmd.tags)
			}

			// Verify the exclude tags was set correctly
			if cmd.excludeTags != tt.expectedExclTags {
				t.Errorf("expected exclude tags %q, got %q", tt.expectedExclTags, cmd.excludeTags)
			}

			// Verify the wildcard was set correctly
			if cmd.wildcard != tt.expectedWildcard {
				t.Errorf("expected wildcard %q, got %q", tt.expectedWildcard, cmd.wildcard)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
