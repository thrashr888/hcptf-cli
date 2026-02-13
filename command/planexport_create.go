package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PlanExportCreateCommand is a command to create a plan export
type PlanExportCreateCommand struct {
	Meta
	planID        string
	dataType      string
	format        string
	planExportSvc planExportCreator
}

// Run executes the planexport create command
func (c *PlanExportCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("planexport create")
	flags.StringVar(&c.planID, "plan-id", "", "Plan ID (required)")
	flags.StringVar(&c.dataType, "data-type", "sentinel-mock-bundle-v0", "Data type for export (default: sentinel-mock-bundle-v0)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.planID == "" {
		c.Ui.Error("Error: -plan-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate data type
	var exportDataType tfe.PlanExportDataType
	switch c.dataType {
	case "sentinel-mock-bundle-v0":
		exportDataType = tfe.PlanExportSentinelMockBundleV0
	default:
		c.Ui.Error(fmt.Sprintf("Error: invalid data-type '%s'. Currently only 'sentinel-mock-bundle-v0' is supported", c.dataType))
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build create options
	options := tfe.PlanExportCreateOptions{
		Plan: &tfe.Plan{
			ID: c.planID,
		},
		DataType: &exportDataType,
	}

	// Create plan export
	planExport, err := c.planExportService(client).Create(client.Context(), options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating plan export: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Plan export '%s' created successfully", planExport.ID))

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
		if !planExport.StatusTimestamps.ExpiredAt.IsZero() {
			data["ExpiredAt"] = planExport.StatusTimestamps.ExpiredAt
		}
	}

	formatter.KeyValue(data)
	return 0
}

func (c *PlanExportCreateCommand) planExportService(client *client.Client) planExportCreator {
	if c.planExportSvc != nil {
		return c.planExportSvc
	}
	return client.PlanExports
}

// Help returns help text for the planexport create command
func (c *PlanExportCreateCommand) Help() string {
	helpText := `
Usage: hcptf planexport create [options]

  Create a plan export request to export plan data for analysis.
  The export process is asynchronous. Use 'planexport read' to check
  the export status, and 'planexport download' once finished.

Options:

  -plan-id=<plan-id>      Plan ID (required)
  -data-type=<type>       Data type for export (default: sentinel-mock-bundle-v0)
                          Currently only 'sentinel-mock-bundle-v0' is supported
  -output=<format>        Output format: table (default) or json

Example:

  hcptf planexport create -plan-id=plan-abc123
  hcptf planexport create -plan-id=plan-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the planexport create command
func (c *PlanExportCreateCommand) Synopsis() string {
	return "Create a plan export request"
}
