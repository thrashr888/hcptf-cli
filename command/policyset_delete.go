package command

import (
	"fmt"
	"strings"
)

// PolicySetDeleteCommand is a command to delete a policy set
type PolicySetDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the policy set delete command
func (c *PolicySetDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset delete")
	flags.StringVar(&c.id, "id", "", "Policy set ID (required)")
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete policy set '%s'? This action cannot be undone! (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete policy set
	err = client.PolicySets.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting policy set: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Policy set '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the policy set delete command
func (c *PolicySetDeleteCommand) Help() string {
	helpText := `
Usage: hcptf policyset delete [options]

  Delete a policy set. This action cannot be undone!

Options:

  -id=<id>  Policy set ID (required)
  -force    Force delete without confirmation

Example:

  hcptf policyset delete -id=polset-12345
  hcptf policyset delete -id=polset-12345 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set delete command
func (c *PolicySetDeleteCommand) Synopsis() string {
	return "Delete a policy set"
}
