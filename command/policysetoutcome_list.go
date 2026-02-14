package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// PolicySetOutcomeListCommand is a command to list policy set outcomes for a policy evaluation
type PolicySetOutcomeListCommand struct {
	Meta
	policyEvaluationID string
	format             string
}

// Run executes the policy set outcome list command
func (c *PolicySetOutcomeListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policysetoutcome list")
	flags.StringVar(&c.policyEvaluationID, "policy-evaluation-id", "", "Policy Evaluation ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.policyEvaluationID == "" {
		c.Ui.Error("Error: -policy-evaluation-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List policy set outcomes
	policySetOutcomes, err := client.PolicySetOutcomes.List(client.Context(), c.policyEvaluationID, &tfe.PolicySetOutcomeListOptions{
		ListOptions: &tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing policy set outcomes: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(policySetOutcomes.Items) == 0 {
		c.Ui.Output("No policy set outcomes found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Policy Set Name", "Overridable", "Passed", "Mandatory Failed", "Advisory Failed", "Errored"}
	var rows [][]string

	for _, pso := range policySetOutcomes.Items {
		overridable := "No"
		if pso.Overridable != nil && *pso.Overridable {
			overridable = "Yes"
		}

		// ResultCount is a struct, not a pointer, so we can access it directly
		passed := pso.ResultCount.Passed
		mandatoryFailed := pso.ResultCount.MandatoryFailed
		advisoryFailed := pso.ResultCount.AdvisoryFailed
		errored := pso.ResultCount.Errored

		rows = append(rows, []string{
			pso.ID,
			pso.PolicySetName,
			overridable,
			fmt.Sprintf("%d", passed),
			fmt.Sprintf("%d", mandatoryFailed),
			fmt.Sprintf("%d", advisoryFailed),
			fmt.Sprintf("%d", errored),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the policy set outcome list command
func (c *PolicySetOutcomeListCommand) Help() string {
	helpText := `
Usage: hcptf policysetoutcome list [options]

  List policy set outcomes for a policy evaluation. Policy set outcomes
  represent the results of evaluating each policy set.

Options:

  -policy-evaluation-id=<id>  Policy Evaluation ID (required)
  -output=<format>            Output format: table (default) or json

Example:

  hcptf policysetoutcome list -policy-evaluation-id=poleval-abc123
  hcptf policysetoutcome list -policy-evaluation-id=poleval-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set outcome list command
func (c *PolicySetOutcomeListCommand) Synopsis() string {
	return "List policy set outcomes for a policy evaluation"
}
