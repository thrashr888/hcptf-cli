package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// AzureoidcReadCommand is a command to read Azure OIDC configuration details
type AzureoidcReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the azureoidc read command
func (c *AzureoidcReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("azureoidc read")
	flags.StringVar(&c.id, "id", "", "Azure OIDC configuration ID (required)")
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

	// Read Azure OIDC configuration
	config, err := client.AzureOIDCConfigurations.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading Azure OIDC configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

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

// Help returns help text for the azureoidc read command
func (c *AzureoidcReadCommand) Help() string {
	helpText := `
Usage: hcptf azureoidc read [options]

  Read Azure OIDC configuration details.

  Displays the configuration details for an Azure OIDC configuration,
  including the application (client) ID, subscription ID, tenant ID,
  and audience settings.

Options:

  -id=<id>          Azure OIDC configuration ID (required)
                    Format: azoidc-XXXXXXXXXX
  -output=<format>  Output format: table (default) or json

Example:

  hcptf azureoidc read -id=azoidc-ABC123
  hcptf azureoidc read -id=azoidc-ABC123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the azureoidc read command
func (c *AzureoidcReadCommand) Synopsis() string {
	return "Read Azure OIDC configuration details"
}
