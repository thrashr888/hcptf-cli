package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// VaultoidcReadCommand is a command to read Vault OIDC configuration details
type VaultoidcReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the vaultoidc read command
func (c *VaultoidcReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("vaultoidc read")
	flags.StringVar(&c.id, "id", "", "Vault OIDC configuration ID (required)")
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

	// Read Vault OIDC configuration
	config, err := client.VaultOIDCConfigurations.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading Vault OIDC configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

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

// Help returns help text for the vaultoidc read command
func (c *VaultoidcReadCommand) Help() string {
	helpText := `
Usage: hcptf vaultoidc read [options]

  Read Vault OIDC configuration details.

  Displays the configuration details for a Vault OIDC configuration,
  including the Vault address, role, namespace, auth path, and audience settings.

Options:

  -id=<id>          Vault OIDC configuration ID (required)
                    Format: voidc-XXXXXXXXXX
  -output=<format>  Output format: table (default) or json

Example:

  hcptf vaultoidc read -id=voidc-ABC123
  hcptf vaultoidc read -id=voidc-ABC123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the vaultoidc read command
func (c *VaultoidcReadCommand) Synopsis() string {
	return "Read Vault OIDC configuration details"
}
