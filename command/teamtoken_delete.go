package command

import (
	"fmt"
	"strings"
)

// TeamTokenDeleteCommand is a command to delete a team token
type TeamTokenDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the team token delete command
func (c *TeamTokenDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("teamtoken delete")
	flags.StringVar(&c.id, "id", "", "Team token ID (required)")
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete team token '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete team token
	err = client.TeamTokens.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting team token: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Team token '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the team token delete command
func (c *TeamTokenDeleteCommand) Help() string {
	helpText := `
Usage: hcptf teamtoken delete [options]

  Delete a team API token.

  This will permanently delete the specified team token. Any applications
  or scripts using this token will no longer be able to authenticate.

Options:

  -id=<id>  Team token ID (required)
  -force    Force delete without confirmation

Example:

  hcptf teamtoken delete -id=at-abc123xyz
  hcptf teamtoken delete -id=at-abc123xyz -force

Security Note:

  Deleting a team token will immediately invalidate it. Any applications
  using this token will need to be updated with a new token.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team token delete command
func (c *TeamTokenDeleteCommand) Synopsis() string {
	return "Delete a team API token"
}
