package command

import (
	"fmt"
	"strings"
)

// AzureoidcDeleteCommand is a command to delete an Azure OIDC configuration
type AzureoidcDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the azureoidc delete command
func (c *AzureoidcDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("azureoidc delete")
	flags.StringVar(&c.id, "id", "", "Azure OIDC configuration ID (required)")
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

	// Read Azure OIDC configuration to get its details for confirmation
	config, err := client.AzureOIDCConfigurations.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading Azure OIDC configuration: %s", err))
		return 1
	}

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete Azure OIDC configuration '%s' (Client ID: %s)? (yes/no): ", c.id, config.ClientID))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete Azure OIDC configuration
	err = client.AzureOIDCConfigurations.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting Azure OIDC configuration: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Azure OIDC configuration '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the azureoidc delete command
func (c *AzureoidcDeleteCommand) Help() string {
	helpText := `
Usage: hcptf azureoidc delete [options]

  Delete an Azure OIDC configuration.

  WARNING: Deleting an Azure OIDC configuration will prevent HCP Terraform
  from generating dynamic Azure credentials using this configuration.
  Ensure no workspaces are actively using this configuration before deletion.

Options:

  -id=<id>  Azure OIDC configuration ID (required)
  -force    Force delete without confirmation

Example:

  hcptf azureoidc delete -id=azoidc-ABC123
  hcptf azureoidc delete -id=azoidc-ABC123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the azureoidc delete command
func (c *AzureoidcDeleteCommand) Synopsis() string {
	return "Delete an Azure OIDC configuration"
}
