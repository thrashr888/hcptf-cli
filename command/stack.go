package command

import "strings"

// StackCommand is a parent command for all stack operations
type StackCommand struct {
	Meta
}

// Run shows help text for stack commands
func (c *StackCommand) Run(args []string) int {
	c.Ui.Output(c.Help())
	return 0
}

// Help returns help text for the stack command
func (c *StackCommand) Help() string {
	helpText := `
Usage: hcptf stack <subcommand> [options]

  Manage Terraform Stacks and their deployments.

  Terraform Stacks enable orchestrating multiple Terraform configurations
  as a cohesive unit, with shared inputs and cross-stack dependencies.

Subcommands:

  Stack Management:
    list                  List stacks
    create                Create a new stack
    read                  Show stack details
    update                Update a stack
    delete                Delete a stack

  Configuration Commands:
    configuration list    List stack configurations
    configuration create  Create a stack configuration
    configuration read    Show configuration details
    configuration update  Update a stack configuration
    configuration delete  Delete a stack configuration

  Deployment Commands:
    deployment list       List stack deployments
    deployment create     Trigger a new deployment
    deployment read       Show deployment details

  State Commands:
    state list           List stack state versions
    state read           Show state version details

Examples:

  # List all stacks in a project
  hcptf stack list -organization=my-org -project=my-project

  # Create a new stack
  hcptf stack create -organization=my-org -project=my-project -name=prod

  # List configurations for a stack
  hcptf stack configuration list -stack-id=stk-abc123

  # Trigger a deployment
  hcptf stack deployment create -stack-id=stk-abc123

For detailed help on any subcommand:
  hcptf stack <subcommand> -help
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack command
func (c *StackCommand) Synopsis() string {
	return "Manage Terraform Stacks"
}
