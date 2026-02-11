package command

import (
	"fmt"
	"strings"
)

// StackDeleteCommand is a command to delete a stack
type StackDeleteCommand struct {
	Meta
	stackID string
	force   bool
}

// Run executes the stack delete command
func (c *StackDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stack delete")
	flags.StringVar(&c.stackID, "id", "", "Stack ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.stackID == "" {
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

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete stack '%s'? (yes/no): ", c.stackID))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete stack
	err = client.Stacks.Delete(client.Context(), c.stackID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting stack: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Stack '%s' deleted successfully", c.stackID))
	return 0
}

// Help returns help text for the stack delete command
func (c *StackDeleteCommand) Help() string {
	helpText := `
Usage: hcptf stack delete [options]

  Delete a stack. This will remove the stack and all associated configurations.

Options:

  -id=<stack-id>  Stack ID (required)
  -force          Force delete without confirmation

Example:

  hcptf stack delete -id=st-abc123
  hcptf stack delete -id=st-old123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack delete command
func (c *StackDeleteCommand) Synopsis() string {
	return "Delete a stack"
}
