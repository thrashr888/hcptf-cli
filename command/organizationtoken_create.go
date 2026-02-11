package command

import (
	"fmt"
	"strings"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// OrganizationTokenCreateCommand is a command to create an organization token
type OrganizationTokenCreateCommand struct {
	Meta
	organization string
	expiredAt    string
	format       string
}

// Run executes the organization token create command
func (c *OrganizationTokenCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organizationtoken create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.expiredAt, "expired-at", "", "Expiration date in ISO 8601 format (e.g., 2024-12-31T23:59:59Z)")
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

	// Parse expiration date if provided
	var expiredAt *time.Time
	if c.expiredAt != "" {
		t, err := time.Parse(time.RFC3339, c.expiredAt)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error parsing expired-at date: %s", err))
			c.Ui.Error("Use ISO 8601 format (e.g., 2024-12-31T23:59:59Z)")
			return 1
		}
		expiredAt = &t
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create organization token
	var token *tfe.OrganizationToken
	if expiredAt != nil {
		options := tfe.OrganizationTokenCreateOptions{
			ExpiredAt: expiredAt,
		}
		token, err = client.OrganizationTokens.CreateWithOptions(client.Context(), c.organization, options)
	} else {
		token, err = client.OrganizationTokens.Create(client.Context(), c.organization)
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating organization token: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output("Organization token created successfully")

	// Show token details including the secret value
	data := map[string]interface{}{
		"ID":        token.ID,
		"Token":     token.Token,
		"CreatedAt": token.CreatedAt,
	}

	if token.CreatedBy != nil && token.CreatedBy.User != nil {
		data["CreatedBy"] = token.CreatedBy.User.Username
	}

	if !token.ExpiredAt.IsZero() {
		data["ExpiredAt"] = token.ExpiredAt.Format("2006-01-02 15:04:05")
	} else {
		data["ExpiredAt"] = "Never"
	}

	formatter.KeyValue(data)

	// Warning about token visibility
	c.Ui.Warn("\nWARNING: This is the only time the token will be displayed. Save it securely.")
	c.Ui.Warn("Creating a new organization token will invalidate the previous token.")

	return 0
}

// Help returns help text for the organization token create command
func (c *OrganizationTokenCreateCommand) Help() string {
	helpText := `
Usage: hcptf organizationtoken create [options]

  Create an organization token (organization-level API token).

  Organization tokens can authenticate as the organization itself, with full
  access to all organization resources. Creating a new token will invalidate
  any existing organization token.

Options:

  -organization=<name>   Organization name (required)
  -org=<name>           Alias for -organization
  -expired-at=<date>    Token expiration date in ISO 8601 format
                        (e.g., 2024-12-31T23:59:59Z). If omitted, token never expires.
  -output=<format>      Output format: table (default) or json

Example:

  # Create organization token that never expires
  hcptf organizationtoken create -org=my-org

  # Create organization token with expiration
  hcptf organizationtoken create -org=my-org -expired-at=2024-12-31T23:59:59Z

Security Note:

  The token value is only displayed once upon creation and cannot be retrieved
  later. Store it securely. Only members of the owners team can create
  organization tokens.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization token create command
func (c *OrganizationTokenCreateCommand) Synopsis() string {
	return "Create an organization token"
}
