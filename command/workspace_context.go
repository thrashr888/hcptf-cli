package command

import (
	"fmt"
	"strings"
)

// WorkspaceContextCommand shows help for workspace context
type WorkspaceContextCommand struct {
	Meta
	organization string
	workspace    string
}

// Run shows workspace-specific subcommands
func (c *WorkspaceContextCommand) Run(args []string) int {
	// Parse the org and workspace from args if provided via flags
	for i, arg := range args {
		if strings.HasPrefix(arg, "-org=") {
			c.organization = strings.TrimPrefix(arg, "-org=")
		} else if arg == "-org" && i+1 < len(args) {
			c.organization = args[i+1]
		}
		if strings.HasPrefix(arg, "-workspace=") {
			c.workspace = strings.TrimPrefix(arg, "-workspace=")
		} else if arg == "-workspace" && i+1 < len(args) {
			c.workspace = args[i+1]
		}
	}

	if c.organization == "" || c.workspace == "" {
		c.Ui.Error("No organization or workspace specified")
		return 1
	}

	help := fmt.Sprintf(`Workspace: %s/%s

Available commands for this workspace:

  Runs:
    hcptf <org> <workspace> runs              List runs
    hcptf <org> <workspace> runs <run-id>     Show run details

  Variables:
    hcptf <org> <workspace> variables         List variables

  State:
    hcptf <org> <workspace> state             List state versions
    hcptf <org> <workspace> state outputs     Show state outputs

  Configuration Versions:
    configversion list -workspace=%s

You can also use traditional command syntax:
    hcptf run list -org=%s -workspace=%s
    hcptf variable list -org=%s -workspace=%s
    hcptf state list -org=%s -workspace=%s
`,
		c.organization, c.workspace,
		c.workspace,
		c.organization, c.workspace,
		c.organization, c.workspace,
		c.organization, c.workspace)

	c.Ui.Output(help)
	return 0
}

// Help returns help text
func (c *WorkspaceContextCommand) Help() string {
	return "Show workspace context help"
}

// Synopsis returns a one-line synopsis
func (c *WorkspaceContextCommand) Synopsis() string {
	return "Show workspace context help"
}
