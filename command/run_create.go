package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunCreateCommand is a command to create a run
type RunCreateCommand struct {
	Meta
	organization string
	workspace    string
	message      string
	destroy      bool
	refreshOnly  bool
	format       string
	workspaceSvc workspaceReader
	runSvc       runCreator
}

// Run executes the run create command
func (c *RunCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("run create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "name", "", "Workspace name (required)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (alias)")
	flags.StringVar(&c.message, "message", "", "Run message")
	flags.BoolVar(&c.destroy, "destroy", false, "Create a destroy plan")
	flags.BoolVar(&c.refreshOnly, "refresh-only", false, "Create a refresh-only run that detects drift only")
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

	// Create run options
	options := tfe.RunCreateOptions{
		Workspace: ws,
		Message:   tfe.String(c.message),
	}

	if c.destroy {
		options.IsDestroy = tfe.Bool(true)
	}

	if c.refreshOnly {
		options.RefreshOnly = tfe.Bool(true)
	}

	// Create run
	run, err := c.runService(client).Create(client.Context(), options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating run: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if c.format != "json" {
		c.Ui.Output("Run created successfully")
	}

	// Show run details
	data := map[string]interface{}{
		"ID":        run.ID,
		"Status":    run.Status,
		"Message":   run.Message,
		"IsDestroy": run.IsDestroy,
		"Source":    run.Source,
		"CreatedAt": run.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *RunCreateCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

func (c *RunCreateCommand) runService(client *client.Client) runCreator {
	if c.runSvc != nil {
		return c.runSvc
	}
	return client.Runs
}

// Help returns help text for the run create command
func (c *RunCreateCommand) Help() string {
	helpText := `
Usage: hcptf workspace run create [options]

  Create a new run for a workspace.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Workspace name (required)
  -workspace=<name>    Alias for -name
  -message=<text>      Run message
  -destroy             Create a destroy plan
  -refresh-only        Create a run that only checks drift (refresh-only)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf workspace run create -org=my-org -name=my-workspace -message="Deploy changes"
  hcptf workspace run create -org=my-org -name=prod -destroy
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run create command
func (c *RunCreateCommand) Synopsis() string {
	return "Create a new run"
}
