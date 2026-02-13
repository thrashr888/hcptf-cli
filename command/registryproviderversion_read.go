package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

type RegistryProviderVersionReadCommand struct {
	Meta
	organization string
	name         string
	namespace    string
	version      string
	format       string
}

// Run executes the registry provider version read command
func (c *RegistryProviderVersionReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registryproviderversion read")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Provider name (required)")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (defaults to organization)")
	flags.StringVar(&c.version, "version", "", "Version string (required)")
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

	// Read provider version
	versionID := tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			RegistryName: tfe.PrivateRegistry,
			Namespace:    c.namespace,
			Name:         c.name,
		},
		Version: c.version,
	}

	providerVersion, err := client.RegistryProviderVersions.Read(client.Context(), versionID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading provider version: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	// Show version details
	data := map[string]interface{}{
		"ID":                 providerVersion.ID,
		"Version":            providerVersion.Version,
		"KeyID":              providerVersion.KeyID,
		"Protocols":          strings.Join(providerVersion.Protocols, ", "),
		"ShasumsUploaded":    providerVersion.ShasumsUploaded,
		"ShasumsSigUploaded": providerVersion.ShasumsSigUploaded,
		"CreatedAt":          providerVersion.CreatedAt,
		"UpdatedAt":          providerVersion.UpdatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the registry provider version read command
func (c *RegistryProviderVersionReadCommand) Help() string {
	helpText := `
Usage: hcptf registryproviderversion read [options]

  Show details of a private registry provider version.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Provider name (required)
  -namespace=<name>    Namespace (defaults to organization)
  -version=<semver>    Version string (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf registryproviderversion read -org=my-org -name=aws -version=3.1.1
  hcptf registryproviderversion read -org=my-org -name=custom -version=1.0.0 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry provider version read command
func (c *RegistryProviderVersionReadCommand) Synopsis() string {
	return "Show details of a private registry provider version"
}
