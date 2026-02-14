package command

import (
	"fmt"
	"strings"
)

// OrganizationTokenListCommand is a command to list organization tokens
type OrganizationTokenListCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the organization token list command
func (c *OrganizationTokenListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organizationtoken list")
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

// Help returns help text for the organization token list command
func (c *OrganizationTokenListCommand) Help() string {
	helpText := `
Usage: hcptf organizationtoken list [options]

  List organization token (organization-level API token).

  This command displays information about the organization token, which can
  be used to authenticate API requests at the organization level. Note that
  the actual token value is only displayed when the token is created.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf organizationtoken list -org=my-org
  hcptf organizationtoken list -organization=my-org -output=json

Note:

  Only members of the owners team can access organization tokens.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization token list command
func (c *OrganizationTokenListCommand) Synopsis() string {
	return "List organization token"
}
