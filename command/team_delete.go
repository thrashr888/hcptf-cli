package command

import (
	"fmt"
	"strings"
)

// TeamDeleteCommand is a command to delete a team
type TeamDeleteCommand struct {
	Meta
	organization string
	name         string
	force        bool
}

// Run executes the team delete command
func (c *TeamDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("team delete")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Team name (required)")
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

	// Read team first to get ID
	team, err := client.Teams.Read(client.Context(), c.name)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading team: %s", err))
		return 1
	}

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete team '%s'? (yes/no): ", c.name))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete team
	err = client.Teams.Delete(client.Context(), team.ID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting team: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Team '%s' deleted successfully", c.name))
	return 0
}

// Help returns help text for the team delete command
func (c *TeamDeleteCommand) Help() string {
	helpText := `
Usage: hcptf team delete [options]

  Delete a team.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Team name (required)
  -force               Force delete without confirmation

Example:

  hcptf team delete -org=my-org -name=old-team
  hcptf team delete -org=my-org -name=deprecated -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team delete command
func (c *TeamDeleteCommand) Synopsis() string {
	return "Delete a team"
}
