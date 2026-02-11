package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunApplyCommand is a command to apply a run
type RunApplyCommand struct {
	Meta
	runID   string
	comment string
	runSvc  runApplier
}

// Run executes the run apply command
func (c *RunApplyCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("run apply")
	flags.StringVar(&c.runID, "id", "", "Run ID (required)")
	flags.StringVar(&c.comment, "comment", "", "Optional comment")

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

	// Apply run
	err = c.runService(client).Apply(client.Context(), c.runID, tfe.RunApplyOptions{
		Comment: &c.comment,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error applying run: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Run %s has been approved and is applying", c.runID))
	return 0
}

func (c *RunApplyCommand) runService(client *client.Client) runApplier {
	if c.runSvc != nil {
		return c.runSvc
	}
	return client.Runs
}

// Help returns help text for the run apply command
func (c *RunApplyCommand) Help() string {
	helpText := `
Usage: hcptf run apply [options]

  Approve and apply a run.

Options:

  -id=<run-id>      Run ID (required)
  -comment=<text>   Optional comment

Example:

  hcptf run apply -id=run-abc123
  hcptf run apply -id=run-abc123 -comment="Approved for deployment"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run apply command
func (c *RunApplyCommand) Synopsis() string {
	return "Approve and apply a run"
}
