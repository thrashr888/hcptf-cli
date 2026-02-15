package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// PublicRegistryPolicyListCommand lists policies from public registry
type PublicRegistryPolicyListCommand struct {
	Meta
	format string
	limit  int
}

// PolicyListResponse represents the registry API response
type PolicyListResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Title     string `json:"title"`
			Name      string `json:"name"`
			Downloads int    `json:"downloads"`
		} `json:"attributes"`
		Relationships struct {
			LatestVersion struct {
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
			} `json:"latest-version"`
		} `json:"relationships"`
	} `json:"data"`
}

// Run executes the command
func (c *PublicRegistryPolicyListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("publicregistry policy list")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")
	flags.IntVar(&c.limit, "limit", 20, "Maximum number of policies to display")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Query registry API for policies
	registryURL := "https://registry.terraform.io/v2/policies?page[size]=100&include=latest-version"

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

	var policyList PolicyListResponse
	if err := json.Unmarshal(body, &policyList); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	if len(policyList.Data) == 0 {
		c.Ui.Output("No policies found")
		return 0
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if c.format == "json" {
		jsonData, _ := json.MarshalIndent(policyList, "", "  ")
		c.Ui.Output(string(jsonData))
		return 0
	}

	// Table format
	count := len(policyList.Data)
	if count > c.limit {
		count = c.limit
	}

	headers := []string{"Name", "Title", "Downloads", "Latest Version"}
	var rows [][]string

	for i := 0; i < count; i++ {
		policy := policyList.Data[i]
		// Extract version from related link (format: /v2/policies/namespace/name/version)
		versionPath := policy.Relationships.LatestVersion.Links.Related
		versionParts := strings.Split(strings.TrimPrefix(versionPath, "/v2/policies/"), "/")
		version := ""
		if len(versionParts) == 3 {
			version = versionParts[2]
		}

		rows = append(rows, []string{
			policy.Attributes.Name,
			policy.Attributes.Title,
			fmt.Sprintf("%d", policy.Attributes.Downloads),
			version,
		})
	}

	c.Ui.Output(fmt.Sprintf("Showing %d of %d policies\n", count, len(policyList.Data)))
	formatter.Table(headers, rows)

	return 0
}

// Help returns help text
func (c *PublicRegistryPolicyListCommand) Help() string {
	helpText := `
Usage: hcptf publicregistry policy list [options]

  List policies available in the public Terraform registry.

  This command queries registry.terraform.io to list public Sentinel and OPA
  policies, showing their latest versions and download counts.

Options:

  -output=<format>  Output format: table (default) or json
  -limit=<number>   Maximum number of policies to display (default: 20)

Examples:

  # List available public policies
  hcptf publicregistry policy list

  # List all policies in JSON format
  hcptf publicregistry policy list -output=json

  # Show first 10 policies
  hcptf publicregistry policy list -limit=10

Output shows:
  - Policy name (identifier)
  - Title (description)
  - Total downloads
  - Latest version available

Use these to find policies for compliance and governance:
  - CIS benchmarks (AWS, Azure, GCP)
  - Security best practices
  - Cost optimization policies
  - Tagging standards
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis
func (c *PublicRegistryPolicyListCommand) Synopsis() string {
	return "List policies from public registry"
}
