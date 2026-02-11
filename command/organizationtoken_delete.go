package command

import (
	"fmt"
	"strings"
)

// OrganizationTokenDeleteCommand is a command to delete an organization token
type OrganizationTokenDeleteCommand struct {
	Meta
	organization string
	force        bool
}

// Run executes the organization token delete command
func (c *OrganizationTokenDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organizationtoken delete")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete the organization token for '%s'? (yes/no): ", c.organization))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete organization token
	err = client.OrganizationTokens.Delete(client.Context(), c.organization)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting organization token: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Organization token for '%s' deleted successfully", c.organization))
	return 0
}

// Help returns help text for the organization token delete command
func (c *OrganizationTokenDeleteCommand) Help() string {
	helpText := `
Usage: hcptf organizationtoken delete [options]

  Delete an organization token (organization-level API token).

  This will permanently delete the organization's API token. Any applications
  or scripts using this token will no longer be able to authenticate.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -force               Force delete without confirmation

Example:

  hcptf organizationtoken delete -org=my-org
  hcptf organizationtoken delete -organization=my-org -force

Security Note:

  Deleting the organization token will immediately invalidate it. Any
  applications using this token will need to be updated with a new token.
  Only members of the owners team can delete organization tokens.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization token delete command
func (c *OrganizationTokenDeleteCommand) Synopsis() string {
	return "Delete an organization token"
}
