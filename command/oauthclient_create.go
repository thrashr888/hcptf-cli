package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// OAuthClientCreateCommand is a command to create an OAuth client
type OAuthClientCreateCommand struct {
	Meta
	organization       string
	serviceProvider    string
	name               string
	httpURL            string
	apiURL             string
	oauthTokenString   string
	key                string
	secret             string
	privateKey         string
	rsaPublicKey       string
	organizationScoped string
	format             string
}

// Run executes the oauthclient create command
func (c *OAuthClientCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("oauthclient create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.serviceProvider, "service-provider", "", "VCS provider: github, github_enterprise, gitlab_hosted, gitlab_community_edition, gitlab_enterprise_edition, ado_server (required)")
	flags.StringVar(&c.name, "name", "", "Display name for the OAuth client")
	flags.StringVar(&c.httpURL, "http-url", "", "VCS provider HTTP URL (required)")
	flags.StringVar(&c.apiURL, "api-url", "", "VCS provider API URL (required)")
	flags.StringVar(&c.oauthTokenString, "oauth-token-string", "", "OAuth token string from VCS provider (required)")
	flags.StringVar(&c.key, "key", "", "OAuth client key")
	flags.StringVar(&c.secret, "secret", "", "OAuth client secret")
	flags.StringVar(&c.privateKey, "private-key", "", "SSH private key (required for Azure DevOps Server)")
	flags.StringVar(&c.rsaPublicKey, "rsa-public-key", "", "RSA public key (required for Bitbucket Data Center)")
	flags.StringVar(&c.organizationScoped, "organization-scoped", "true", "Whether OAuth client is scoped to all projects (true/false)")
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

	if c.serviceProvider == "" {
		c.Ui.Error("Error: -service-provider flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.httpURL == "" {
		c.Ui.Error("Error: -http-url flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.apiURL == "" {
		c.Ui.Error("Error: -api-url flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.oauthTokenString == "" {
		c.Ui.Error("Error: -oauth-token-string flag is required")
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
	sp := tfe.ServiceProviderType(c.serviceProvider)
	options := tfe.OAuthClientCreateOptions{
		ServiceProvider: &sp,
		HTTPURL:         tfe.String(c.httpURL),
		APIURL:          tfe.String(c.apiURL),
		OAuthToken:      tfe.String(c.oauthTokenString),
	}

	if c.name != "" {
		options.Name = tfe.String(c.name)
	}

	if c.key != "" {
		options.Key = tfe.String(c.key)
	}

	if c.secret != "" {
		options.Secret = tfe.String(c.secret)
	}

	if c.privateKey != "" {
		options.PrivateKey = tfe.String(c.privateKey)
	}

	if c.rsaPublicKey != "" {
		options.RSAPublicKey = tfe.String(c.rsaPublicKey)
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

	// Create OAuth client
	oauthClient, err := client.OAuthClients.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating OAuth client: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("OAuth client '%s' created successfully", oauthClient.ID))

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
		"CallbackURL":     oauthClient.CallbackURL,
		"ConnectPath":     oauthClient.ConnectPath,
		"CreatedAt":       oauthClient.CreatedAt,
	}

	if oauthClient.OrganizationScoped != nil {
		data["OrganizationScoped"] = *oauthClient.OrganizationScoped
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the oauthclient create command
func (c *OAuthClientCreateCommand) Help() string {
	helpText := `
Usage: hcptf oauthclient create [options]

  Create a new OAuth client for VCS integration. This establishes
  a connection between the organization and a VCS provider
  (GitHub, GitLab, Azure DevOps, etc.).

Options:

  -organization=<name>         Organization name (required)
  -org=<name>                 Alias for -organization
  -service-provider=<type>    VCS provider type (required)
                              Valid values: github, github_enterprise,
                              gitlab_hosted, gitlab_community_edition,
                              gitlab_enterprise_edition, ado_server
  -http-url=<url>             VCS provider HTTP URL (required)
                              Example: https://github.com
  -api-url=<url>              VCS provider API URL (required)
                              Example: https://api.github.com
  -oauth-token-string=<token> OAuth token from VCS provider (required)
  -name=<name>                Display name for the OAuth client
  -key=<key>                  OAuth client key
  -secret=<secret>            OAuth client secret
  -private-key=<key>          SSH private key (for Azure DevOps Server)
  -rsa-public-key=<key>       RSA public key (for Bitbucket Data Center)
  -organization-scoped=<bool> Scope to all projects (default: true)
  -output=<format>            Output format: table (default) or json

Example:

  # Create GitHub OAuth client
  hcptf oauthclient create \
    -org=my-org \
    -service-provider=github \
    -http-url=https://github.com \
    -api-url=https://api.github.com \
    -oauth-token-string=ghp_xxxxxxxxxxxx

  # Create GitHub Enterprise OAuth client
  hcptf oauthclient create \
    -org=my-org \
    -service-provider=github_enterprise \
    -http-url=https://github.example.com \
    -api-url=https://github.example.com/api/v3 \
    -oauth-token-string=ghp_xxxxxxxxxxxx \
    -name="GitHub Enterprise"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the oauthclient create command
func (c *OAuthClientCreateCommand) Synopsis() string {
	return "Create a new OAuth client for VCS integration"
}
