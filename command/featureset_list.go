package command

import (
	"fmt"
	"net/http"
	"strings"
)

// FeatureSetListCommand lists feature sets.
type FeatureSetListCommand struct {
	Meta
	format string
}

func (c *FeatureSetListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("featureset list")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	responseBody, status, err := executeAPIRequest(client, http.MethodGet, "/api/v2/feature-sets", nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error requesting feature sets: %s", err))
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
		return 1
	}

	formatter := c.Meta.NewFormatter(c.format)
	printAPIResponse(formatter, payload)
	return 0
}

func (c *FeatureSetListCommand) Help() string {
	helpText := `
Usage: hcptf featureset list [options]

  List available feature sets.

Options:

  -output=<format>   Output format: table (default) or json
`
	return strings.TrimSpace(helpText)
}

func (c *FeatureSetListCommand) Synopsis() string {
	return "List feature sets"
}
