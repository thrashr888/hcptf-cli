package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// StackStateListCommand is a command to list stack states
type StackStateListCommand struct {
	Meta
	stackID string
	format  string
}

// Run executes the stack state list command
func (c *StackStateListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stackstate list")
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

	// List stack states
	states, err := client.StackStates.List(client.Context(), c.stackID, &tfe.StackStateListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing stack states: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(states.Items) == 0 {
		c.Ui.Output("No stack states found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Generation", "Deployment", "Status", "Current", "Resources"}
	var rows [][]string

	for _, state := range states.Items {
		isCurrent := "no"
		if state.IsCurrent {
			isCurrent = "yes"
		}

		rows = append(rows, []string{
			state.ID,
			fmt.Sprintf("%d", state.Generation),
			state.Deployment,
			state.Status,
			isCurrent,
			fmt.Sprintf("%d", state.ResourceInstanceCount),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the stack state list command
func (c *StackStateListCommand) Help() string {
	helpText := `
Usage: hcptf stackstate list [options]

  List state versions for a stack. Stack states provide insight into the
  state of a particular stack deployment.

Options:

  -stack-id=<id>    Stack ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf stackstate list -stack-id=st-abc123
  hcptf stackstate list -stack-id=st-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack state list command
func (c *StackStateListCommand) Synopsis() string {
	return "List state versions for a stack"
}
