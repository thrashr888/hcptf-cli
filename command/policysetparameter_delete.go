package command

import (
	"context"
	"fmt"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"strings"
)

type policySetParameterDeleter interface {
	Delete(ctx context.Context, policySetID, parameterID string) error
}

// PolicySetParameterDeleteCommand is a command to delete a policy set parameter
type PolicySetParameterDeleteCommand struct {
	Meta
	policySetID           string
	parameterID           string
	force                 bool
	yes                   bool
	policySetParameterSvc policySetParameterDeleter
}

// Run executes the policy set parameter delete command
func (c *PolicySetParameterDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policysetparameter delete")
	flags.StringVar(&c.policySetID, "policy-set-id", "", "Policy Set ID (required)")
	flags.StringVar(&c.parameterID, "id", "", "Parameter ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")
	flags.BoolVar(&c.force, "f", false, "Shorthand for -force")
	flags.BoolVar(&c.yes, "y", false, "Confirm delete without prompt")

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

	// Confirm deletion unless force or -y is set
	if !c.force && !c.yes {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete parameter '%s' from policy set '%s'? (yes/no): ", c.parameterID, c.policySetID))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.TrimSpace(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete policy set parameter
	err = c.policySetParameterService(client).Delete(client.Context(), c.policySetID, c.parameterID)
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
	Usage: hcptf policyset parameter delete [options]

  Delete a policy set parameter. This action cannot be undone.

Options:

  -policy-set-id=<id>  Policy Set ID (required)
  -id=<parameter-id>   Parameter ID (required)
  -force               Force delete without confirmation
  -f                   Shorthand for -force
  -y                   Confirm delete without prompt

Example:

  hcptf policysetparameter delete -policy-set-id=polset-abc123 -id=var-xyz789
  hcptf policysetparameter delete -policy-set-id=polset-abc123 -id=var-xyz789 -force
  hcptf policysetparameter delete -policy-set-id=polset-abc123 -id=var-xyz789 -y
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set parameter delete command
func (c *PolicySetParameterDeleteCommand) Synopsis() string {
	return "Delete a policy set parameter"
}

func (c *PolicySetParameterDeleteCommand) policySetParameterService(client *client.Client) policySetParameterDeleter {
	if c.policySetParameterSvc != nil {
		return c.policySetParameterSvc
	}
	return client.PolicySetParameters
}
