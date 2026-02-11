package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// TeamCreateCommand is a command to create a team
type TeamCreateCommand struct {
	Meta
	organization string
	name         string
	visibility   string
	format       string
}

// Run executes the team create command
func (c *TeamCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("team create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Team name (required)")
	flags.StringVar(&c.visibility, "visibility", "secret", "Team visibility: secret or organization (default: secret)")
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

	// Validate visibility
	if c.visibility != "secret" && c.visibility != "organization" {
		c.Ui.Error("Error: -visibility must be 'secret' or 'organization'")
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create team
	options := tfe.TeamCreateOptions{
		Name:       tfe.String(c.name),
		Visibility: tfe.String(c.visibility),
	}

	team, err := client.Teams.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating team: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Team '%s' created successfully", team.Name))

	// Show team details
	data := map[string]interface{}{
		"ID":         team.ID,
		"Name":       team.Name,
		"Visibility": team.Visibility,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the team create command
func (c *TeamCreateCommand) Help() string {
	helpText := `
Usage: hcptf team create [options]

  Create a new team.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Team name (required)
  -visibility=<type>   Team visibility: secret or organization (default: secret)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf team create -org=my-org -name=developers
  hcptf team create -org=my-org -name=admins -visibility=organization
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team create command
func (c *TeamCreateCommand) Synopsis() string {
	return "Create a new team"
}
