package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// RunTaskUpdateCommand is a command to update a run task
type RunTaskUpdateCommand struct {
	Meta
	id          string
	name        string
	url         string
	hmacKey     string
	category    string
	enabled     string
	description string
	format      string
}

// Run executes the run task update command
func (c *RunTaskUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("runtask update")
	flags.StringVar(&c.id, "id", "", "Run task ID (required)")
	flags.StringVar(&c.name, "name", "", "Run task name")
	flags.StringVar(&c.url, "url", "", "Run task URL")
	flags.StringVar(&c.hmacKey, "hmac-key", "", "HMAC key for request authentication")
	flags.StringVar(&c.category, "category", "", "Run task category: task or advisory")
	flags.StringVar(&c.enabled, "enabled", "", "Enable run task (true/false)")
	flags.StringVar(&c.description, "description", "", "Run task description")
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

	// Build update options
	options := tfe.RunTaskUpdateOptions{}

	if c.name != "" {
		options.Name = tfe.String(c.name)
	}

	if c.url != "" {
		options.URL = tfe.String(c.url)
	}

	if c.hmacKey != "" {
		options.HMACKey = tfe.String(c.hmacKey)
	}

	if c.category != "" {
		if c.category != "task" && c.category != "advisory" {
			c.Ui.Error("Error: -category must be 'task' or 'advisory'")
			return 1
		}
		options.Category = tfe.String(c.category)
	}

	if c.enabled != "" {
		if c.enabled == "true" {
			options.Enabled = tfe.Bool(true)
		} else if c.enabled == "false" {
			options.Enabled = tfe.Bool(false)
		} else {
			c.Ui.Error("Error: -enabled must be 'true' or 'false'")
			return 1
		}
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	// Update run task
	runTask, err := client.RunTasks.Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating run task: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Run task '%s' updated successfully", runTask.Name))

	// Show run task details
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

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the run task update command
func (c *RunTaskUpdateCommand) Help() string {
	helpText := `
Usage: hcptf runtask update [options]

  Update run task settings.

Options:

  -id=<id>             Run task ID (required)
  -name=<name>         Run task name
  -url=<url>           Run task URL
  -hmac-key=<key>      HMAC key for request authentication
  -category=<cat>      Run task category: task or advisory
                       task: Can block runs if they fail
                       advisory: Results are informational only
  -enabled=<bool>      Enable run task (true/false)
  -description=<text>  Run task description
  -output=<format>     Output format: table (default) or json

Example:

  hcptf runtask update -id=task-ABC123 -enabled=false
  hcptf runtask update -id=task-ABC123 -name="Updated Security Scan" \
    -url=https://new-scanner.example.com/webhook
  hcptf runtask update -id=task-ABC123 -category=advisory
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run task update command
func (c *RunTaskUpdateCommand) Synopsis() string {
	return "Update run task settings"
}
