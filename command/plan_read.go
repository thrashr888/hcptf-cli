package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// PlanReadCommand is a command to read plan details
type PlanReadCommand struct {
	Meta
	planID  string
	format  string
	planSvc planReader
}

// Run executes the plan read command
func (c *PlanReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("plan read")
	flags.StringVar(&c.planID, "id", "", "Plan ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.planID == "" {
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

	// Read plan
	plan, err := c.planService(client).Read(client.Context(), c.planID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading plan: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

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

  Show plan details.

Options:

  -id=<plan-id>     Plan ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf plan read -id=plan-abc123
  hcptf plan read -id=plan-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the plan read command
func (c *PlanReadCommand) Synopsis() string {
	return "Show plan details"
}
