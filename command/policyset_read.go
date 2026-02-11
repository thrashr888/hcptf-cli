package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// PolicySetReadCommand is a command to read policy set details
type PolicySetReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the policy set read command
func (c *PolicySetReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset read")
	flags.StringVar(&c.id, "id", "", "Policy set ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.id == "" {
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

	// Read policy set
	policySet, err := client.PolicySets.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading policy set: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                policySet.ID,
		"Name":              policySet.Name,
		"Description":       policySet.Description,
		"Global":            policySet.Global,
		"PolicyCount":       policySet.PolicyCount,
		"WorkspaceCount":    policySet.WorkspaceCount,
		"ProjectCount":      policySet.ProjectCount,
		"VCSRepo":           policySet.VCSRepo,
		"PoliciesPath":      policySet.PoliciesPath,
		"AgentEnabled":      policySet.AgentEnabled,
		"PolicyToolVersion": policySet.PolicyToolVersion,
		"CreatedAt":         policySet.CreatedAt,
		"UpdatedAt":         policySet.UpdatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the policy set read command
func (c *PolicySetReadCommand) Help() string {
	helpText := `
Usage: hcptf policyset read [options]

  Read policy set details.

Options:

  -id=<id>          Policy set ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf policyset read -id=polset-12345
  hcptf policyset read -id=polset-12345 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set read command
func (c *PolicySetReadCommand) Synopsis() string {
	return "Read policy set details"
}
