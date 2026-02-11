package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunTriggerDeleteCommand is a command to delete a run trigger
type RunTriggerDeleteCommand struct {
	Meta
	id            string
	force         bool
	runTriggerSvc runTriggerDeleter
}

// Run executes the run trigger delete command
func (c *RunTriggerDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("runtrigger delete")
	flags.StringVar(&c.id, "id", "", "Run trigger ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.id == "" {
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

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete run trigger '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete run trigger
	err = c.runTriggerService(client).Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting run trigger: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Run trigger '%s' deleted successfully", c.id))
	return 0
}

func (c *RunTriggerDeleteCommand) runTriggerService(client *client.Client) runTriggerDeleter {
	if c.runTriggerSvc != nil {
		return c.runTriggerSvc
	}
	return client.RunTriggers
}

// Help returns help text for the run trigger delete command
func (c *RunTriggerDeleteCommand) Help() string {
	helpText := `
Usage: hcptf runtrigger delete [options]

  Delete a run trigger. This removes the automatic orchestration link
  between workspaces.

Options:

  -id=<id>  Run trigger ID (required)
  -force    Force delete without confirmation

Examples:

  # Delete a run trigger with confirmation
  hcptf runtrigger delete -id=rt-3yVQZvHzf5j3WRJ1

  # Force delete without confirmation
  hcptf runtrigger delete -id=rt-3yVQZvHzf5j3WRJ1 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run trigger delete command
func (c *RunTriggerDeleteCommand) Synopsis() string {
	return "Delete a run trigger"
}
