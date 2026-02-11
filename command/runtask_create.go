package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// RunTaskCreateCommand is a command to create a run task
type RunTaskCreateCommand struct {
	Meta
	organization string
	name         string
	url          string
	hmacKey      string
	category     string
	enabled      bool
	description  string
	format       string
}

// Run executes the run task create command
func (c *RunTaskCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("runtask create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Run task name (required)")
	flags.StringVar(&c.url, "url", "", "Run task URL (required)")
	flags.StringVar(&c.hmacKey, "hmac-key", "", "HMAC key for request authentication (optional)")
	flags.StringVar(&c.category, "category", "task", "Run task category: task or advisory (default: task)")
	flags.BoolVar(&c.enabled, "enabled", true, "Enable run task (default: true)")
	flags.StringVar(&c.description, "description", "", "Run task description (optional)")
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

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.url == "" {
		c.Ui.Error("Error: -url flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate category
	if c.category != "task" && c.category != "advisory" {
		c.Ui.Error("Error: -category must be 'task' or 'advisory'")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build create options
	options := tfe.RunTaskCreateOptions{
		Name:     c.name,
		URL:      c.url,
		Category: c.category,
		Enabled:  tfe.Bool(c.enabled),
	}

	if c.hmacKey != "" {
		options.HMACKey = tfe.String(c.hmacKey)
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	// Create run task
	runTask, err := client.RunTasks.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating run task: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Run task '%s' created successfully", runTask.Name))

	// Show run task details
	data := map[string]interface{}{
		"ID":          runTask.ID,
		"Name":        runTask.Name,
		"URL":         runTask.URL,
		"Category":    runTask.Category,
		"Enabled":     runTask.Enabled,
		"Description": runTask.Description,
		"HMACKey":     "***",
	}

	if runTask.HMACKey != nil && *runTask.HMACKey != "" {
		data["HMACKey"] = "set"
	} else {
		data["HMACKey"] = "not set"
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the run task create command
func (c *RunTaskCreateCommand) Help() string {
	helpText := `
Usage: hcptf runtask create [options]

  Create a new run task in an organization.

  Run tasks integrate external systems into the Terraform run workflow. They allow
  you to trigger external systems at specific points during a run and receive
  results that can affect the run's status.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Run task name (required)
  -url=<url>           Run task URL (required)
  -hmac-key=<key>      HMAC key for request authentication (optional)
  -category=<cat>      Run task category: task or advisory (default: task)
                       task: Can block runs if they fail
                       advisory: Results are informational only
  -enabled             Enable run task (default: true)
  -description=<text>  Run task description (optional)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf runtask create -org=my-org -name="Security Scan" \
    -url=https://security-scanner.example.com/webhook \
    -hmac-key=secret123 -category=task

  hcptf runtask create -org=my-org -name="Cost Estimation" \
    -url=https://cost-estimator.example.com/webhook \
    -category=advisory -description="Estimates infrastructure costs"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run task create command
func (c *RunTaskCreateCommand) Synopsis() string {
	return "Create a new run task"
}
