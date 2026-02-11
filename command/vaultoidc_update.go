package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// VaultoidcUpdateCommand is a command to update a Vault OIDC configuration
type VaultoidcUpdateCommand struct {
	Meta
	id               string
	address          string
	roleName         string
	namespace        string
	jwtAuthPath      string
	tlsCACertificate string
	format           string
}

// Run executes the vaultoidc update command
func (c *VaultoidcUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("vaultoidc update")
	flags.StringVar(&c.id, "id", "", "Vault OIDC configuration ID (required)")
	flags.StringVar(&c.address, "address", "", "Vault instance address")
	flags.StringVar(&c.roleName, "role", "", "Vault JWT auth role name")
	flags.StringVar(&c.namespace, "namespace", "", "Vault namespace")
	flags.StringVar(&c.jwtAuthPath, "auth-path", "", "Vault JWT auth mount path")
	flags.StringVar(&c.tlsCACertificate, "encoded-cacert", "", "Base64-encoded CA certificate")
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
	options := tfe.VaultOIDCConfigurationUpdateOptions{}

	if c.address != "" {
		options.Address = &c.address
	}

	if c.roleName != "" {
		options.RoleName = &c.roleName
	}

	if c.namespace != "" {
		options.Namespace = &c.namespace
	}

	if c.jwtAuthPath != "" {
		options.JWTAuthPath = &c.jwtAuthPath
	}

	if c.tlsCACertificate != "" {
		options.TLSCACertificate = &c.tlsCACertificate
	}

	// Update Vault OIDC configuration
	config, err := client.VaultOIDCConfigurations.Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating Vault OIDC configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Vault OIDC configuration '%s' updated successfully", config.ID))

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

// Help returns help text for the vaultoidc update command
func (c *VaultoidcUpdateCommand) Help() string {
	helpText := `
Usage: hcptf vaultoidc update [options]

  Update Vault OIDC configuration settings.

  Updates the Vault address, role, namespace, auth path, CA certificate,
  or audience settings for an existing Vault OIDC configuration.

Options:

  -id=<id>                  Vault OIDC configuration ID (required)
  -address=<url>           Vault instance address
                           Format: https://vault.example.com:8200
  -role=<name>             Vault JWT auth role name
  -namespace=<name>        Vault namespace
  -auth-path=<path>        Vault JWT auth mount path
  -encoded-cacert=<cert>   Base64-encoded CA certificate
  -output=<format>         Output format: table (default) or json

Example:

  hcptf vaultoidc update -id=voidc-ABC123 \
    -address=https://vault.example.com:8200

  hcptf vaultoidc update -id=voidc-ABC123 \
    -role=new-terraform-role \
    -namespace=admin

  hcptf vaultoidc update -id=voidc-ABC123 \
    -auth-path=jwt-auth
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the vaultoidc update command
func (c *VaultoidcUpdateCommand) Synopsis() string {
	return "Update Vault OIDC configuration settings"
}
