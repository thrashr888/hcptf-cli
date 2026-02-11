package command

import (
	"fmt"
	"strings"
)

// VariableSetVariableDeleteCommand is a command to delete a variable from a variable set
type VariableSetVariableDeleteCommand struct {
	Meta
	variableSetID string
	variableID    string
}

// Run executes the variable set variable delete command
func (c *VariableSetVariableDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset variable delete")
	flags.StringVar(&c.variableSetID, "variableset-id", "", "Variable set ID (required)")
	flags.StringVar(&c.variableID, "variable-id", "", "Variable ID (required)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.variableSetID == "" {
		c.Ui.Error("Error: -variableset-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.variableID == "" {
		c.Ui.Error("Error: -variable-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Delete variable
	err = client.VariableSetVariables.Delete(client.Context(), c.variableSetID, c.variableID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting variable: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Variable '%s' deleted successfully from variable set", c.variableID))
	return 0
}

// Help returns help text for the variable set variable delete command
func (c *VariableSetVariableDeleteCommand) Help() string {
	helpText := `
Usage: hcptf variableset variable delete [options]

  Delete a variable from a variable set.

Options:

  -variableset-id=<id>  Variable set ID (required)
  -variable-id=<id>     Variable ID (required)

Example:

  hcptf variableset variable delete -variableset-id=varset-12345 -variable-id=var-abc123
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable set variable delete command
func (c *VariableSetVariableDeleteCommand) Synopsis() string {
	return "Delete a variable from a variable set"
}
