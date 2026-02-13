package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// OrganizationMembershipCreateCommand is a command to create an organization membership
type OrganizationMembershipCreateCommand struct {
	Meta
	organization string
	email        string
	teamIDs      string
	format       string
}

// Run executes the organization membership create command
func (c *OrganizationMembershipCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organizationmembership create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.email, "email", "", "Email address of user to invite (required)")
	flags.StringVar(&c.teamIDs, "team-ids", "", "Comma-separated list of team IDs (required)")
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

	if c.email == "" {
		c.Ui.Error("Error: -email flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.teamIDs == "" {
		c.Ui.Error("Error: -team-ids flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Parse team IDs
	teamIDList := strings.Split(c.teamIDs, ",")
	var teams []*tfe.Team
	for _, id := range teamIDList {
		teams = append(teams, &tfe.Team{ID: strings.TrimSpace(id)})
	}

	// Build create options
	options := tfe.OrganizationMembershipCreateOptions{
		Email: tfe.String(c.email),
		Teams: teams,
	}

	// Create organization membership
	membership, err := client.OrganizationMemberships.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating organization membership: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Organization membership created successfully (Status: %s)", membership.Status))

	// Show membership details
	data := map[string]interface{}{
		"ID":     membership.ID,
		"Email":  c.email,
		"Status": membership.Status,
	}

	if membership.User != nil {
		data["UserID"] = membership.User.ID
		if membership.User.Username != "" {
			data["Username"] = membership.User.Username
		}
	}

	if membership.Organization != nil {
		data["Organization"] = membership.Organization.Name
	}

	// Show team IDs
	var teamIDStrings []string
	for _, team := range teams {
		teamIDStrings = append(teamIDStrings, team.ID)
	}
	data["Teams"] = strings.Join(teamIDStrings, ", ")

	formatter.KeyValue(data)

	if membership.Status == "invited" {
		c.Ui.Info("\nNote: An invitation email has been sent to the user.")
	}

	return 0
}

// Help returns help text for the organization membership create command
func (c *OrganizationMembershipCreateCommand) Help() string {
	helpText := `
Usage: hcptf organizationmembership create [options]

  Invite a user to join an organization.

  This command invites a user to join an organization and assigns them to
  one or more teams. If the user doesn't have an account, they can create
  one using the invitation email.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -email=<email>       Email address of user to invite (required)
  -team-ids=<ids>      Comma-separated list of team IDs (required)
                       Users must be added to at least one team
  -output=<format>     Output format: table (default) or json

Example:

  # Invite user to a single team
  hcptf organizationmembership create -org=my-org \
    -email=user@example.com -team-ids=team-abc123

  # Invite user to multiple teams
  hcptf organizationmembership create -org=my-org \
    -email=user@example.com -team-ids=team-abc123,team-def456

Note:

  Only members with team management permissions or owners can invite users
  to an organization. Users must be added to at least one team.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization membership create command
func (c *OrganizationMembershipCreateCommand) Synopsis() string {
	return "Invite a user to join an organization"
}
