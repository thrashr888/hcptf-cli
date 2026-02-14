package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PlanReadCommand is a command to read plan details
type PlanReadCommand struct {
	Meta
	planID  string
	runID   string
	format  string
	planSvc planReader
	runSvc  runReader
}

// Run executes the plan read command
func (c *PlanReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("plan read")
	flags.StringVar(&c.planID, "id", "", "Plan ID or Run ID")
	flags.StringVar(&c.runID, "run-id", "", "Run ID (alternative to -id)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags - need either planID or runID
	id := c.planID
	if id == "" {
		id = c.runID
	}
	if id == "" {
		c.Ui.Error("Error: -id or -run-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// If ID starts with "run-", get the plan ID from the run
	planID := id
	if strings.HasPrefix(id, "run-") {
		run, err := c.runService(client).Read(client.Context(), id)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading run: %s", err))
			return 1
		}
		if run.Plan == nil {
			c.Ui.Error("Error: run has no plan")
			return 1
		}
		planID = run.Plan.ID
	}

	// Read plan
	plan, err := c.planService(client).Read(client.Context(), planID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading plan: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                     plan.ID,
		"Status":                 plan.Status,
		"HasChanges":             plan.HasChanges,
		"ResourceAdditions":      plan.ResourceAdditions,
		"ResourceChanges":        plan.ResourceChanges,
		"ResourceDestructions":   plan.ResourceDestructions,
		"ResourceImports":        plan.ResourceImports,
		"GeneratedConfiguration": plan.GeneratedConfiguration,
	}

	if plan.StatusTimestamps != nil {
		if !plan.StatusTimestamps.QueuedAt.IsZero() {
			data["QueuedAt"] = plan.StatusTimestamps.QueuedAt
		}
		if !plan.StatusTimestamps.StartedAt.IsZero() {
			data["StartedAt"] = plan.StatusTimestamps.StartedAt
		}
		if !plan.StatusTimestamps.FinishedAt.IsZero() {
			data["FinishedAt"] = plan.StatusTimestamps.FinishedAt
		}
	}

	if plan.LogReadURL != "" {
		data["LogReadURL"] = plan.LogReadURL
	}

	formatter.KeyValue(data)
	return 0
}

func (c *PlanReadCommand) runService(client *client.Client) runReader {
	if c.runSvc != nil {
		return c.runSvc
	}
	return client.Runs
}

func (c *PlanReadCommand) planService(client *client.Client) planReader {
	if c.planSvc != nil {
		return c.planSvc
	}
	return client.Plans
}

// Help returns help text for the plan read command
func (c *PlanReadCommand) Help() string {
	helpText := `
Usage: hcptf plan read [options]

  Show plan details. You can provide either a plan ID or a run ID.
  If you provide a run ID, the command will automatically look up the
  associated plan.

Options:

  -id=<id>          Plan ID (plan-xxx) or Run ID (run-xxx) (required)
  -run-id=<id>      Run ID (alternative to -id)
  -output=<format>  Output format: table (default) or json

Examples:

  # Using plan ID
  hcptf plan read -id=plan-abc123

  # Using run ID
  hcptf plan read -id=run-xyz789
  hcptf plan read -run-id=run-xyz789

  # URL-style (via run show routing)
  hcptf my-org my-workspace runs run-xyz789 plan
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the plan read command
func (c *PlanReadCommand) Synopsis() string {
	return "Show plan details"
}
