package command

import (
	"context"
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type registryProviderPlatformDeleter interface {
	Delete(ctx context.Context, platformID tfe.RegistryProviderPlatformID) error
}

type RegistryProviderPlatformDeleteCommand struct {
	Meta
	organization                string
	name                        string
	namespace                   string
	version                     string
	os                          string
	arch                        string
	force                       bool
	yes                         bool
	registryProviderPlatformSvc registryProviderPlatformDeleter
}

// Run executes the registry provider platform delete command
func (c *RegistryProviderPlatformDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registryproviderplatform delete")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Provider name (required)")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (defaults to organization)")
	flags.StringVar(&c.version, "version", "", "Version string (required)")
	flags.StringVar(&c.os, "os", "", "Operating system (required)")
	flags.StringVar(&c.arch, "arch", "", "Architecture (required)")
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

	if !c.force && !c.yes {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete provider platform '%s/%s:%s (%s/%s)'? (yes/no): ", c.namespace, c.name, c.version, c.os, c.arch))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}
		if strings.TrimSpace(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete provider platform
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

	err = c.registryProviderPlatformService(client).Delete(client.Context(), platformID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting provider platform: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Provider platform '%s/%s:%s (%s/%s)' deleted successfully", c.namespace, c.name, c.version, c.os, c.arch))
	return 0
}

// Help returns help text for the registry provider platform delete command
func (c *RegistryProviderPlatformDeleteCommand) Help() string {
	helpText := `
	Usage: hcptf registry provider platform delete [options]

  Delete a specific platform binary of a private registry provider version.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Provider name (required)
  -namespace=<name>    Namespace (defaults to organization)
  -version=<semver>    Version string (required)
  -os=<os>             Operating system (required)
  -arch=<arch>         Architecture (required)
  -force               Force delete without confirmation
  -f                   Shorthand for -force
  -y                   Confirm delete without prompt

Example:

  hcptf registryproviderplatform delete -org=my-org -name=aws -version=3.1.1 -os=linux -arch=amd64
  hcptf registryproviderplatform delete -org=my-org -name=custom -version=1.0.0 -os=darwin -arch=arm64
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry provider platform delete command
func (c *RegistryProviderPlatformDeleteCommand) Synopsis() string {
	return "Delete a platform binary of a private registry provider version"
}

func (c *RegistryProviderPlatformDeleteCommand) registryProviderPlatformService(client *client.Client) registryProviderPlatformDeleter {
	if c.registryProviderPlatformSvc != nil {
		return c.registryProviderPlatformSvc
	}
	return client.RegistryProviderPlatforms
}
