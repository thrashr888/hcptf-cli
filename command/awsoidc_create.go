package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// AWSoidcCreateCommand is a command to create an AWS OIDC configuration
type AWSoidcCreateCommand struct {
	Meta
	organization string
	roleArn      string
	format       string
}

// Run executes the awsoidc create command
func (c *AWSoidcCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("awsoidc create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.roleArn, "role-arn", "", "AWS IAM role ARN (required)")
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

	if c.roleArn == "" {
		c.Ui.Error("Error: -role-arn flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build create options
	options := tfe.AWSOIDCConfigurationCreateOptions{
		RoleARN: c.roleArn,
	}

	// Create AWS OIDC configuration
	config, err := client.AWSOIDCConfigurations.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating AWS OIDC configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("AWS OIDC configuration created successfully with ID: %s", config.ID))

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

// Help returns help text for the awsoidc create command
func (c *AWSoidcCreateCommand) Help() string {
	helpText := `
Usage: hcptf awsoidc create [options]

  Create a new AWS OIDC configuration for dynamic AWS credentials.

  AWS OIDC configurations enable HCP Terraform to dynamically generate
  AWS credentials using OpenID Connect. This eliminates the need to store
  static AWS access keys in HCP Terraform.

  Prerequisites:
  - AWS IAM role configured with a trust relationship for HCP Terraform's OIDC provider
  - Role must have permissions to perform the required AWS operations

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -role-arn=<arn>      AWS IAM role ARN to assume (required)
                       Format: arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME
  -output=<format>     Output format: table (default) or json

Example:

  hcptf awsoidc create -org=my-org \
    -role-arn=arn:aws:iam::123456789012:role/terraform-role
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the awsoidc create command
func (c *AWSoidcCreateCommand) Synopsis() string {
	return "Create an AWS OIDC configuration for dynamic credentials"
}
