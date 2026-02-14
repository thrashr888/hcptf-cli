package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PolicyReadCommand is a command to read policy details
type PolicyReadCommand struct {
	Meta
	policyID    string
	format      string
	policySvc   policyReader
	downloadSvc policyDownloader
}

// Run executes the policy read command
func (c *PolicyReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policy read")
	flags.StringVar(&c.policyID, "id", "", "Policy ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.policyID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read policy
	policy, err := c.policyService(client).Read(client.Context(), c.policyID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading policy: %s", err))
		return 1
	}

	// Download policy content
	policyContent, err := c.policyDownloadService(client).Download(client.Context(), c.policyID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error downloading policy content: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":               policy.ID,
		"Name":             policy.Name,
		"Description":      policy.Description,
		"EnforcementLevel": string(policy.EnforcementLevel),
		"PolicySetCount":   policy.PolicySetCount,
		"UpdatedAt":        policy.UpdatedAt,
		"Content":          string(policyContent),
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the policy read command
func (c *PolicyReadCommand) Help() string {
	helpText := `
Usage: hcptf policy read [options]

  Read policy details.

Options:

  -id=<policy-id>   Policy ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf policy read -id=pol-abc123
  hcptf policy read -id=pol-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

func (c *PolicyReadCommand) policyService(client *client.Client) policyReader {
	if c.policySvc != nil {
		return c.policySvc
	}
	return client.Policies
}

func (c *PolicyReadCommand) policyDownloadService(client *client.Client) policyDownloader {
	if c.downloadSvc != nil {
		return c.downloadSvc
	}
	return client.Policies
}

// Synopsis returns a short synopsis for the policy read command
func (c *PolicyReadCommand) Synopsis() string {
	return "Read policy details"
}
