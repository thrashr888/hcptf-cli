package command

import (
	"fmt"
	"strings"
)

// ReservedTagKeyDeleteCommand is a command to delete a reserved tag key
type ReservedTagKeyDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the reservedtagkey delete command
func (c *ReservedTagKeyDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("reservedtagkey delete")
	flags.StringVar(&c.id, "id", "", "Reserved tag key ID (required)")
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete reserved tag key '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete reserved tag key
	err = client.ReservedTagKeys.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting reserved tag key: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Reserved tag key '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the reservedtagkey delete command
func (c *ReservedTagKeyDeleteCommand) Help() string {
	helpText := `
Usage: hcptf reservedtagkey delete [options]

  Delete a reserved tag key from an organization.

Options:

  -id=<id>  Reserved tag key ID (required)
  -force    Force delete without confirmation

Example:

  hcptf reservedtagkey delete -id=rtk-ABC123
  hcptf reservedtagkey delete -id=rtk-ABC123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the reservedtagkey delete command
func (c *ReservedTagKeyDeleteCommand) Synopsis() string {
	return "Delete a reserved tag key"
}
