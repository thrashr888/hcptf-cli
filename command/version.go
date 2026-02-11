package command

import (
	"fmt"
)

// VersionCommand is a Command implementation that prints the version
type VersionCommand struct {
	Meta
}

// Run executes the version command
func (c *VersionCommand) Run(args []string) int {
	version := "0.1.0-dev"
	c.Ui.Output(fmt.Sprintf("hcptf version %s", version))
	return 0
}

// Help returns help text for the version command
func (c *VersionCommand) Help() string {
	return `Usage: hcptf version

  Prints the version of this CLI.
`
}

// Synopsis returns a short synopsis for the version command
func (c *VersionCommand) Synopsis() string {
	return "Prints the version"
}
