package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// PolicySetAddPolicyCommand is a command to add a policy to a policy set
type PolicySetAddPolicyCommand struct {
	Meta
	policySetID string
	policyID    string
}

// Run executes the policy set add-policy command
func (c *PolicySetAddPolicyCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset add-policy")
	flags.StringVar(&c.policySetID, "policyset-id", "", "Policy set ID (required)")
	flags.StringVar(&c.policyID, "policy-id", "", "Policy ID to add (required)")

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

	// Add policy to policy set
	options := tfe.PolicySetAddPoliciesOptions{
		Policies: []*tfe.Policy{
			{ID: c.policyID},
		},
	}

	err = client.PolicySets.AddPolicies(client.Context(), c.policySetID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error adding policy to policy set: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Policy '%s' added to policy set '%s' successfully", c.policyID, c.policySetID))
	return 0
}

// Help returns help text for the policy set add-policy command
func (c *PolicySetAddPolicyCommand) Help() string {
	helpText := `
Usage: hcptf policyset add-policy [options]

  Add a policy to a policy set. This only works for policy sets that are not
  managed by VCS.

Options:

  -policyset-id=<id>  Policy set ID (required)
  -policy-id=<id>     Policy ID to add (required)

Example:

  hcptf policyset add-policy -policyset-id=polset-12345 -policy-id=pol-67890
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set add-policy command
func (c *PolicySetAddPolicyCommand) Synopsis() string {
	return "Add a policy to a policy set"
}
