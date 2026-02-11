package command

import (
	"fmt"
	"strings"
)

// PolicyDeleteCommand is a command to delete a policy
type PolicyDeleteCommand struct {
	Meta
	policyID string
	force    bool
}

// Run executes the policy delete command
func (c *PolicyDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policy delete")
	flags.StringVar(&c.policyID, "id", "", "Policy ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")

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

	// Read policy to get the name for confirmation
	policy, err := client.Policies.Read(client.Context(), c.policyID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading policy: %s", err))
		return 1
	}

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete policy '%s' (%s)? (yes/no): ", policy.Name, c.policyID))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete policy
	err = client.Policies.Delete(client.Context(), c.policyID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting policy: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Policy '%s' (%s) deleted successfully", policy.Name, c.policyID))
	return 0
}

// Help returns help text for the policy delete command
func (c *PolicyDeleteCommand) Help() string {
	helpText := `
Usage: hcptf policy delete [options]

  Delete a policy.

Options:

  -id=<policy-id>  Policy ID (required)
  -force           Force delete without confirmation

Example:

  hcptf policy delete -id=pol-abc123
  hcptf policy delete -id=pol-abc123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy delete command
func (c *PolicyDeleteCommand) Synopsis() string {
	return "Delete a policy"
}
