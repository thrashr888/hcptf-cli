package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// PolicySetRemovePolicyCommand is a command to remove a policy from a policy set
type PolicySetRemovePolicyCommand struct {
	Meta
	policySetID string
	policyID    string
}

// Run executes the policy set remove-policy command
func (c *PolicySetRemovePolicyCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset remove-policy")
	flags.StringVar(&c.policySetID, "policyset-id", "", "Policy set ID (required)")
	flags.StringVar(&c.policyID, "policy-id", "", "Policy ID to remove (required)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.policySetID == "" {
		c.Ui.Error("Error: -policyset-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.policyID == "" {
		c.Ui.Error("Error: -policy-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Remove policy from policy set
	options := tfe.PolicySetRemovePoliciesOptions{
		Policies: []*tfe.Policy{
			{ID: c.policyID},
		},
	}

	err = client.PolicySets.RemovePolicies(client.Context(), c.policySetID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error removing policy from policy set: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Policy '%s' removed from policy set '%s' successfully", c.policyID, c.policySetID))
	return 0
}

// Help returns help text for the policy set remove-policy command
func (c *PolicySetRemovePolicyCommand) Help() string {
	helpText := `
Usage: hcptf policyset remove-policy [options]

  Remove a policy from a policy set. This only works for policy sets that are
  not managed by VCS.

Options:

  -policyset-id=<id>  Policy set ID (required)
  -policy-id=<id>     Policy ID to remove (required)

Example:

  hcptf policyset remove-policy -policyset-id=polset-12345 -policy-id=pol-67890
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set remove-policy command
func (c *PolicySetRemovePolicyCommand) Synopsis() string {
	return "Remove a policy from a policy set"
}
