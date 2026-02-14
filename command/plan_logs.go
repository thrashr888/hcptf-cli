package command

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PlanLogsCommand is a command to get plan logs
type PlanLogsCommand struct {
	Meta
	planID     string
	runID      string
	format     string
	planLogSvc planLogReader
}

// Run executes the plan logs command
func (c *PlanLogsCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("plan logs")
	flags.StringVar(&c.planID, "id", "", "Plan ID or Run ID")
	flags.StringVar(&c.runID, "run-id", "", "Run ID (alternative to -id)")
	flags.StringVar(&c.format, "output", "raw", "Output format: raw or json")

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
		run, err := client.Runs.Read(client.Context(), id)
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

	// Get plan logs
	logs, err := c.planLogService(client).Logs(client.Context(), planID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading plan logs: %s", err))
		return 1
	}

	// Read logs
	logData, err := io.ReadAll(logs)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading log data: %s", err))
		return 1
	}

	// Format output
	if c.format == "json" {
		output := map[string]interface{}{
			"plan_id": planID,
			"logs":    string(logData),
		}
		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error formatting JSON: %s", err))
			return 1
		}
		c.Ui.Output(string(jsonData))
	} else {
		c.Ui.Output(string(logData))
	}

	return 0
}

func (c *PlanLogsCommand) planLogService(client *client.Client) planLogReader {
	if c.planLogSvc != nil {
		return c.planLogSvc
	}
	return client.Plans
}

// Help returns help text for the plan logs command
func (c *PlanLogsCommand) Help() string {
	helpText := `
Usage: hcptf plan logs [options]

  Get plan logs. You can provide either a plan ID or a run ID.
  If you provide a run ID, the command will automatically look up the
  associated plan.

Options:

  -id=<id>          Plan ID (plan-xxx) or Run ID (run-xxx) (required)
  -run-id=<id>      Run ID (alternative to -id)
  -output=<format>  Output format: raw (default) or json

Examples:

  # Using plan ID
  hcptf plan logs -id=plan-abc123

  # Using run ID
  hcptf plan logs -id=run-xyz789
  hcptf plan logs -run-id=run-xyz789

  # URL-style
  hcptf my-org my-workspace runs run-xyz789 logs
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the plan logs command
func (c *PlanLogsCommand) Synopsis() string {
	return "Get plan logs"
}
