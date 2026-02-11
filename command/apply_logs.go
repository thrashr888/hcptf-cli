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
	format      string
	applyLogSvc applyLogReader
}

// Run executes the apply logs command
func (c *ApplyLogsCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("apply logs")
	flags.StringVar(&c.applyID, "id", "", "Apply ID (required)")
	flags.StringVar(&c.format, "output", "raw", "Output format: raw or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.applyID == "" {
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

	// Get apply logs
	logs, err := c.applyLogService(client).Logs(client.Context(), c.applyID)
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
			"apply_id": c.applyID,
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

  Get apply logs (JSON or raw output).

Options:

  -id=<apply-id>    Apply ID (required)
  -output=<format>  Output format: raw (default) or json

Example:

  hcptf apply logs -id=apply-abc123
  hcptf apply logs -id=apply-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the apply logs command
func (c *ApplyLogsCommand) Synopsis() string {
	return "Get apply logs"
}
