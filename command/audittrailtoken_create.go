package command

import (
	"fmt"
	"strings"
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

// AuditTrailTokenCreateCommand is a command to create an audit trail token
type AuditTrailTokenCreateCommand struct {
	Meta
	organization string
	expiredAt    string
	format       string
}

// Run executes the audit trail token create command
func (c *AuditTrailTokenCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("audittrailtoken create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.expiredAt, "expired-at", "", "Token expiration date (ISO8601 format: YYYY-MM-DDTHH:MM:SS.SSSZ)")
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
			c.Ui.Error(fmt.Sprintf("Error: invalid date format for -expired-at: %s", err))
			c.Ui.Error("Use ISO8601 format: YYYY-MM-DDTHH:MM:SS.SSSZ")
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

	// Create audit trail token
	tokenType := tfe.AuditTrailToken
	options := tfe.OrganizationTokenCreateOptions{
		ExpiredAt: expiredAt,
		TokenType: &tokenType,
	}

	token, err := client.OrganizationTokens.CreateWithOptions(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating audit trail token: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output("Audit trail token created successfully")

	// Show token details
	data := map[string]interface{}{
		"ID":        token.ID,
		"Token":     token.Token,
		"CreatedAt": token.CreatedAt.Format("2006-01-02 15:04:05 MST"),
	}

	if !token.ExpiredAt.IsZero() {
		data["ExpiredAt"] = token.ExpiredAt.Format("2006-01-02 15:04:05 MST")
	} else {
		data["ExpiredAt"] = "Never"
	}

	if token.CreatedBy != nil && token.CreatedBy.User != nil {
		data["CreatedBy"] = token.CreatedBy.User.Username
	}

	formatter.KeyValue(data)

	// Warning about token visibility
	c.Ui.Warn("\nWARNING: This is the only time the token will be displayed. Save it securely.")
	c.Ui.Warn("This token can be used to access your organization's audit trail data.")

	return 0
}

// Help returns help text for the audit trail token create command
func (c *AuditTrailTokenCreateCommand) Help() string {
	helpText := `
Usage: hcptf audittrailtoken create [options]

  Create an audit trail token for an organization. This generates a new
  audit trail token, replacing any existing token. The token is used to
  authenticate integrations pulling audit trail data.

  Note: Only owners team members can access this command.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -expired-at=<date>   Token expiration date (ISO8601 format). If omitted, token never expires.
  -output=<format>     Output format: table (default) or json

Example:

  hcptf audittrailtoken create -org=my-org
  hcptf audittrailtoken create -org=my-org -expired-at=2025-12-31T23:59:59.000Z
  hcptf audittrailtoken create -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the audit trail token create command
func (c *AuditTrailTokenCreateCommand) Synopsis() string {
	return "Create an audit trail token for an organization"
}
