package command

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

// OrganizationTagCreateCommand creates an organization tag.
type OrganizationTagCreateCommand struct {
	Meta
	organization string
	name         string
	format       string
}

// Run executes the organizationtag create command.
func (c *OrganizationTagCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("organizationtag create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Tag name (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}
	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	endpoint := fmt.Sprintf("/api/v2/organizations/%s/tags", c.organization)
	body := fmt.Sprintf(`{"data":{"type":"tags","attributes":{"name":%q}}}`, c.name)
	responseBody, status, err := executeAPIRequest(client, http.MethodPost, endpoint, bytes.NewBufferString(body))
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating organization tag: %s", err))
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

// Help returns help text for the organizationtag create command.
func (c *OrganizationTagCreateCommand) Help() string {
	helpText := `
Usage: hcptf organizationtag create [options]

  Create a tag in an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>           Alias for -organization
  -name=<name>          Tag name (required)
  -output=<format>      Output format: table (default) or json

Example:

  hcptf organizationtag create -org=my-org -name=platform
  hcptf organizationtag create -org=my-org -name=platform -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the organizationtag create command.
func (c *OrganizationTagCreateCommand) Synopsis() string {
	return "Create an organization tag"
}
