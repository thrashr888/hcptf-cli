package command

import (
	"fmt"
	"strings"
)

// OAuthClientDeleteCommand is a command to delete an OAuth client
type OAuthClientDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the oauthclient delete command
func (c *OAuthClientDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("oauthclient delete")
	flags.StringVar(&c.id, "id", "", "OAuth client ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")

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

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete OAuth client '%s'? This will unlink all workspaces using this connection. (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete OAuth client
	err = client.OAuthClients.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting OAuth client: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("OAuth client '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the oauthclient delete command
func (c *OAuthClientDeleteCommand) Help() string {
	helpText := `
Usage: hcptf oauthclient delete [options]

  Delete an OAuth client. This removes the VCS connection and will
  unlink all workspaces that use this OAuth client from their
  repositories.

  WARNING: Workspaces using this OAuth client will need to be
  manually linked to another VCS connection.

Options:

  -id=<id>  OAuth client ID (required)
  -force    Force delete without confirmation

Example:

  hcptf oauthclient delete -id=oc-XKFwG6ggfA9n7t1K
  hcptf oauthclient delete -id=oc-XKFwG6ggfA9n7t1K -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the oauthclient delete command
func (c *OAuthClientDeleteCommand) Synopsis() string {
	return "Delete an OAuth client"
}
