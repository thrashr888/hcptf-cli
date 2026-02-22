package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PolicyDownloadCommand downloads policy content.
type PolicyDownloadCommand struct {
	Meta
	policyID   string
	outputFile string
	policySvc  policyDownloader
}

// Run executes the policy download command.
func (c *PolicyDownloadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policy download")
	flags.StringVar(&c.policyID, "id", "", "Policy ID (required)")
	flags.StringVar(&c.outputFile, "output", "", "Output file path (required)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.policyID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}
	if c.outputFile == "" {
		c.Ui.Error("Error: -output flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	content, err := c.policyService(client).Download(client.Context(), c.policyID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error downloading policy content: %s", err))
		return 1
	}

	if err := os.WriteFile(c.outputFile, content, 0o600); err != nil {
		c.Ui.Error(fmt.Sprintf("Error writing output file: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Policy content downloaded to %s", c.outputFile))
	return 0
}

func (c *PolicyDownloadCommand) policyService(client *client.Client) policyDownloader {
	if c.policySvc != nil {
		return c.policySvc
	}
	return client.Policies
}

// Help returns help text for the policy download command.
func (c *PolicyDownloadCommand) Help() string {
	helpText := `
Usage: hcptf policy download [options]

  Download content for an existing policy.

Options:

  -id=<policy-id>   Policy ID (required)
  -output=<path>    Output file path (required)

Example:

  hcptf policy download -id=pol-abc123 -output=policy.sentinel
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy download command.
func (c *PolicyDownloadCommand) Synopsis() string {
	return "Download policy content"
}
