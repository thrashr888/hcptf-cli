package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// VariableSetDeleteCommand is a command to delete a variable set
type VariableSetDeleteCommand struct {
	Meta
	id             string
	force          bool
	yes            bool
	variableSetSvc variableSetDeleter
}

// Run executes the variable set delete command
func (c *VariableSetDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset delete")
	flags.StringVar(&c.id, "id", "", "Variable set ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")
	flags.BoolVar(&c.force, "f", false, "Shorthand for -force")
	flags.BoolVar(&c.yes, "y", false, "Confirm delete without prompt")

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

	if !c.force && !c.yes {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete variable set '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}
		if strings.TrimSpace(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete variable set
	err = c.variableSetService(client).Delete(client.Context(), c.id)
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
  -force     Force delete without confirmation
  -f         Shorthand for -force
  -y         Confirm delete without prompt

Example:

  hcptf variableset delete -id=varset-12345
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable set delete command
func (c *VariableSetDeleteCommand) Synopsis() string {
	return "Delete a variable set"
}

func (c *VariableSetDeleteCommand) variableSetService(client *client.Client) variableSetDeleter {
	if c.variableSetSvc != nil {
		return c.variableSetSvc
	}
	return client.VariableSets
}
