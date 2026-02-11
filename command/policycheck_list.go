package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// PolicyCheckListCommand is a command to list policy checks for a run
type PolicyCheckListCommand struct {
	Meta
	runID  string
	format string
}

// Run executes the policy check list command
func (c *PolicyCheckListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policycheck list")
	flags.StringVar(&c.runID, "run-id", "", "Run ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.runID == "" {
		c.Ui.Error("Error: -run-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List policy checks
	policyChecks, err := client.PolicyChecks.List(client.Context(), c.runID, &tfe.PolicyCheckListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing policy checks: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(policyChecks.Items) == 0 {
		c.Ui.Output("No policy checks found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Status", "Scope", "Overridable", "Passed", "Failed", "Soft Failed"}
	var rows [][]string

	for _, pc := range policyChecks.Items {
		overridable := "No"
		if pc.Actions.IsOverridable {
			overridable = "Yes"
		}

		passed := 0
		totalFailed := 0
		softFailed := 0
		if pc.Result != nil {
			passed = pc.Result.Passed
			totalFailed = pc.Result.TotalFailed
			softFailed = pc.Result.SoftFailed
		}

		rows = append(rows, []string{
			pc.ID,
			string(pc.Status),
			string(pc.Scope),
			overridable,
			fmt.Sprintf("%d", passed),
			fmt.Sprintf("%d", totalFailed),
			fmt.Sprintf("%d", softFailed),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the policy check list command
func (c *PolicyCheckListCommand) Help() string {
	helpText := `
Usage: hcptf policycheck list [options]

  List policy checks for a run.

Options:

  -run-id=<id>      Run ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf policycheck list -run-id=run-abc123
  hcptf policycheck list -run-id=run-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy check list command
func (c *PolicyCheckListCommand) Synopsis() string {
	return "List policy checks for a run"
}
