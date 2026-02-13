package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunTriggerReadCommand is a command to read run trigger details
type RunTriggerReadCommand struct {
	Meta
	id            string
	format        string
	runTriggerSvc runTriggerReader
}

// Run executes the run trigger read command
func (c *RunTriggerReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("runtrigger read")
	flags.StringVar(&c.id, "id", "", "Run trigger ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read run trigger
	runTrigger, err := c.runTriggerService(client).Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading run trigger: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":             runTrigger.ID,
		"WorkspaceName":  runTrigger.WorkspaceName,
		"SourceableName": runTrigger.SourceableName,
		"CreatedAt":      runTrigger.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if runTrigger.Workspace != nil {
		data["WorkspaceID"] = runTrigger.Workspace.ID
	}

	if runTrigger.Sourceable != nil {
		data["SourceableID"] = runTrigger.Sourceable.ID
	}

	formatter.KeyValue(data)
	return 0
}

func (c *RunTriggerReadCommand) runTriggerService(client *client.Client) runTriggerReader {
	if c.runTriggerSvc != nil {
		return c.runTriggerSvc
	}
	return client.RunTriggers
}

// Help returns help text for the run trigger read command
func (c *RunTriggerReadCommand) Help() string {
	helpText := `
Usage: hcptf runtrigger read [options]

  Show details of a run trigger.

Options:

  -id=<id>          Run trigger ID (required)
  -output=<format>  Output format: table (default) or json

Examples:

  # Show run trigger details
  hcptf runtrigger read -id=rt-3yVQZvHzf5j3WRJ1

  # Output as JSON
  hcptf runtrigger read -id=rt-3yVQZvHzf5j3WRJ1 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run trigger read command
func (c *RunTriggerReadCommand) Synopsis() string {
	return "Show run trigger details"
}
