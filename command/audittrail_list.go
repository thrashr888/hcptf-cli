package command

import (
	"fmt"
	"strings"
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

// AuditTrailListCommand is a command to list audit trail events
type AuditTrailListCommand struct {
	Meta
	organization string
	since        string
	pageNumber   int
	pageSize     int
	format       string
}

// Run executes the audit trail list command
func (c *AuditTrailListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("audittrail list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.since, "since", "", "Return audit events since this date (ISO8601 format: YYYY-MM-DDTHH:MM:SS.SSSZ)")
	flags.IntVar(&c.pageNumber, "page-number", 1, "Page number")
	flags.IntVar(&c.pageSize, "page-size", 100, "Number of items per page (max 1000)")
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

	// Validate since date format if provided
	var sinceTime time.Time
	if c.since != "" {
		t, err := time.Parse(time.RFC3339, c.since)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error: invalid date format for -since: %s", err))
			c.Ui.Error("Use ISO8601 format: YYYY-MM-DDTHH:MM:SS.SSSZ")
			return 1
		}
		sinceTime = t
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build list options
	options := &tfe.AuditTrailListOptions{
		Since: sinceTime,
		ListOptions: &tfe.ListOptions{
			PageNumber: c.pageNumber,
			PageSize:   c.pageSize,
		},
	}

	// List audit trail events
	auditTrails, err := client.AuditTrails.List(client.Context(), options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing audit trail events: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(auditTrails.Items) == 0 {
		c.Ui.Output("No audit trail events found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Timestamp", "Type", "Resource Type", "Action", "Actor"}
	var rows [][]string

	for _, at := range auditTrails.Items {
		actor := at.Auth.Description
		if actor == "" {
			actor = at.Auth.AccessorID
		}

		rows = append(rows, []string{
			at.ID,
			at.Timestamp.Format("2006-01-02 15:04:05"),
			at.Type,
			at.Resource.Type,
			at.Resource.Action,
			actor,
		})
	}

	formatter.Table(headers, rows)

	// Show pagination info
	if auditTrails.AuditTrailPagination != nil {
		c.Ui.Output(fmt.Sprintf("\nPage %d of %d (Total: %d events)",
			auditTrails.AuditTrailPagination.CurrentPage,
			auditTrails.AuditTrailPagination.TotalPages,
			auditTrails.AuditTrailPagination.TotalCount))
	}

	return 0
}

// Help returns help text for the audit trail list command
func (c *AuditTrailListCommand) Help() string {
	helpText := `
Usage: hcptf audittrail list [options]

  List audit trail events for an organization. Audit trails provide compliance
  and security monitoring by logging all API actions. HCP Terraform retains
  14 days of audit log information.

  Note: Requires an organization token or audit trail token.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -since=<datetime>    Return audit events since this date (ISO8601 format)
  -page-number=<num>   Page number (default: 1)
  -page-size=<num>     Number of items per page (default: 100, max: 1000)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf audittrail list -org=my-org
  hcptf audittrail list -org=my-org -since=2024-01-01T00:00:00.000Z
  hcptf audittrail list -org=my-org -page-size=50 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the audit trail list command
func (c *AuditTrailListCommand) Synopsis() string {
	return "List audit trail events for an organization"
}
