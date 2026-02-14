package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// AuditTrailTokenListCommand is a command to list audit trail tokens
type AuditTrailTokenListCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the audit trail token list command
func (c *AuditTrailTokenListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("audittrailtoken list")
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

	// Read audit trail token (there's only one per organization)
	tokenType := tfe.AuditTrailToken
	options := tfe.OrganizationTokenReadOptions{
		TokenType: &tokenType,
	}

	token, err := client.OrganizationTokens.ReadWithOptions(client.Context(), c.organization, options)
	if err != nil {
		// Check if it's a 404 - no token exists
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			c.Ui.Output("No audit trail token found")
			return 0
		}
		c.Ui.Error(fmt.Sprintf("Error reading audit trail token: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	// Prepare table data
	headers := []string{"ID", "Description", "Created At", "Last Used At", "Expired At"}
	var rows [][]string

	lastUsed := "Never"
	if !token.LastUsedAt.IsZero() {
		lastUsed = token.LastUsedAt.Format("2006-01-02 15:04:05")
	}

	expiredAt := "Never"
	if !token.ExpiredAt.IsZero() {
		expiredAt = token.ExpiredAt.Format("2006-01-02 15:04:05")
	}

	description := token.Description
	if description == "" {
		description = "-"
	}

	rows = append(rows, []string{
		token.ID,
		description,
		token.CreatedAt.Format("2006-01-02 15:04:05"),
		lastUsed,
		expiredAt,
	})

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the audit trail token list command
func (c *AuditTrailTokenListCommand) Help() string {
	helpText := `
Usage: hcptf audittrailtoken list [options]

  Show the audit trail token for an organization. Organizations can have
  only one audit trail token at a time. Audit trail tokens are used to
  authenticate integrations pulling audit trail data, such as the
  HCP Terraform for Splunk app.

  Note: Only owners team members can access this command.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf audittrailtoken list -org=my-org
  hcptf audittrailtoken list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the audit trail token list command
func (c *AuditTrailTokenListCommand) Synopsis() string {
	return "Show the audit trail token for an organization"
}
