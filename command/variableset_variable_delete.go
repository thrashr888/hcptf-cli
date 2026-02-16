package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

type variableSetVariableDeleter interface {
	Delete(ctx context.Context, variableSetID, variableID string) error
}

// VariableSetVariableDeleteCommand is a command to delete a variable from a variable set
type VariableSetVariableDeleteCommand struct {
	Meta
	variableSetID          string
	variableID             string
	force                  bool
	yes                    bool
	variableSetVariableSvc variableSetVariableDeleter
}

// Run executes the variable set variable delete command
func (c *VariableSetVariableDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset variable delete")
	flags.StringVar(&c.variableSetID, "variableset-id", "", "Variable set ID (required)")
	flags.StringVar(&c.variableID, "variable-id", "", "Variable ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")
	flags.BoolVar(&c.force, "f", false, "Shorthand for -force")
	flags.BoolVar(&c.yes, "y", false, "Confirm delete without prompt")

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

	if !c.force && !c.yes {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete variable '%s' from variable set '%s'? (yes/no): ", c.variableID, c.variableSetID))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}
		if strings.TrimSpace(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete variable
	err = c.variableSetVariableService(client).Delete(client.Context(), c.variableSetID, c.variableID)
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
  -force                Force delete without confirmation
  -f                    Shorthand for -force
  -y                    Confirm delete without prompt

Example:

  hcptf variableset variable delete -variableset-id=varset-12345 -variable-id=var-abc123
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable set variable delete command
func (c *VariableSetVariableDeleteCommand) Synopsis() string {
	return "Delete a variable from a variable set"
}

func (c *VariableSetVariableDeleteCommand) variableSetVariableService(client *client.Client) variableSetVariableDeleter {
	if c.variableSetVariableSvc != nil {
		return c.variableSetVariableSvc
	}
	return client.VariableSetVariables
}
