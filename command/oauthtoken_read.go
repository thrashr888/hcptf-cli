package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// OAuthTokenReadCommand is a command to read OAuth token details
type OAuthTokenReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the oauthtoken read command
func (c *OAuthTokenReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("oauthtoken read")
	flags.StringVar(&c.id, "id", "", "OAuth token ID (required)")
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

	// Read OAuth token
	oauthToken, err := client.OAuthTokens.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading OAuth token: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                  oauthToken.ID,
		"ServiceProviderUser": oauthToken.ServiceProviderUser,
		"HasSSHKey":           oauthToken.HasSSHKey,
		"CreatedAt":           oauthToken.CreatedAt,
	}

	if oauthToken.OAuthClient != nil {
		data["OAuthClientID"] = oauthToken.OAuthClient.ID
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the oauthtoken read command
func (c *OAuthTokenReadCommand) Help() string {
	helpText := `
Usage: hcptf oauthtoken read [options]

  Read OAuth token details. Shows information about an authenticated
  VCS connection including the service provider user and SSH key status.

Options:

  -id=<id>          OAuth token ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf oauthtoken read -id=ot-hmAyP66qk2AMVdbJ
  hcptf oauthtoken read -id=ot-hmAyP66qk2AMVdbJ -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the oauthtoken read command
func (c *OAuthTokenReadCommand) Synopsis() string {
	return "Read OAuth token details"
}
