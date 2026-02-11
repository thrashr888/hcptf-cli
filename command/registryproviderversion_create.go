package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type RegistryProviderVersionCreateCommand struct {
	Meta
	organization string
	name         string
	namespace    string
	version      string
	keyID        string
	protocols    string
	format       string
}

// Run executes the registry provider version create command
func (c *RegistryProviderVersionCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("registryproviderversion create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Provider name (required)")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (defaults to organization)")
	flags.StringVar(&c.version, "version", "", "Version string, e.g., 1.0.0 (required)")
	flags.StringVar(&c.keyID, "key-id", "", "GPG key ID for signing (required)")
	flags.StringVar(&c.protocols, "protocols", "5.0,6.0", "Comma-separated protocol versions (default: 5.0,6.0)")
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

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.version == "" {
		c.Ui.Error("Error: -version flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.keyID == "" {
		c.Ui.Error("Error: -key-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Default namespace to organization
	if c.namespace == "" {
		c.namespace = c.organization
	}

	// Parse protocols
	protocolList := strings.Split(c.protocols, ",")

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create provider version
	providerID := tfe.RegistryProviderID{
		RegistryName: tfe.PrivateRegistry,
		Namespace:    c.namespace,
		Name:         c.name,
	}

	options := tfe.RegistryProviderVersionCreateOptions{
		Version:   c.version,
		KeyID:     c.keyID,
		Protocols: protocolList,
	}

	providerVersion, err := client.RegistryProviderVersions.Create(client.Context(), providerID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating provider version: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Provider version '%s' created successfully", providerVersion.Version))
	c.Ui.Output("\nNext steps:")
	c.Ui.Output("1. Upload SHA256SUMS file to the shasums-upload URL")
	c.Ui.Output("2. Upload SHA256SUMS.sig file to the shasums-sig-upload URL")
	c.Ui.Output("3. Create platform binaries with 'registryproviderplatform create'")

	if shasumsURL, ok := providerVersion.Links["shasums-upload"]; ok {
		c.Ui.Output(fmt.Sprintf("\nSHA256SUMS upload URL: %s", shasumsURL))
	}
	if shasumsSignURL, ok := providerVersion.Links["shasums-sig-upload"]; ok {
		c.Ui.Output(fmt.Sprintf("SHA256SUMS.sig upload URL: %s", shasumsSignURL))
	}

	// Show version details
	data := map[string]interface{}{
		"ID":        providerVersion.ID,
		"Version":   providerVersion.Version,
		"KeyID":     providerVersion.KeyID,
		"Protocols": strings.Join(providerVersion.Protocols, ", "),
		"CreatedAt": providerVersion.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the registry provider version create command
func (c *RegistryProviderVersionCreateCommand) Help() string {
	helpText := `
Usage: hcptf registryproviderversion create [options]

  Create a new version for a private registry provider.

  After creating the version, you must:
  1. Upload SHA256SUMS file
  2. Upload SHA256SUMS.sig file
  3. Create platform binaries

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Provider name (required)
  -namespace=<name>    Namespace (defaults to organization)
  -version=<semver>    Version string, e.g., 1.0.0 (required)
  -key-id=<id>         GPG key ID for signing (required)
  -protocols=<list>    Comma-separated protocol versions (default: 5.0,6.0)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf registryproviderversion create -org=my-org -name=aws -version=3.1.1 -key-id=32966F3FB5AC1129
  hcptf registryproviderversion create -org=my-org -name=custom -version=1.0.0 -key-id=ABCD1234 -protocols=5.0,6.0
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the registry provider version create command
func (c *RegistryProviderVersionCreateCommand) Synopsis() string {
	return "Create a new version for a private registry provider"
}
