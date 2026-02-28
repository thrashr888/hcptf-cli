package command

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunLogsCommand is a command to get plan or apply logs by run ID
type RunLogsCommand struct {
	Meta
	runID       string
	phase       string
	format      string
	runSvc      runReader
	planLogSvc  planLogReader
	applyLogSvc applyLogReader
}

// Run executes the run logs command
func (c *RunLogsCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("run logs")
	flags.StringVar(&c.runID, "id", "", "Run ID (required)")
	flags.StringVar(&c.phase, "phase", "auto", "Phase to show logs for: plan, apply, or auto (default: auto)")
	flags.StringVar(&c.format, "output", "raw", "Output format: raw or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.runID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate phase flag
	switch c.phase {
	case "auto", "plan", "apply":
		// valid
	default:
		c.Ui.Error(fmt.Sprintf("Error: invalid -phase value %q, must be plan, apply, or auto", c.phase))
		return 1
	}

	// Get API client
	cl, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read the run to determine phase and get plan/apply IDs
	run, err := c.runService(cl).Read(cl.Context(), c.runID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading run: %s", err))
		return 1
	}

	// Determine which phase to show
	phase := c.phase
	if phase == "auto" {
		if run.Apply != nil {
			phase = "apply"
		} else {
			phase = "plan"
		}
	}

	switch phase {
	case "apply":
		if run.Apply == nil {
			c.Ui.Error("Error: run has no apply (may not have been applied yet)")
			return 1
		}
		return c.fetchApplyLogs(cl, run.Apply.ID)
	default:
		if run.Plan == nil {
			c.Ui.Error("Error: run has no plan")
			return 1
		}
		return c.fetchPlanLogs(cl, run.Plan.ID)
	}
}

func (c *RunLogsCommand) fetchPlanLogs(cl *client.Client, planID string) int {
	logs, err := c.planLogService(cl).Logs(cl.Context(), planID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading plan logs: %s", err))
		return 1
	}

	logData, err := io.ReadAll(logs)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading log data: %s", err))
		return 1
	}

	if c.format == "json" {
		output := map[string]interface{}{
			"run_id":  c.runID,
			"phase":   "plan",
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

func (c *RunLogsCommand) fetchApplyLogs(cl *client.Client, applyID string) int {
	logs, err := c.applyLogService(cl).Logs(cl.Context(), applyID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading apply logs: %s", err))
		return 1
	}

	logData, err := io.ReadAll(logs)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading log data: %s", err))
		return 1
	}

	if c.format == "json" {
		output := map[string]interface{}{
			"run_id":   c.runID,
			"phase":    "apply",
			"apply_id": applyID,
			"logs":     string(logData),
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

func (c *RunLogsCommand) runService(cl *client.Client) runReader {
	if c.runSvc != nil {
		return c.runSvc
	}
	return cl.Runs
}

func (c *RunLogsCommand) planLogService(cl *client.Client) planLogReader {
	if c.planLogSvc != nil {
		return c.planLogSvc
	}
	return cl.Plans
}

func (c *RunLogsCommand) applyLogService(cl *client.Client) applyLogReader {
	if c.applyLogSvc != nil {
		return c.applyLogSvc
	}
	return cl.Applies
}

// Help returns help text for the run logs command
func (c *RunLogsCommand) Help() string {
	helpText := `
Usage: hcptf run logs [options]

  Get plan or apply logs for a run. By default, auto-detects the appropriate
  phase: if the run has an apply, shows apply logs; otherwise shows plan logs.
  Use -phase to explicitly select plan or apply logs.

Options:

  -id=<run-id>        Run ID (required)
  -phase=<phase>      Phase to show: plan, apply, or auto (default: auto)
  -output=<format>    Output format: raw (default) or json

Examples:

  # Auto-detect phase (apply if available, otherwise plan)
  hcptf run logs -id=run-abc123

  # Explicitly show plan logs
  hcptf run logs -id=run-abc123 -phase=plan

  # Explicitly show apply logs
  hcptf run logs -id=run-abc123 -phase=apply

  # JSON output
  hcptf run logs -id=run-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run logs command
func (c *RunLogsCommand) Synopsis() string {
	return "Get plan or apply logs for a run"
}
