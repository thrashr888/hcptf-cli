package command

import (
	"fmt"
	"strings"
)

// PolicySetParameterDeleteCommand is a command to delete a policy set parameter
type PolicySetParameterDeleteCommand struct {
	Meta
	policySetID string
	parameterID string
	autoApprove bool
}

// Run executes the policy set parameter delete command
func (c *PolicySetParameterDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policysetparameter delete")
	flags.StringVar(&c.policySetID, "policy-set-id", "", "Policy Set ID (required)")
	flags.StringVar(&c.parameterID, "id", "", "Parameter ID (required)")
	flags.BoolVar(&c.autoApprove, "auto-approve", false, "Skip confirmation prompt")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.policySetID == "" {
		c.Ui.Error("Error: -policy-set-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.parameterID == "" {
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

	// Confirm deletion unless auto-approve is set
	if !c.autoApprove {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete parameter '%s' from policy set '%s'? (yes/no): ", c.parameterID, c.policySetID))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete policy set parameter
	err = client.PolicySetParameters.Delete(client.Context(), c.policySetID, c.parameterID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting policy set parameter: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Policy set parameter %s deleted successfully", c.parameterID))
	return 0
}

// Help returns help text for the policy set parameter delete command
func (c *PolicySetParameterDeleteCommand) Help() string {
	helpText := `
Usage: hcptf policysetparameter delete [options]

  Delete a policy set parameter. This action cannot be undone.

Options:

  -policy-set-id=<id>  Policy Set ID (required)
  -id=<parameter-id>   Parameter ID (required)
  -auto-approve        Skip confirmation prompt

Example:

  hcptf policysetparameter delete -policy-set-id=polset-abc123 -id=var-xyz789
  hcptf policysetparameter delete -policy-set-id=polset-abc123 -id=var-xyz789 -auto-approve
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set parameter delete command
func (c *PolicySetParameterDeleteCommand) Synopsis() string {
	return "Delete a policy set parameter"
}
