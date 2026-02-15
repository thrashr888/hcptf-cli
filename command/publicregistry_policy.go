package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// PublicRegistryPolicyCommand queries the public Terraform registry for policy info
type PublicRegistryPolicyCommand struct {
	Meta
	policy  string
	version string
	format  string
}

// PolicyInfo represents public registry policy information
type PolicyInfo struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Readme string `json:"readme"`
		} `json:"attributes"`
	} `json:"data"`
	Included []struct {
		Type       string `json:"type"`
		Attributes struct {
			Name   string `json:"name"`
			Shasum string `json:"shasum"`
		} `json:"attributes"`
	} `json:"included"`
}

// Run executes the command
func (c *PublicRegistryPolicyCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("publicregistry policy")
	flags.StringVar(&c.policy, "name", "", "Policy name (e.g., hashicorp/CIS-Policy-Set-for-AWS-Terraform)")
	flags.StringVar(&c.version, "version", "", "Policy version (e.g., 1.0.1). If not specified, uses 'latest'")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.policy == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Parse namespace/name
	parts := strings.Split(c.policy, "/")
	if len(parts) != 2 {
		c.Ui.Error("Error: policy name must be in format namespace/name (e.g., hashicorp/CIS-Policy-Set-for-AWS-Terraform)")
		return 1
	}
	namespace := parts[0]
	name := parts[1]

	// If no version specified, search for the policy to get latest version
	version := c.version
	if version == "" {
		// Search for policy to get latest version
		searchURL := "https://registry.terraform.io/v2/policies?page%5Bsize%5D=100&include=latest-version"
		searchResp, err := http.Get(searchURL)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error searching for policy: %s", err))
			return 1
		}
		defer searchResp.Body.Close()

		if searchResp.StatusCode != http.StatusOK {
			c.Ui.Error(fmt.Sprintf("Error searching for policy: status %d", searchResp.StatusCode))
			return 1
		}

		searchBody, err := io.ReadAll(searchResp.Body)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading search response: %s", err))
			return 1
		}

		var searchResult PolicyListResponse
		if err := json.Unmarshal(searchBody, &searchResult); err != nil {
			c.Ui.Error(fmt.Sprintf("Error parsing search response: %s", err))
			return 1
		}

		// Find matching policy
		found := false
		for _, policy := range searchResult.Data {
			if policy.Attributes.Name == name {
				// Extract version from related link
				versionPath := policy.Relationships.LatestVersion.Links.Related
				versionParts := strings.Split(strings.TrimPrefix(versionPath, "/v2/policies/"), "/")
				if len(versionParts) == 3 {
					version = versionParts[2]
					found = true
					break
				}
			}
		}

		if !found {
			c.Ui.Error(fmt.Sprintf("Policy %s/%s not found in registry", namespace, name))
			return 1
		}
	}

	// Query public registry API
	registryURL := fmt.Sprintf("https://registry.terraform.io/v2/policies/%s/%s/%s?include=policies,policy-modules,policy-library", namespace, name, version)

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

	var policyInfo PolicyInfo
	if err := json.Unmarshal(body, &policyInfo); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	// Extract policies and modules from included section
	var policies []string
	var modules []string
	for _, item := range policyInfo.Included {
		if item.Type == "policies" {
			policies = append(policies, item.Attributes.Name)
		} else if item.Type == "policy-modules" {
			modules = append(modules, item.Attributes.Name)
		}
	}

	data := map[string]interface{}{
		"Name":         fmt.Sprintf("%s/%s", namespace, name),
		"Version":      version,
		"Policies":     strings.Join(policies, ", "),
		"Modules":      strings.Join(modules, ", "),
		"PolicyCount":  len(policies),
		"ModuleCount":  len(modules),
		"DocsURL":      fmt.Sprintf("https://registry.terraform.io/policies/%s/%s/%s", namespace, name, version),
		"VersionsURL":  fmt.Sprintf("https://registry.terraform.io/policies/%s/%s", namespace, name),
	}

	// Add readme if present (truncated for table view)
	if policyInfo.Data.Attributes.Readme != "" {
		readme := policyInfo.Data.Attributes.Readme
		if c.format == "table" && len(readme) > 200 {
			readme = readme[:197] + "..."
		}
		data["Readme"] = readme
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text
func (c *PublicRegistryPolicyCommand) Help() string {
	helpText := `
Usage: hcptf publicregistry policy [options]

  Get information about a policy from the public Terraform registry.

  This command queries registry.terraform.io to get details about a public
  Sentinel or OPA policy, including the policies and modules it contains.

Options:

  -name=<policy>    Policy name in format namespace/name (required)
                    Examples: hashicorp/CIS-Policy-Set-for-AWS-Terraform,
                             hashicorp/CIS-Policy-Set-for-Azure-Terraform
  -version=<ver>    Policy version (default: latest)
  -output=<format>  Output format: table (default) or json

Examples:

  # Get AWS CIS policy set info
  hcptf publicregistry policy -name=hashicorp/CIS-Policy-Set-for-AWS-Terraform

  # Get specific version of Azure CIS policy
  hcptf publicregistry policy -name=hashicorp/CIS-Policy-Set-for-Azure-Terraform -version=1.0.0

  # Get GCP CIS policy in JSON format
  hcptf publicregistry policy -name=hashicorp/CIS-Policy-Set-for-GCP-Terraform -output=json

Output includes:
  - Policy set version
  - Number of policies included
  - Number of modules included
  - Policy and module names
  - Documentation URL
  - Versions page URL

Use these URLs to:
  - Review policy details and configuration
  - Check version history and changelogs
  - Find usage instructions and examples
  - Review compliance mappings (CIS benchmarks, etc.)
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis
func (c *PublicRegistryPolicyCommand) Synopsis() string {
	return "Get policy info from public registry"
}
