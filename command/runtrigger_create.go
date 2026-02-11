package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// RunTriggerCreateCommand is a command to create a run trigger
type RunTriggerCreateCommand struct {
	Meta
	organization      string
	workspace         string
	sourceWorkspace   string
	format            string
	workspaceSvc      workspaceReader
	runTriggerSvc     runTriggerCreator
}

// Run executes the run trigger create command
func (c *RunTriggerCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("runtrigger create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Target workspace name (required)")
	flags.StringVar(&c.sourceWorkspace, "source-workspace", "", "Source workspace name (required)")
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

	if c.sourceWorkspace == "" {
		c.Ui.Error("Error: -source-workspace flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Get target workspace
	targetWs, err := c.workspaceService(client).Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading target workspace: %s", err))
		return 1
	}

	// Get source workspace
	sourceWs, err := c.workspaceService(client).Read(client.Context(), c.organization, c.sourceWorkspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading source workspace: %s", err))
		return 1
	}

	// Create run trigger
	runTrigger, err := c.runTriggerService(client).Create(client.Context(), targetWs.ID, tfe.RunTriggerCreateOptions{
		Sourceable: &tfe.Workspace{
			ID: sourceWs.ID,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating run trigger: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":             runTrigger.ID,
		"WorkspaceName":  runTrigger.WorkspaceName,
		"SourceableName": runTrigger.SourceableName,
		"CreatedAt":      runTrigger.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	c.Ui.Output(fmt.Sprintf("Run trigger created successfully"))
	formatter.KeyValue(data)
	return 0
}

func (c *RunTriggerCreateCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

func (c *RunTriggerCreateCommand) runTriggerService(client *client.Client) runTriggerCreator {
	if c.runTriggerSvc != nil {
		return c.runTriggerSvc
	}
	return client.RunTriggers
}

// Help returns help text for the run trigger create command
func (c *RunTriggerCreateCommand) Help() string {
	helpText := `
Usage: hcptf runtrigger create [options]

  Create a run trigger to link two workspaces for automatic orchestration.
  When the source workspace completes a run successfully, the target workspace
  will automatically start a new run.

Options:

  -organization=<name>       Organization name (required)
  -org=<name>               Alias for -organization
  -workspace=<name>         Target workspace name - will start runs automatically (required)
  -source-workspace=<name>  Source workspace name - triggers runs when complete (required)
  -output=<format>          Output format: table (default) or json

Examples:

  # Create a run trigger so prod-app runs when prod-networking completes
  hcptf runtrigger create -org=my-org \
    -workspace=prod-app \
    -source-workspace=prod-networking

  # Output as JSON
  hcptf runtrigger create -org=my-org \
    -workspace=prod-app \
    -source-workspace=prod-networking \
    -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run trigger create command
func (c *RunTriggerCreateCommand) Synopsis() string {
	return "Create a run trigger to connect workspaces"
}
