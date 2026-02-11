package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/config"
)

// LogoutCommand is a command to remove saved credentials
type LogoutCommand struct {
	Meta
	hostname string
}

// Run executes the logout command
func (c *LogoutCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("logout")
	flags.StringVar(&c.hostname, "hostname", "app.terraform.io", "HCP Terraform hostname")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Check if credentials exist
	credsPath := config.GetTerraformCredentialsPath()
	creds, err := config.LoadTerraformCredentialsFile()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("No credentials found at %s", credsPath))
		return 0
	}

	if _, exists := creds.Credentials[c.hostname]; !exists {
		c.Ui.Output(fmt.Sprintf("No credentials found for %s", c.hostname))
		c.Ui.Output(fmt.Sprintf("Credentials file: %s", credsPath))
		return 0
	}

	// Remove credentials for the specified hostname
	if err := config.RemoveCredential(c.hostname); err != nil {
		c.Ui.Error(fmt.Sprintf("Error removing credentials: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Successfully logged out of %s", c.hostname))
	c.Ui.Output(fmt.Sprintf("Credentials removed from %s", credsPath))

	return 0
}

// Help returns help text for the logout command
func (c *LogoutCommand) Help() string {
	helpText := `
Usage: hcptf logout [options]

  Remove saved credentials for HCP Terraform.

  This command removes stored API tokens from
  ~/.terraform.d/credentials.tfrc.json for the specified hostname.

Options:

  -hostname=<hostname>  HCP Terraform hostname (default: app.terraform.io)

Example:

  # Logout from HCP Terraform
  hcptf logout

  # Logout from Terraform Enterprise
  hcptf logout -hostname=tfe.example.com

Note:

  This only removes credentials from the credentials file. It does not
  revoke the token on the server. To revoke tokens, use the web UI:
  https://app.terraform.io/app/settings/tokens
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the logout command
func (c *LogoutCommand) Synopsis() string {
	return "Remove saved credentials"
}
