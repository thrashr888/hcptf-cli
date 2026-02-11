package command

import (
	"fmt"
	"strings"
)

// AWSoidcDeleteCommand is a command to delete an AWS OIDC configuration
type AWSoidcDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the awsoidc delete command
func (c *AWSoidcDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("awsoidc delete")
	flags.StringVar(&c.id, "id", "", "AWS OIDC configuration ID (required)")
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

	// Read AWS OIDC configuration to get its details for confirmation
	config, err := client.AWSOIDCConfigurations.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading AWS OIDC configuration: %s", err))
		return 1
	}

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete AWS OIDC configuration '%s' (Role ARN: %s)? (yes/no): ", c.id, config.RoleARN))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete AWS OIDC configuration
	err = client.AWSOIDCConfigurations.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting AWS OIDC configuration: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("AWS OIDC configuration '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the awsoidc delete command
func (c *AWSoidcDeleteCommand) Help() string {
	helpText := `
Usage: hcptf awsoidc delete [options]

  Delete an AWS OIDC configuration.

  WARNING: Deleting an AWS OIDC configuration will prevent HCP Terraform
  from generating dynamic AWS credentials using this configuration.
  Ensure no workspaces are actively using this configuration before deletion.

Options:

  -id=<id>  AWS OIDC configuration ID (required)
  -force    Force delete without confirmation

Example:

  hcptf awsoidc delete -id=awsoidc-ABC123
  hcptf awsoidc delete -id=awsoidc-ABC123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the awsoidc delete command
func (c *AWSoidcDeleteCommand) Synopsis() string {
	return "Delete an AWS OIDC configuration"
}
