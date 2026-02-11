package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// RunTaskListCommand is a command to list run tasks
type RunTaskListCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the run task list command
func (c *RunTaskListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("runtask list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List run tasks
	runTasks, err := client.RunTasks.List(client.Context(), c.organization, &tfe.RunTaskListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing run tasks: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(runTasks.Items) == 0 {
		c.Ui.Output("No run tasks found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "URL", "Category", "Enabled", "HMAC Key"}
	var rows [][]string

	for _, rt := range runTasks.Items {
		enabled := "false"
		if rt.Enabled {
			enabled = "true"
		}

		hmacKey := "not set"
		if rt.HMACKey != nil && *rt.HMACKey != "" {
			hmacKey = "set"
		}

		rows = append(rows, []string{
			rt.ID,
			rt.Name,
			rt.URL,
			rt.Category,
			enabled,
			hmacKey,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the run task list command
func (c *RunTaskListCommand) Help() string {
	helpText := `
Usage: hcptf runtask list [options]

  List run tasks in an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf runtask list -org=my-org
  hcptf runtask list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run task list command
func (c *RunTaskListCommand) Synopsis() string {
	return "List run tasks in an organization"
}
