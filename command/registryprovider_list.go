package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type RegistryProviderListCommand struct {
	Meta
	organization        string
	format              string
	registryProviderSvc registryProviderLister
}

// Run executes the registry provider list command
func (c *RegistryProviderListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registryprovider list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List registry providers
	providers, err := c.registryProviderService(client).List(client.Context(), c.organization, &tfe.RegistryProviderListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing registry providers: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(providers.Items) == 0 {
		c.Ui.Output("No registry providers found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Namespace", "Registry", "Created At"}
	var rows [][]string

	for _, prov := range providers.Items {
		rows = append(rows, []string{
			prov.ID,
			prov.Name,
			prov.Namespace,
			string(prov.RegistryName),
			prov.CreatedAt,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *RegistryProviderListCommand) registryProviderService(client *client.Client) registryProviderLister {
	if c.registryProviderSvc != nil {
		return c.registryProviderSvc
	}
	return client.RegistryProviders
}

// Help returns help text for the registry provider list command
func (c *RegistryProviderListCommand) Help() string {
	helpText := `
Usage: hcptf registryprovider list [options]

  List private registry providers in an organization.

  The private registry enables organizations to publish and share
  custom Terraform providers within their organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf registryprovider list -organization=my-org
  hcptf registryprovider list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry provider list command
func (c *RegistryProviderListCommand) Synopsis() string {
	return "List private registry providers in an organization"
}
