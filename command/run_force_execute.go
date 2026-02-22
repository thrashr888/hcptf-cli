package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunForceExecuteCommand is a command to force-execute a run.
type RunForceExecuteCommand struct {
	Meta
	runID  string
	runSvc runForceExecutor
}

// Run executes the run force-execute command.
func (c *RunForceExecuteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("run force-execute")
	flags.StringVar(&c.runID, "id", "", "Run ID (required)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.runID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	if err := c.runService(client).ForceExecute(client.Context(), c.runID); err != nil {
		c.Ui.Error(fmt.Sprintf("Error force-executing run: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Run %s has been force-executed", c.runID))
	return 0
}

func (c *RunForceExecuteCommand) runService(client *client.Client) runForceExecutor {
	if c.runSvc != nil {
		return c.runSvc
	}
	return client.Runs
}

// Help returns help text for the run force-execute command.
func (c *RunForceExecuteCommand) Help() string {
	helpText := `
Usage: hcptf run force-execute [options]

  Force-execute a run.

Options:

  -id=<run-id>  Run ID (required)

Example:

  hcptf run force-execute -id=run-abc123
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run force-execute command.
func (c *RunForceExecuteCommand) Synopsis() string {
	return "Force-execute a run"
}
