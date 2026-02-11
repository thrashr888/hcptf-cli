package command

import (
	"strings"
)

// StackConfigurationUpdateCommand is a command to update a stack configuration
// Note: The API does not support direct updates to stack configurations
type StackConfigurationUpdateCommand struct {
	Meta
	configID string
	format   string
}

// Run executes the stack configuration update command
func (c *StackConfigurationUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stackconfiguration update")
	flags.StringVar(&c.configID, "id", "", "Stack configuration ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	c.Ui.Error("Error: Stack configurations cannot be updated directly")
	c.Ui.Error("Stack configurations are immutable snapshots. To make changes:")
	c.Ui.Error("  1. Update the stack's VCS repository, or")
	c.Ui.Error("  2. Create a new stack configuration")
	c.Ui.Error(c.Help())
	return 1
}

// Help returns help text for the stack configuration update command
func (c *StackConfigurationUpdateCommand) Help() string {
	helpText := `
Usage: hcptf stackconfiguration update [options]

  Stack configurations are immutable and cannot be updated directly.
  They represent point-in-time snapshots of your stack setup.

  To make changes:
    - For VCS-backed stacks: Update the VCS repository
    - For manual stacks: Create a new stack configuration

Alternatives:

  hcptf stackconfiguration create -stack-id=<id>  # Create new configuration
  hcptf stack update -id=<id>                     # Update stack settings
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack configuration update command
func (c *StackConfigurationUpdateCommand) Synopsis() string {
	return "Update stack configuration (not supported - configurations are immutable)"
}
