package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// OAuthClientUpdateCommand is a command to update an OAuth client
type OAuthClientUpdateCommand struct {
	Meta
	id                 string
	name               string
	key                string
	secret             string
	rsaPublicKey       string
	oauthTokenString   string
	organizationScoped string
	format             string
}

// Run executes the oauthclient update command
func (c *OAuthClientUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("oauthclient update")
	flags.StringVar(&c.id, "id", "", "OAuth client ID (required)")
	flags.StringVar(&c.name, "name", "", "Display name for the OAuth client")
	flags.StringVar(&c.key, "key", "", "OAuth client key")
	flags.StringVar(&c.secret, "secret", "", "OAuth client secret")
	flags.StringVar(&c.rsaPublicKey, "rsa-public-key", "", "RSA public key")
	flags.StringVar(&c.oauthTokenString, "oauth-token-string", "", "New OAuth token string (for credential rotation)")
	flags.StringVar(&c.organizationScoped, "organization-scoped", "", "Whether OAuth client is scoped to all projects (true/false)")
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
	options := tfe.OAuthClientUpdateOptions{}

	if c.name != "" {
		options.Name = tfe.String(c.name)
	}

	if c.key != "" {
		options.Key = tfe.String(c.key)
	}

	if c.secret != "" {
		options.Secret = tfe.String(c.secret)
	}

	if c.rsaPublicKey != "" {
		options.RSAPublicKey = tfe.String(c.rsaPublicKey)
	}

	if c.oauthTokenString != "" {
		options.OAuthToken = tfe.String(c.oauthTokenString)
	}

	if c.organizationScoped != "" {
		if c.organizationScoped == "true" {
			options.OrganizationScoped = tfe.Bool(true)
		} else if c.organizationScoped == "false" {
			options.OrganizationScoped = tfe.Bool(false)
		} else {
			c.Ui.Error("Error: -organization-scoped must be 'true' or 'false'")
			return 1
		}
	}

	// Update OAuth client
	oauthClient, err := client.OAuthClients.Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating OAuth client: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("OAuth client '%s' updated successfully", oauthClient.ID))

	name := ""
	if oauthClient.Name != nil && *oauthClient.Name != "" {
		name = *oauthClient.Name
	} else {
		name = oauthClient.ServiceProviderName
	}

	// Show OAuth client details
	data := map[string]interface{}{
		"ID":              oauthClient.ID,
		"Name":            name,
		"ServiceProvider": string(oauthClient.ServiceProvider),
		"HTTPURL":         oauthClient.HTTPURL,
		"APIURL":          oauthClient.APIURL,
		"UpdatedAt":       oauthClient.CreatedAt,
	}

	if oauthClient.OrganizationScoped != nil {
		data["OrganizationScoped"] = *oauthClient.OrganizationScoped
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the oauthclient update command
func (c *OAuthClientUpdateCommand) Help() string {
	helpText := `
Usage: hcptf oauthclient update [options]

  Update OAuth client settings. Use this to rotate credentials,
  change display name, or modify scope settings.

Options:

  -id=<id>                    OAuth client ID (required)
  -name=<name>                Display name for the OAuth client
  -key=<key>                  OAuth client key
  -secret=<secret>            OAuth client secret
  -rsa-public-key=<key>       RSA public key
  -oauth-token-string=<token> New OAuth token (for credential rotation)
  -organization-scoped=<bool> Scope to all projects (true/false)
  -output=<format>            Output format: table (default) or json

Example:

  # Update OAuth client name
  hcptf oauthclient update -id=oc-XKFwG6ggfA9n7t1K -name="GitHub Production"

  # Rotate OAuth token
  hcptf oauthclient update -id=oc-XKFwG6ggfA9n7t1K -oauth-token-string=ghp_newtoken

  # Change organization scope
  hcptf oauthclient update -id=oc-XKFwG6ggfA9n7t1K -organization-scoped=false
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the oauthclient update command
func (c *OAuthClientUpdateCommand) Synopsis() string {
	return "Update OAuth client settings"
}
