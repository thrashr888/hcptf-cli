package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// OrganizationMembershipListCommand is a command to list organization memberships
type OrganizationMembershipListCommand struct {
	Meta
	organization string
	status       string
	email        string
	format       string
}

// Run executes the organization membership list command
func (c *OrganizationMembershipListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organizationmembership list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.status, "status", "", "Filter by status (invited, active)")
	flags.StringVar(&c.email, "email", "", "Filter by email address")
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

	// Build list options
	options := &tfe.OrganizationMembershipListOptions{
		Include: []tfe.OrgMembershipIncludeOpt{tfe.OrgMembershipUser},
	}

	if c.status != "" {
		options.Status = tfe.OrganizationMembershipStatus(c.status)
	}

	if c.email != "" {
		options.Emails = []string{c.email}
	}

	// List organization memberships
	memberships, err := client.OrganizationMemberships.List(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing organization memberships: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(memberships.Items) == 0 {
		c.Ui.Output("No organization memberships found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "User ID", "Email", "Username", "Status"}
	var rows [][]string

	for _, membership := range memberships.Items {
		userID := ""
		email := ""
		username := ""

		if membership.User != nil {
			userID = membership.User.ID
			email = membership.User.Email
			username = membership.User.Username
			if username == "" {
				username = "-"
			}
		}

		rows = append(rows, []string{
			membership.ID,
			userID,
			email,
			username,
			string(membership.Status),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the organization membership list command
func (c *OrganizationMembershipListCommand) Help() string {
	helpText := `
Usage: hcptf organizationmembership list [options]

  List organization memberships.

  This command displays all members of an organization, including their
  membership status (invited or active).

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -status=<status>     Filter by status (invited, active)
  -email=<email>       Filter by email address
  -output=<format>     Output format: table (default) or json

Example:

  # List all memberships
  hcptf organizationmembership list -org=my-org

  # List only active memberships
  hcptf organizationmembership list -org=my-org -status=active

  # List invited memberships
  hcptf organizationmembership list -org=my-org -status=invited

  # Filter by email
  hcptf organizationmembership list -org=my-org -email=user@example.com

  # Output as JSON
  hcptf organizationmembership list -org=my-org -output=json

Note:

  Only members with team management permissions or owners can list
  organization memberships.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization membership list command
func (c *OrganizationMembershipListCommand) Synopsis() string {
	return "List organization memberships"
}
