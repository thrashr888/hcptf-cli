package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// OrganizationCreateCommand is a command to create an organization
type OrganizationCreateCommand struct {
	Meta
	name   string
	email  string
	format string
}

// Run executes the organization create command
func (c *OrganizationCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organization create")
	flags.StringVar(&c.name, "name", "", "Organization name (required)")
	flags.StringVar(&c.email, "email", "", "Admin email address (required)")
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

	if c.email == "" {
		c.Ui.Error("Error: -email flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create organization
	options := tfe.OrganizationCreateOptions{
		Name:  tfe.String(c.name),
		Email: tfe.String(c.email),
	}

	org, err := client.Organizations.Create(client.Context(), options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating organization: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Organization '%s' created successfully", org.Name))

	// Show organization details
	data := map[string]interface{}{
		"Name":                   org.Name,
		"Email":                  org.Email,
		"CollaboratorAuthPolicy": org.CollaboratorAuthPolicy,
		"CostEstimationEnabled":  org.CostEstimationEnabled,
		"SessionTimeout":         org.SessionTimeout,
		"SessionRemember":        org.SessionRemember,
		"TwoFactorConformant":    org.TwoFactorConformant,
		"CreatedAt":              org.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the organization create command
func (c *OrganizationCreateCommand) Help() string {
	helpText := `
Usage: hcptf organization create [options]

  Create a new organization.

Options:

  -name=<name>      Organization name (required)
  -email=<email>    Admin email address (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf organization create -name=my-org -email=admin@example.com
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organization create command
func (c *OrganizationCreateCommand) Synopsis() string {
	return "Create a new organization"
}
