package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// OrganizationMemberReadCommand is a command to read detailed organization member information
type OrganizationMemberReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the organization member read command
func (c *OrganizationMemberReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organizationmember read")
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
		c.Ui.Error(fmt.Sprintf("Error reading organization member: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	// Show membership details
	data := map[string]interface{}{
		"ID":     membership.ID,
		"Status": membership.Status,
	}

	// User information
	if membership.User != nil {
		data["UserID"] = membership.User.ID
		data["Email"] = membership.User.Email
		if membership.User.Username != "" {
			data["Username"] = membership.User.Username
		}

		// User permissions
		if membership.User.Permissions != nil {
			data["CanCreateOrganizations"] = membership.User.Permissions.CanCreateOrganizations
			data["CanChangeEmail"] = membership.User.Permissions.CanChangeEmail
			data["CanChangeUsername"] = membership.User.Permissions.CanChangeUsername
		}

		// Two-factor authentication status
		if membership.User.TwoFactor != nil {
			data["TwoFactorEnabled"] = membership.User.TwoFactor.Enabled
			data["TwoFactorVerified"] = membership.User.TwoFactor.Verified
		}
	}

	if membership.Organization != nil {
		data["Organization"] = membership.Organization.Name
	}

	// Show detailed team information
	if len(membership.Teams) > 0 {
		var teamList []string
		var teamIDs []string
		for _, team := range membership.Teams {
			teamList = append(teamList, team.Name)
			teamIDs = append(teamIDs, team.ID)
		}
		data["Teams"] = strings.Join(teamList, ", ")
		data["TeamIDs"] = strings.Join(teamIDs, ", ")
		data["TeamCount"] = len(membership.Teams)
	} else {
		data["Teams"] = "None"
		data["TeamCount"] = 0
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the organization member read command
func (c *OrganizationMemberReadCommand) Help() string {
	helpText := `
Usage: hcptf organizationmember read [options]

  Show detailed organization member information including teams and permissions.

  This command displays comprehensive information about a specific organization
  membership, including the user's status, team memberships, and account permissions.

Options:

  -id=<id>          Organization membership ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf organizationmember read -id=ou-abc123xyz
  hcptf organizationmember read -id=ou-abc123xyz -output=json

Note:

  Only members with team management permissions or owners can read
  detailed organization membership information.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization member read command
func (c *OrganizationMemberReadCommand) Synopsis() string {
	return "Show detailed organization member information"
}
