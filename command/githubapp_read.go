package command

import (
	"fmt"
	"net/http"
	"strings"
)

// GitHubAppReadCommand reads a GitHub App installation by ID.
type GitHubAppReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the githubapp read command.
func (c *GitHubAppReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("githubapp read")
	flags.StringVar(&c.id, "id", "", "GitHub App installation ID (required)")
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

	responseBody, status, err := executeAPIRequest(client, http.MethodGet, fmt.Sprintf("/api/v2/github-app-installations/%s", c.id), nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error requesting GitHub App installation: %s", err))
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

// Help returns help text for the githubapp read command.
func (c *GitHubAppReadCommand) Help() string {
	helpText := `
Usage: hcptf githubapp read [options]

  Read a GitHub App installation by ID.

Options:

  -id=<id>           GitHub App installation ID (required)
  -output=<format>   Output format: table (default) or json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the githubapp read command.
func (c *GitHubAppReadCommand) Synopsis() string {
	return "Read a GitHub App installation"
}
