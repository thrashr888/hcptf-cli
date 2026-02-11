package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// OAuthClientReadCommand is a command to read OAuth client details
type OAuthClientReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the oauthclient read command
func (c *OAuthClientReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("oauthclient read")
	flags.StringVar(&c.id, "id", "", "OAuth client ID (required)")
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

	// Read OAuth client
	oauthClient, err := client.OAuthClients.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading OAuth client: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	name := ""
	if oauthClient.Name != nil && *oauthClient.Name != "" {
		name = *oauthClient.Name
	} else {
		name = oauthClient.ServiceProviderName
	}

	data := map[string]interface{}{
		"ID":                  oauthClient.ID,
		"Name":                name,
		"ServiceProvider":     string(oauthClient.ServiceProvider),
		"ServiceProviderName": oauthClient.ServiceProviderName,
		"HTTPURL":             oauthClient.HTTPURL,
		"APIURL":              oauthClient.APIURL,
		"CallbackURL":         oauthClient.CallbackURL,
		"ConnectPath":         oauthClient.ConnectPath,
		"CreatedAt":           oauthClient.CreatedAt,
	}

	if oauthClient.OrganizationScoped != nil {
		data["OrganizationScoped"] = *oauthClient.OrganizationScoped
	}

	if oauthClient.Organization != nil {
		data["Organization"] = oauthClient.Organization.Name
	}

	if oauthClient.Key != "" {
		data["Key"] = oauthClient.Key
	}

	if oauthClient.RSAPublicKey != "" {
		data["RSAPublicKey"] = oauthClient.RSAPublicKey
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the oauthclient read command
func (c *OAuthClientReadCommand) Help() string {
	helpText := `
Usage: hcptf oauthclient read [options]

  Read OAuth client details. Shows information about a VCS connection
  including service provider, URLs, and configuration.

Options:

  -id=<id>          OAuth client ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf oauthclient read -id=oc-XKFwG6ggfA9n7t1K
  hcptf oauthclient read -id=oc-XKFwG6ggfA9n7t1K -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the oauthclient read command
func (c *OAuthClientReadCommand) Synopsis() string {
	return "Read OAuth client details"
}
