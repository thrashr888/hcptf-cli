package command

import (
	"fmt"
	"strings"
)

// PolicyCheckOverrideCommand is a command to override a soft-mandatory policy check
type PolicyCheckOverrideCommand struct {
	Meta
	policyCheckID string
	format        string
	autoApprove   bool
}

// Run executes the policy check override command
func (c *PolicyCheckOverrideCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policycheck override")
	flags.StringVar(&c.policyCheckID, "id", "", "Policy Check ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")
	flags.BoolVar(&c.autoApprove, "auto-approve", false, "Skip confirmation prompt")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.policyCheckID == "" {
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

	// Read policy check first to verify it can be overridden
	policyCheck, err := client.PolicyChecks.Read(client.Context(), c.policyCheckID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading policy check: %s", err))
		return 1
	}

	if !policyCheck.Actions.IsOverridable {
		c.Ui.Error("Error: This policy check cannot be overridden (not soft-mandatory)")
		return 1
	}

	if !policyCheck.Permissions.CanOverride {
		c.Ui.Error("Error: You do not have permission to override this policy check")
		return 1
	}

	// Confirm override unless auto-approve is set
	if !c.autoApprove {
		c.Ui.Output(fmt.Sprintf("Policy Check: %s", policyCheck.ID))
		c.Ui.Output(fmt.Sprintf("Status: %s", policyCheck.Status))
		c.Ui.Output(fmt.Sprintf("Scope: %s", policyCheck.Scope))
		if policyCheck.Result != nil {
			c.Ui.Output(fmt.Sprintf("Failed Policies: %d", policyCheck.Result.SoftFailed))
		}
		c.Ui.Output("")
		c.Ui.Output("Are you sure you want to override this policy check?")
		c.Ui.Output("Only 'yes' will be accepted to confirm.")
		c.Ui.Output("")

		response, err := c.Ui.Ask("Enter a value: ")
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading input: %s", err))
			return 1
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "yes" {
			c.Ui.Output("Override cancelled.")
			return 0
		}
	}

	// Override policy check
	policyCheck, err = client.PolicyChecks.Override(client.Context(), c.policyCheckID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error overriding policy check: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output("Policy check overridden successfully")

	data := map[string]interface{}{
		"ID":     policyCheck.ID,
		"Status": string(policyCheck.Status),
		"Scope":  string(policyCheck.Scope),
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the policy check override command
func (c *PolicyCheckOverrideCommand) Help() string {
	helpText := `
Usage: hcptf policycheck override [options]

  Override a soft-mandatory policy check. This allows a run to proceed
  even when a soft-mandatory policy has failed.

  Note: You must have the appropriate permissions to override policy checks.
  Only soft-mandatory (advisory) policy checks can be overridden.

Options:

  -id=<policy-check-id>  Policy Check ID (required)
  -output=<format>       Output format: table (default) or json
  -auto-approve          Skip confirmation prompt

Example:

  hcptf policycheck override -id=polchk-abc123
  hcptf policycheck override -id=polchk-abc123 -auto-approve
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy check override command
func (c *PolicyCheckOverrideCommand) Synopsis() string {
	return "Override a soft-mandatory policy check"
}
