package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// AzureoidcCreateCommand is a command to create an Azure OIDC configuration
type AzureoidcCreateCommand struct {
	Meta
	organization   string
	clientID       string
	subscriptionID string
	tenantID       string
	format         string
}

// Run executes the azureoidc create command
func (c *AzureoidcCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("azureoidc create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.clientID, "client-id", "", "Azure application (client) ID (required)")
	flags.StringVar(&c.subscriptionID, "subscription-id", "", "Azure subscription ID (required)")
	flags.StringVar(&c.tenantID, "tenant-id", "", "Azure tenant (directory) ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.clientID == "" {
		c.Ui.Error("Error: -client-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.subscriptionID == "" {
		c.Ui.Error("Error: -subscription-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.tenantID == "" {
		c.Ui.Error("Error: -tenant-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build create options
	options := tfe.AzureOIDCConfigurationCreateOptions{
		ClientID:       c.clientID,
		SubscriptionID: c.subscriptionID,
		TenantID:       c.tenantID,
	}

	// Create Azure OIDC configuration
	config, err := client.AzureOIDCConfigurations.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating Azure OIDC configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Azure OIDC configuration created successfully with ID: %s", config.ID))

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

// Help returns help text for the azureoidc create command
func (c *AzureoidcCreateCommand) Help() string {
	helpText := `
Usage: hcptf azureoidc create [options]

  Create a new Azure OIDC configuration for dynamic Azure credentials.

  Azure OIDC configurations enable HCP Terraform to dynamically generate
  Azure credentials using OpenID Connect. This eliminates the need to store
  static Azure service principal credentials in HCP Terraform.

  Prerequisites:
  - Azure Entra ID application with federated credentials configured for HCP Terraform
  - Application must have appropriate permissions for required Azure operations

Options:

  -organization=<name>      Organization name (required)
  -org=<name>              Alias for -organization
  -client-id=<id>          Azure application (client) ID (required)
  -subscription-id=<id>    Azure subscription ID (required)
  -tenant-id=<id>          Azure tenant (directory) ID (required)
  -output=<format>         Output format: table (default) or json

Example:

  hcptf azureoidc create -org=my-org \
    -client-id=12345678-1234-1234-1234-123456789012 \
    -subscription-id=87654321-4321-4321-4321-210987654321 \
    -tenant-id=abcdefab-abcd-abcd-abcd-abcdefabcdef
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the azureoidc create command
func (c *AzureoidcCreateCommand) Synopsis() string {
	return "Create an Azure OIDC configuration for dynamic credentials"
}
