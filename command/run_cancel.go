package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunCancelCommand is a command to cancel a run
type RunCancelCommand struct {
	Meta
	runID   string
	comment string
	force   bool
	runSvc  runCanceler
}

// Run executes the run cancel command
func (c *RunCancelCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("run cancel")
	flags.StringVar(&c.runID, "id", "", "Run ID (required)")
	flags.StringVar(&c.comment, "comment", "", "Optional comment")
	flags.BoolVar(&c.force, "force", false, "Force cancel the run")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.runID == "" {
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

	// Cancel run
	if c.force {
		err = c.runService(client).ForceCancel(client.Context(), c.runID, tfe.RunForceCancelOptions{
			Comment: &c.comment,
		})
	} else {
		err = c.runService(client).Cancel(client.Context(), c.runID, tfe.RunCancelOptions{
			Comment: &c.comment,
		})
	}

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error canceling run: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Run %s has been canceled", c.runID))
	return 0
}

func (c *RunCancelCommand) runService(client *client.Client) runCanceler {
	if c.runSvc != nil {
		return c.runSvc
	}
	return client.Runs
}

// Help returns help text for the run cancel command
func (c *RunCancelCommand) Help() string {
	helpText := `
Usage: hcptf workspace run cancel [options]

  Cancel a run.

Options:

  -id=<run-id>      Run ID (required)
  -comment=<text>   Optional comment
  -force            Force cancel the run

Example:

  hcptf workspace run cancel -id=run-abc123
  hcptf workspace run cancel -id=run-abc123 -force -comment="Emergency stop"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run cancel command
func (c *RunCancelCommand) Synopsis() string {
	return "Cancel a run"
}
