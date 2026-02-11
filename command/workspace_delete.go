package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// WorkspaceDeleteCommand is a command to delete a workspace
type WorkspaceDeleteCommand struct {
	Meta
	organization string
	name         string
	force        bool
	workspaceSvc workspaceDeleter
}

// Run executes the workspace delete command
func (c *WorkspaceDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspace delete")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Workspace name (required)")
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

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete workspace '%s'? (yes/no): ", c.name))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete workspace
	err = c.workspaceService(client).Delete(client.Context(), c.organization, c.name)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting workspace: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Workspace '%s' deleted successfully", c.name))
	return 0
}

func (c *WorkspaceDeleteCommand) workspaceService(client *client.Client) workspaceDeleter {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

// Help returns help text for the workspace delete command
func (c *WorkspaceDeleteCommand) Help() string {
	helpText := `
Usage: hcptf workspace delete [options]

  Delete a workspace.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Workspace name (required)
  -force               Force delete without confirmation

Example:

  hcptf workspace delete -org=my-org -name=my-workspace
  hcptf workspace delete -org=my-org -name=old-workspace -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace delete command
func (c *WorkspaceDeleteCommand) Synopsis() string {
	return "Delete a workspace"
}
