package command

import (
	"fmt"
	"strings"
)

// HYOKKeyDeleteCommand is a command to revoke a HYOK customer key version
type HYOKKeyDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the HYOK key delete command
func (c *HYOKKeyDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("hyokkey delete")
	flags.StringVar(&c.id, "id", "", "HYOK customer key version ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force revocation without confirmation")

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

	// Confirm revocation unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to revoke HYOK customer key version '%s'? This action cannot be undone! (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Revocation cancelled")
			return 0
		}
	}

	// Revoke HYOK customer key version
	err = client.HYOKCustomerKeyVersions.Revoke(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error revoking HYOK customer key version: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("HYOK customer key version '%s' revocation initiated successfully", c.id))
	return 0
}

// Help returns help text for the HYOK key delete command
func (c *HYOKKeyDeleteCommand) Help() string {
	helpText := `
Usage: hcptf hyokkey delete [options]

  Revoke a HYOK customer key version. This action cannot be undone!

  Revoking a key version prevents it from being used for new encryption
  operations. Existing data encrypted with this key version can still be
  decrypted until the key is removed from your KMS.

  HYOK is an Enterprise-only feature that lets you manage encryption keys
  using your own key management system (AWS KMS, Azure Key Vault, or GCP Cloud KMS).

Options:

  -id=<id>  HYOK customer key version ID (required)
  -force    Force revocation without confirmation

Example:

  hcptf hyokkey delete -id=keyv-123456
  hcptf hyokkey delete -id=keyv-123456 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the HYOK key delete command
func (c *HYOKKeyDeleteCommand) Synopsis() string {
	return "Revoke a HYOK customer key version"
}
