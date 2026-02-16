package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// NoCodeCreateCommand creates a no-code provisioning policy.
type NoCodeCreateCommand struct {
	Meta
	organization string
	payload      string
	format       string
}

func (c *NoCodeCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("nocode create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.payload, "payload", "", "JSON payload for the create request")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.payload == "" {
		c.Ui.Error("Error: -payload flag is required")
		c.Ui.Error("Provide a JSON object representing the provisioning request")
		c.Ui.Error(c.Help())
		return 1
	}

	if !json.Valid([]byte(c.payload)) {
		c.Ui.Error("Error: -payload must be valid JSON")
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	endpoint := fmt.Sprintf("/api/v2/organizations/%s/no-code-provisioning", url.PathEscape(c.organization))
	responseBody, status, err := executeAPIRequest(client, http.MethodPost, endpoint, strings.NewReader(c.payload))
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating no-code provisioning: %s", err))
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

func (c *NoCodeCreateCommand) Help() string {
	helpText := `
Usage: hcptf nocode create [options]

  Create no-code provisioning configuration for an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>           Alias for -organization
  -payload=<json>       JSON payload (required)
  -output=<format>      Output format: table (default) or json

Example:

  hcptf nocode create -org=my-org -payload='{"enabled":true}'
`
	return strings.TrimSpace(helpText)
}

func (c *NoCodeCreateCommand) Synopsis() string {
	return "Create no-code provisioning"
}
