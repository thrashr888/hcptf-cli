package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

type ProjectTeamAccessDeleteCommand struct {
	Meta
	id                     string
	force                  bool
	projectTeamAccessSvc   projectTeamAccessDeleter
}

// Run executes the project team access delete command
func (c *ProjectTeamAccessDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("projectteamaccess delete")
	flags.StringVar(&c.id, "id", "", "Project team access ID (required)")
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to remove project team access '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete project team access
	err = c.projectTeamAccessService(client).Remove(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error removing project team access: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Project team access '%s' removed successfully", c.id))
	return 0
}

func (c *ProjectTeamAccessDeleteCommand) projectTeamAccessService(client *client.Client) projectTeamAccessDeleter {
	if c.projectTeamAccessSvc != nil {
		return c.projectTeamAccessSvc
	}
	return client.TeamProjectAccess
}

// Help returns help text for the project team access delete command
func (c *ProjectTeamAccessDeleteCommand) Help() string {
	helpText := `
Usage: hcptf projectteamaccess delete [options]

  Remove team access from a project.

Options:

  -id=<id>  Project team access ID (required)
  -force    Force delete without confirmation

Example:

  hcptf projectteamaccess delete -id=tprj-123abc
  hcptf projectteamaccess delete -id=tprj-123abc -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the project team access delete command
func (c *ProjectTeamAccessDeleteCommand) Synopsis() string {
	return "Remove team access from a project"
}
