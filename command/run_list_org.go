package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunListOrgCommand lists runs across an organization.
type RunListOrgCommand struct {
	Meta
	organization string
	user         string
	commit       string
	search       string
	status       string
	source       string
	operation    string
	agentPool    string
	statusGroup  string
	timeframe    string
	workspace    string
	include      string
	format       string
	runSvc       runOrgLister
}

// Run executes the run list-org command.
func (c *RunListOrgCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("run list-org")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.user, "user", "", "Filter by VCS username")
	flags.StringVar(&c.commit, "commit", "", "Filter by commit SHA")
	flags.StringVar(&c.search, "search", "", "Search runs by username, commit, run ID, or message")
	flags.StringVar(&c.status, "status", "", "Filter by run status (comma-separated)")
	flags.StringVar(&c.source, "source", "", "Filter by run source (comma-separated)")
	flags.StringVar(&c.operation, "operation", "", "Filter by run operation (comma-separated)")
	flags.StringVar(&c.agentPool, "agent-pool", "", "Filter by agent pool names (comma-separated)")
	flags.StringVar(&c.statusGroup, "status-group", "", "Filter by run status group (comma-separated)")
	flags.StringVar(&c.timeframe, "timeframe", "", "Filter by timeframe (comma-separated)")
	flags.StringVar(&c.workspace, "workspace", "", "Filter by workspace names (comma-separated)")
	flags.StringVar(&c.include, "include", "", "Comma-separated related resources to include")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	options := &tfe.RunListForOrganizationOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 50,
		},
		User:           c.user,
		Commit:         c.commit,
		Basic:          c.search,
		Status:         c.status,
		Source:         c.source,
		Operation:      c.operation,
		AgentPoolNames: c.agentPool,
		StatusGroup:    c.statusGroup,
		Timeframe:      c.timeframe,
		WorkspaceNames: c.workspace,
	}
	if c.include != "" {
		for _, include := range splitCommaList(c.include) {
			if include == "" {
				continue
			}
			options.Include = append(options.Include, tfe.RunIncludeOpt(include))
		}
	}

	runs, err := c.runService(client).ListForOrganization(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing organization runs: %s", err))
		return 1
	}

	formatter := c.Meta.NewFormatter(c.format)
	if len(runs.Items) == 0 {
		c.Ui.Output("No runs found")
		return 0
	}

	headers := []string{"ID", "Workspace", "Status", "Source", "Message", "Created At"}
	var rows [][]string
	for _, run := range runs.Items {
		workspaceName := ""
		if run.Workspace != nil {
			workspaceName = run.Workspace.Name
		}

		message := run.Message
		if len(message) > 50 {
			message = message[:47] + "..."
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

func (c *RunListOrgCommand) runService(client *client.Client) runOrgLister {
	if c.runSvc != nil {
		return c.runSvc
	}
	return client.Runs
}

// Help returns help text for the run list-org command.
func (c *RunListOrgCommand) Help() string {
	helpText := `
Usage: hcptf run list-org [options]

  List runs across an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -user=<name>         Filter by VCS username
  -commit=<sha>        Filter by commit SHA
  -search=<query>      Search username/commit/run/message
  -status=<values>     Filter by status (comma-separated)
  -source=<values>     Filter by source (comma-separated)
  -operation=<values>  Filter by operation (comma-separated)
  -agent-pool=<values> Filter by agent pool names (comma-separated)
  -status-group=<vals> Filter by status groups (comma-separated)
  -timeframe=<vals>    Filter by timeframe values (comma-separated)
  -workspace=<values>  Filter by workspace names (comma-separated)
  -include=<values>    Comma-separated include values
  -output=<format>     Output format: table (default) or json

Example:

  hcptf run list-org -org=my-org
  hcptf run list-org -org=my-org -status=planned,applied -workspace=prod
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run list-org command.
func (c *RunListOrgCommand) Synopsis() string {
	return "List runs across an organization"
}
