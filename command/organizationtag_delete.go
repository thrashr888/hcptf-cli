package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// OrganizationTagDeleteCommand is a command to delete organization tags
type OrganizationTagDeleteCommand struct {
	Meta
	organization string
	id           string
	force        bool
}

// Run executes the organizationtag delete command
func (c *OrganizationTagDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organizationtag delete")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.id, "id", "", "Tag ID (required)")
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete tag '%s'? This will remove it from all workspaces. (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete organization tag
	options := tfe.OrganizationTagsDeleteOptions{
		IDs: []string{c.id},
	}

	err = client.OrganizationTags.Delete(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting organization tag: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Organization tag '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the organizationtag delete command
func (c *OrganizationTagDeleteCommand) Help() string {
	helpText := `
Usage: hcptf organizationtag delete [options]

  Delete an organization tag. This will remove the tag from
  all workspaces in the organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -id=<id>             Tag ID (required)
  -force               Force delete without confirmation

Example:

  hcptf organizationtag delete -org=my-org -id=tag-ABC123
  hcptf organizationtag delete -org=my-org -id=tag-ABC123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organizationtag delete command
func (c *OrganizationTagDeleteCommand) Synopsis() string {
	return "Delete an organization tag"
}
