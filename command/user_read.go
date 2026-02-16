package command

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// UserReadCommand reads a user by ID.
type UserReadCommand struct {
	Meta
	userID string
	format string
}

func (c *UserReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("user read")
	flags.StringVar(&c.userID, "id", "", "User ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.userID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	responseBody, status, err := executeAPIRequest(client, http.MethodGet, fmt.Sprintf("/api/v2/users/%s", url.PathEscape(c.userID)), nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error requesting user: %s", err))
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

func (c *UserReadCommand) Help() string {
	helpText := `
Usage: hcptf user read [options]

  Read a user by user ID.

Options:

  -id=<id>           User ID (required)
  -output=<format>   Output format: table (default) or json
`
	return strings.TrimSpace(helpText)
}

func (c *UserReadCommand) Synopsis() string {
	return "Read user details"
}
