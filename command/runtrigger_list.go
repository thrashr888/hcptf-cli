package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// RunTriggerListCommand is a command to list run triggers
type RunTriggerListCommand struct {
	Meta
	organization   string
	workspace      string
	triggerType    string
	format         string
	workspaceSvc   workspaceReader
	runTriggerSvc  runTriggerLister
}

// Run executes the run trigger list command
func (c *RunTriggerListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("runtrigger list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
	flags.StringVar(&c.triggerType, "type", "inbound", "Run trigger type: inbound or outbound (default: inbound)")
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

	// Validate trigger type
	if c.triggerType != "inbound" && c.triggerType != "outbound" {
		c.Ui.Error("Error: -type must be either 'inbound' or 'outbound'")
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

	// List run triggers
	runTriggers, err := c.runTriggerService(client).List(client.Context(), ws.ID, &tfe.RunTriggerListOptions{
		RunTriggerType: tfe.RunTriggerFilterOp(c.triggerType),
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing run triggers: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(runTriggers.Items) == 0 {
		c.Ui.Output(fmt.Sprintf("No %s run triggers found", c.triggerType))
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Workspace", "Sourceable", "Created At"}
	var rows [][]string

	for _, rt := range runTriggers.Items {
		rows = append(rows, []string{
			rt.ID,
			rt.WorkspaceName,
			rt.SourceableName,
			rt.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *RunTriggerListCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

func (c *RunTriggerListCommand) runTriggerService(client *client.Client) runTriggerLister {
	if c.runTriggerSvc != nil {
		return c.runTriggerSvc
	}
	return client.RunTriggers
}

// Help returns help text for the run trigger list command
func (c *RunTriggerListCommand) Help() string {
	helpText := `
Usage: hcptf runtrigger list [options]

  List run triggers for a workspace. Run triggers link workspaces to create
  automatic orchestration - when a source workspace completes a run, the
  target workspace automatically starts a new run.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -type=<type>         Run trigger type: inbound or outbound (default: inbound)
                       - inbound: Triggers that cause runs in this workspace
                       - outbound: Triggers that this workspace causes in other workspaces
  -output=<format>     Output format: table (default) or json

Examples:

  # List inbound triggers (workspaces that trigger this workspace)
  hcptf runtrigger list -org=my-org -workspace=my-workspace

  # List outbound triggers (workspaces this workspace triggers)
  hcptf runtrigger list -org=my-org -workspace=my-workspace -type=outbound

  # Output as JSON
  hcptf runtrigger list -org=my-org -workspace=prod -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run trigger list command
func (c *RunTriggerListCommand) Synopsis() string {
	return "List run triggers for a workspace"
}
