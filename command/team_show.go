package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// TeamShowCommand is a command to show team details
type TeamShowCommand struct {
	Meta
	organization string
	name         string
	format       string
	teamSvc      teamReader
}

// Run executes the team show command
func (c *TeamShowCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("team show")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Team name (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read team
	team, err := c.teamService(client).Read(client.Context(), c.name)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading team: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                                 team.ID,
		"Name":                               team.Name,
		"Visibility":                         team.Visibility,
		"OrganizationAccessManageWorkspaces": false,
		"OrganizationAccessManagePolicies":   false,
	}

	if team.OrganizationAccess != nil {
		data["OrganizationAccessManageWorkspaces"] = team.OrganizationAccess.ManageWorkspaces
		data["OrganizationAccessManagePolicies"] = team.OrganizationAccess.ManagePolicies
		data["OrganizationAccessManageVCSSettings"] = team.OrganizationAccess.ManageVCSSettings
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the team show command
func (c *TeamShowCommand) Help() string {
	helpText := `
Usage: hcptf team show [options]

  Show team details.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Team name (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf team show -org=my-org -name=developers
  hcptf team show -org=my-org -name=admins -output=json
`
	return strings.TrimSpace(helpText)
}

func (c *TeamShowCommand) teamService(client *client.Client) teamReader {
	if c.teamSvc != nil {
		return c.teamSvc
	}
	return client.Teams
}

// Synopsis returns a short synopsis for the team show command
func (c *TeamShowCommand) Synopsis() string {
	return "Show team details"
}
