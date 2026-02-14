package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// OAuthTokenUpdateCommand is a command to update an OAuth token
type OAuthTokenUpdateCommand struct {
	Meta
	id     string
	sshKey string
	format string
}

// Run executes the oauthtoken update command
func (c *OAuthTokenUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("oauthtoken update")
	flags.StringVar(&c.id, "id", "", "OAuth token ID (required)")
	flags.StringVar(&c.sshKey, "ssh-key", "", "SSH private key content")
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

	if c.sshKey == "" {
		c.Ui.Error("Error: -ssh-key flag is required")
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
	options := tfe.OAuthTokenUpdateOptions{
		PrivateSSHKey: tfe.String(c.sshKey),
	}

	// Update OAuth token
	oauthToken, err := client.OAuthTokens.Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating OAuth token: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("OAuth token '%s' updated successfully", oauthToken.ID))

	// Show OAuth token details
	data := map[string]interface{}{
		"ID":                  oauthToken.ID,
		"ServiceProviderUser": oauthToken.ServiceProviderUser,
		"HasSSHKey":           oauthToken.HasSSHKey,
		"CreatedAt":           oauthToken.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the oauthtoken update command
func (c *OAuthTokenUpdateCommand) Help() string {
	helpText := `
Usage: hcptf oauthtoken update [options]

  Update OAuth token settings. This is typically used to add or
  update the SSH key associated with the OAuth token for accessing
  private Git submodules.

Options:

  -id=<id>       OAuth token ID (required)
  -ssh-key=<key> SSH private key content (required)
  -output=<fmt>  Output format: table (default) or json

Example:

  hcptf oauthtoken update -id=ot-hmAyP66qk2AMVdbJ -ssh-key="$(cat ~/.ssh/id_rsa)"
  hcptf oauthtoken update -id=ot-hmAyP66qk2AMVdbJ -ssh-key="-----BEGIN RSA PRIVATE KEY-----..."
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the oauthtoken update command
func (c *OAuthTokenUpdateCommand) Synopsis() string {
	return "Update OAuth token settings"
}
