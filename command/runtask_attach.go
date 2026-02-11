package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// RunTaskAttachCommand is a command to attach a run task to a workspace
type RunTaskAttachCommand struct {
	Meta
	organization     string
	workspace        string
	runTaskID        string
	enforcementLevel string
	stage            string
	format           string
}

// Run executes the run task attach command
func (c *RunTaskAttachCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("runtask attach")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
	flags.StringVar(&c.runTaskID, "runtask-id", "", "Run task ID (required)")
	flags.StringVar(&c.enforcementLevel, "enforcement-level", "advisory", "Enforcement level: advisory or mandatory (default: advisory)")
	flags.StringVar(&c.stage, "stage", "post_plan", "Stage: post_plan, pre_plan, or pre_apply (default: post_plan)")
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

	if c.workspace == "" {
		c.Ui.Error("Error: -workspace flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.runTaskID == "" {
		c.Ui.Error("Error: -runtask-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate enforcement level
	var enforcementLevel tfe.TaskEnforcementLevel
	switch c.enforcementLevel {
	case "advisory":
		enforcementLevel = tfe.Advisory
	case "mandatory":
		enforcementLevel = tfe.Mandatory
	default:
		c.Ui.Error("Error: -enforcement-level must be 'advisory' or 'mandatory'")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate stage
	var stage tfe.Stage
	switch c.stage {
	case "post_plan":
		stage = tfe.PostPlan
	case "pre_plan":
		stage = tfe.PrePlan
	case "pre_apply":
		stage = tfe.PreApply
	default:
		c.Ui.Error("Error: -stage must be 'post_plan', 'pre_plan', or 'pre_apply'")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Get workspace to obtain its ID
	workspace, err := client.Workspaces.Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	// Attach run task to workspace
	options := tfe.WorkspaceRunTaskCreateOptions{
		EnforcementLevel: enforcementLevel,
		RunTask:          &tfe.RunTask{ID: c.runTaskID},
		Stage:            &stage,
	}

	workspaceRunTask, err := client.WorkspaceRunTasks.Create(client.Context(), workspace.ID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error attaching run task to workspace: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Run task attached to workspace '%s' successfully", c.workspace))

	// Show workspace run task details
	data := map[string]interface{}{
		"ID":               workspaceRunTask.ID,
		"EnforcementLevel": workspaceRunTask.EnforcementLevel,
		"Stage":            workspaceRunTask.Stage,
	}

	if workspaceRunTask.RunTask != nil {
		data["RunTaskID"] = workspaceRunTask.RunTask.ID
		data["RunTaskName"] = workspaceRunTask.RunTask.Name
	}

	if workspaceRunTask.Workspace != nil {
		data["WorkspaceID"] = workspaceRunTask.Workspace.ID
		data["WorkspaceName"] = workspaceRunTask.Workspace.Name
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the run task attach command
func (c *RunTaskAttachCommand) Help() string {
	helpText := `
Usage: hcptf runtask attach [options]

  Attach a run task to a workspace. This creates a workspace-task association
  that configures how and when the run task executes during a Terraform run.

Options:

  -organization=<name>       Organization name (required)
  -org=<name>               Alias for -organization
  -workspace=<name>         Workspace name (required)
  -runtask-id=<id>          Run task ID (required)
  -enforcement-level=<lvl>  Enforcement level (default: advisory)
                            advisory: Results are informational, runs continue
                            mandatory: Failed checks block run progression
  -stage=<stage>            Run stage (default: post_plan)
                            post_plan: After plan, before apply
                            pre_plan: Before plan operation
                            pre_apply: Before apply operation
  -output=<format>          Output format: table (default) or json

Example:

  hcptf runtask attach -org=my-org -workspace=prod \
    -runtask-id=task-ABC123 -enforcement-level=mandatory -stage=post_plan

  hcptf runtask attach -org=my-org -workspace=dev \
    -runtask-id=task-XYZ789 -enforcement-level=advisory -stage=pre_plan
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run task attach command
func (c *RunTaskAttachCommand) Synopsis() string {
	return "Attach a run task to a workspace"
}
