package command

import (
	"fmt"
	"os"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// AuditTrailReadCommand is a command to read audit trail event details
type AuditTrailReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the audit trail read command
func (c *AuditTrailReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("audittrail read")
	flags.StringVar(&c.id, "id", "", "Audit trail event ID (required)")
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

	// Note: The audit trails API doesn't have a direct "read by ID" endpoint,
	// so we need to list and filter. We'll fetch multiple pages if needed.
	options := &tfe.AuditTrailListOptions{
		ListOptions: &tfe.ListOptions{
			PageNumber: 1,
			PageSize:   1000,
		},
	}

	var auditTrail *tfe.AuditTrail
	for {
		auditTrails, err := client.AuditTrails.List(client.Context(), options)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error listing audit trail events: %s", err))
			return 1
		}

		// Search for the ID
		for _, at := range auditTrails.Items {
			if at.ID == c.id {
				auditTrail = at
				break
			}
		}

		if auditTrail != nil {
			break
		}

		// Check if there are more pages
		if auditTrails.AuditTrailPagination == nil || auditTrails.AuditTrailPagination.NextPage == 0 {
			break
		}

		options.ListOptions.PageNumber = auditTrails.AuditTrailPagination.NextPage
	}

	if auditTrail == nil {
		c.Ui.Error(fmt.Sprintf("Audit trail event with ID '%s' not found", c.id))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)
	if c.Meta.OutputWriter == nil && c.Meta.ErrorWriter == nil {
		formatter = output.NewFormatterWithWriters(c.format, os.Stdout, os.Stderr)
	}

	data := map[string]interface{}{
		"ID":        auditTrail.ID,
		"Version":   auditTrail.Version,
		"Type":      auditTrail.Type,
		"Timestamp": auditTrail.Timestamp.Format("2006-01-02 15:04:05 MST"),
	}

	// Auth information
	data["Auth.AccessorID"] = auditTrail.Auth.AccessorID
	data["Auth.Description"] = auditTrail.Auth.Description
	data["Auth.Type"] = auditTrail.Auth.Type
	data["Auth.OrganizationID"] = auditTrail.Auth.OrganizationID
	if auditTrail.Auth.ImpersonatorID != nil && *auditTrail.Auth.ImpersonatorID != "" {
		data["Auth.ImpersonatorID"] = *auditTrail.Auth.ImpersonatorID
	}

	// Request information
	data["Request.ID"] = auditTrail.Request.ID

	// Resource information
	data["Resource.ID"] = auditTrail.Resource.ID
	data["Resource.Type"] = auditTrail.Resource.Type
	data["Resource.Action"] = auditTrail.Resource.Action
	if auditTrail.Resource.Meta != nil {
		data["Resource.Meta"] = auditTrail.Resource.Meta
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the audit trail read command
func (c *AuditTrailReadCommand) Help() string {
	helpText := `
Usage: hcptf audit trail read [options]

  Read detailed information about a specific audit trail event by searching
  through recent audit logs. Note: This searches through the available audit
  trail history (14 days) to find the event.

  Note: Requires an organization token or audit trail token.

Options:

  -id=<id>          Audit trail event ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf audit trail read -id=ae66e491-db59-457c-8445-9c908ee726ae
  hcptf audit trail read -id=ae66e491-db59-457c-8445-9c908ee726ae -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the audit trail read command
func (c *AuditTrailReadCommand) Synopsis() string {
	return "Read audit trail event details"
}
