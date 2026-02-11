package command

import (
	"fmt"
	"strings"
)

// OrganizationDeleteCommand is a command to delete an organization
type OrganizationDeleteCommand struct {
	Meta
	name  string
	force bool
}

// Run executes the organization delete command
func (c *OrganizationDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organization delete")
	flags.StringVar(&c.name, "name", "", "Organization name (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete organization '%s'? This action cannot be undone! (yes/no): ", c.name))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete organization
	err = client.Organizations.Delete(client.Context(), c.name)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting organization: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Organization '%s' deleted successfully", c.name))
	return 0
}

// Help returns help text for the organization delete command
func (c *OrganizationDeleteCommand) Help() string {
	helpText := `
Usage: hcptf organization delete [options]

  Delete an organization. This action cannot be undone!

Options:

  -name=<name>  Organization name (required)
  -force        Force delete without confirmation

Example:

  hcptf organization delete -name=my-org
  hcptf organization delete -name=old-org -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization delete command
func (c *OrganizationDeleteCommand) Synopsis() string {
	return "Delete an organization"
}
