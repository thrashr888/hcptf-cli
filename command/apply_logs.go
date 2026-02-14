package command

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// ApplyLogsCommand is a command to get apply logs
type ApplyLogsCommand struct {
	Meta
	applyID     string
	runID       string
	format      string
	applyLogSvc applyLogReader
}

// Run executes the apply logs command
func (c *ApplyLogsCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("apply logs")
	flags.StringVar(&c.applyID, "id", "", "Apply ID or Run ID")
	flags.StringVar(&c.runID, "run-id", "", "Run ID (alternative to -id)")
	flags.StringVar(&c.format, "output", "raw", "Output format: raw or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags - need either applyID or runID
	id := c.applyID
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

	// If ID starts with "run-", get the apply ID from the run
	applyID := id
	if strings.HasPrefix(id, "run-") {
		run, err := client.Runs.Read(client.Context(), id)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading run: %s", err))
			return 1
		}
		if run.Apply == nil {
			c.Ui.Error("Error: run has no apply (may not have been applied yet)")
			return 1
		}
		applyID = run.Apply.ID
	}

	// Get apply logs
	logs, err := c.applyLogService(client).Logs(client.Context(), applyID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading apply logs: %s", err))
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

func (c *ApplyLogsCommand) applyLogService(client *client.Client) applyLogReader {
	if c.applyLogSvc != nil {
		return c.applyLogSvc
	}
	return client.Applies
}

// Help returns help text for the apply logs command
func (c *ApplyLogsCommand) Help() string {
	helpText := `
Usage: hcptf apply logs [options]

  Get apply logs. You can provide either an apply ID or a run ID.
  If you provide a run ID, the command will automatically look up the
  associated apply.

Options:

  -id=<id>          Apply ID (apply-xxx) or Run ID (run-xxx) (required)
  -run-id=<id>      Run ID (alternative to -id)
  -output=<format>  Output format: raw (default) or json

Examples:

  # Using apply ID
  hcptf apply logs -id=apply-abc123

  # Using run ID
  hcptf apply logs -id=run-xyz789
  hcptf apply logs -run-id=run-xyz789

  # URL-style
  hcptf my-org my-workspace runs run-xyz789 apply logs
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the apply logs command
func (c *ApplyLogsCommand) Synopsis() string {
	return "Get apply logs"
}
