package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type RegistryModuleReadCommand struct {
	Meta
	organization      string
	name              string
	provider          string
	namespace         string
	registryName      string
	format            string
	registryModuleSvc registryModuleReader
}

// Run executes the registry module read command
func (c *RegistryModuleReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registrymodule read")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Module name (required)")
	flags.StringVar(&c.provider, "provider", "", "Provider name (required)")
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

	if c.provider == "" {
		c.Ui.Error("Error: -provider flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Default namespace to organization for private modules
	if c.namespace == "" {
		c.namespace = c.organization
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read registry module
	moduleID := tfe.RegistryModuleID{
		Organization: c.organization,
		Name:         c.name,
		Provider:     c.provider,
		Namespace:    c.namespace,
		RegistryName: tfe.RegistryName(c.registryName),
	}

	module, err := c.registryModuleService(client).Read(client.Context(), moduleID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading registry module: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	// Show module details
	data := map[string]interface{}{
		"ID":           module.ID,
		"Name":         module.Name,
		"Namespace":    module.Namespace,
		"Provider":     module.Provider,
		"RegistryName": string(module.RegistryName),
		"Status":       module.Status,
		"CreatedAt":    module.CreatedAt,
		"UpdatedAt":    module.UpdatedAt,
	}

	// Add VCS repo info if available
	if module.VCSRepo != nil {
		data["VCSIdentifier"] = module.VCSRepo.Identifier
		data["VCSBranch"] = module.VCSRepo.Branch
	}

	formatter.KeyValue(data)
	return 0
}

func (c *RegistryModuleReadCommand) registryModuleService(client *client.Client) registryModuleReader {
	if c.registryModuleSvc != nil {
		return c.registryModuleSvc
	}
	return client.RegistryModules
}

// Help returns help text for the registry module read command
func (c *RegistryModuleReadCommand) Help() string {
	helpText := `
Usage: hcptf registrymodule read [options]

  Show details of a private registry module.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Module name (required)
  -provider=<name>     Provider name (required)
  -namespace=<name>    Namespace (defaults to organization)
  -registry-name=<val> Registry name: public or private (default: private)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf registrymodule read -org=my-org -name=vpc -provider=aws
  hcptf registrymodule read -org=my-org -name=vpc -provider=aws -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry module read command
func (c *RegistryModuleReadCommand) Synopsis() string {
	return "Show details of a private registry module"
}
