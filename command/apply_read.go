package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// ApplyReadCommand is a command to read apply details
type ApplyReadCommand struct {
	Meta
	applyID  string
	runID    string
	format   string
	applySvc applyReader
}

// Run executes the apply read command
func (c *ApplyReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("apply read")
	flags.StringVar(&c.applyID, "id", "", "Apply ID or Run ID")
	flags.StringVar(&c.runID, "run-id", "", "Run ID (alternative to -id)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags - need either applyID or runID
	id := c.applyID
	if id == "" {
		id = c.runID
	}
	if id == "" {
		c.Ui.Error("Error: -id or -run-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// If ID starts with "run-", get the apply ID from the run
	applyID := id
	if strings.HasPrefix(id, "run-") {
		run, err := client.Runs.Read(client.Context(), id)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading run: %s", err))
			return 1
		}
		if run.Apply == nil {
			c.Ui.Error("Error: run has no apply (may not have been applied yet)")
			return 1
		}
		applyID = run.Apply.ID
	}

	// Read apply
	apply, err := c.applyService(client).Read(client.Context(), applyID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading apply: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                   apply.ID,
		"Status":               apply.Status,
		"ResourceAdditions":    apply.ResourceAdditions,
		"ResourceChanges":      apply.ResourceChanges,
		"ResourceDestructions": apply.ResourceDestructions,
		"ResourceImports":      apply.ResourceImports,
	}

	if apply.StatusTimestamps != nil {
		if !apply.StatusTimestamps.QueuedAt.IsZero() {
			data["QueuedAt"] = apply.StatusTimestamps.QueuedAt
		}
		if !apply.StatusTimestamps.StartedAt.IsZero() {
			data["StartedAt"] = apply.StatusTimestamps.StartedAt
		}
		if !apply.StatusTimestamps.FinishedAt.IsZero() {
			data["FinishedAt"] = apply.StatusTimestamps.FinishedAt
		}
	}

	if apply.LogReadURL != "" {
		data["LogReadURL"] = apply.LogReadURL
	}

	formatter.KeyValue(data)
	return 0
}

func (c *ApplyReadCommand) applyService(client *client.Client) applyReader {
	if c.applySvc != nil {
		return c.applySvc
	}
	return client.Applies
}

// Help returns help text for the apply read command
func (c *ApplyReadCommand) Help() string {
	helpText := `
Usage: hcptf apply read [options]

  Show apply details. You can provide either an apply ID or a run ID.
  If you provide a run ID, the command will automatically look up the
  associated apply.

Options:

  -id=<id>          Apply ID (apply-xxx) or Run ID (run-xxx) (required)
  -run-id=<id>      Run ID (alternative to -id)
  -output=<format>  Output format: table (default) or json

Examples:

  # Using apply ID
  hcptf apply read -id=apply-abc123

  # Using run ID
  hcptf apply read -id=run-xyz789
  hcptf apply read -run-id=run-xyz789

  # URL-style
  hcptf my-org my-workspace runs run-xyz789 applyread
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the apply read command
func (c *ApplyReadCommand) Synopsis() string {
	return "Show apply details"
}
