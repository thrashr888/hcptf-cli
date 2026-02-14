package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// OrganizationShowCommand is a command to show organization details
type OrganizationShowCommand struct {
	Meta
	name   string
	format string
	orgSvc organizationReader
}

// Run executes the organization show command
func (c *OrganizationShowCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organization show")
	flags.StringVar(&c.name, "name", "", "Organization name (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read organization
	org, err := c.orgService(client).Read(client.Context(), c.name)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading organization: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"Name":                   org.Name,
		"Email":                  org.Email,
		"ExternalID":             org.ExternalID,
		"CollaboratorAuthPolicy": org.CollaboratorAuthPolicy,
		"CostEstimationEnabled":  org.CostEstimationEnabled,
		"SessionTimeout":         org.SessionTimeout,
		"SessionRemember":        org.SessionRemember,
		"TwoFactorConformant":    org.TwoFactorConformant,
		"Permissions":            org.Permissions,
		"SAMLEnabled":            org.SAMLEnabled,
		"CreatedAt":              org.CreatedAt,
		"TrialExpiresAt":         org.TrialExpiresAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the organization show command
func (c *OrganizationShowCommand) Help() string {
	helpText := `
Usage: hcptf organization show [options]

  Show organization details.

Options:

  -name=<name>      Organization name (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf organization show -name=my-org
  hcptf organization show -name=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

func (c *OrganizationShowCommand) orgService(client *client.Client) organizationReader {
	if c.orgSvc != nil {
		return c.orgSvc
	}
	return client.Organizations
}

// Synopsis returns a short synopsis for the organization show command
func (c *OrganizationShowCommand) Synopsis() string {
	return "Show organization details"
}
