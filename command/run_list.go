package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunListCommand is a command to list runs
type RunListCommand struct {
	Meta
	organization string
	workspace    string
	format       string
	workspaceSvc workspaceReader
	runSvc       runLister
}

// Run executes the run list command
func (c *RunListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("run list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
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
		c.Ui.Error("Error: -workspace flag is required")
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

	// List runs
	runs, err := c.runService(client).List(client.Context(), ws.ID, &tfe.RunListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 50,
		},
	})
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
	headers := []string{"ID", "Status", "Source", "Message", "Created At"}
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
			run.CreatedAt.Format("2006-01-02 15:04:05"),
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
Usage: hcptf run list [options]

  List runs for a workspace.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf run list -org=my-org -workspace=my-workspace
  hcptf run list -org=my-org -workspace=prod -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run list command
func (c *RunListCommand) Synopsis() string {
	return "List runs for a workspace"
}
