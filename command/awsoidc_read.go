package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// AWSoidcReadCommand is a command to read AWS OIDC configuration details
type AWSoidcReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the awsoidc read command
func (c *AWSoidcReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("awsoidc read")
	flags.StringVar(&c.id, "id", "", "AWS OIDC configuration ID (required)")
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

	// Read AWS OIDC configuration
	config, err := client.AWSOIDCConfigurations.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading AWS OIDC configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":      config.ID,
		"RoleARN": config.RoleARN,
	}

	if config.Organization != nil {
		data["Organization"] = config.Organization.Name
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the awsoidc read command
func (c *AWSoidcReadCommand) Help() string {
	helpText := `
Usage: hcptf awsoidc read [options]

  Read AWS OIDC configuration details.

  Displays the configuration details for an AWS OIDC configuration,
  including the IAM role ARN and audience settings.

Options:

  -id=<id>          AWS OIDC configuration ID (required)
                    Format: awsoidc-XXXXXXXXXX
  -output=<format>  Output format: table (default) or json

Example:

  hcptf awsoidc read -id=awsoidc-ABC123
  hcptf awsoidc read -id=awsoidc-ABC123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the awsoidc read command
func (c *AWSoidcReadCommand) Synopsis() string {
	return "Read AWS OIDC configuration details"
}
