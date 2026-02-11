package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type RegistryProviderReadCommand struct {
	Meta
	organization        string
	name                string
	namespace           string
	registryName        string
	format              string
	registryProviderSvc registryProviderReader
}

// Run executes the registry provider read command
func (c *RegistryProviderReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registryprovider read")
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

	// Read registry provider
	providerID := tfe.RegistryProviderID{
		RegistryName: tfe.RegistryName(c.registryName),
		Namespace:    c.namespace,
		Name:         c.name,
	}

	provider, err := c.registryProviderService(client).Read(client.Context(), providerID, nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading registry provider: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	// Show provider details
	data := map[string]interface{}{
		"ID":           provider.ID,
		"Name":         provider.Name,
		"Namespace":    provider.Namespace,
		"RegistryName": string(provider.RegistryName),
		"CreatedAt":    provider.CreatedAt,
		"UpdatedAt":    provider.UpdatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *RegistryProviderReadCommand) registryProviderService(client *client.Client) registryProviderReader {
	if c.registryProviderSvc != nil {
		return c.registryProviderSvc
	}
	return client.RegistryProviders
}

// Help returns help text for the registry provider read command
func (c *RegistryProviderReadCommand) Help() string {
	helpText := `
Usage: hcptf registryprovider read [options]

  Show details of a private registry provider.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Provider name (required)
  -namespace=<name>    Namespace (defaults to organization)
  -registry-name=<val> Registry name: public or private (default: private)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf registryprovider read -org=my-org -name=aws
  hcptf registryprovider read -org=my-org -name=custom -namespace=hashicorp -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry provider read command
func (c *RegistryProviderReadCommand) Synopsis() string {
	return "Show details of a private registry provider"
}
