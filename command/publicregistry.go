package command

import "strings"

// PublicRegistryCommand is a parent command for querying the public Terraform registry
type PublicRegistryCommand struct {
	Meta
}

// Run shows help text for publicregistry commands
func (c *PublicRegistryCommand) Run(args []string) int {
	c.Ui.Output(c.Help())
	return 0
}

// Help returns help text for the publicregistry command
func (c *PublicRegistryCommand) Help() string {
	helpText := `
Usage: hcptf publicregistry <subcommand> [options]

  Query the public Terraform registry at registry.terraform.io.

  Use these commands to find information about publicly available
  Terraform providers and modules, check latest versions, and view
  documentation links.

Subcommands:

  Provider Commands:
    provider                 Get provider information
    provider versions        List all available provider versions

  Module Commands:
    module                   Get module information

Examples:

  # Get information about the AWS provider
  hcptf publicregistry provider -name=hashicorp/aws

  # List all available versions of the AWS provider
  hcptf publicregistry provider versions -name=hashicorp/aws

  # Get information about the VPC module
  hcptf publicregistry module -name=terraform-aws-modules/vpc/aws

  # Check latest version of Random provider
  hcptf publicregistry provider -name=hashicorp/random

For detailed help on any subcommand:
  hcptf publicregistry <subcommand> -help
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the publicregistry command
func (c *PublicRegistryCommand) Synopsis() string {
	return "Query public Terraform registry"
}
