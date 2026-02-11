package command

import (
	"fmt"
	"strings"
)

// OrganizationMembershipDeleteCommand is a command to delete an organization membership
type OrganizationMembershipDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the organization membership delete command
func (c *OrganizationMembershipDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organizationmembership delete")
	flags.StringVar(&c.id, "id", "", "Organization membership ID (required)")
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to remove organization membership '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete organization membership
	err = client.OrganizationMemberships.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting organization membership: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Organization membership '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the organization membership delete command
func (c *OrganizationMembershipDeleteCommand) Help() string {
	helpText := `
Usage: hcptf organizationmembership delete [options]

  Remove a user from an organization.

  This command removes a user's membership from an organization. The user
  will lose access to all workspaces and resources in the organization.

Options:

  -id=<id>  Organization membership ID (required)
  -force    Force delete without confirmation

Example:

  hcptf organizationmembership delete -id=ou-abc123xyz
  hcptf organizationmembership delete -id=ou-abc123xyz -force

Security Note:

  Removing a user from an organization will immediately revoke their access
  to all resources in the organization. You cannot remove yourself from an
  organization that you own.

Note:

  Only members with team management permissions or owners can remove users
  from an organization.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization membership delete command
func (c *OrganizationMembershipDeleteCommand) Synopsis() string {
	return "Remove a user from an organization"
}
