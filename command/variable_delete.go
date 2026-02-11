package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// VariableDeleteCommand is a command to delete a variable
type VariableDeleteCommand struct {
	Meta
	organization string
	workspace    string
	id           string
	force        bool
	workspaceSvc workspaceReader
	variableSvc  variableDeleter
}

// Run executes the variable delete command
func (c *VariableDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variable delete")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
	flags.StringVar(&c.id, "id", "", "Variable ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.workspace == "" {
		c.Ui.Error("Error: -workspace flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

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

	// Get workspace first
	ws, err := c.workspaceService(client).Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete variable '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete variable
	err = c.variableService(client).Delete(client.Context(), ws.ID, c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting variable: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Variable '%s' deleted successfully", c.id))
	return 0
}

func (c *VariableDeleteCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

func (c *VariableDeleteCommand) variableService(client *client.Client) variableDeleter {
	if c.variableSvc != nil {
		return c.variableSvc
	}
	return client.Variables
}

// Help returns help text for the variable delete command
func (c *VariableDeleteCommand) Help() string {
	helpText := `
Usage: hcptf variable delete [options]

  Delete a variable.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -id=<id>             Variable ID (required)
  -force               Force delete without confirmation

Example:

  hcptf variable delete -org=my-org -workspace=prod -id=var-123
  hcptf variable delete -org=my-org -workspace=prod -id=var-456 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable delete command
func (c *VariableDeleteCommand) Synopsis() string {
	return "Delete a variable"
}
