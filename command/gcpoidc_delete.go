package command

import (
	"fmt"
	"strings"
)

// GCPoidcDeleteCommand is a command to delete a GCP OIDC configuration
type GCPoidcDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the gcpoidc delete command
func (c *GCPoidcDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("gcpoidc delete")
	flags.StringVar(&c.id, "id", "", "GCP OIDC configuration ID (required)")
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

	// Read GCP OIDC configuration to get its details for confirmation
	config, err := client.GCPOIDCConfigurations.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading GCP OIDC configuration: %s", err))
		return 1
	}

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete GCP OIDC configuration '%s' (Service Account: %s)? (yes/no): ", c.id, config.ServiceAccountEmail))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete GCP OIDC configuration
	err = client.GCPOIDCConfigurations.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting GCP OIDC configuration: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("GCP OIDC configuration '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the gcpoidc delete command
func (c *GCPoidcDeleteCommand) Help() string {
	helpText := `
Usage: hcptf gcpoidc delete [options]

  Delete a GCP OIDC configuration.

  WARNING: Deleting a GCP OIDC configuration will prevent HCP Terraform
  from generating dynamic GCP credentials using this configuration.
  Ensure no workspaces are actively using this configuration before deletion.

Options:

  -id=<id>  GCP OIDC configuration ID (required)
  -force    Force delete without confirmation

Example:

  hcptf gcpoidc delete -id=gcpoidc-ABC123
  hcptf gcpoidc delete -id=gcpoidc-ABC123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the gcpoidc delete command
func (c *GCPoidcDeleteCommand) Synopsis() string {
	return "Delete a GCP OIDC configuration"
}
