package command

import (
	"fmt"
	"strings"
)

// SSHKeyDeleteCommand is a command to delete an SSH key
type SSHKeyDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the SSH key delete command
func (c *SSHKeyDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("sshkey delete")
	flags.StringVar(&c.id, "id", "", "SSH key ID (required)")
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete SSH key '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete SSH key
	err = client.SSHKeys.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting SSH key: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("SSH key '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the SSH key delete command
func (c *SSHKeyDeleteCommand) Help() string {
	helpText := `
Usage: hcptf sshkey delete [options]

  Delete an SSH key.

Options:

  -id=<id>             SSH key ID (required)
  -force               Force delete without confirmation

Example:

  hcptf sshkey delete -id=sshkey-123abc
  hcptf sshkey delete -id=sshkey-456def -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the SSH key delete command
func (c *SSHKeyDeleteCommand) Synopsis() string {
	return "Delete an SSH key"
}
