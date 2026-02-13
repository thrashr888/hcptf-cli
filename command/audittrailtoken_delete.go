package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// AuditTrailTokenDeleteCommand is a command to delete an audit trail token
type AuditTrailTokenDeleteCommand struct {
	Meta
	organization string
	force        bool
	orgTokenSvc  auditTrailTokenDeleter
}

// Run executes the audit trail token delete command
func (c *AuditTrailTokenDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("audittrailtoken delete")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")

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

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete the audit trail token for organization '%s'? (yes/no): ", c.organization))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete audit trail token
	tokenType := tfe.AuditTrailToken
	options := tfe.OrganizationTokenDeleteOptions{
		TokenType: &tokenType,
	}

	err = c.auditTrailTokenService(client).DeleteWithOptions(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting audit trail token: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Audit trail token for organization '%s' deleted successfully", c.organization))
	return 0
}

func (c *AuditTrailTokenDeleteCommand) auditTrailTokenService(client *client.Client) auditTrailTokenDeleter {
	if c.orgTokenSvc != nil {
		return c.orgTokenSvc
	}
	return client.OrganizationTokens
}

// Help returns help text for the audit trail token delete command
func (c *AuditTrailTokenDeleteCommand) Help() string {
	helpText := `
Usage: hcptf audittrailtoken delete [options]

  Delete the audit trail token for an organization. This will revoke access
  for any integrations using this token to pull audit trail data.

  Note: Only owners team members can access this command.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -force               Force delete without confirmation

Example:

  hcptf audittrailtoken delete -org=my-org
  hcptf audittrailtoken delete -org=my-org -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the audit trail token delete command
func (c *AuditTrailTokenDeleteCommand) Synopsis() string {
	return "Delete an audit trail token"
}
