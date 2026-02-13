package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type RegistryProviderCreateCommand struct {
	Meta
	organization        string
	name                string
	namespace           string
	registryName        string
	format              string
	registryProviderSvc registryProviderCreator
}

// Run executes the registry provider create command
func (c *RegistryProviderCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registryprovider create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Provider name (required)")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (defaults to organization)")
	flags.StringVar(&c.registryName, "registry-name", "private", "Registry name: public or private (default: private)")
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

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Default namespace to organization for private providers
	if c.namespace == "" {
		c.namespace = c.organization
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create registry provider
	options := tfe.RegistryProviderCreateOptions{
		Name:         c.name,
		Namespace:    c.namespace,
		RegistryName: tfe.RegistryName(c.registryName),
	}

	provider, err := c.registryProviderService(client).Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating registry provider: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if c.format != "json" {
		c.Ui.Output(fmt.Sprintf("Registry provider '%s/%s' created successfully", provider.Namespace, provider.Name))
	}

	if c.registryName == "private" {
		if c.format != "json" {
			c.Ui.Output("\nNext steps:")
			c.Ui.Output("1. Create a provider version using 'registryproviderversion create'")
			c.Ui.Output("2. Upload shasums and signature files")
			c.Ui.Output("3. Create platform binaries using 'registryproviderplatform create'")
		}
	}

	// Show provider details
	data := map[string]interface{}{
		"ID":           provider.ID,
		"Name":         provider.Name,
		"Namespace":    provider.Namespace,
		"RegistryName": string(provider.RegistryName),
		"CreatedAt":    provider.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *RegistryProviderCreateCommand) registryProviderService(client *client.Client) registryProviderCreator {
	if c.registryProviderSvc != nil {
		return c.registryProviderSvc
	}
	return client.RegistryProviders
}

// Help returns help text for the registry provider create command
func (c *RegistryProviderCreateCommand) Help() string {
	helpText := `
Usage: hcptf registryprovider create [options]

  Create a new private registry provider.

  After creating a provider, you must:
  1. Create a version with 'registryproviderversion create'
  2. Upload SHA256SUMS and SHA256SUMS.sig files
  3. Create platform binaries with 'registryproviderplatform create'

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Provider name (required)
  -namespace=<name>    Namespace (defaults to organization)
  -registry-name=<val> Registry name: public or private (default: private)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf registryprovider create -org=my-org -name=aws
  hcptf registryprovider create -org=my-org -name=custom -namespace=hashicorp
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry provider create command
func (c *RegistryProviderCreateCommand) Synopsis() string {
	return "Create a new private registry provider"
}
