package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// ApplyReadCommand is a command to read apply details
type ApplyReadCommand struct {
	Meta
	applyID  string
	format   string
	applySvc applyReader
}

// Run executes the apply read command
func (c *ApplyReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("apply read")
	flags.StringVar(&c.applyID, "id", "", "Apply ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.applyID == "" {
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

	// Read apply
	apply, err := c.applyService(client).Read(client.Context(), c.applyID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading apply: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

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

  Show apply details.

Options:

  -id=<apply-id>    Apply ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf apply read -id=apply-abc123
  hcptf apply read -id=apply-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the apply read command
func (c *ApplyReadCommand) Synopsis() string {
	return "Show apply details"
}
