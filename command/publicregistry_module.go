package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// PublicRegistryModuleCommand queries the public Terraform registry for module info
type PublicRegistryModuleCommand struct {
	Meta
	module string
	format string
}

// ModuleInfo represents public registry module information
type ModuleInfo struct {
	ID          string `json:"id"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Downloads   int    `json:"downloads"`
	Published   string `json:"published_at"`
	Verified    bool   `json:"verified"`
}

// Run executes the command
func (c *PublicRegistryModuleCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("publicregistry module")
	flags.StringVar(&c.module, "name", "", "Module name (e.g., terraform-aws-modules/s3-bucket/aws)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.module == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Parse namespace/name/system
	parts := strings.Split(c.module, "/")
	if len(parts) != 3 {
		c.Ui.Error("Error: module name must be in format namespace/name/system (e.g., terraform-aws-modules/vpc/aws)")
		return 1
	}
	namespace := parts[0]
	name := parts[1]
	system := parts[2]

	// Query public registry API
	registryURL := fmt.Sprintf("https://registry.terraform.io/v1/modules/%s/%s/%s", namespace, name, system)

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

	var moduleInfo ModuleInfo
	if err := json.Unmarshal(body, &moduleInfo); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"Name":        fmt.Sprintf("%s/%s/%s", namespace, name, system),
		"Version":     moduleInfo.Version,
		"Description": moduleInfo.Description,
		"Source":      moduleInfo.Source,
		"Downloads":   moduleInfo.Downloads,
		"Published":   moduleInfo.Published,
		"Verified":    moduleInfo.Verified,
		"DocsURL":     fmt.Sprintf("https://registry.terraform.io/modules/%s/%s/%s/latest", namespace, name, system),
		"VersionsURL": fmt.Sprintf("https://registry.terraform.io/modules/%s/%s/%s?tab=versions", namespace, name, system),
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text
func (c *PublicRegistryModuleCommand) Help() string {
	helpText := `
Usage: hcptf publicregistry module [options]

  Get information about a module from the public Terraform registry.

  This command queries registry.terraform.io to get the latest version,
  description, download count, and documentation links for a public module.

Options:

  -name=<module>    Module name in format namespace/name/system (required)
                    Examples: terraform-aws-modules/vpc/aws,
                             hashicorp/dir/template,
                             terraform-aws-modules/s3-bucket/aws
  -output=<format>  Output format: table (default) or json

Examples:

  # Get VPC module info
  hcptf publicregistry module -name=terraform-aws-modules/vpc/aws

  # Get S3 bucket module info
  hcptf publicregistry module -name=terraform-aws-modules/s3-bucket/aws

  # Get HashiCorp dir template module
  hcptf publicregistry module -name=hashicorp/dir/template

Output includes:
  - Latest version available
  - Module description
  - Source repository
  - Total downloads
  - Publication date
  - Verified status (official/partner modules)
  - Documentation URL (with inputs, outputs, examples)
  - Versions page URL

Use these URLs to:
  - Review module changelog and version history
  - Check inputs/outputs before upgrading
  - Find usage examples
  - Review README and documentation
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis
func (c *PublicRegistryModuleCommand) Synopsis() string {
	return "Get module info from public registry"
}
