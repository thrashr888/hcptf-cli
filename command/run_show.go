package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// RunShowCommand is a command to show run details
type RunShowCommand struct {
	Meta
	runID  string
	format string
	runSvc runReader
}

// Run executes the run show command
func (c *RunShowCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("run show")
	flags.StringVar(&c.runID, "id", "", "Run ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

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

	// Read run
	run, err := c.runService(client).Read(client.Context(), c.runID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading run: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                   run.ID,
		"Status":               run.Status,
		"Message":              run.Message,
		"IsDestroy":            run.IsDestroy,
		"Source":               run.Source,
		"AutoApply":            run.AutoApply,
		"HasChanges":           run.HasChanges,
		"CreatedAt":            run.CreatedAt,
		"StatusTimestamps":     run.StatusTimestamps,
		"TargetAddrs":          run.TargetAddrs,
		"ReplaceAddrs":         run.ReplaceAddrs,
		"RefreshOnly":          run.RefreshOnly,
		"AllowEmptyApply":      run.AllowEmptyApply,
		"PlanOnly":             run.PlanOnly,
		"TerraformVersion":     run.TerraformVersion,
		"PositionInQueue":      run.PositionInQueue,
		"ResourceAdditions":    run.Plan.ResourceAdditions,
		"ResourceChanges":      run.Plan.ResourceChanges,
		"ResourceDestructions": run.Plan.ResourceDestructions,
	}

	if run.ConfigurationVersion != nil {
		data["ConfigurationVersionID"] = run.ConfigurationVersion.ID
	}

	if run.Workspace != nil {
		data["WorkspaceID"] = run.Workspace.ID
		data["WorkspaceName"] = run.Workspace.Name
	}

	formatter.KeyValue(data)
	return 0
}

func (c *RunShowCommand) runService(client *client.Client) runReader {
	if c.runSvc != nil {
		return c.runSvc
	}
	return client.Runs
}

// Help returns help text for the run show command
func (c *RunShowCommand) Help() string {
	helpText := `
Usage: hcptf run show [options]

  Show run details.

Options:

  -id=<run-id>      Run ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf run show -id=run-abc123
  hcptf run show -id=run-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run show command
func (c *RunShowCommand) Synopsis() string {
	return "Show run details"
}
