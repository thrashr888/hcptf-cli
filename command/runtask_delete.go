package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunTaskDeleteCommand is a command to delete a run task
type RunTaskDeleteCommand struct {
	Meta
	id         string
	force      bool
	runTaskSvc runTaskDeleterReader
}

// Run executes the run task delete command
func (c *RunTaskDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("runtask delete")
	flags.StringVar(&c.id, "id", "", "Run task ID (required)")
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

	// Read run task to get its name for confirmation
	runTask, err := c.runTaskService(client).Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading run task: %s", err))
		return 1
	}

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete run task '%s' (%s)? (yes/no): ", runTask.Name, c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete run task
	err = c.runTaskService(client).Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting run task: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Run task '%s' (%s) deleted successfully", runTask.Name, c.id))
	return 0
}

func (c *RunTaskDeleteCommand) runTaskService(client *client.Client) runTaskDeleterReader {
	if c.runTaskSvc != nil {
		return c.runTaskSvc
	}
	return client.RunTasks
}

// Help returns help text for the run task delete command
func (c *RunTaskDeleteCommand) Help() string {
	helpText := `
Usage: hcptf runtask delete [options]

  Delete a run task.

  Warning: Deleting a run task will remove it from all workspaces it is
  attached to. This action cannot be undone.

Options:

  -id=<id>  Run task ID (required)
  -force    Force delete without confirmation

Example:

  hcptf runtask delete -id=task-ABC123
  hcptf runtask delete -id=task-ABC123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run task delete command
func (c *RunTaskDeleteCommand) Synopsis() string {
	return "Delete a run task"
}
