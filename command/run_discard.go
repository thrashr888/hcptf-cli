package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunDiscardCommand is a command to discard a run
type RunDiscardCommand struct {
	Meta
	runID   string
	comment string
	runSvc  runDiscarder
}

// Run executes the run discard command
func (c *RunDiscardCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("run discard")
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

	// Discard run
	err = c.runService(client).Discard(client.Context(), c.runID, tfe.RunDiscardOptions{
		Comment: &c.comment,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error discarding run: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Run %s has been discarded", c.runID))
	return 0
}

func (c *RunDiscardCommand) runService(client *client.Client) runDiscarder {
	if c.runSvc != nil {
		return c.runSvc
	}
	return client.Runs
}

// Help returns help text for the run discard command
func (c *RunDiscardCommand) Help() string {
	helpText := `
Usage: hcptf run discard [options]

  Discard a run.

Options:

  -id=<run-id>      Run ID (required)
  -comment=<text>   Optional comment

Example:

  hcptf run discard -id=run-abc123
  hcptf run discard -id=run-abc123 -comment="No longer needed"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run discard command
func (c *RunDiscardCommand) Synopsis() string {
	return "Discard a run"
}
