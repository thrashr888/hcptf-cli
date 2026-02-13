package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// VaultoidcCreateCommand is a command to create a Vault OIDC configuration
type VaultoidcCreateCommand struct {
	Meta
	organization     string
	address          string
	roleName         string
	namespace        string
	jwtAuthPath      string
	tlsCACertificate string
	format           string
}

// Run executes the vaultoidc create command
func (c *VaultoidcCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("vaultoidc create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.address, "address", "", "Vault instance address (required)")
	flags.StringVar(&c.roleName, "role", "", "Vault JWT auth role name (required)")
	flags.StringVar(&c.namespace, "namespace", "", "Vault namespace (required)")
	flags.StringVar(&c.jwtAuthPath, "auth-path", "jwt", "Vault JWT auth mount path (default: jwt)")
	flags.StringVar(&c.tlsCACertificate, "encoded-cacert", "", "Base64-encoded CA certificate (optional)")
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

	if c.address == "" {
		c.Ui.Error("Error: -address flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.roleName == "" {
		c.Ui.Error("Error: -role flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.namespace == "" {
		c.Ui.Error("Error: -namespace flag is required")
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
	options := tfe.VaultOIDCConfigurationCreateOptions{
		Address:   c.address,
		RoleName:  c.roleName,
		Namespace: c.namespace,
	}

	if c.jwtAuthPath != "" {
		options.JWTAuthPath = c.jwtAuthPath
	}

	if c.tlsCACertificate != "" {
		options.TLSCACertificate = c.tlsCACertificate
	}

	// Create Vault OIDC configuration
	config, err := client.VaultOIDCConfigurations.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating Vault OIDC configuration: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Vault OIDC configuration created successfully with ID: %s", config.ID))

	// Show configuration details
	data := map[string]interface{}{
		"ID":          config.ID,
		"Address":     config.Address,
		"RoleName":    config.RoleName,
		"Namespace":   config.Namespace,
		"JWTAuthPath": config.JWTAuthPath,
	}

	if config.TLSCACertificate != "" {
		data["TLSCACertificate"] = "[set]"
	}

	if config.Organization != nil {
		data["Organization"] = config.Organization.Name
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the vaultoidc create command
func (c *VaultoidcCreateCommand) Help() string {
	helpText := `
Usage: hcptf vaultoidc create [options]

  Create a new Vault OIDC configuration for dynamic Vault credentials.

  Vault OIDC configurations enable HCP Terraform to dynamically generate
  Vault tokens using OpenID Connect. This eliminates the need to store
  static Vault tokens in HCP Terraform.

  Prerequisites:
  - Vault instance with JWT auth method enabled
  - JWT auth role configured with appropriate policies
  - Vault must be accessible from HCP Terraform

Options:

  -organization=<name>      Organization name (required)
  -org=<name>              Alias for -organization
  -address=<url>           Vault instance address (required)
                           Format: https://vault.example.com:8200
  -role=<name>             Vault JWT auth role name (required)
  -namespace=<name>        Vault namespace (required)
  -auth-path=<path>        Vault JWT auth mount path (default: jwt)
  -encoded-cacert=<cert>   Base64-encoded CA certificate (optional)
                           Only needed for self-signed certificates
  -output=<format>         Output format: table (default) or json

Example:

  hcptf vaultoidc create -org=my-org \
    -address=https://vault.example.com:8200 \
    -role=terraform-role \
    -namespace=admin \
    -auth-path=jwt

  hcptf vaultoidc create -org=my-org \
    -address=https://my-vault-cluster.vault.cloud:8200 \
    -role=terraform-role \
    -namespace=admin \
    -auth-path=jwt-auth
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the vaultoidc create command
func (c *VaultoidcCreateCommand) Synopsis() string {
	return "Create a Vault OIDC configuration for dynamic credentials"
}
