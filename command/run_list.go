package command

import (
	"fmt"
	"strings"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunListCommand is a command to list runs
type RunListCommand struct {
	Meta
	organization string
	workspace    string
	user         string
	commit       string
	search       string
	status       string
	source       string
	operation    string
	include      string
	format       string
	workspaceSvc workspaceReader
	runSvc       runLister
}

// Run executes the run list command
func (c *RunListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("run list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "name", "", "Workspace name (required)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (alias)")
	flags.StringVar(&c.user, "user", "", "Filter runs by VCS username")
	flags.StringVar(&c.commit, "commit", "", "Filter runs by commit SHA")
	flags.StringVar(&c.search, "search", "", "Basic search across username/commit/run/message")
	flags.StringVar(&c.status, "status", "", "Filter by run status (comma-separated)")
	flags.StringVar(&c.source, "source", "", "Filter by run source (comma-separated)")
	flags.StringVar(&c.operation, "operation", "", "Filter by run operation type (comma-separated)")
	flags.StringVar(&c.include, "include", "", "Comma-separated related resources to include")
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

	if c.workspace == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Get workspace first
	ws, err := c.workspaceService(client).Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	options := &tfe.RunListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 50,
		},
		User:      c.user,
		Commit:    c.commit,
		Search:    c.search,
		Status:    c.status,
		Source:    c.source,
		Operation: c.operation,
	}
	if c.include != "" {
		for _, include := range splitCommaList(c.include) {
			if include == "" {
				continue
			}
			options.Include = append(options.Include, tfe.RunIncludeOpt(include))
		}
	}

	// List runs
	runs, err := c.runService(client).List(client.Context(), ws.ID, options)
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
	headers := []string{"ID", "Status", "Source", "Message", "CreatedAt"}
	var rows [][]string

	for _, run := range runs.Items {
		message := run.Message
		if len(message) > 50 {
			message = message[:47] + "..."
		}

		rows = append(rows, []string{
			run.ID,
			string(run.Status),
			string(run.Source),
			message,
			run.CreatedAt.Format(time.RFC3339),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *RunListCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

func (c *RunListCommand) runService(client *client.Client) runLister {
	if c.runSvc != nil {
		return c.runSvc
	}
	return client.Runs
}

// Help returns help text for the run list command
func (c *RunListCommand) Help() string {
	helpText := `
Usage: hcptf workspace run list [options]

  List runs for a workspace.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Workspace name (required)
  -workspace=<name>    Alias for -name
  -user=<name>         Filter by VCS username
  -commit=<sha>        Filter by commit SHA
  -search=<query>      Search username/commit/run/message
  -status=<values>     Filter by status (comma-separated)
  -source=<values>     Filter by source (comma-separated)
  -operation=<values>  Filter by operation type (comma-separated)
  -include=<values>    Comma-separated include values
  -output=<format>     Output format: table (default) or json

Example:

  hcptf workspace run list -org=my-org -name=my-workspace
  hcptf workspace run list -org=my-org -name=my-workspace -status=planned,applied -include=plan
  hcptf workspace run list -org=my-org -name=prod -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run list command
func (c *RunListCommand) Synopsis() string {
	return "List runs for a workspace"
}
