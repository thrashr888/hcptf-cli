package command

import (
	"fmt"
	"net/http"
	"strings"
)

// GitHubAppListCommand lists GitHub App installations for an organization.
type GitHubAppListCommand struct {
	Meta
	organization string
	format       string
}

func (c *GitHubAppListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("githubapp list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	endpoint := fmt.Sprintf("/api/v2/organizations/%s/github-app-installations", c.organization)
	responseBody, status, err := executeAPIRequest(client, http.MethodGet, endpoint, nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error requesting GitHub App installations: %s", err))
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
		c.Ui.Output(string(responseBody))
		return 1
	}

	formatter := c.Meta.NewFormatter(c.format)
	printAPIResponse(formatter, payload)
	return 0
}

func (c *GitHubAppListCommand) Help() string {
	helpText := `
Usage: hcptf githubapp list [options]

  List GitHub App installations for an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>           Alias for -organization
  -output=<format>      Output format: table (default) or json
`
	return strings.TrimSpace(helpText)
}

func (c *GitHubAppListCommand) Synopsis() string {
	return "List GitHub App installations"
}
