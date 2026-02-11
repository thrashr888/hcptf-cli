package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// PolicySetOutcomeReadCommand is a command to read policy set outcome details
type PolicySetOutcomeReadCommand struct {
	Meta
	policySetOutcomeID string
	format             string
}

// Run executes the policy set outcome read command
func (c *PolicySetOutcomeReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policysetoutcome read")
	flags.StringVar(&c.policySetOutcomeID, "id", "", "Policy Set Outcome ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.policySetOutcomeID == "" {
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

	// Read policy set outcome
	policySetOutcome, err := client.PolicySetOutcomes.Read(client.Context(), c.policySetOutcomeID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading policy set outcome: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	overridable := "Unknown"
	if policySetOutcome.Overridable != nil {
		if *policySetOutcome.Overridable {
			overridable = "Yes"
		} else {
			overridable = "No"
		}
	}

	data := map[string]interface{}{
		"ID":                   policySetOutcome.ID,
		"PolicySetName":        policySetOutcome.PolicySetName,
		"PolicySetDescription": policySetOutcome.PolicySetDescription,
		"Overridable":          overridable,
		"Error":                policySetOutcome.Error,
	}

	// ResultCount is a struct, not a pointer, so we can access it directly
	data["Passed"] = policySetOutcome.ResultCount.Passed
	data["MandatoryFailed"] = policySetOutcome.ResultCount.MandatoryFailed
	data["AdvisoryFailed"] = policySetOutcome.ResultCount.AdvisoryFailed
	data["Errored"] = policySetOutcome.ResultCount.Errored

	// Show individual policy outcomes
	if len(policySetOutcome.Outcomes) > 0 {
		c.Ui.Output("")
		c.Ui.Output("Individual Policy Outcomes:")
		c.Ui.Output("")

		outcomeHeaders := []string{"Policy Name", "Status", "Enforcement Level", "Query", "Description"}
		var outcomeRows [][]string

		for _, outcome := range policySetOutcome.Outcomes {
			outcomeRows = append(outcomeRows, []string{
				outcome.PolicyName,
				outcome.Status,
				string(outcome.EnforcementLevel),
				outcome.Query,
				outcome.Description,
			})
		}

		outcomeFormatter := output.NewFormatter("table")
		outcomeFormatter.Table(outcomeHeaders, outcomeRows)
	}

	c.Ui.Output("")
	c.Ui.Output("Policy Set Outcome Summary:")
	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the policy set outcome read command
func (c *PolicySetOutcomeReadCommand) Help() string {
	helpText := `
Usage: hcptf policysetoutcome read [options]

  Read policy set outcome details. Shows the results of evaluating a policy set,
  including individual policy outcomes within the set.

Options:

  -id=<policy-set-outcome-id>  Policy Set Outcome ID (required)
  -output=<format>             Output format: table (default) or json

Example:

  hcptf policysetoutcome read -id=psout-abc123
  hcptf policysetoutcome read -id=psout-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set outcome read command
func (c *PolicySetOutcomeReadCommand) Synopsis() string {
	return "Read policy set outcome details"
}
