package command

import (
	"fmt"
	"strings"
)

// VaultoidcDeleteCommand is a command to delete a Vault OIDC configuration
type VaultoidcDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the vaultoidc delete command
func (c *VaultoidcDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("vaultoidc delete")
	flags.StringVar(&c.id, "id", "", "Vault OIDC configuration ID (required)")
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

	// Read Vault OIDC configuration to get its details for confirmation
	config, err := client.VaultOIDCConfigurations.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading Vault OIDC configuration: %s", err))
		return 1
	}

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete Vault OIDC configuration '%s' (Address: %s)? (yes/no): ", c.id, config.Address))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete Vault OIDC configuration
	err = client.VaultOIDCConfigurations.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting Vault OIDC configuration: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Vault OIDC configuration '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the vaultoidc delete command
func (c *VaultoidcDeleteCommand) Help() string {
	helpText := `
Usage: hcptf vaultoidc delete [options]

  Delete a Vault OIDC configuration.

  WARNING: Deleting a Vault OIDC configuration will prevent HCP Terraform
  from generating dynamic Vault tokens using this configuration.
  Ensure no workspaces are actively using this configuration before deletion.

Options:

  -id=<id>  Vault OIDC configuration ID (required)
  -force    Force delete without confirmation

Example:

  hcptf vaultoidc delete -id=voidc-ABC123
  hcptf vaultoidc delete -id=voidc-ABC123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the vaultoidc delete command
func (c *VaultoidcDeleteCommand) Synopsis() string {
	return "Delete a Vault OIDC configuration"
}
