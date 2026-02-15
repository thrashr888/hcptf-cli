package command

import (
	"fmt"
	"strings"
)

// NamespaceCommand is a lightweight parent command for namespace-level command groups.
type NamespaceCommand struct {
	Meta
	name     string
	synopsis string
}

// Run shows help text for the namespace.
func (c *NamespaceCommand) Run(args []string) int {
	c.Ui.Output(c.Help())
	return 0
}

// Help returns help text for the namespace command.
func (c *NamespaceCommand) Help() string {
	return strings.TrimSpace(fmt.Sprintf(`
Usage: hcptf %s <subcommand> [options]

  %s

For available subcommands, run:
  hcptf %s -help
`, c.name, c.synopsis, c.name))
}

// Synopsis returns a short synopsis for the namespace command.
func (c *NamespaceCommand) Synopsis() string {
	return c.synopsis
}
