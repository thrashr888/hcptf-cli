package command

import "strings"

// RegistryCommand is a parent command for all registry operations
type RegistryCommand struct {
	Meta
}

// Run shows help text for registry commands
func (c *RegistryCommand) Run(args []string) int {
	c.Ui.Output(c.Help())
	return 0
}

// Help returns help text for the registry command
func (c *RegistryCommand) Help() string {
	helpText := `
Usage: hcptf registry <subcommand> [options]

  Manage private registry resources (modules and providers).

  The HCP Terraform private registry allows you to publish and share
  Terraform modules and providers within your organization.

Subcommands:

  Module Commands:
    module list              List private modules
    module create            Publish a new module
    module read              Show module details
    module delete            Delete a module
    module version create    Create a new module version
    module version delete    Delete a module version

  Provider Commands:
    provider list            List private providers
    provider create          Publish a new provider
    provider read            Show provider details
    provider delete          Delete a provider
    provider version create  Create a new provider version
    provider version read    Show provider version details
    provider version delete  Delete a provider version
    provider platform create Add a platform to a provider version
    provider platform read   Show provider platform details
    provider platform delete Delete a provider platform

Examples:

  # List all modules in the private registry
  hcptf registry module list -organization=my-org

  # Create a new private provider
  hcptf registry provider create -organization=my-org -name=my-provider

  # Show details of a specific module
  hcptf registry module read -organization=my-org -name=vpc -namespace=my-org

For detailed help on any subcommand:
  hcptf registry <subcommand> -help
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry command
func (c *RegistryCommand) Synopsis() string {
	return "Manage private registry resources"
}
