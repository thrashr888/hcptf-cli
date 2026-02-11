package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// TeamRemoveMemberCommand is a command to remove a member from a team
type TeamRemoveMemberCommand struct {
	Meta
	organization string
	teamName     string
	username     string
}

// Run executes the team remove-member command
func (c *TeamRemoveMemberCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("team remove-member")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.teamName, "team", "", "Team name (required)")
	flags.StringVar(&c.username, "username", "", "Username to remove (required)")

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

	// Remove team member
	options := tfe.TeamMemberRemoveOptions{
		Usernames: []string{c.username},
	}

	err = client.TeamMembers.Remove(client.Context(), team.ID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error removing team member: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("User '%s' removed from team '%s' successfully", c.username, c.teamName))
	return 0
}

// Help returns help text for the team remove-member command
func (c *TeamRemoveMemberCommand) Help() string {
	helpText := `
Usage: hcptf team remove-member [options]

  Remove a member from a team.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -team=<name>         Team name (required)
  -username=<user>     Username to remove (required)

Example:

  hcptf team remove-member -org=my-org -team=developers -username=alice
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team remove-member command
func (c *TeamRemoveMemberCommand) Synopsis() string {
	return "Remove a member from a team"
}
