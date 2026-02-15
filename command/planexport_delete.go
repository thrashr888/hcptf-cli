package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PlanExportDeleteCommand deletes a plan export.
type PlanExportDeleteCommand struct {
	Meta
	planExportID  string
	force         bool
	yes           bool
	planExportSvc planExportDeleter
}

// Run executes the planexport delete command.
func (c *PlanExportDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("planexport delete")
	flags.StringVar(&c.planExportID, "id", "", "Plan export ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")
	flags.BoolVar(&c.yes, "y", false, "Confirm delete without prompt")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.planExportID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	apiClient, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	if !c.force && !c.yes {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete plan export '%s'? (yes/no): ", c.planExportID))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}
		if strings.ToLower(strings.TrimSpace(confirmation)) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	if err := c.planExportService(apiClient).Delete(apiClient.Context(), c.planExportID); err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting plan export: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Plan export '%s' deleted successfully", c.planExportID))
	return 0
}

func (c *PlanExportDeleteCommand) planExportService(client *client.Client) planExportDeleter {
	if c.planExportSvc != nil {
		return c.planExportSvc
	}
	return client.PlanExports
}

// Help returns help text for the planexport delete command.
func (c *PlanExportDeleteCommand) Help() string {
	helpText := `
Usage: hcptf planexport delete [options]

  Delete a plan export.

Options:

  -id=<export-id>     Plan export ID (required)
  -force              Force delete without confirmation
  -y                  Confirm delete without prompt

Example:

  hcptf planexport delete -id=pe-abc123
  hcptf planexport delete -id=pe-abc123 -force
  hcptf planexport delete -id=pe-abc123 -y
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the planexport delete command.
func (c *PlanExportDeleteCommand) Synopsis() string {
	return "Delete a plan export"
}
