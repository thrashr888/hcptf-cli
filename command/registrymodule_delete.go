package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

type RegistryModuleDeleteCommand struct {
	Meta
	organization      string
	name              string
	force             bool
	yes               bool
	registryModuleSvc registryModuleDeleter
}

// Run executes the registry module delete command
func (c *RegistryModuleDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registrymodule delete")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Module name (required)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	if !c.force && !c.yes {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete module '%s/%s'? (yes/no): ", c.organization, c.name))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}
		if strings.TrimSpace(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete registry module
	err = c.registryModuleService(client).Delete(client.Context(), c.organization, c.name)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting registry module: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Registry module '%s/%s' deleted successfully", c.organization, c.name))
	return 0
}

func (c *RegistryModuleDeleteCommand) registryModuleService(client *client.Client) registryModuleDeleter {
	if c.registryModuleSvc != nil {
		return c.registryModuleSvc
	}
	return client.RegistryModules
}

// Help returns help text for the registry module delete command
func (c *RegistryModuleDeleteCommand) Help() string {
	helpText := `
	Usage: hcptf registry module delete [options]

  Delete a private registry module.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Module name (required)
  -force               Force delete without confirmation
  -f                   Shorthand for -force
  -y                   Confirm delete without prompt

Example:

  hcptf registry module delete -org=my-org -name=vpc
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry module delete command
func (c *RegistryModuleDeleteCommand) Synopsis() string {
	return "Delete a private registry module"
}
