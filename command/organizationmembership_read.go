package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// OrganizationMembershipReadCommand is a command to read an organization membership
type OrganizationMembershipReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the organization membership read command
func (c *OrganizationMembershipReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organizationmembership read")
	flags.StringVar(&c.id, "id", "", "Organization membership ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read organization membership with included resources
	options := tfe.OrganizationMembershipReadOptions{
		Include: []tfe.OrgMembershipIncludeOpt{
			tfe.OrgMembershipUser,
		},
	}

	membership, err := client.OrganizationMemberships.ReadWithOptions(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading organization membership: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	// Show membership details
	data := map[string]interface{}{
		"ID":     membership.ID,
		"Status": membership.Status,
	}

	if membership.User != nil {
		data["UserID"] = membership.User.ID
		data["Email"] = membership.User.Email
		if membership.User.Username != "" {
			data["Username"] = membership.User.Username
		}
	}

	if membership.Organization != nil {
		data["Organization"] = membership.Organization.Name
	}

	// Show team information
	if len(membership.Teams) > 0 {
		var teamInfo []string
		for _, team := range membership.Teams {
			teamInfo = append(teamInfo, fmt.Sprintf("%s (%s)", team.Name, team.ID))
		}
		data["Teams"] = strings.Join(teamInfo, ", ")
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the organization membership read command
func (c *OrganizationMembershipReadCommand) Help() string {
	helpText := `
Usage: hcptf organizationmembership read [options]

  Show organization membership details.

  This command displays detailed information about a specific organization
  membership, including the user's status and team memberships.

Options:

  -id=<id>          Organization membership ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf organizationmembership read -id=ou-abc123xyz
  hcptf organizationmembership read -id=ou-abc123xyz -output=json

Note:

  Only members with team management permissions or owners can read
  organization membership details.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization membership read command
func (c *OrganizationMembershipReadCommand) Synopsis() string {
	return "Show organization membership details"
}
