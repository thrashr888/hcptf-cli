package command

import (
	"fmt"
	"strings"
)

// OrganizationContextCommand shows help for organization context
type OrganizationContextCommand struct {
	Meta
	organization string
}

// Run shows organization-specific subcommands
func (c *OrganizationContextCommand) Run(args []string) int {
	// Parse the org from args if provided via flag
	for i, arg := range args {
		if strings.HasPrefix(arg, "-org=") {
			c.organization = strings.TrimPrefix(arg, "-org=")
		} else if arg == "-org" && i+1 < len(args) {
			c.organization = args[i+1]
		}
	}

	if c.organization == "" {
		c.Ui.Error("No organization specified")
		return 1
	}

	help := fmt.Sprintf(`Organization: %s

Available commands for this organization:

  Workspaces:
    hcptf <org> workspaces          List workspaces
    hcptf <org> <workspace>         Show workspace details
    hcptf <org> <workspace> -h      Show workspace commands

  Projects:
    hcptf <org> projects            List projects

  Teams:
    hcptf <org> teams               List teams

  Policies:
    hcptf <org> policies            List policies
    hcptf <org> policysets          List policy sets

You can also use traditional command syntax:
    hcptf workspace list -org=%s
    hcptf project list -org=%s
    hcptf team list -org=%s
`, c.organization,
		c.organization, c.organization, c.organization)

	c.Ui.Output(help)
	return 0
}

// Help returns help text
func (c *OrganizationContextCommand) Help() string {
	return "Show organization context help"
}

// Synopsis returns a one-line synopsis
func (c *OrganizationContextCommand) Synopsis() string {
	return "Show organization context help"
}
