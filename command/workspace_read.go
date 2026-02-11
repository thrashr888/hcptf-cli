package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// WorkspaceReadCommand is a command to read workspace details
type WorkspaceReadCommand struct {
	Meta
	organization string
	name         string
	format       string
	workspaceSvc workspaceReader
}

// Run executes the workspace read command
func (c *WorkspaceReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspace read")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Workspace name (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read workspace
	workspace, err := c.workspaceService(client).Read(client.Context(), c.organization, c.name)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                   workspace.ID,
		"Name":                 workspace.Name,
		"Organization":         c.organization,
		"TerraformVersion":     workspace.TerraformVersion,
		"AutoApply":            workspace.AutoApply,
		"AllowDestroyPlan":     workspace.AllowDestroyPlan,
		"Description":          workspace.Description,
		"Environment":          workspace.Environment,
		"ExecutionMode":        workspace.ExecutionMode,
		"FileTriggersEnabled":  workspace.FileTriggersEnabled,
		"GlobalRemoteState":    workspace.GlobalRemoteState,
		"Locked":               workspace.Locked,
		"QueueAllRuns":         workspace.QueueAllRuns,
		"SpeculativeEnabled":   workspace.SpeculativeEnabled,
		"SourceName":           workspace.SourceName,
		"SourceURL":            workspace.SourceURL,
		"TriggerPrefixes":      workspace.TriggerPrefixes,
		"WorkingDirectory":     workspace.WorkingDirectory,
		"CreatedAt":            workspace.CreatedAt,
		"UpdatedAt":            workspace.UpdatedAt,
		"ResourceCount":        workspace.ResourceCount,
		"ApplyDurationAverage": workspace.ApplyDurationAverage,
		"PlanDurationAverage":  workspace.PlanDurationAverage,
		"RunsCount":            workspace.RunsCount,
	}

	if workspace.CurrentRun != nil {
		data["CurrentRunID"] = workspace.CurrentRun.ID
		data["CurrentRunStatus"] = workspace.CurrentRun.Status
	}

	formatter.KeyValue(data)
	return 0
}

func (c *WorkspaceReadCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

// Help returns help text for the workspace read command
func (c *WorkspaceReadCommand) Help() string {
	helpText := `
Usage: hcptf workspace read [options]

  Read workspace details.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Workspace name (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf workspace read -org=my-org -name=my-workspace
  hcptf workspace read -org=my-org -name=prod -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace read command
func (c *WorkspaceReadCommand) Synopsis() string {
	return "Read workspace details"
}
