package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type RegistryProviderDeleteCommand struct {
	Meta
	organization        string
	name                string
	namespace           string
	registryName        string
	force               bool
	yes                 bool
	registryProviderSvc registryProviderDeleter
}

// Run executes the registry provider delete command
func (c *RegistryProviderDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registryprovider delete")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Provider name (required)")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (defaults to organization)")
	flags.StringVar(&c.registryName, "registry-name", "private", "Registry name: public or private (default: private)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")
	flags.BoolVar(&c.force, "f", false, "Shorthand for -force")
	flags.BoolVar(&c.yes, "y", false, "Confirm delete without prompt")

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

	if !c.force && !c.yes {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete provider '%s/%s'? (yes/no): ", c.namespace, c.name))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}
		if strings.TrimSpace(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete registry provider
	providerID := tfe.RegistryProviderID{
		RegistryName: tfe.RegistryName(c.registryName),
		Namespace:    c.namespace,
		Name:         c.name,
	}

	err = c.registryProviderService(client).Delete(client.Context(), providerID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting registry provider: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Registry provider '%s/%s' deleted successfully", c.namespace, c.name))
	return 0
}

func (c *RegistryProviderDeleteCommand) registryProviderService(client *client.Client) registryProviderDeleter {
	if c.registryProviderSvc != nil {
		return c.registryProviderSvc
	}
	return client.RegistryProviders
}

// Help returns help text for the registry provider delete command
func (c *RegistryProviderDeleteCommand) Help() string {
	helpText := `
	Usage: hcptf registry provider delete [options]

  Delete a private registry provider and all its versions.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Provider name (required)
  -namespace=<name>    Namespace (defaults to organization)
  -registry-name=<val> Registry name: public or private (default: private)
  -force               Force delete without confirmation
  -f                   Shorthand for -force
  -y                   Confirm delete without prompt

Example:

  hcptf registryprovider delete -org=my-org -name=aws
  hcptf registryprovider delete -org=my-org -name=custom -namespace=hashicorp
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry provider delete command
func (c *RegistryProviderDeleteCommand) Synopsis() string {
	return "Delete a private registry provider"
}
