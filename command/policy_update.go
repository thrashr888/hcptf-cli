package command

import (
	"fmt"
	"os"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// PolicyUpdateCommand is a command to update a policy
type PolicyUpdateCommand struct {
	Meta
	policyID    string
	description string
	enforce     string
	policyFile  string
	format      string
}

// Run executes the policy update command
func (c *PolicyUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policy update")
	flags.StringVar(&c.policyID, "id", "", "Policy ID (required)")
	flags.StringVar(&c.description, "description", "", "Policy description")
	flags.StringVar(&c.enforce, "enforce", "", "Enforcement level: advisory, soft-mandatory, or hard-mandatory")
	flags.StringVar(&c.policyFile, "policy-file", "", "Path to policy file")
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

	// Build update options
	options := tfe.PolicyUpdateOptions{}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	if c.enforce != "" {
		var enforcementLevel tfe.EnforcementLevel
		switch c.enforce {
		case "advisory":
			enforcementLevel = tfe.EnforcementAdvisory
		case "soft-mandatory":
			enforcementLevel = tfe.EnforcementSoft
		case "hard-mandatory":
			enforcementLevel = tfe.EnforcementHard
		default:
			c.Ui.Error("Error: -enforce must be 'advisory', 'soft-mandatory', or 'hard-mandatory'")
			return 1
		}
		options.EnforcementLevel = &enforcementLevel
	}

	// Update policy
	policy, err := client.Policies.Update(client.Context(), c.policyID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating policy: %s", err))
		return 1
	}

	// Upload new policy content if provided
	if c.policyFile != "" {
		policyContent, err := os.ReadFile(c.policyFile)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading policy file: %s", err))
			return 1
		}

		err = client.Policies.Upload(client.Context(), c.policyID, policyContent)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error uploading policy content: %s", err))
			return 1
		}
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Policy '%s' updated successfully", policy.Name))

	// Show policy details
	data := map[string]interface{}{
		"ID":               policy.ID,
		"Name":             policy.Name,
		"Description":      policy.Description,
		"EnforcementLevel": string(policy.EnforcementLevel),
		"PolicySetCount":   policy.PolicySetCount,
		"UpdatedAt":        policy.UpdatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the policy update command
func (c *PolicyUpdateCommand) Help() string {
	helpText := `
Usage: hcptf policy update [options]

  Update policy settings.

Options:

  -id=<policy-id>       Policy ID (required)
  -description=<text>   Policy description
  -enforce=<level>      Enforcement level: advisory, soft-mandatory, or hard-mandatory
  -policy-file=<path>   Path to policy file
  -output=<format>      Output format: table (default) or json

Example:

  hcptf policy update -id=pol-abc123 -enforce=hard-mandatory
  hcptf policy update -id=pol-abc123 -description="Updated policy" -policy-file=updated-policy.sentinel
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy update command
func (c *PolicyUpdateCommand) Synopsis() string {
	return "Update policy settings"
}
