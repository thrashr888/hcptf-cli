package command

import (
	"fmt"
	"net/http"
	"strings"
)

// IPRangeReadCommand reads cloud IP range information.
type IPRangeReadCommand struct {
	Meta
	format string
}

func (c *IPRangeReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("iprange read")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	responseBody, status, err := executeAPIRequest(client, http.MethodGet, "/api/v2/meta/ip-ranges", nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error requesting IP ranges: %s", err))
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

func (c *IPRangeReadCommand) Help() string {
	helpText := `
Usage: hcptf iprange read [options]

  Read current HCP Terraform IP ranges.

Options:

  -output=<format>   Output format: table (default) or json
`
	return strings.TrimSpace(helpText)
}

func (c *IPRangeReadCommand) Synopsis() string {
	return "Read public IP ranges"
}
