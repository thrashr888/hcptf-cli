package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// ExplorerQueryCommand queries the Explorer API
type ExplorerQueryCommand struct {
	Meta
	organization string
	queryType    string
	sort         string
	filter       string
	fields       string
	limit        int
	pageNumber   int
	pageSize     int
	format       string
	exportCSV    bool
}

// Run executes the explorer query command
func (c *ExplorerQueryCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("explorer query")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.queryType, "type", "", "Query type: workspaces, tf_versions, providers, modules (required)")
	flags.StringVar(&c.sort, "sort", "", "Sort field (prefix with - for descending)")
	flags.StringVar(&c.filter, "filter", "", "Filter conditions")
	flags.StringVar(&c.fields, "fields", "", "Comma-separated fields to return")
	flags.IntVar(&c.limit, "limit", 0, "Maximum number of results to return (overrides page-size)")
	flags.IntVar(&c.pageNumber, "page", 1, "Page number")
	flags.IntVar(&c.pageSize, "page-size", 20, "Page size")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")
	flags.BoolVar(&c.exportCSV, "csv", false, "Export as CSV instead of paginated JSON")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.queryType == "" {
		c.Ui.Error("Error: -type flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Apply limit if set (overrides page-size)
	if c.limit > 0 {
		c.pageSize = c.limit
	}

	// Validate query type
	validTypes := []string{"workspaces", "tf_versions", "providers", "modules"}
	valid := false
	for _, t := range validTypes {
		if c.queryType == t {
			valid = true
			break
		}
	}
	if !valid {
		c.Ui.Error(fmt.Sprintf("Error: invalid type %q, must be one of: %s", c.queryType, strings.Join(validTypes, ", ")))
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build URL - Explorer API is at /api/v2/organizations/:org/explorer
	baseURL := client.BaseURL()
	var endpoint string
	if c.exportCSV {
		endpoint = fmt.Sprintf("%s/api/v2/organizations/%s/explorer/export/csv", baseURL, c.organization)
	} else {
		endpoint = fmt.Sprintf("%s/api/v2/organizations/%s/explorer", baseURL, c.organization)
	}

	// Build query parameters
	params := url.Values{}
	params.Add("type", c.queryType)
	if c.sort != "" {
		params.Add("sort", c.sort)
	}
	if c.filter != "" {
		params.Add("filter", c.filter)
	}
	if c.fields != "" {
		params.Add("fields", c.fields)
	}
	if !c.exportCSV {
		params.Add("page[number]", fmt.Sprintf("%d", c.pageNumber))
		params.Add("page[size]", fmt.Sprintf("%d", c.pageSize))
	}

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	// Make HTTP request
	req, err := http.NewRequestWithContext(client.Context(), "GET", fullURL, nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating request: %s", err))
		return 1
	}

	// Add auth token
	token := client.Token()
	req.Header.Set("Authorization", "Bearer "+token)
	if !c.exportCSV {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error making request: %s", err))
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusNotFound {
			c.Ui.Error("Error: Explorer API not found (404)")
			c.Ui.Error(fmt.Sprintf("Response: %s", string(body)))
			c.Ui.Error("")
			c.Ui.Error("Possible causes:")
			c.Ui.Error(fmt.Sprintf("  - Organization '%s' does not exist", c.organization))
			c.Ui.Error("  - Explorer API is not available for this organization")
			c.Ui.Error("  - Explorer API may require HCP Terraform Plus or higher")
			c.Ui.Error("")
			c.Ui.Error(fmt.Sprintf("Verify your organization name with: hcptf organization show -name=%s", c.organization))
		} else {
			c.Ui.Error(fmt.Sprintf("Error: API returned %d: %s", resp.StatusCode, string(body)))
		}
		return 1
	}

	// Handle CSV export
	if c.exportCSV {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading response: %s", err))
			return 1
		}
		c.Ui.Output(string(data))
		return 0
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	// Format output
	if c.format == "json" {
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error formatting JSON: %s", err))
			return 1
		}
		c.Ui.Output(string(jsonData))
		return 0
	}

	// Table format
	formatter := output.NewFormatter(c.format)

	// Extract data array from JSON API response
	data, ok := result["data"].([]interface{})
	if !ok {
		c.Ui.Error("Error: unexpected response format")
		return 1
	}

	if len(data) == 0 {
		c.Ui.Output(fmt.Sprintf("No %s found", c.queryType))
		return 0
	}

	// Dynamic table based on fields
	var headers []string
	var rows [][]string

	for i, item := range data {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Get attributes
		attrs, ok := itemMap["attributes"].(map[string]interface{})
		if !ok {
			continue
		}

		// First item: build headers
		if i == 0 {
			if c.fields != "" {
				headers = strings.Split(c.fields, ",")
			} else {
				// Use all attribute keys as headers
				for k := range attrs {
					headers = append(headers, k)
				}
			}
		}

		// Build row
		var row []string
		for _, header := range headers {
			val := attrs[header]
			row = append(row, fmt.Sprintf("%v", val))
		}
		rows = append(rows, row)
	}

	c.Ui.Output(fmt.Sprintf("Showing %d of %s (page %d)\n", len(data), c.queryType, c.pageNumber))
	formatter.Table(headers, rows)
	return 0
}

// Help returns help text
func (c *ExplorerQueryCommand) Help() string {
	helpText := `
Usage: hcptf explorer query [options]

  Query resources across your organization using the Explorer API.
  Supports querying workspaces, Terraform versions, providers, and modules.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -type=<type>         Query type: workspaces, tf_versions, providers, modules (required)
  -sort=<field>        Sort by field (prefix with - for descending, e.g., -created_at)
  -filter=<expr>       Filter expression (e.g., "name:prod*")
  -fields=<fields>     Comma-separated fields to return
  -limit=<number>      Maximum number of results to return (overrides page-size)
  -page=<number>       Page number (default: 1)
  -page-size=<size>    Page size (default: 20)
  -output=<format>     Output format: table (default) or json
  -csv                 Export results as CSV (unpaged)

Examples:

  # List all workspaces
  hcptf explorer query -org=my-org -type=workspaces

  # List first 10 workspaces
  hcptf explorer query -org=my-org -type=workspaces -limit=10

  # List workspaces sorted by creation date
  hcptf explorer query -org=my-org -type=workspaces -sort=-created_at

  # Query Terraform versions in use
  hcptf explorer query -org=my-org -type=tf_versions

  # Query providers
  hcptf explorer query -org=my-org -type=providers -fields=name,version,source

  # Export to CSV
  hcptf explorer query -org=my-org -type=workspaces -csv > workspaces.csv

  # URL-style
  hcptf my-org explorer workspaces
  hcptf my-org explorer providers
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis
func (c *ExplorerQueryCommand) Synopsis() string {
	return "Query resources using Explorer API"
}
