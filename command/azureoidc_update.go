package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// AzureoidcUpdateCommand is a command to update an Azure OIDC configuration
type AzureoidcUpdateCommand struct {
	Meta
	id             string
	clientID       string
	subscriptionID string
	tenantID       string
	format         string
}

// Run executes the azureoidc update command
func (c *AzureoidcUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("azureoidc update")
	flags.StringVar(&c.id, "id", "", "Azure OIDC configuration ID (required)")
	flags.StringVar(&c.clientID, "client-id", "", "Azure application (client) ID")
	flags.StringVar(&c.subscriptionID, "subscription-id", "", "Azure subscription ID")
	flags.StringVar(&c.tenantID, "tenant-id", "", "Azure tenant (directory) ID")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

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

	// Build update options
	options := tfe.AzureOIDCConfigurationUpdateOptions{}

	if c.clientID != "" {
		options.ClientID = &c.clientID
	}

	if c.subscriptionID != "" {
		options.SubscriptionID = &c.subscriptionID
	}

	if c.tenantID != "" {
		options.TenantID = &c.tenantID
	}

	// Update Azure OIDC configuration
	config, err := client.AzureOIDCConfigurations.Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating Azure OIDC configuration: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Azure OIDC configuration '%s' updated successfully", config.ID))

	// Show configuration details
	data := map[string]interface{}{
		"ID":             config.ID,
		"ClientID":       config.ClientID,
		"SubscriptionID": config.SubscriptionID,
		"TenantID":       config.TenantID,
	}

	if config.Organization != nil {
		data["Organization"] = config.Organization.Name
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the azureoidc update command
func (c *AzureoidcUpdateCommand) Help() string {
	helpText := `
Usage: hcptf azureoidc update [options]

  Update Azure OIDC configuration settings.

  Updates the Azure application (client) ID, subscription ID, tenant ID,
  or audience settings for an existing Azure OIDC configuration.

Options:

  -id=<id>                 Azure OIDC configuration ID (required)
  -client-id=<id>          Azure application (client) ID
  -subscription-id=<id>    Azure subscription ID
  -tenant-id=<id>          Azure tenant (directory) ID
  -output=<format>         Output format: table (default) or json

Example:

  hcptf azureoidc update -id=azoidc-ABC123 \
    -client-id=12345678-1234-1234-1234-123456789012

  hcptf azureoidc update -id=azoidc-ABC123 \
    -subscription-id=87654321-4321-4321-4321-210987654321 \
    -tenant-id=abcdefab-abcd-abcd-abcd-abcdefabcdef
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the azureoidc update command
func (c *AzureoidcUpdateCommand) Synopsis() string {
	return "Update Azure OIDC configuration settings"
}
