package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

type RegistryProviderPlatformReadCommand struct {
	Meta
	organization string
	name         string
	namespace    string
	version      string
	os           string
	arch         string
	format       string
}

// Run executes the registry provider platform read command
func (c *RegistryProviderPlatformReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registryproviderplatform read")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Provider name (required)")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (defaults to organization)")
	flags.StringVar(&c.version, "version", "", "Version string (required)")
	flags.StringVar(&c.os, "os", "", "Operating system (required)")
	flags.StringVar(&c.arch, "arch", "", "Architecture (required)")
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

	if c.os == "" {
		c.Ui.Error("Error: -os flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.arch == "" {
		c.Ui.Error("Error: -arch flag is required")
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

	// Read provider platform
	platformID := tfe.RegistryProviderPlatformID{
		RegistryProviderVersionID: tfe.RegistryProviderVersionID{
			RegistryProviderID: tfe.RegistryProviderID{
				RegistryName: tfe.PrivateRegistry,
				Namespace:    c.namespace,
				Name:         c.name,
			},
			Version: c.version,
		},
		OS:   c.os,
		Arch: c.arch,
	}

	platform, err := client.RegistryProviderPlatforms.Read(client.Context(), platformID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading provider platform: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	// Show platform details
	data := map[string]interface{}{
		"ID":                     platform.ID,
		"OS":                     platform.OS,
		"Arch":                   platform.Arch,
		"Filename":               platform.Filename,
		"Shasum":                 platform.Shasum,
		"ProviderBinaryUploaded": platform.ProviderBinaryUploaded,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the registry provider platform read command
func (c *RegistryProviderPlatformReadCommand) Help() string {
	helpText := `
Usage: hcptf registryproviderplatform read [options]

  Show details of a private registry provider platform.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Provider name (required)
  -namespace=<name>    Namespace (defaults to organization)
  -version=<semver>    Version string (required)
  -os=<os>             Operating system (required)
  -arch=<arch>         Architecture (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf registryproviderplatform read -org=my-org -name=aws -version=3.1.1 -os=linux -arch=amd64
  hcptf registryproviderplatform read -org=my-org -name=custom -version=1.0.0 -os=darwin -arch=arm64 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry provider platform read command
func (c *RegistryProviderPlatformReadCommand) Synopsis() string {
	return "Show details of a private registry provider platform"
}
