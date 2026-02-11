package command

import (
	"fmt"
	"strings"
)

// VariableSetDeleteCommand is a command to delete a variable set
type VariableSetDeleteCommand struct {
	Meta
	id string
}

// Run executes the variable set delete command
func (c *VariableSetDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset delete")
	flags.StringVar(&c.id, "id", "", "Variable set ID (required)")

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

	// Delete variable set
	err = client.VariableSets.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting variable set: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Variable set '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the variable set delete command
func (c *VariableSetDeleteCommand) Help() string {
	helpText := `
Usage: hcptf variableset delete [options]

  Delete a variable set.

Options:

  -id=<id>  Variable set ID (required)

Example:

  hcptf variableset delete -id=varset-12345
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable set delete command
func (c *VariableSetDeleteCommand) Synopsis() string {
	return "Delete a variable set"
}
