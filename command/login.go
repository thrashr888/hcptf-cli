package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/config"
)

// LoginCommand is a command to authenticate and save credentials
type LoginCommand struct {
	Meta
	hostname  string
	showToken bool
}

// Run executes the login command
func (c *LoginCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("login")
	flags.StringVar(&c.hostname, "hostname", "app.terraform.io", "HCP Terraform hostname")
	flags.BoolVar(&c.showToken, "show-token", false, "Show token after successful login")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.showToken {
		cfg, err := c.Config()
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error loading config: %s", err))
			return 1
		}

		token := cfg.GetToken(c.hostname)
		if token == "" {
			c.Ui.Error("No token found for this hostname. Run 'hcptf login' first.")
			return 1
		}

		c.Ui.Output(token)
		return 0
	}

	c.Ui.Output(fmt.Sprintf("Authenticating to %s", c.hostname))
	c.Ui.Output("")
	c.Ui.Output("This command will store an API token in ~/.terraform.d/credentials.tfrc.json")
	c.Ui.Output("The token will be shared with the Terraform CLI.")
	c.Ui.Output("")

	// Get token from user
	c.Ui.Output("Generate a token at:")
	if c.hostname == "app.terraform.io" {
		c.Ui.Output("  https://app.terraform.io/app/settings/tokens")
	} else {
		c.Ui.Output(fmt.Sprintf("  https://%s/app/settings/tokens", c.hostname))
	}
	c.Ui.Output("")

	token, err := c.Ui.Ask("Enter your API token:")
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading token: %s", err))
		return 1
	}

	token = strings.TrimSpace(token)
	if token == "" {
		c.Ui.Error("Error: token cannot be empty")
		return 1
	}

	// Validate token by making a test API call
	c.Ui.Output("")
	c.Ui.Output("Validating token...")

	if err := config.ValidateToken(c.hostname, token); err != nil {
		c.Ui.Error(fmt.Sprintf("Error: token validation failed: %s", err))
		c.Ui.Error("")
		c.Ui.Error("Please verify:")
		c.Ui.Error("  1. The token is correct")
		c.Ui.Error("  2. The token has not expired")
		c.Ui.Error(fmt.Sprintf("  3. You can access %s", c.hostname))
		return 1
	}

	// Save token to credentials file
	if err := config.SaveCredential(c.hostname, token); err != nil {
		c.Ui.Error(fmt.Sprintf("Error saving credentials: %s", err))
		return 1
	}

	c.Ui.Output("")
	c.Ui.Output("Success! Credentials saved.")
	c.Ui.Output("")
	c.Ui.Output(fmt.Sprintf("Token stored in %s", config.GetTerraformCredentialsPath()))
	c.Ui.Output("")
	c.Ui.Output("You can now use hcptf commands without setting TFE_TOKEN or HCPTF_TOKEN.")
	c.Ui.Output("The Terraform CLI will also use these credentials.")

	return 0
}

// Help returns help text for the login command
func (c *LoginCommand) Help() string {
	helpText := `
Usage: hcptf login [options]

  Authenticate to HCP Terraform and save credentials.

  This command will prompt you for an API token and store it in
  ~/.terraform.d/credentials.tfrc.json, making it available to both
  hcptf and the Terraform CLI.

Options:

  -hostname=<hostname>  HCP Terraform hostname (default: app.terraform.io)
  -show-token          Show token and exit without prompting (default: false)

Example:

  # Login to HCP Terraform
  hcptf login

  # Login to Terraform Enterprise
  hcptf login -hostname=tfe.example.com

Note:

  The stored credentials are compatible with 'terraform login' and will
  be automatically used by both the Terraform CLI and hcptf.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the login command
func (c *LoginCommand) Synopsis() string {
	return "Authenticate to HCP Terraform"
}
