package command

import (
	"context"
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// WorkspaceReadCommand is a command to read workspace details
type WorkspaceReadCommand struct {
	Meta
	organization string
	name         string
	include      string
	format       string
	workspaceSvc workspaceReader
}

type workspaceReaderWithOptions interface {
	ReadWithOptions(ctx context.Context, organization, workspace string, options *tfe.WorkspaceReadOptions) (*tfe.Workspace, error)
}

// Run executes the workspace read command
func (c *WorkspaceReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspace read")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Workspace name (required)")
	flags.StringVar(&c.include, "include", "", "Comma-separated related resources to include (e.g. project,current_run)")
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

	workspaceSvc := c.workspaceService(client)
	var workspace *tfe.Workspace
	includes := parseWorkspaceReadIncludes(c.include)
	if !containsWorkspaceInclude(includes, tfe.WSProject) {
		includes = append(includes, tfe.WSProject)
	}

	if withOptions, ok := any(workspaceSvc).(workspaceReaderWithOptions); ok {
		workspace, err = withOptions.ReadWithOptions(client.Context(), c.organization, c.name, &tfe.WorkspaceReadOptions{
			Include: includes,
		})
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
			return 1
		}
	} else {
		workspace, err = workspaceSvc.Read(client.Context(), c.organization, c.name)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
			return 1
		}
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	projectID := ""
	projectName := ""
	if workspace.Project != nil {
		projectID = workspace.Project.ID
		projectName = workspace.Project.Name
	}

	data := map[string]interface{}{
		"ID":                   workspace.ID,
		"Name":                 workspace.Name,
		"Organization":         c.organization,
		"ProjectID":            projectID,
		"ProjectName":          projectName,
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

func parseWorkspaceReadIncludes(raw string) []tfe.WSIncludeOpt {
	if raw == "" {
		return nil
	}

	parts := splitCommaList(raw)
	includes := make([]tfe.WSIncludeOpt, 0, len(parts))
	for _, include := range parts {
		if include == "" {
			continue
		}
		includes = append(includes, tfe.WSIncludeOpt(include))
	}
	return includes
}

func containsWorkspaceInclude(includes []tfe.WSIncludeOpt, include tfe.WSIncludeOpt) bool {
	for _, current := range includes {
		if current == include {
			return true
		}
	}
	return false
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
  -include=<values>    Comma-separated include values (optional)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf workspace read -org=my-org -name=my-workspace
  hcptf workspace read -org=my-org -name=my-workspace -include=project,current_run
  hcptf workspace read -org=my-org -name=prod -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace read command
func (c *WorkspaceReadCommand) Synopsis() string {
	return "Read workspace details"
}
