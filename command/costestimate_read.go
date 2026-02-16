package command

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// CostEstimateReadCommand reads a cost estimate by ID.
type CostEstimateReadCommand struct {
	Meta
	id     string
	format string
}

func (c *CostEstimateReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("costestimate read")
	flags.StringVar(&c.id, "id", "", "Cost estimate ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	responseBody, status, err := executeAPIRequest(client, http.MethodGet, fmt.Sprintf("/api/v2/cost-estimates/%s", url.PathEscape(c.id)), nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error requesting cost estimate: %s", err))
		return 1
	}
	if status < 200 || status >= 300 {
		c.Ui.Error(fmt.Sprintf("API request failed with status %d: %s", status, string(responseBody)))
		return 1
	}

	payload, err := parseAPIResponse(responseBody)
	if err != nil {
		if c.format == "json" {
			c.Ui.Output(string(responseBody))
			return 0
		}
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		c.Ui.Error(string(responseBody))
		return 1
	}

	formatter := c.Meta.NewFormatter(c.format)
	printAPIResponse(formatter, payload)
	return 0
}

func (c *CostEstimateReadCommand) Help() string {
	helpText := `
Usage: hcptf costestimate read [options]

  Read a cost estimate by ID.

Options:

  -id=<id>           Cost estimate ID (required)
  -output=<format>   Output format: table (default) or json

Example:

  hcptf costestimate read -id=ce-12345
  hcptf costestimate read -id=ce-12345 -output=json
`
	return strings.TrimSpace(helpText)
}

func (c *CostEstimateReadCommand) Synopsis() string {
	return "Read a cost estimate"
}
