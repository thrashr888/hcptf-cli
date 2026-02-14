package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// StackConfigurationListCommand is a command to list stack configurations
type StackConfigurationListCommand struct {
	Meta
	stackID string
	format  string
}

// Run executes the stack configuration list command
func (c *StackConfigurationListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stackconfiguration list")
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

	// List stack configurations
	configs, err := client.StackConfigurations.List(client.Context(), c.stackID, &tfe.StackConfigurationListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing stack configurations: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(configs.Items) == 0 {
		c.Ui.Output("No stack configurations found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Sequence", "Status", "Speculative", "Created"}
	var rows [][]string

	for _, config := range configs.Items {
		speculative := "false"
		if config.Speculative {
			speculative = "true"
		}

		rows = append(rows, []string{
			config.ID,
			fmt.Sprintf("%d", config.SequenceNumber),
			string(config.Status),
			speculative,
			config.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the stack configuration list command
func (c *StackConfigurationListCommand) Help() string {
	helpText := `
Usage: hcptf stackconfiguration list [options]

  List stack configurations for a stack. A stack configuration represents
  a snapshot of all the pieces that make up your stack.

Options:

  -stack-id=<id>    Stack ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf stackconfiguration list -stack-id=st-abc123
  hcptf stackconfiguration list -stack-id=st-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack configuration list command
func (c *StackConfigurationListCommand) Synopsis() string {
	return "List stack configurations for a stack"
}
