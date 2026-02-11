package command

import (
	"fmt"
	"strings"
)

// HYOKDeleteCommand is a command to delete a HYOK configuration
type HYOKDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the HYOK delete command
func (c *HYOKDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("hyok delete")
	flags.StringVar(&c.id, "id", "", "HYOK configuration ID (required)")
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete HYOK configuration '%s'? This action cannot be undone! (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete HYOK configuration
	err = client.HYOKConfigurations.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting HYOK configuration: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("HYOK configuration '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the HYOK delete command
func (c *HYOKDeleteCommand) Help() string {
	helpText := `
Usage: hcptf hyok delete [options]

  Delete a HYOK (Hold Your Own Key) configuration. This action cannot be undone!

  HYOK is an Enterprise-only feature that lets you manage encryption keys
  using your own key management system (AWS KMS, Azure Key Vault, or GCP Cloud KMS).

Options:

  -id=<id>  HYOK configuration ID (required)
  -force    Force delete without confirmation

Example:

  hcptf hyok delete -id=hyokc-123456
  hcptf hyok delete -id=hyokc-123456 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the HYOK delete command
func (c *HYOKDeleteCommand) Synopsis() string {
	return "Delete a HYOK configuration"
}
