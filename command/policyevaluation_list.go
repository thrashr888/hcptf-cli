package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// PolicyEvaluationListCommand is a command to list policy evaluations for a task stage
type PolicyEvaluationListCommand struct {
	Meta
	taskStageID string
	format      string
}

// Run executes the policy evaluation list command
func (c *PolicyEvaluationListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyevaluation list")
	flags.StringVar(&c.taskStageID, "task-stage-id", "", "Task Stage ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.taskStageID == "" {
		c.Ui.Error("Error: -task-stage-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List policy evaluations
	policyEvaluations, err := client.PolicyEvaluations.List(client.Context(), c.taskStageID, &tfe.PolicyEvaluationListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing policy evaluations: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(policyEvaluations.Items) == 0 {
		c.Ui.Output("No policy evaluations found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Status", "Policy Kind", "Passed", "Mandatory Failed", "Advisory Failed", "Errored"}
	var rows [][]string

	for _, pe := range policyEvaluations.Items {
		passed := 0
		mandatoryFailed := 0
		advisoryFailed := 0
		errored := 0

		// ResultCount is a struct, not a pointer, so we can access it directly
		passed = pe.ResultCount.Passed
		mandatoryFailed = pe.ResultCount.MandatoryFailed
		advisoryFailed = pe.ResultCount.AdvisoryFailed
		errored = pe.ResultCount.Errored

		rows = append(rows, []string{
			pe.ID,
			string(pe.Status),
			string(pe.PolicyKind),
			fmt.Sprintf("%d", passed),
			fmt.Sprintf("%d", mandatoryFailed),
			fmt.Sprintf("%d", advisoryFailed),
			fmt.Sprintf("%d", errored),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the policy evaluation list command
func (c *PolicyEvaluationListCommand) Help() string {
	helpText := `
Usage: hcptf policyevaluation list [options]

  List policy evaluations for a task stage. Policy evaluations represent
  individual policy executions within the task stage. This is primarily
  used for OPA policies.

  Note: Policy evaluations are part of the newer task-based policy workflow,
  not the legacy policy check workflow. Use 'policycheck' commands for
  Sentinel policy checks.

Options:

  -task-stage-id=<id>  Task Stage ID (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf policyevaluation list -task-stage-id=ts-abc123
  hcptf policyevaluation list -task-stage-id=ts-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy evaluation list command
func (c *PolicyEvaluationListCommand) Synopsis() string {
	return "List policy evaluations for a task stage"
}
