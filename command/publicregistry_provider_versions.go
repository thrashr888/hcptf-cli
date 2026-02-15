package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// PublicRegistryProviderVersionsCommand lists provider versions from public registry
type PublicRegistryProviderVersionsCommand struct {
	Meta
	provider string
	format   string
}

// ProviderVersionsResponse represents the registry API response
type ProviderVersionsResponse struct {
	Versions []struct {
		Version   string   `json:"version"`
		Protocols []string `json:"protocols"`
		Platforms []struct {
			OS   string `json:"os"`
			Arch string `json:"arch"`
		} `json:"platforms"`
	} `json:"versions"`
}

// Run executes the command
func (c *PublicRegistryProviderVersionsCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("publicregistry provider versions")
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

	// Query registry API for versions
	registryURL := fmt.Sprintf("https://registry.terraform.io/v1/providers/%s/%s/versions", namespace, name)

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

	var versionsResp ProviderVersionsResponse
	if err := json.Unmarshal(body, &versionsResp); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	if len(versionsResp.Versions) == 0 {
		c.Ui.Output("No versions found")
		return 0
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if c.format == "json" {
		jsonData, _ := json.MarshalIndent(versionsResp, "", "  ")
		c.Ui.Output(string(jsonData))
		return 0
	}

	// Table format - show latest 20 versions
	headers := []string{"Version", "Protocols", "Platforms"}
	var rows [][]string

	// Show newest first (registry returns newest first)
	count := len(versionsResp.Versions)
	if count > 20 {
		count = 20
	}

	for i := 0; i < count; i++ {
		ver := versionsResp.Versions[i]
		protocols := strings.Join(ver.Protocols, ", ")

		platformStrs := make([]string, 0, len(ver.Platforms))
		for _, p := range ver.Platforms {
			platformStrs = append(platformStrs, fmt.Sprintf("%s/%s", p.OS, p.Arch))
		}
		platforms := strings.Join(platformStrs, ", ")
		if len(platforms) > 50 {
			platforms = platforms[:47] + "..."
		}

		rows = append(rows, []string{ver.Version, protocols, platforms})
	}

	c.Ui.Output(fmt.Sprintf("Showing %d of %d versions for %s/%s\n", count, len(versionsResp.Versions), namespace, name))
	formatter.Table(headers, rows)

	// Show latest version highlight
	if len(versionsResp.Versions) > 0 {
		c.Ui.Output(fmt.Sprintf("\nLatest version: %s", versionsResp.Versions[0].Version))
		c.Ui.Output(fmt.Sprintf("Documentation: https://registry.terraform.io/providers/%s/%s/latest/docs", namespace, name))
	}

	return 0
}

// Help returns help text
func (c *PublicRegistryProviderVersionsCommand) Help() string {
	helpText := `
Usage: hcptf publicregistry provider versions [options]

  List all available versions of a provider from the public Terraform registry.

  Versions are listed newest first, showing supported protocols and platforms.
  Use this to find the latest version when planning upgrades.

Options:

  -name=<provider>  Provider name in format namespace/name (required)
                    Examples: hashicorp/aws, hashicorp/google, integrations/github
  -output=<format>  Output format: table (default) or json

Examples:

  # List AWS provider versions
  hcptf publicregistry provider versions -name=hashicorp/aws

  # List all versions in JSON format
  hcptf publicregistry provider versions -name=hashicorp/random -output=json

  # Check available versions for GitHub provider
  hcptf publicregistry provider versions -name=integrations/github

Output shows:
  - Version number (semantic versioning)
  - Supported protocols (5.0, 6.0, etc.)
  - Available platforms (linux/amd64, darwin/arm64, etc.)
  - Latest version highlighted
  - Documentation URL

Use with version upgrade workflow:
  1. Check current version: hcptf explorer query -type=providers
  2. Find latest version: hcptf publicregistry provider versions -name=hashicorp/aws
  3. Review docs and upgrade guides at the documentation URL
  4. Update version constraint in code
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis
func (c *PublicRegistryProviderVersionsCommand) Synopsis() string {
	return "List provider versions from public registry"
}
