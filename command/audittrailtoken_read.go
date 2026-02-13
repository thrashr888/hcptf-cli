package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// AuditTrailTokenReadCommand is a command to read audit trail token details
type AuditTrailTokenReadCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the audit trail token read command
func (c *AuditTrailTokenReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("audittrailtoken read")
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

	// Read audit trail token
	tokenType := tfe.AuditTrailToken
	options := tfe.OrganizationTokenReadOptions{
		TokenType: &tokenType,
	}

	token, err := client.OrganizationTokens.ReadWithOptions(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading audit trail token: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":          token.ID,
		"Description": token.Description,
		"CreatedAt":   token.CreatedAt.Format("2006-01-02 15:04:05 MST"),
	}

	if !token.LastUsedAt.IsZero() {
		data["LastUsedAt"] = token.LastUsedAt.Format("2006-01-02 15:04:05 MST")
	} else {
		data["LastUsedAt"] = "Never"
	}

	if !token.ExpiredAt.IsZero() {
		data["ExpiredAt"] = token.ExpiredAt.Format("2006-01-02 15:04:05 MST")
	} else {
		data["ExpiredAt"] = "Never"
	}

	if token.CreatedBy != nil && token.CreatedBy.User != nil {
		data["CreatedBy"] = token.CreatedBy.User.Username
	}

	// Note: Token value is not returned by Read endpoint
	c.Ui.Info("Note: Token value is only visible when the token is created.")

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the audit trail token read command
func (c *AuditTrailTokenReadCommand) Help() string {
	helpText := `
Usage: hcptf audittrailtoken read [options]

  Read audit trail token details for an organization. This shows information
  about the token but not the actual token value (only visible at creation).

  Note: Only owners team members can access this command.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf audittrailtoken read -org=my-org
  hcptf audittrailtoken read -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the audit trail token read command
func (c *AuditTrailTokenReadCommand) Synopsis() string {
	return "Read audit trail token details"
}
