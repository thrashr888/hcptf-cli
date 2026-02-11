package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

type RegistryProviderVersionDeleteCommand struct {
	Meta
	organization string
	name         string
	namespace    string
	version      string
}

// Run executes the registry provider version delete command
func (c *RegistryProviderVersionDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registryproviderversion delete")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Provider name (required)")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (defaults to organization)")
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

	// Delete provider version
	versionID := tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			RegistryName: tfe.PrivateRegistry,
			Namespace:    c.namespace,
			Name:         c.name,
		},
		Version: c.version,
	}

	err = client.RegistryProviderVersions.Delete(client.Context(), versionID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting provider version: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Provider version '%s/%s:%s' deleted successfully", c.namespace, c.name, c.version))
	return 0
}

// Help returns help text for the registry provider version delete command
func (c *RegistryProviderVersionDeleteCommand) Help() string {
	helpText := `
Usage: hcptf registryproviderversion delete [options]

  Delete a specific version of a private registry provider.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Provider name (required)
  -namespace=<name>    Namespace (defaults to organization)
  -version=<semver>    Version string (required)

Example:

  hcptf registryproviderversion delete -org=my-org -name=aws -version=3.1.1
  hcptf registryproviderversion delete -org=my-org -name=custom -version=1.0.0
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry provider version delete command
func (c *RegistryProviderVersionDeleteCommand) Synopsis() string {
	return "Delete a version of a private registry provider"
}
