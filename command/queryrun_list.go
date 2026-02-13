package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// QueryRunListCommand is a command to search runs across an organization
type QueryRunListCommand struct {
	Meta
	organization string
	status       string
	operation    string
	source       string
	workspace    string
	agentPool    string
	statusGroup  string
	searchUser   string
	searchCommit string
	searchBasic  string
	format       string
}

// Run executes the query run list command
func (c *QueryRunListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("queryrun list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.status, "status", "", "Filter by run status (comma-separated)")
	flags.StringVar(&c.operation, "operation", "", "Filter by operation type (comma-separated)")
	flags.StringVar(&c.source, "source", "", "Filter by run source (comma-separated)")
	flags.StringVar(&c.workspace, "workspace", "", "Filter by workspace name (comma-separated)")
	flags.StringVar(&c.agentPool, "agent-pool", "", "Filter by agent pool name (comma-separated)")
	flags.StringVar(&c.statusGroup, "status-group", "", "Filter by status group (final, non_final, discardable)")
	flags.StringVar(&c.searchUser, "search-user", "", "Search by VCS username")
	flags.StringVar(&c.searchCommit, "search-commit", "", "Search by commit SHA")
	flags.StringVar(&c.searchBasic, "search-basic", "", "Basic search (username, commit, run ID, or message)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build list options
	options := &tfe.RunListForOrganizationOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 50,
		},
	}

	if c.status != "" {
		options.Status = c.status
	}

	if c.operation != "" {
		options.Operation = c.operation
	}

	if c.source != "" {
		options.Source = c.source
	}

	if c.workspace != "" {
		options.WorkspaceNames = c.workspace
	}

	if c.agentPool != "" {
		options.AgentPoolNames = c.agentPool
	}

	if c.statusGroup != "" {
		options.StatusGroup = c.statusGroup
	}

	if c.searchUser != "" {
		options.User = c.searchUser
	}

	if c.searchCommit != "" {
		options.Commit = c.searchCommit
	}

	if c.searchBasic != "" {
		options.Basic = c.searchBasic
	}

	// List runs across organization
	runs, err := client.Runs.ListForOrganization(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing runs: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(runs.Items) == 0 {
		c.Ui.Output("No runs found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Workspace", "Status", "Source", "Message", "Created At"}
	var rows [][]string

	for _, run := range runs.Items {
		message := run.Message
		if len(message) > 50 {
			message = message[:47] + "..."
		}

		workspaceName := ""
		if run.Workspace != nil {
			workspaceName = run.Workspace.Name
		}

		rows = append(rows, []string{
			run.ID,
			workspaceName,
			string(run.Status),
			string(run.Source),
			message,
			run.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the query run list command
func (c *QueryRunListCommand) Help() string {
	helpText := `
Usage: hcptf queryrun list [options]

  Search and list runs across an organization with filters.

  This command allows you to search for runs across all workspaces in an
  organization using various filters and search criteria.

Options:

  -organization=<name>   Organization name (required)
  -org=<name>           Alias for -organization
  -status=<status>      Filter by run status (comma-separated)
                        Examples: pending, planning, planned, applying, applied
  -operation=<op>       Filter by operation (comma-separated)
                        Options: plan_only, plan_and_apply, save_plan,
                        refresh_only, destroy, empty_apply
  -source=<source>      Filter by run source (comma-separated)
                        Options: tfe-ui, tfe-api, tfe-configuration-version
  -workspace=<name>     Filter by workspace name (comma-separated)
  -agent-pool=<name>    Filter by agent pool name (comma-separated)
  -status-group=<grp>   Filter by status group
                        Options: final, non_final, discardable
  -search-user=<user>   Search by VCS username
  -search-commit=<sha>  Search by commit SHA
  -search-basic=<term>  Basic search (username, commit, run ID, or message)
  -output=<format>      Output format: table (default) or json

Example:

  # List all runs in organization
  hcptf queryrun list -org=my-org

  # List runs by workspace
  hcptf queryrun list -org=my-org -workspace=prod-app

  # List runs by status
  hcptf queryrun list -org=my-org -status=applied,applying

  # List runs from API
  hcptf queryrun list -org=my-org -source=tfe-api

  # Search by commit
  hcptf queryrun list -org=my-org -search-commit=abc123def

  # Output as JSON
  hcptf queryrun list -org=my-org -output=json

Note:

  This endpoint is rate-limited to 30 requests per minute.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the query run list command
func (c *QueryRunListCommand) Synopsis() string {
	return "Search runs across organization"
}
