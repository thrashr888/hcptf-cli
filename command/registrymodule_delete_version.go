package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type RegistryModuleDeleteVersionCommand struct {
	Meta
	organization      string
	name              string
	provider          string
	namespace         string
	registryName      string
	version           string
	registryModuleSvc registryModuleVersionDeleter
}

// Run executes the registry module delete-version command
func (c *RegistryModuleDeleteVersionCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registrymodule delete-version")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Module name (required)")
	flags.StringVar(&c.provider, "provider", "", "Provider name (required)")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (defaults to organization)")
	flags.StringVar(&c.registryName, "registry-name", "private", "Registry name: must be 'private'")
	flags.StringVar(&c.version, "version", "", "Version string (required)")

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

	if c.version == "" {
		c.Ui.Error("Error: -version flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Default namespace to organization
	if c.namespace == "" {
		c.namespace = c.organization
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Delete module version
	moduleID := tfe.RegistryModuleID{
		RegistryName: tfe.RegistryName(c.registryName),
		Namespace:    c.namespace,
		Name:         c.name,
		Provider:     c.provider,
	}

	err = c.registryModuleService(client).DeleteVersion(client.Context(), moduleID, c.version)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting module version: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Module version '%s/%s/%s:%s' deleted successfully", c.organization, c.name, c.provider, c.version))
	return 0
}

func (c *RegistryModuleDeleteVersionCommand) registryModuleService(client *client.Client) registryModuleVersionDeleter {
	if c.registryModuleSvc != nil {
		return c.registryModuleSvc
	}
	return client.RegistryModules
}

// Help returns help text for the registry module delete-version command
func (c *RegistryModuleDeleteVersionCommand) Help() string {
	helpText := `
Usage: hcptf registrymodule delete-version [options]

  Delete a specific version of a private registry module.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Module name (required)
  -provider=<name>     Provider name (required)
  -version=<semver>    Version string (required)

Example:

  hcptf registrymodule delete-version -org=my-org -name=vpc -provider=aws -version=1.2.3
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry module delete-version command
func (c *RegistryModuleDeleteVersionCommand) Synopsis() string {
	return "Delete a version of a private registry module"
}
