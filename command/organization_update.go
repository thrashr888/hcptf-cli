package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// OrganizationUpdateCommand is a command to update an organization
type OrganizationUpdateCommand struct {
	Meta
	name                  string
	email                 string
	sessionTimeout        int
	sessionRemember       int
	costEstimationEnabled string
	format                string
	orgSvc                organizationUpdater
}

// Run executes the organization update command
func (c *OrganizationUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organization update")
	flags.StringVar(&c.name, "name", "", "Organization name (required)")
	flags.StringVar(&c.email, "email", "", "Admin email address")
	flags.IntVar(&c.sessionTimeout, "session-timeout", 0, "Session timeout in minutes")
	flags.IntVar(&c.sessionRemember, "session-remember", 0, "Session remember duration in minutes")
	flags.StringVar(&c.costEstimationEnabled, "cost-estimation", "", "Enable cost estimation (true/false)")
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

	// Build update options
	options := tfe.OrganizationUpdateOptions{}

	if c.email != "" {
		options.Email = tfe.String(c.email)
	}

	if c.sessionTimeout > 0 {
		options.SessionTimeout = tfe.Int(c.sessionTimeout)
	}

	if c.sessionRemember > 0 {
		options.SessionRemember = tfe.Int(c.sessionRemember)
	}

	if c.costEstimationEnabled != "" {
		if c.costEstimationEnabled == "true" {
			options.CostEstimationEnabled = tfe.Bool(true)
		} else if c.costEstimationEnabled == "false" {
			options.CostEstimationEnabled = tfe.Bool(false)
		} else {
			c.Ui.Error("Error: -cost-estimation must be 'true' or 'false'")
			return 1
		}
	}

	// Update organization
	org, err := c.organizationService(client).Update(client.Context(), c.name, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating organization: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Organization '%s' updated successfully", org.Name))

	// Show organization details
	data := map[string]interface{}{
		"Name":                  org.Name,
		"Email":                 org.Email,
		"SessionTimeout":        org.SessionTimeout,
		"SessionRemember":       org.SessionRemember,
		"CostEstimationEnabled": org.CostEstimationEnabled,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *OrganizationUpdateCommand) organizationService(client *client.Client) organizationUpdater {
	if c.orgSvc != nil {
		return c.orgSvc
	}
	return client.Organizations
}

// Help returns help text for the organization update command
func (c *OrganizationUpdateCommand) Help() string {
	helpText := `
Usage: hcptf organization update [options]

  Update organization settings.

Options:

  -name=<name>               Organization name (required)
  -email=<email>             Admin email address
  -session-timeout=<mins>    Session timeout in minutes
  -session-remember=<mins>   Session remember duration in minutes
  -cost-estimation=<bool>    Enable cost estimation (true/false)
  -output=<format>           Output format: table (default) or json

Example:

  hcptf organization update -name=my-org -email=newemail@example.com
  hcptf organization update -name=my-org -cost-estimation=true
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization update command
func (c *OrganizationUpdateCommand) Synopsis() string {
	return "Update organization settings"
}
