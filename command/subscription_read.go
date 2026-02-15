package command

import (
	"fmt"
	"net/http"
	"strings"
)

// SubscriptionReadCommand reads a subscription by ID.
type SubscriptionReadCommand struct {
	Meta
	subscriptionID string
	format         string
}

func (c *SubscriptionReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("subscription read")
	flags.StringVar(&c.subscriptionID, "id", "", "Subscription ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.subscriptionID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	responseBody, status, err := executeAPIRequest(client, http.MethodGet, fmt.Sprintf("/api/v2/subscriptions/%s", c.subscriptionID), nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error requesting subscription: %s", err))
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

func (c *SubscriptionReadCommand) Help() string {
	helpText := `
Usage: hcptf subscription read [options]

  Read a subscription by ID.

Options:

  -id=<id>           Subscription ID (required)
  -output=<format>   Output format: table (default) or json
`
	return strings.TrimSpace(helpText)
}

func (c *SubscriptionReadCommand) Synopsis() string {
	return "Read a subscription"
}
