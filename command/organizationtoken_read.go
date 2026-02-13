package command

import (
	"fmt"
	"strings"

)

// OrganizationTokenReadCommand is a command to read an organization token
type OrganizationTokenReadCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the organization token read command
func (c *OrganizationTokenReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organizationtoken read")
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

	// Read organization token
	token, err := client.OrganizationTokens.Read(client.Context(), c.organization)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading organization token: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if token == nil {
		c.Ui.Output("No organization token found")
		return 0
	}

	// Show token details (without the secret)
	data := map[string]interface{}{
		"ID":        token.ID,
		"CreatedAt": token.CreatedAt,
	}

	if token.CreatedBy != nil && token.CreatedBy.User != nil {
		data["CreatedBy"] = token.CreatedBy.User.Username
	}

	if !token.LastUsedAt.IsZero() {
		data["LastUsedAt"] = token.LastUsedAt.Format("2006-01-02 15:04:05")
	} else {
		data["LastUsedAt"] = "Never"
	}

	if !token.ExpiredAt.IsZero() {
		data["ExpiredAt"] = token.ExpiredAt.Format("2006-01-02 15:04:05")
	} else {
		data["ExpiredAt"] = "Never"
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the organization token read command
func (c *OrganizationTokenReadCommand) Help() string {
	helpText := `
Usage: hcptf organizationtoken read [options]

  Show organization token details (organization-level API token).

  This command displays metadata about the organization token. Note that
  the actual token value is only displayed when the token is created and
  cannot be retrieved later.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf organizationtoken read -org=my-org
  hcptf organizationtoken read -organization=my-org -output=json

Note:

  Only members of the owners team can access organization tokens.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization token read command
func (c *OrganizationTokenReadCommand) Synopsis() string {
	return "Show organization token details"
}
