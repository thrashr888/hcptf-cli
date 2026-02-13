package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// AWSoidcUpdateCommand is a command to update an AWS OIDC configuration
type AWSoidcUpdateCommand struct {
	Meta
	id      string
	roleArn string
	format  string
}

// Run executes the awsoidc update command
func (c *AWSoidcUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("awsoidc update")
	flags.StringVar(&c.id, "id", "", "AWS OIDC configuration ID (required)")
	flags.StringVar(&c.roleArn, "role-arn", "", "AWS IAM role ARN")
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

	// Build update options
	options := tfe.AWSOIDCConfigurationUpdateOptions{}

	if c.roleArn != "" {
		options.RoleARN = c.roleArn
	}

	// Update AWS OIDC configuration
	config, err := client.AWSOIDCConfigurations.Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating AWS OIDC configuration: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("AWS OIDC configuration '%s' updated successfully", config.ID))

	// Show configuration details
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

// Help returns help text for the awsoidc update command
func (c *AWSoidcUpdateCommand) Help() string {
	helpText := `
Usage: hcptf awsoidc update [options]

  Update AWS OIDC configuration settings.

  Updates the AWS IAM role ARN or audience settings for an existing
  AWS OIDC configuration.

Options:

  -id=<id>             AWS OIDC configuration ID (required)
  -role-arn=<arn>      AWS IAM role ARN to assume
                       Format: arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME
  -output=<format>     Output format: table (default) or json

Example:

  hcptf awsoidc update -id=awsoidc-ABC123 \
    -role-arn=arn:aws:iam::123456789012:role/new-terraform-role
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the awsoidc update command
func (c *AWSoidcUpdateCommand) Synopsis() string {
	return "Update AWS OIDC configuration settings"
}
