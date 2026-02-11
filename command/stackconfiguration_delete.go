package command

import (
	"strings"
)

// StackConfigurationDeleteCommand is a command to delete a stack configuration
// Note: The API does not support direct deletion of stack configurations
type StackConfigurationDeleteCommand struct {
	Meta
	configID string
	force    bool
}

// Run executes the stack configuration delete command
func (c *StackConfigurationDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stackconfiguration delete")
	flags.StringVar(&c.configID, "id", "", "Stack configuration ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	c.Ui.Error("Error: Stack configurations cannot be deleted directly")
	c.Ui.Error("Stack configurations are managed by HCP Terraform automatically.")
	c.Ui.Error("To remove a stack's configurations, delete the stack itself:")
	c.Ui.Error("  hcptf stack delete -id=<stack-id>")
	c.Ui.Error(c.Help())
	return 1
}

// Help returns help text for the stack configuration delete command
func (c *StackConfigurationDeleteCommand) Help() string {
	helpText := `
Usage: hcptf stackconfiguration delete [options]

  Stack configurations cannot be deleted individually. They are managed
  automatically by HCP Terraform as part of the stack lifecycle.

  To remove all configurations:
    hcptf stack delete -id=<stack-id>

Alternative:

  hcptf stack delete -id=<stack-id>  # Delete stack and all configurations
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack configuration delete command
func (c *StackConfigurationDeleteCommand) Synopsis() string {
	return "Delete stack configuration (not supported - managed by HCP Terraform)"
}
