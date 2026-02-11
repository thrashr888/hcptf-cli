package command

import (
	"fmt"
	"strings"
)

// UserTokenDeleteCommand is a command to delete a user token
type UserTokenDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the user token delete command
func (c *UserTokenDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("usertoken delete")
	flags.StringVar(&c.id, "id", "", "User token ID (required)")
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete user token '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete user token
	err = client.UserTokens.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting user token: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("User token '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the user token delete command
func (c *UserTokenDeleteCommand) Help() string {
	helpText := `
Usage: hcptf usertoken delete [options]

  Delete a user API token.

  This will permanently delete the specified user token. Any applications
  or scripts using this token will no longer be able to authenticate.

Options:

  -id=<id>  User token ID (required)
  -force    Force delete without confirmation

Example:

  hcptf usertoken delete -id=at-abc123xyz
  hcptf usertoken delete -id=at-abc123xyz -force

Security Note:

  Deleting a user token will immediately invalidate it. Any applications
  using this token will need to be updated with a new token. This command
  can only delete tokens for the currently authenticated user.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the user token delete command
func (c *UserTokenDeleteCommand) Synopsis() string {
	return "Delete a user API token"
}
