package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PolicyUploadCommand uploads policy content for an existing policy.
type PolicyUploadCommand struct {
	Meta
	policyID   string
	policyFile string
	policySvc  policyUploader
}

// Run executes the policy upload command.
func (c *PolicyUploadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policy upload")
	flags.StringVar(&c.policyID, "id", "", "Policy ID (required)")
	flags.StringVar(&c.policyFile, "policy-file", "", "Path to policy file (required)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.policyID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}
	if c.policyFile == "" {
		c.Ui.Error("Error: -policy-file flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	content, err := os.ReadFile(c.policyFile)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading policy file: %s", err))
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	if err := c.policyService(client).Upload(client.Context(), c.policyID, content); err != nil {
		c.Ui.Error(fmt.Sprintf("Error uploading policy content: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Policy content uploaded for %s", c.policyID))
	return 0
}

func (c *PolicyUploadCommand) policyService(client *client.Client) policyUploader {
	if c.policySvc != nil {
		return c.policySvc
	}
	return client.Policies
}

// Help returns help text for the policy upload command.
func (c *PolicyUploadCommand) Help() string {
	helpText := `
Usage: hcptf policy upload [options]

  Upload content for an existing policy.

Options:

  -id=<policy-id>      Policy ID (required)
  -policy-file=<path>  Path to policy file (required)

Example:

  hcptf policy upload -id=pol-abc123 -policy-file=policy.sentinel
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy upload command.
func (c *PolicyUploadCommand) Synopsis() string {
	return "Upload policy content"
}
