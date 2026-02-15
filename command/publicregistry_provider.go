package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// PublicRegistryProviderCommand queries the public Terraform registry for provider info
type PublicRegistryProviderCommand struct {
	Meta
	provider string
	format   string
}

// ProviderInfo represents public registry provider information
type ProviderInfo struct {
	ID          string `json:"id"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Published   string `json:"published_at"`
}

// Run executes the public registry provider command
func (c *PublicRegistryProviderCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("publicregistry provider")
	flags.StringVar(&c.provider, "name", "", "Provider name (e.g., hashicorp/aws)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.provider == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Parse namespace/name
	parts := strings.Split(c.provider, "/")
	if len(parts) != 2 {
		c.Ui.Error("Error: provider name must be in format namespace/name (e.g., hashicorp/aws)")
		return 1
	}
	namespace := parts[0]
	name := parts[1]

	// Query public registry API
	registryURL := fmt.Sprintf("https://registry.terraform.io/v1/providers/%s/%s", namespace, name)

	resp, err := http.Get(registryURL)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error querying registry: %s", err))
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.Ui.Error(fmt.Sprintf("Registry API returned %d: %s", resp.StatusCode, string(body)))
		return 1
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading response: %s", err))
		return 1
	}

	var providerInfo ProviderInfo
	if err := json.Unmarshal(body, &providerInfo); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"Name":        fmt.Sprintf("%s/%s", namespace, name),
		"Version":     providerInfo.Version,
		"Description": providerInfo.Description,
		"Source":      providerInfo.Source,
		"Published":   providerInfo.Published,
		"DocsURL":     fmt.Sprintf("https://registry.terraform.io/providers/%s/%s/latest/docs", namespace, name),
		"VersionsURL": fmt.Sprintf("https://registry.terraform.io/providers/%s/%s/versions", namespace, name),
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text
func (c *PublicRegistryProviderCommand) Help() string {
	helpText := `
Usage: hcptf publicregistry provider [options]

  Get information about a provider from the public Terraform registry.

  This command queries registry.terraform.io to get the latest version,
  description, and documentation links for a public provider.

Options:

  -name=<provider>  Provider name in format namespace/name (required)
                    Examples: hashicorp/aws, hashicorp/google, integrations/github
  -output=<format>  Output format: table (default) or json

Examples:

  # Get latest AWS provider info
  hcptf publicregistry provider -name=hashicorp/aws

  # Get Google Cloud provider info
  hcptf publicregistry provider -name=hashicorp/google

  # Get GitHub provider info
  hcptf publicregistry provider -name=integrations/github

Output includes:
  - Latest version available
  - Provider description
  - Source repository
  - Documentation URL
  - Versions page URL

Use these URLs to:
  - Review upgrade guides and changelogs
  - Check for breaking changes before upgrading
  - Find documentation for specific resources
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis
func (c *PublicRegistryProviderCommand) Synopsis() string {
	return "Get provider info from public registry"
}
