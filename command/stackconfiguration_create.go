package command

import (
	"strings"
)

// StackConfigurationCreateCommand is a command to create a stack configuration
type StackConfigurationCreateCommand struct {
	Meta
	stackID     string
	speculative bool
	format      string
}

// Run executes the stack configuration create command
func (c *StackConfigurationCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stackconfiguration create")
	flags.StringVar(&c.stackID, "stack-id", "", "Stack ID (required)")
	flags.BoolVar(&c.speculative, "speculative", false, "Create a speculative configuration (plan-only)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.stackID == "" {
		c.Ui.Error("Error: -stack-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Note: CreateAndUpload requires a path to upload. For manual configurations,
	// users should use the upload URL workflow separately.
	c.Ui.Error("Error: Creating empty stack configurations is not directly supported")
	c.Ui.Error("Stack configurations must be created with content via CreateAndUpload")
	c.Ui.Error("or fetched from VCS using: hcptf stackdeployment create")
	return 1
}

// Help returns help text for the stack configuration create command
func (c *StackConfigurationCreateCommand) Help() string {
	helpText := `
Usage: hcptf stackconfiguration create [options]

  Create a new stack configuration. This creates an empty configuration that
  you can upload a configuration file to. Note: This cannot be used if the
  stack is connected to a VCS repository.

Options:

  -stack-id=<id>   Stack ID (required)
  -speculative     Create a speculative configuration (plan-only)
  -output=<format> Output format: table (default) or json

Example:

  hcptf stackconfiguration create -stack-id=st-abc123
  hcptf stackconfiguration create -stack-id=st-abc123 -speculative
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack configuration create command
func (c *StackConfigurationCreateCommand) Synopsis() string {
	return "Create a new stack configuration"
}
