package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

type TeamAccessDeleteCommand struct {
	Meta
	id            string
	force         bool
	teamAccessSvc teamAccessDeleter
}

// Run executes the team access delete command
func (c *TeamAccessDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("teamaccess delete")
	flags.StringVar(&c.id, "id", "", "Team access ID (required)")
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to remove team access '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete team access
	err = c.teamAccessService(client).Remove(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error removing team access: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Team access '%s' removed successfully", c.id))
	return 0
}

func (c *TeamAccessDeleteCommand) teamAccessService(client *client.Client) teamAccessDeleter {
	if c.teamAccessSvc != nil {
		return c.teamAccessSvc
	}
	return client.TeamAccess
}

// Help returns help text for the team access delete command
func (c *TeamAccessDeleteCommand) Help() string {
	helpText := `
Usage: hcptf teamaccess delete [options]

  Remove team access from a workspace.

Options:

  -id=<id>    Team access ID (required)
  -force      Force delete without confirmation

Example:

  hcptf teamaccess delete -id=tws-123abc
  hcptf teamaccess delete -id=tws-123abc -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team access delete command
func (c *TeamAccessDeleteCommand) Synopsis() string {
	return "Remove team access from a workspace"
}
