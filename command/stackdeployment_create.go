package command

import (
	"fmt"
	"strings"

)

// StackDeploymentCreateCommand is a command to create/trigger a stack deployment
type StackDeploymentCreateCommand struct {
	Meta
	stackID string
	format  string
}

// Run executes the stack deployment create command
func (c *StackDeploymentCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stackdeployment create")
	flags.StringVar(&c.stackID, "stack-id", "", "Stack ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.stackID == "" {
		c.Ui.Error("Error: -stack-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// For VCS-backed stacks, fetch the latest configuration from VCS
	// This will trigger a new deployment
	_, err = client.Stacks.FetchLatestFromVcs(client.Context(), c.stackID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error triggering stack deployment: %s", err))
		c.Ui.Error("")
		c.Ui.Error("Note: For manually configured stacks (non-VCS), create a stack")
		c.Ui.Error("configuration first using:")
		c.Ui.Error("  hcptf stackconfiguration create -stack-id=<id>")
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Stack deployment triggered for stack '%s'", c.stackID))
	c.Ui.Output("Fetching latest configuration from VCS...")

	data := map[string]interface{}{
		"StackID": c.stackID,
		"Status":  "Configuration fetch initiated",
		"Message": "Check stack configurations for deployment progress",
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the stack deployment create command
func (c *StackDeploymentCreateCommand) Help() string {
	helpText := `
Usage: hcptf stackdeployment create [options]

  Trigger a new stack deployment by fetching the latest configuration from VCS.
  This command is for VCS-backed stacks only.

  For manually configured stacks, create a stack configuration first:
    hcptf stackconfiguration create -stack-id=<id>

Options:

  -stack-id=<id>    Stack ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf stackdeployment create -stack-id=st-abc123
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack deployment create command
func (c *StackDeploymentCreateCommand) Synopsis() string {
	return "Trigger a new stack deployment"
}
