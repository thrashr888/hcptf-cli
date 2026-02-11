package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// TeamListCommand is a command to list teams
type TeamListCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the team list command
func (c *TeamListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("team list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List teams
	teams, err := client.Teams.List(client.Context(), c.organization, &tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing teams: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(teams.Items) == 0 {
		c.Ui.Output("No teams found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Visibility"}
	var rows [][]string

	for _, team := range teams.Items {
		rows = append(rows, []string{
			team.ID,
			team.Name,
			team.Visibility,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the team list command
func (c *TeamListCommand) Help() string {
	helpText := `
Usage: hcptf team list [options]

  List teams in an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf team list -org=my-org
  hcptf team list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team list command
func (c *TeamListCommand) Synopsis() string {
	return "List teams in an organization"
}
