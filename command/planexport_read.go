package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PlanExportReadCommand is a command to read plan export details
type PlanExportReadCommand struct {
	Meta
	planExportID  string
	format        string
	planExportSvc planExportReader
}

// Run executes the planexport read command
func (c *PlanExportReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("planexport read")
	flags.StringVar(&c.planExportID, "id", "", "Plan export ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.planExportID == "" {
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

	// Read plan export
	planExport, err := c.planExportService(client).Read(client.Context(), c.planExportID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading plan export: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":       planExport.ID,
		"DataType": planExport.DataType,
		"Status":   planExport.Status,
	}

	if planExport.StatusTimestamps != nil {
		if !planExport.StatusTimestamps.QueuedAt.IsZero() {
			data["QueuedAt"] = planExport.StatusTimestamps.QueuedAt
		}
		if !planExport.StatusTimestamps.FinishedAt.IsZero() {
			data["FinishedAt"] = planExport.StatusTimestamps.FinishedAt
		}
		if !planExport.StatusTimestamps.ErroredAt.IsZero() {
			data["ErroredAt"] = planExport.StatusTimestamps.ErroredAt
		}
		if !planExport.StatusTimestamps.CanceledAt.IsZero() {
			data["CanceledAt"] = planExport.StatusTimestamps.CanceledAt
		}
		if !planExport.StatusTimestamps.ExpiredAt.IsZero() {
			data["ExpiredAt"] = planExport.StatusTimestamps.ExpiredAt
		}
	}

	formatter.KeyValue(data)
	return 0
}

func (c *PlanExportReadCommand) planExportService(client *client.Client) planExportReader {
	if c.planExportSvc != nil {
		return c.planExportSvc
	}
	return client.PlanExports
}

// Help returns help text for the planexport read command
func (c *PlanExportReadCommand) Help() string {
	helpText := `
Usage: hcptf planexport read [options]

  Show plan export details and status.

  The status field shows the current state:
  - queued: Export is waiting to be processed
  - pending: Export is being processed
  - finished: Export is complete and ready to download
  - errored: Export failed
  - canceled: Export was canceled
  - expired: Export data has expired (1 hour after completion)

Options:

  -id=<export-id>     Plan export ID (required)
  -output=<format>    Output format: table (default) or json

Example:

  hcptf planexport read -id=pe-abc123
  hcptf planexport read -id=pe-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the planexport read command
func (c *PlanExportReadCommand) Synopsis() string {
	return "Show plan export details and status"
}
