package command

import (
	"fmt"
	"strings"

)

// RunTaskReadCommand is a command to read run task details
type RunTaskReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the run task read command
func (c *RunTaskReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("runtask read")
	flags.StringVar(&c.id, "id", "", "Run task ID (required)")
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

	// Read run task
	runTask, err := client.RunTasks.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading run task: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":          runTask.ID,
		"Name":        runTask.Name,
		"URL":         runTask.URL,
		"Category":    runTask.Category,
		"Enabled":     runTask.Enabled,
		"Description": runTask.Description,
	}

	if runTask.HMACKey != nil && *runTask.HMACKey != "" {
		data["HMACKey"] = "set"
	} else {
		data["HMACKey"] = "not set"
	}

	// Add organization info if available
	if runTask.Organization != nil {
		data["Organization"] = runTask.Organization.Name
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the run task read command
func (c *RunTaskReadCommand) Help() string {
	helpText := `
Usage: hcptf runtask read [options]

  Read run task details by ID.

Options:

  -id=<id>          Run task ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf runtask read -id=task-ABC123
  hcptf runtask read -id=task-ABC123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run task read command
func (c *RunTaskReadCommand) Synopsis() string {
	return "Read run task details by ID"
}
