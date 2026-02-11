package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// TeamAddMemberCommand is a command to add a member to a team
type TeamAddMemberCommand struct {
	Meta
	organization string
	teamName     string
	username     string
}

// Run executes the team add-member command
func (c *TeamAddMemberCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("team add-member")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.teamName, "team", "", "Team name (required)")
	flags.StringVar(&c.username, "username", "", "Username to add (required)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.teamName == "" {
		c.Ui.Error("Error: -team flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.username == "" {
		c.Ui.Error("Error: -username flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read team to get ID
	team, err := client.Teams.Read(client.Context(), c.teamName)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading team: %s", err))
		return 1
	}

	// Add team member
	options := tfe.TeamMemberAddOptions{
		Usernames: []string{c.username},
	}

	err = client.TeamMembers.Add(client.Context(), team.ID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error adding team member: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("User '%s' added to team '%s' successfully", c.username, c.teamName))
	return 0
}

// Help returns help text for the team add-member command
func (c *TeamAddMemberCommand) Help() string {
	helpText := `
Usage: hcptf team add-member [options]

  Add a member to a team.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -team=<name>         Team name (required)
  -username=<user>     Username to add (required)

Example:

  hcptf team add-member -org=my-org -team=developers -username=alice
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team add-member command
func (c *TeamAddMemberCommand) Synopsis() string {
	return "Add a member to a team"
}
